// Copyright © 2018 The Homeport Team
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

package bunt

import (
	"bytes"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

var boldMarker = regexp.MustCompile(`\*([^*]+?)\*`)
var italicMarker = regexp.MustCompile(`_([^_]+?)_`)
var underlineMarker = regexp.MustCompile(`~([^~]+?)~`)
var colorMarker = regexp.MustCompile(`(\w+)\{([^}]+?)\}`)

type ansiString struct {
	text string
	seqs []marker
}

type marker struct {
	position int
	codes    []uint8
	fgColor  *colorful.Color
}

func (a *ansiString) String() string {
	if !UseColors() {
		return a.text
	}

	result := a.text
	offset := 0

	for _, marker := range a.seqs {
		seq := markerToSeq(marker)
		pos := marker.position + offset
		result = result[:pos] + seq + result[pos:]
		offset += len(seq)
	}

	return result
}

func markerToSeq(m marker) string {
	values := []uint8{}

	values = append(values, m.codes...)

	if m.fgColor != nil {
		r, g, b := (*m.fgColor).RGB255()
		values = append(values, 38, 2, r, g, b)
	}

	return ansiSeq(values...)
}

func ansiSeq(values ...uint8) string {
	var buf bytes.Buffer
	buf.WriteString(seq)
	buf.WriteString("[")
	buf.WriteString(strings.Join(uintSliceToStringSlice(values), ";"))
	buf.WriteString("m")

	return buf.String()
}

func uintSliceToStringSlice(uintSlice []uint8) []string {
	result := make([]string, len(uintSlice))
	for i, x := range uintSlice {
		result[i] = strconv.Itoa(int(x))
	}

	return result
}

func stringSliceToUintSlice(strSlice []string) ([]uint8, error) {
	result := make([]uint8, len(strSlice))
	for i, x := range strSlice {
		num, err := strconv.Atoi(x)
		if err != nil {
			return nil, err
		}

		result[i] = uint8(num)
	}

	return result, nil
}

func parseSettings(input []string) ([]uint8, *colorful.Color, error) {
	settings, err := stringSliceToUintSlice(input)
	if err != nil {
		return nil, nil, err
	}

	var codes []uint8
	var color *colorful.Color

	for len(settings) > 0 {
		switch {
		case len(settings) >= 5 && settings[0] == 38 && settings[1] == 2:
			color = &colorful.Color{
				R: float64(settings[2]) / 255.0,
				G: float64(settings[3]) / 255.0,
				B: float64(settings[4]) / 255.0,
			}
			settings = settings[5:]

		default:
			codes = append(codes, settings[0])
			settings = settings[1:]
		}
	}

	return codes, color, nil
}

func blendIn(orig []marker, oldLength int, newLength int, new ...marker) []marker {
	if len(orig) == 0 {
		return append(orig, new...)
	}

	if len(new) != 2 {
		panic("blendIn function only supports two new markers, one start and one end")
	}

	posA, posB := new[0].position, new[1].position
	for idx := 0; idx < len(orig); idx += 2 {
		if orig[idx].position >= posA && posB <= orig[idx+1].position {
			textOffset := oldLength - newLength
			orig[idx+1].position -= textOffset

			insertPoint := idx + 1
			return append(orig[:insertPoint], append(new, orig[insertPoint:]...)...)
		}
	}

	return append(orig, new...)
}

func parseString(input string) (*ansiString, error) {
	text := RemoveAllEscapeSequences(input)
	sequences := []marker{}

	// --- --- ---

	type textMarker struct {
		search  string
		replace string
		posA    int
		posB    int
		codes   []uint8
		fgColor *colorful.Color
	}

	positionsToTextMarker := func(list *[]textMarker, str string, positions []int, code uint8) {
		if len(positions) > 0 {
			*list = append(*list, textMarker{
				search:  str[positions[0]:positions[1]],
				replace: str[positions[2]:positions[3]],
				posA:    positions[0],
				posB:    positions[3] - 1,
				codes:   []uint8{code},
			})
		}
	}

	positionsToTextMarker2 := func(list *[]textMarker, str string, positions []int) {
		if len(positions) == 6 {
			colorName := str[positions[2]:positions[3]]
			if color := lookupColorByName(colorName); color != nil {
				*list = append(*list, textMarker{
					search:  str[positions[0]:positions[1]],
					replace: str[positions[4]:positions[5]],
					posA:    positions[0],
					posB:    positions[5] - len(colorName) - 1,
					fgColor: color,
				})
			}
		}
	}

	next := func(str string) *textMarker {
		results := []textMarker{}
		positionsToTextMarker2(&results, str, colorMarker.FindStringSubmatchIndex(str))
		positionsToTextMarker(&results, str, boldMarker.FindStringSubmatchIndex(str), 1)
		positionsToTextMarker(&results, str, italicMarker.FindStringSubmatchIndex(str), 3)
		positionsToTextMarker(&results, str, underlineMarker.FindStringSubmatchIndex(str), 4)

		if len(results) == 0 {
			return nil
		}

		sort.Slice(results, func(i, j int) bool {
			return results[i].posA < results[j].posA
		})

		return &results[0]
	}

	for x := next(text); x != nil; x = next(text) {
		text = strings.Replace(text, x.search, x.replace, 1)
		sequences = blendIn(sequences, len(x.search), len(x.replace),
			marker{position: x.posA, codes: x.codes, fgColor: x.fgColor},
			marker{position: x.posB, codes: []uint8{0}},
		)
	}

	// --- --- ---

	escapeSeqFinderRegExp := regexp.MustCompile(seq + `\[(\d+(;\d+)*)m`)

	offset := 0
	for _, submatch := range escapeSeqFinderRegExp.FindAllStringSubmatchIndex(input, -1) {
		fullMatchStart, fullMatchEnd := submatch[0], submatch[1]
		settingsStart, settingsEnd := submatch[2], submatch[3]

		settings := strings.Split(input[settingsStart:settingsEnd], ";")
		codes, color, err := parseSettings(settings)
		if err != nil {
			return nil, err
		}

		sequences = append(sequences, marker{position: fullMatchStart - offset, codes: codes, fgColor: color})
		offset += fullMatchEnd - fullMatchStart
	}

	return &ansiString{text, sequences}, nil
}
