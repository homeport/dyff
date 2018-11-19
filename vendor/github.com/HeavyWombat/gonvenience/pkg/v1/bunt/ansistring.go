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

package bunt

import (
	"bytes"
	"regexp"
	"sort"
	"strconv"
	"strings"
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
	setting  uint32
}

func (a *ansiString) String() string {
	if !UseColors() {
		return a.text
	}

	result := a.text
	offset := 0

	for _, marker := range a.seqs {
		seq := markerToSeq(marker.setting)
		pos := marker.position + offset
		result = result[:pos] + seq + result[pos:]
		offset += len(seq)
	}

	return result
}

func markerToSeq(setting uint32) string {
	switch setting {
	case 0x00000000, 0x01000000, 0x02000000, 0x03000000, 0x04000000:
		return ansiSeq(uint8(setting >> 24))

	default:
		values := []uint8{}
		if value := uint8(setting >> 24); value != 0 {
			values = append(values, value)
		}

		values = append(values, 38, 2,
			uint8(setting>>16)&0xFF,
			uint8(setting>>8)&0xFF,
			uint8(setting>>0)&0xFF)

		return ansiSeq(values...)
	}
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
	result := make([]string, len(uintSlice), len(uintSlice))
	for i, x := range uintSlice {
		result[i] = strconv.Itoa(int(x))
	}

	return result
}

func stringSliceToUintSlice(strSlice []string) ([]uint8, error) {
	result := make([]uint8, len(strSlice), len(strSlice))
	for i, x := range strSlice {
		num, err := strconv.Atoi(x)
		if err != nil {
			return nil, err
		}

		result[i] = uint8(num)
	}

	return result, nil
}

func parseSettings(input []string) (uint32, error) {
	settings, err := stringSliceToUintSlice(input)
	if err != nil {
		return 0, err
	}

	var result uint32
	for len(settings) > 0 {
		switch {
		case len(settings) >= 5 && settings[0] == 38 && settings[1] == 2:
			result |= uint32(settings[2])<<16 |
				uint32(settings[3])<<8 |
				uint32(settings[4])<<0
			settings = settings[5:]

		default:
			result |= uint32(settings[0]) << 24
			settings = settings[1:]
		}
	}

	return result, nil
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
		setting uint32
	}

	positionsToTextMarker := func(list *[]textMarker, str string, positions []int, setting uint32) {
		if len(positions) > 0 {
			*list = append(*list, textMarker{
				search:  str[positions[0]:positions[1]],
				replace: str[positions[2]:positions[3]],
				posA:    positions[0],
				posB:    positions[3] - 1,
				setting: setting,
			})
		}
	}

	positionsToTextMarker2 := func(list *[]textMarker, str string, positions []int) {
		if len(positions) == 6 {
			colorName := str[positions[2]:positions[3]]
			if color := lookupColorByName(colorName); color != nil {
				r, g, b := color.RGB255()
				setting := uint32((uint32(r) << 16) |
					(uint32(g) << 8) |
					(uint32(b) << 0))

				*list = append(*list, textMarker{
					search:  str[positions[0]:positions[1]],
					replace: str[positions[4]:positions[5]],
					posA:    positions[0],
					posB:    positions[5] - len(colorName) - 1,
					setting: setting,
				})
			}
		}
	}

	next := func(str string) *textMarker {
		results := []textMarker{}
		positionsToTextMarker2(&results, str, colorMarker.FindStringSubmatchIndex(str))
		positionsToTextMarker(&results, str, boldMarker.FindStringSubmatchIndex(str), 0x01000000)
		positionsToTextMarker(&results, str, italicMarker.FindStringSubmatchIndex(str), 0x03000000)
		positionsToTextMarker(&results, str, underlineMarker.FindStringSubmatchIndex(str), 0x04000000)

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
			marker{position: x.posA, setting: x.setting},
			marker{position: x.posB, setting: 0x00000000},
		)
	}

	// --- --- ---

	escapeSeqFinderRegExp := regexp.MustCompile(seq + `\[(\d+(;\d+)*)m`)

	offset := 0
	for _, submatch := range escapeSeqFinderRegExp.FindAllStringSubmatchIndex(input, -1) {
		fullMatchStart, fullMatchEnd := submatch[0], submatch[1]
		settingsStart, settingsEnd := submatch[2], submatch[3]

		settings := strings.Split(input[settingsStart:settingsEnd], ";")
		parsed, err := parseSettings(settings)
		if err != nil {
			return nil, err
		}

		sequences = append(sequences, marker{position: fullMatchStart - offset, setting: parsed})
		offset += fullMatchEnd - fullMatchStart
	}

	return &ansiString{text, sequences}, nil
}
