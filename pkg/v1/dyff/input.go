// Copyright © 2018 Matthias Diester
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package dyff

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/HeavyWombat/dyff/pkg/v1/bunt"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"

	ordered "github.com/virtuald/go-ordered-json"
)

// PreserveKeyOrderInJSON specifies whether a special library is used to decode
// JSON input to preserve the order of keys in maps even though that is not part
// of the JSON specification.
var PreserveKeyOrderInJSON = false

// DecoderProxy can either be used with the standard JSON Decoder, or the
// specialised JSON library fork that supports preserving key order
type DecoderProxy struct {
	standard *json.Decoder
	ordered  *ordered.Decoder
}

// InputFile represents the actual input file (either local, or fetched remotely) that needs to be processed. It can contain multiple documents, where a document is a map or a list of things.
type InputFile struct {
	Location  string
	Note      string
	Documents []interface{}
}

// NewDecoderProxy creates a new decoder proxy which either works in ordered
// mode or standard mode.
func NewDecoderProxy(keepOrder bool, r io.Reader) *DecoderProxy {
	if keepOrder {
		decoder := ordered.NewDecoder(r)
		decoder.UseOrderedObject()
		return &DecoderProxy{ordered: decoder}
	}

	return &DecoderProxy{standard: json.NewDecoder(r)}
}

// Decode is a delegate function that calls JSON Decoder `Decode`
func (d *DecoderProxy) Decode(v interface{}) error {
	if d.ordered != nil {
		return d.ordered.Decode(v)
	}

	return d.standard.Decode(v)
}

// HumanReadableLocationInformation create a nicely decorated information about the provided input location. It will output the absolut path of the file (rather than the possibly relative location), or it will show the URL in the usual look-and-feel of URIs.
func HumanReadableLocationInformation(inputFile InputFile) string {
	var buf bytes.Buffer

	// Start with a nice location output
	buf.WriteString(HumanReadableLocation(inputFile.Location))

	// Add additional note if it is set
	if inputFile.Note != "" {
		buf.WriteString(", ")
		buf.WriteString(bunt.Colorize(inputFile.Note, bunt.Orange))
	}

	// Add an information about how many documents are in the provided input file
	if documents := len(inputFile.Documents); documents > 1 {
		buf.WriteString(", ")
		buf.WriteString(bunt.Colorize(Plural(documents, "document"), bunt.Aquamarine, bunt.Bold))
	}

	return buf.String()
}

// HumanReadableLocation returns a human readable location with proper coloring
func HumanReadableLocation(location string) string {
	var buf bytes.Buffer

	if IsStdin(location) {
		buf.WriteString(bunt.Style("<STDIN>", bunt.Italic))

	} else if _, err := os.Stat(location); err == nil {
		if abs, err := filepath.Abs(location); err == nil {
			buf.WriteString(bunt.Style(abs, bunt.Bold))
		} else {
			buf.WriteString(bunt.Style(location, bunt.Bold))
		}

	} else if _, err := url.ParseRequestURI(location); err == nil {
		buf.WriteString(bunt.Colorize(location, bunt.CornflowerBlue, bunt.Underline))
	}

	return buf.String()
}

// LoadFiles concurrently loads two files from the provided locations
func LoadFiles(locationA string, locationB string) (InputFile, InputFile, error) {
	type resultPair struct {
		result InputFile
		err    error
	}

	fromChan := make(chan resultPair, 1)
	toChan := make(chan resultPair, 1)

	go func() {
		result, err := LoadFile(locationA)
		fromChan <- resultPair{result, err}
	}()

	go func() {
		result, err := LoadFile(locationB)
		toChan <- resultPair{result, err}
	}()

	from := <-fromChan
	if from.err != nil {
		return InputFile{}, InputFile{}, from.err
	}

	to := <-toChan
	if to.err != nil {
		return InputFile{}, InputFile{}, to.err
	}

	return from.result, to.result, nil
}

