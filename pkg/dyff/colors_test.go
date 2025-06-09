// Copyright © 2019 The Homeport Team
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

package dyff_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/homeport/dyff/pkg/dyff"
	"github.com/lucasb-eyer/go-colorful"
)

var _ = Describe("Color theme functionality", func() {

	Context("DefaultColorTheme", func() {
		It("should return the default color values", func() {
			theme := dyff.DefaultColorTheme()
			Expect(theme).ToNot(BeNil())
			
			// Verify default colors match the original hardcoded values
			Expect(theme.Addition.Hex()).To(Equal("#58bf38"))
			Expect(theme.Modification.Hex()).To(Equal("#c7c43f"))
			Expect(theme.Removal.Hex()).To(Equal("#b9311b"))
		})
	})

	Context("Theme-aware color methods", func() {
		var customTheme *dyff.ColorTheme

		BeforeEach(func() {
			customTheme = &dyff.ColorTheme{
				Addition:     colorful.Color{R: 0.0, G: 1.0, B: 0.0}, // Pure green
				Modification: colorful.Color{R: 1.0, G: 0.5, B: 0.0}, // Orange
				Removal:      colorful.Color{R: 1.0, G: 0.0, B: 0.0}, // Pure red
			}
		})

		It("should properly store custom colors", func() {
			// Verify that custom theme stores the colors correctly
			Expect(customTheme.Addition.R).To(BeNumerically("~", 0.0, 0.01))
			Expect(customTheme.Addition.G).To(BeNumerically("~", 1.0, 0.01))
			Expect(customTheme.Addition.B).To(BeNumerically("~", 0.0, 0.01))
			
			Expect(customTheme.Modification.R).To(BeNumerically("~", 1.0, 0.01))
			Expect(customTheme.Modification.G).To(BeNumerically("~", 0.5, 0.01))
			Expect(customTheme.Modification.B).To(BeNumerically("~", 0.0, 0.01))
			
			Expect(customTheme.Removal.R).To(BeNumerically("~", 1.0, 0.01))
			Expect(customTheme.Removal.G).To(BeNumerically("~", 0.0, 0.01))
			Expect(customTheme.Removal.B).To(BeNumerically("~", 0.0, 0.01))
		})
	})

	Context("Integration with report structures", func() {
		It("should allow nil ColorTheme in HumanReport", func() {
			report := &dyff.HumanReport{
				Report: dyff.Report{},
				ColorTheme: nil, // Should use default theme
			}
			Expect(report.ColorTheme).To(BeNil())
		})

		It("should allow custom ColorTheme in HumanReport", func() {
			customTheme := &dyff.ColorTheme{
				Addition:     colorful.Color{R: 0.1, G: 0.9, B: 0.1},
				Modification: colorful.Color{R: 0.9, G: 0.9, B: 0.1},
				Removal:      colorful.Color{R: 0.9, G: 0.1, B: 0.1},
			}
			
			report := &dyff.HumanReport{
				Report: dyff.Report{},
				ColorTheme: customTheme,
			}
			Expect(report.ColorTheme).To(Equal(customTheme))
		})

		It("should allow nil ColorTheme in BriefReport", func() {
			report := &dyff.BriefReport{
				Report: dyff.Report{},
				ColorTheme: nil, // Should use default theme
			}
			Expect(report.ColorTheme).To(BeNil())
		})

		It("should allow custom ColorTheme in BriefReport", func() {
			customTheme := &dyff.ColorTheme{
				Addition:     colorful.Color{R: 0.2, G: 0.8, B: 0.2},
				Modification: colorful.Color{R: 0.8, G: 0.8, B: 0.2},
				Removal:      colorful.Color{R: 0.8, G: 0.2, B: 0.2},
			}
			
			report := &dyff.BriefReport{
				Report: dyff.Report{},
				ColorTheme: customTheme,
			}
			Expect(report.ColorTheme).To(Equal(customTheme))
		})
	})
})