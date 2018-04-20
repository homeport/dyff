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
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"

	"github.com/HeavyWombat/yaml"
	"github.com/pkg/errors"
)

type InputFile struct {
	Documents []interface{}
}

// LoadFile processes the provided input location to load a YAML (or JSON, or raw text)
func LoadFile(location string) (InputFile, error) {
	var (
		documents []interface{}
		data      []byte
		err       error
	)

	if data, err = GetBytesFromLocation(location); err != nil {
		return InputFile{}, errors.Wrap(err, fmt.Sprintf("Unable to load data from %s", location))
	}

	if documents, err = LoadDocuments(data); err != nil {
		return InputFile{}, errors.Wrap(err, fmt.Sprintf("Unable to parse data from %s", location))
	}

	return InputFile{Documents: documents}, nil
}

func LoadDocuments(input []byte) ([]interface{}, error) {
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
		}

		valueType := reflect.TypeOf(value).Kind()
		switch valueType {
		case reflect.Slice:
			if isComplexSlice(value.([]interface{})) {
				types = append(types, "complex-slice")

			} else {
				types = append(types, valueType.String())
			}

		default:
			types = append(types, valueType.String())
		}
	}

	// Second pass: Based on the types, initialise a proper variable to unmarshal data into
	values = make([]interface{}, len(types))
	decoder = yaml.NewDecoder(bytes.NewReader(input))
	for i := 0; i < len(types); i++ {
		switch types[i] {
		case "map":
			var value yaml.MapSlice
			decoder.Decode(&value)
			values[i] = value

		case "slice":
			var value []interface{}
			decoder.Decode(&value)
			values[i] = value

		case "complex-slice":
			var value []yaml.MapSlice
			decoder.Decode(&value)
			values[i] = value

		default:
			return nil, fmt.Errorf("Unsupported type %s in load document function", types[i])
		}
	}

	return values, nil
}

func GetBytesFromLocation(location string) ([]byte, error) {
	var data []byte
	var err error

	// Handle special location "-" which referes to STDIN stream
	if location == "-" {
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
		data = buf.Bytes()

		return data, nil
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
		case map[interface{}]interface{}, yaml.MapSlice:
			counter++
		}
	}

	return counter == len(slice)
}