// LoadFile processes the provided input location to load a YAML (or JSON, or raw text)
func LoadFile(location string) (InputFile, error) {
	start := time.Now()

	var (
		documents []interface{}
		data      []byte
		err       error
	)

	if data, err = getBytesFromLocation(location); err != nil {
		return InputFile{}, errors.Wrap(err, fmt.Sprintf("Unable to load data from %s", location))
	}

	if documents, err = LoadDocuments(data); err != nil {
		return InputFile{}, errors.Wrap(err, fmt.Sprintf("Unable to parse data from %s", location))
	}

	DebugLogger.Printf("Loaded %s (%d byte) with %s in %s", location, len(data), Plural(len(documents), "document"), time.Since(start))
	return InputFile{Location: location, Documents: documents}, nil
}

// LoadDocuments reads the provided input data slice as a YAML or JSON file with
// potential multiple documents. It only acts as a dispatcher and depending on
// the input will either use `LoadJSONDocuments` or `LoadYAMLDocuments`.
func LoadDocuments(input []byte) ([]interface{}, error) {
	switch input[0] {
	case '{', '[':
		return LoadJSONDocuments(input)

	default:
		return LoadYAMLDocuments(input)
	}
}

// LoadJSONDocuments reads the provided input data slice as a YAML file with
// potential multiple documents. Each document in the JSON stream results in an
// entry of the result slice. This function performs two decoding passes over
// the input data slice, the first one to detect the respective types in use.
// And a second one to properly unmarshal the data in the most suitable Go types
// available. JSON does not support key orders in maps.
func LoadJSONDocuments(input []byte) ([]interface{}, error) {
	var (
		types   []string
		values  []interface{}
		decoder *DecoderProxy
	)

	// First pass: decode all documents and save the actual types
	types = make([]string, 0)
	decoder = NewDecoderProxy(false, bytes.NewReader(input))
	for {
		var value interface{}

		if err := decoder.Decode(&value); err == io.EOF {
			break

		} else if err != nil {
			return nil, err
		}

		types = append(types, getType(value))
	}

	DebugLogger.Printf("load JSON input, first pass complete, types detected: %s", strings.Join(types, ", "))

	// Second pass: Based on the types, initialise a proper variable to unmarshal data into
	values = make([]interface{}, len(types))
	decoder = NewDecoderProxy(PreserveKeyOrderInJSON, bytes.NewReader(input))
	for i := 0; i < len(types); i++ {
		switch types[i] {
		case typeMap:
			var value interface{}
			decoder.Decode(&value)
			values[i] = mapSlicify(value)

		case typeSimpleList, typeComplexList:
			var value []interface{}
			decoder.Decode(&value)
			values[i] = mapSlicify(value)

		case typeString:
			var value string
			decoder.Decode(&value)
			values[i] = value

		default:
			return nil, fmt.Errorf("Unsupported type %s in load document function", types[i])
		}
	}

	return values, nil
}

// LoadYAMLDocuments reads the provided input data slice as a YAML file with
// potential multiple documents. Each document in the YAML stream results in an
// entry of the result slice. This function performs two decoding passes over
// the input data slice, the first one to detect the respective types in use.
// And a second one to properly unmarshal the data in the most suitable Go types
// available so that key orders in hashes are preserved.
func LoadYAMLDocuments(input []byte) ([]interface{}, error) {
	var (
		types   []string
		values  []interface{}
		decoder *yaml.Decoder
	)

	// First pass: decode all documents and save the actual types
	types = make([]string, 0)
	decoder = yaml.NewDecoder(bytes.NewReader(input))
	for {
		var value interface{}

		if err := decoder.Decode(&value); err == io.EOF {
			break

		} else if err != nil {
			return nil, err
		}

		types = append(types, getType(value))
	}

	DebugLogger.Printf("load YAML input, first pass complete, types detected: %s", strings.Join(types, ", "))

	// Second pass: Based on the types, initialise a proper variable to unmarshal data into
	values = make([]interface{}, len(types))
	decoder = yaml.NewDecoder(bytes.NewReader(input))
	for i := 0; i < len(types); i++ {
		switch types[i] {
		case typeMap:
			var value yaml.MapSlice
			decoder.Decode(&value)
			values[i] = value

		case typeSimpleList:
			var value []interface{}
			decoder.Decode(&value)
			values[i] = value

		case typeComplexList:
			var value []yaml.MapSlice
			decoder.Decode(&value)
			values[i] = value

		case typeString:
			var value string
			decoder.Decode(&value)
			values[i] = value

		default:
			return nil, fmt.Errorf("Unsupported type %s in load document function", types[i])
		}
	}

	return values, nil
}

// mapSlicify makes sure that each occurrence of a map in the provided structure
// is changed to a YAML MapSlice.
//
// Please note: In case the input data were decoded by the default standard JSON
// parser, there will be no preservation of the order of keys, because JSON does
// not support such thing as an order of keys. Therfore, the keys are sorted to
// have a consistent and testable output structure.
//
// This function supports `OrderedObjects` from the JSON library fork
// `github.com/virtuald/go-ordered-json` and will translate this structure into
// the compatible YAML structure.
func mapSlicify(obj interface{}) interface{} {
	switch obj.(type) {
	case ordered.OrderedObject:
		orderedObj := obj.(ordered.OrderedObject)
		result := make(yaml.MapSlice, 0, len(orderedObj))
		for _, member := range orderedObj {
			result = append(result, yaml.MapItem{Key: member.Key, Value: mapSlicify(member.Value)})
		}

		return result

	case map[string]interface{}:
		hash := obj.(map[string]interface{})
		keys := make([]string, 0, len(hash))
		for key := range hash {
			keys = append(keys, key)
		}

		sort.Strings(keys)

		result := make(yaml.MapSlice, 0, len(hash))
		for _, key := range keys {
			result = append(result, yaml.MapItem{Key: key, Value: mapSlicify(hash[key])})
		}

		return result

	case []interface{}:
		list := obj.([]interface{})
		result := make([]interface{}, len(list))
		for idx, entry := range list {
			result[idx] = mapSlicify(entry)
		}

		return result

	default:
		return obj
	}
}

func getBytesFromLocation(location string) ([]byte, error) {
	var data []byte
	var err error

	// Handle special location "-" which referes to STDIN stream
	if IsStdin(location) {
		if data, err = ioutil.ReadAll(os.Stdin); err != nil {
			return nil, err
		}

		return data, nil
	}

	// Handle location as local file if there is a file at that location
	if _, err = os.Stat(location); err == nil {
		if data, err = ioutil.ReadFile(location); err != nil {
			return nil, err
		}

		return data, nil
	}

	// Handle location as a URI if it looks like one
	if _, err = url.ParseRequestURI(location); err == nil {
		var response *http.Response
		response, err = http.Get(location)
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		buf := new(bytes.Buffer)
		buf.ReadFrom(response.Body)

		if response.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to load from location: %s", buf.Bytes())
		}

		return buf.Bytes(), nil
	}

	// In any other case, bail out ...
	return nil, fmt.Errorf("Unable to get any content using location %s: it is not a file or usable URI", location)
}

func isComplexSlice(slice []interface{}) bool {
	// This is kind of a weird case, but by definition an empty list is a simple slice
	if len(slice) == 0 {
		return false
	}

	// Count the number of entries which are maps or YAML MapSlices
	counter := 0
	for _, entry := range slice {
		switch entry.(type) {
		case map[string]interface{}, map[interface{}]interface{}, yaml.MapSlice:
			counter++
		}
	}

	return counter == len(slice)
}

// IsStdin checks whether the provided input location refers to the dash
// character which usually serves as the replacement to point to STDIN rather
// than a file.
func IsStdin(location string) bool {
	return strings.TrimSpace(location) == "-"
}
