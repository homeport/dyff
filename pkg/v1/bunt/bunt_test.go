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

package bunt_test

import (
	"fmt"

	. "github.com/HeavyWombat/dyff/pkg/v1/bunt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bunt tests", func() {
	Context("Helper functions", func() {
		BeforeEach(func() {
			ColorSetting = ON
			TrueColorSetting = ON
		})

		It("should be able to break up an existing string with colors in its respective parts", func() {
			sample := fmt.Sprintf("[%s] This is a %s with various %s combined in one %s:\n(%s)",
				Colorize("Info", Green),
				Style("text", Italic),
				Style("styles", Bold),
				Colorize(`"string"`, Magenta, Bold),
				Colorize("https://github.com/HeavyWombat/dyff/pkg/bunt", CornflowerBlue, Italic, Underline))

			result, err := BreakUpStringIntoColorSegments(sample)
			Expect(err).To(BeNil())

			expected := String{Segment{Data: "["},
				Segment{Data: "Info", Attributes: []Attribute{38, 2, 0, 128, 0}},
				Segment{Data: "] This is a "},
				Segment{Data: "text", Attributes: []Attribute{Italic}},
				Segment{Data: " with various "},
				Segment{Data: "styles", Attributes: []Attribute{Bold}},
				Segment{Data: " combined in one "},
				Segment{Data: `"string"`, Attributes: []Attribute{1, 38, 2, 255, 0, 255}},
				Segment{Data: ":\n"},
				Segment{Data: "("},
				Segment{Data: "https://github.com/HeavyWombat/dyff/pkg/bunt", Attributes: []Attribute{3, 4, 38, 2, 100, 149, 237}},
				Segment{Data: ")"},
			}

			Expect(result).To(BeEquivalentTo(expected))
		})

		It("should translate a 24-bit color into the nearest 4-bit color attribute", func() {
			Expect(Get4bitEquivalentColorAttribute(Black)).To(BeEquivalentTo(Attribute(30)))                // Black matches Black (#30)
			Expect(Get4bitEquivalentColorAttribute(Brown)).To(BeEquivalentTo(Attribute(31)))                // Brown matches Red (#31)
			Expect(Get4bitEquivalentColorAttribute(DarkRed)).To(BeEquivalentTo(Attribute(31)))              // DarkRed matches Red (#31)
			Expect(Get4bitEquivalentColorAttribute(FireBrick)).To(BeEquivalentTo(Attribute(31)))            // FireBrick matches Red (#31)
			Expect(Get4bitEquivalentColorAttribute(Maroon)).To(BeEquivalentTo(Attribute(31)))               // Maroon matches Red (#31)
			Expect(Get4bitEquivalentColorAttribute(RemovalRed)).To(BeEquivalentTo(Attribute(31)))           // RemovalRed matches Red (#31)
			Expect(Get4bitEquivalentColorAttribute(SaddleBrown)).To(BeEquivalentTo(Attribute(31)))          // SaddleBrown matches Red (#31)
			Expect(Get4bitEquivalentColorAttribute(Sienna)).To(BeEquivalentTo(Attribute(31)))               // Sienna matches Red (#31)
			Expect(Get4bitEquivalentColorAttribute(AdditionGreen)).To(BeEquivalentTo(Attribute(32)))        // AdditionGreen matches Green (#32)
			Expect(Get4bitEquivalentColorAttribute(DarkGreen)).To(BeEquivalentTo(Attribute(32)))            // DarkGreen matches Green (#32)
			Expect(Get4bitEquivalentColorAttribute(DarkSeaGreen)).To(BeEquivalentTo(Attribute(32)))         // DarkSeaGreen matches Green (#32)
			Expect(Get4bitEquivalentColorAttribute(ForestGreen)).To(BeEquivalentTo(Attribute(32)))          // ForestGreen matches Green (#32)
			Expect(Get4bitEquivalentColorAttribute(Green)).To(BeEquivalentTo(Attribute(32)))                // Green matches Green (#32)
			Expect(Get4bitEquivalentColorAttribute(LimeGreen)).To(BeEquivalentTo(Attribute(32)))            // LimeGreen matches Green (#32)
			Expect(Get4bitEquivalentColorAttribute(MediumSeaGreen)).To(BeEquivalentTo(Attribute(32)))       // MediumSeaGreen matches Green (#32)
			Expect(Get4bitEquivalentColorAttribute(Olive)).To(BeEquivalentTo(Attribute(32)))                // Olive matches Green (#32)
			Expect(Get4bitEquivalentColorAttribute(OliveDrab)).To(BeEquivalentTo(Attribute(32)))            // OliveDrab matches Green (#32)
			Expect(Get4bitEquivalentColorAttribute(SeaGreen)).To(BeEquivalentTo(Attribute(32)))             // SeaGreen matches Green (#32)
			Expect(Get4bitEquivalentColorAttribute(Gold)).To(BeEquivalentTo(Attribute(33)))                 // Gold matches Yellow (#33)
			Expect(Get4bitEquivalentColorAttribute(Yellow)).To(BeEquivalentTo(Attribute(33)))               // Yellow matches Yellow (#33)
			Expect(Get4bitEquivalentColorAttribute(Blue)).To(BeEquivalentTo(Attribute(34)))                 // Blue matches Blue (#34)
			Expect(Get4bitEquivalentColorAttribute(DarkBlue)).To(BeEquivalentTo(Attribute(34)))             // DarkBlue matches Blue (#34)
			Expect(Get4bitEquivalentColorAttribute(DarkSlateBlue)).To(BeEquivalentTo(Attribute(34)))        // DarkSlateBlue matches Blue (#34)
			Expect(Get4bitEquivalentColorAttribute(Indigo)).To(BeEquivalentTo(Attribute(34)))               // Indigo matches Blue (#34)
			Expect(Get4bitEquivalentColorAttribute(MediumBlue)).To(BeEquivalentTo(Attribute(34)))           // MediumBlue matches Blue (#34)
			Expect(Get4bitEquivalentColorAttribute(MidnightBlue)).To(BeEquivalentTo(Attribute(34)))         // MidnightBlue matches Blue (#34)
			Expect(Get4bitEquivalentColorAttribute(Navy)).To(BeEquivalentTo(Attribute(34)))                 // Navy matches Blue (#34)
			Expect(Get4bitEquivalentColorAttribute(BlueViolet)).To(BeEquivalentTo(Attribute(35)))           // BlueViolet matches Magenta (#35)
			Expect(Get4bitEquivalentColorAttribute(DarkMagenta)).To(BeEquivalentTo(Attribute(35)))          // DarkMagenta matches Magenta (#35)
			Expect(Get4bitEquivalentColorAttribute(DarkOrchid)).To(BeEquivalentTo(Attribute(35)))           // DarkOrchid matches Magenta (#35)
			Expect(Get4bitEquivalentColorAttribute(DarkViolet)).To(BeEquivalentTo(Attribute(35)))           // DarkViolet matches Magenta (#35)
			Expect(Get4bitEquivalentColorAttribute(MediumVioletRed)).To(BeEquivalentTo(Attribute(35)))      // MediumVioletRed matches Magenta (#35)
			Expect(Get4bitEquivalentColorAttribute(Purple)).To(BeEquivalentTo(Attribute(35)))               // Purple matches Magenta (#35)
			Expect(Get4bitEquivalentColorAttribute(CadetBlue)).To(BeEquivalentTo(Attribute(36)))            // CadetBlue matches Cyan (#36)
			Expect(Get4bitEquivalentColorAttribute(DarkCyan)).To(BeEquivalentTo(Attribute(36)))             // DarkCyan matches Cyan (#36)
			Expect(Get4bitEquivalentColorAttribute(DarkTurquoise)).To(BeEquivalentTo(Attribute(36)))        // DarkTurquoise matches Cyan (#36)
			Expect(Get4bitEquivalentColorAttribute(DeepSkyBlue)).To(BeEquivalentTo(Attribute(36)))          // DeepSkyBlue matches Cyan (#36)
			Expect(Get4bitEquivalentColorAttribute(LightSeaGreen)).To(BeEquivalentTo(Attribute(36)))        // LightSeaGreen matches Cyan (#36)
			Expect(Get4bitEquivalentColorAttribute(MediumAquamarine)).To(BeEquivalentTo(Attribute(36)))     // MediumAquamarine matches Cyan (#36)
			Expect(Get4bitEquivalentColorAttribute(Teal)).To(BeEquivalentTo(Attribute(36)))                 // Teal matches Cyan (#36)
			Expect(Get4bitEquivalentColorAttribute(BurlyWood)).To(BeEquivalentTo(Attribute(37)))            // BurlyWood matches White (#37)
			Expect(Get4bitEquivalentColorAttribute(DarkGoldenrod)).To(BeEquivalentTo(Attribute(37)))        // DarkGoldenrod matches White (#37)
			Expect(Get4bitEquivalentColorAttribute(DarkGray)).To(BeEquivalentTo(Attribute(37)))             // DarkGray matches White (#37)
			Expect(Get4bitEquivalentColorAttribute(Gray)).To(BeEquivalentTo(Attribute(37)))                 // Gray matches White (#37)
			Expect(Get4bitEquivalentColorAttribute(LightPink)).To(BeEquivalentTo(Attribute(37)))            // LightPink matches White (#37)
			Expect(Get4bitEquivalentColorAttribute(LightSkyBlue)).To(BeEquivalentTo(Attribute(37)))         // LightSkyBlue matches White (#37)
			Expect(Get4bitEquivalentColorAttribute(LightSlateGray)).To(BeEquivalentTo(Attribute(37)))       // LightSlateGray matches White (#37)
			Expect(Get4bitEquivalentColorAttribute(LightSteelBlue)).To(BeEquivalentTo(Attribute(37)))       // LightSteelBlue matches White (#37)
			Expect(Get4bitEquivalentColorAttribute(Pink)).To(BeEquivalentTo(Attribute(37)))                 // Pink matches White (#37)
			Expect(Get4bitEquivalentColorAttribute(RosyBrown)).To(BeEquivalentTo(Attribute(37)))            // RosyBrown matches White (#37)
			Expect(Get4bitEquivalentColorAttribute(SandyBrown)).To(BeEquivalentTo(Attribute(37)))           // SandyBrown matches White (#37)
			Expect(Get4bitEquivalentColorAttribute(Silver)).To(BeEquivalentTo(Attribute(37)))               // Silver matches White (#37)
			Expect(Get4bitEquivalentColorAttribute(Tan)).To(BeEquivalentTo(Attribute(37)))                  // Tan matches White (#37)
			Expect(Get4bitEquivalentColorAttribute(Thistle)).To(BeEquivalentTo(Attribute(37)))              // Thistle matches White (#37)
			Expect(Get4bitEquivalentColorAttribute(DarkOliveGreen)).To(BeEquivalentTo(Attribute(90)))       // DarkOliveGreen matches BrightBlack (#90)
			Expect(Get4bitEquivalentColorAttribute(DarkSlateGray)).To(BeEquivalentTo(Attribute(90)))        // DarkSlateGray matches BrightBlack (#90)
			Expect(Get4bitEquivalentColorAttribute(DimGray)).To(BeEquivalentTo(Attribute(90)))              // DimGray matches BrightBlack (#90)
			Expect(Get4bitEquivalentColorAttribute(Chocolate)).To(BeEquivalentTo(Attribute(91)))            // Chocolate matches BrightRed (#91)
			Expect(Get4bitEquivalentColorAttribute(Coral)).To(BeEquivalentTo(Attribute(91)))                // Coral matches BrightRed (#91)
			Expect(Get4bitEquivalentColorAttribute(Crimson)).To(BeEquivalentTo(Attribute(91)))              // Crimson matches BrightRed (#91)
			Expect(Get4bitEquivalentColorAttribute(DarkOrange)).To(BeEquivalentTo(Attribute(91)))           // DarkOrange matches BrightRed (#91)
			Expect(Get4bitEquivalentColorAttribute(DarkSalmon)).To(BeEquivalentTo(Attribute(91)))           // DarkSalmon matches BrightRed (#91)
			Expect(Get4bitEquivalentColorAttribute(IndianRed)).To(BeEquivalentTo(Attribute(91)))            // IndianRed matches BrightRed (#91)
			Expect(Get4bitEquivalentColorAttribute(LightCoral)).To(BeEquivalentTo(Attribute(91)))           // LightCoral matches BrightRed (#91)
			Expect(Get4bitEquivalentColorAttribute(LightSalmon)).To(BeEquivalentTo(Attribute(91)))          // LightSalmon matches BrightRed (#91)
			Expect(Get4bitEquivalentColorAttribute(OrangeRed)).To(BeEquivalentTo(Attribute(91)))            // OrangeRed matches BrightRed (#91)
			Expect(Get4bitEquivalentColorAttribute(PaleVioletRed)).To(BeEquivalentTo(Attribute(91)))        // PaleVioletRed matches BrightRed (#91)
			Expect(Get4bitEquivalentColorAttribute(Peru)).To(BeEquivalentTo(Attribute(91)))                 // Peru matches BrightRed (#91)
			Expect(Get4bitEquivalentColorAttribute(Red)).To(BeEquivalentTo(Attribute(91)))                  // Red matches BrightRed (#91)
			Expect(Get4bitEquivalentColorAttribute(Salmon)).To(BeEquivalentTo(Attribute(91)))               // Salmon matches BrightRed (#91)
			Expect(Get4bitEquivalentColorAttribute(Tomato)).To(BeEquivalentTo(Attribute(91)))               // Tomato matches BrightRed (#91)
			Expect(Get4bitEquivalentColorAttribute(Chartreuse)).To(BeEquivalentTo(Attribute(92)))           // Chartreuse matches BrightGreen (#92)
			Expect(Get4bitEquivalentColorAttribute(GreenYellow)).To(BeEquivalentTo(Attribute(92)))          // GreenYellow matches BrightGreen (#92)
			Expect(Get4bitEquivalentColorAttribute(LawnGreen)).To(BeEquivalentTo(Attribute(92)))            // LawnGreen matches BrightGreen (#92)
			Expect(Get4bitEquivalentColorAttribute(LightGreen)).To(BeEquivalentTo(Attribute(92)))           // LightGreen matches BrightGreen (#92)
			Expect(Get4bitEquivalentColorAttribute(Lime)).To(BeEquivalentTo(Attribute(92)))                 // Lime matches BrightGreen (#92)
			Expect(Get4bitEquivalentColorAttribute(MediumSpringGreen)).To(BeEquivalentTo(Attribute(92)))    // MediumSpringGreen matches BrightGreen (#92)
			Expect(Get4bitEquivalentColorAttribute(PaleGreen)).To(BeEquivalentTo(Attribute(92)))            // PaleGreen matches BrightGreen (#92)
			Expect(Get4bitEquivalentColorAttribute(SpringGreen)).To(BeEquivalentTo(Attribute(92)))          // SpringGreen matches BrightGreen (#92)
			Expect(Get4bitEquivalentColorAttribute(YellowGreen)).To(BeEquivalentTo(Attribute(92)))          // YellowGreen matches BrightGreen (#92)
			Expect(Get4bitEquivalentColorAttribute(DarkKhaki)).To(BeEquivalentTo(Attribute(93)))            // DarkKhaki matches BrightYellow (#93)
			Expect(Get4bitEquivalentColorAttribute(Goldenrod)).To(BeEquivalentTo(Attribute(93)))            // Goldenrod matches BrightYellow (#93)
			Expect(Get4bitEquivalentColorAttribute(Khaki)).To(BeEquivalentTo(Attribute(93)))                // Khaki matches BrightYellow (#93)
			Expect(Get4bitEquivalentColorAttribute(ModificationYellow)).To(BeEquivalentTo(Attribute(93)))   // ModificationYellow matches BrightYellow (#93)
			Expect(Get4bitEquivalentColorAttribute(Orange)).To(BeEquivalentTo(Attribute(93)))               // Orange matches BrightYellow (#93)
			Expect(Get4bitEquivalentColorAttribute(PaleGoldenrod)).To(BeEquivalentTo(Attribute(93)))        // PaleGoldenrod matches BrightYellow (#93)
			Expect(Get4bitEquivalentColorAttribute(CornflowerBlue)).To(BeEquivalentTo(Attribute(94)))       // CornflowerBlue matches BrightBlue (#94)
			Expect(Get4bitEquivalentColorAttribute(DodgerBlue)).To(BeEquivalentTo(Attribute(94)))           // DodgerBlue matches BrightBlue (#94)
			Expect(Get4bitEquivalentColorAttribute(MediumPurple)).To(BeEquivalentTo(Attribute(94)))         // MediumPurple matches BrightBlue (#94)
			Expect(Get4bitEquivalentColorAttribute(MediumSlateBlue)).To(BeEquivalentTo(Attribute(94)))      // MediumSlateBlue matches BrightBlue (#94)
			Expect(Get4bitEquivalentColorAttribute(RoyalBlue)).To(BeEquivalentTo(Attribute(94)))            // RoyalBlue matches BrightBlue (#94)
			Expect(Get4bitEquivalentColorAttribute(SlateBlue)).To(BeEquivalentTo(Attribute(94)))            // SlateBlue matches BrightBlue (#94)
			Expect(Get4bitEquivalentColorAttribute(SlateGray)).To(BeEquivalentTo(Attribute(94)))            // SlateGray matches BrightBlue (#94)
			Expect(Get4bitEquivalentColorAttribute(SteelBlue)).To(BeEquivalentTo(Attribute(94)))            // SteelBlue matches BrightBlue (#94)
			Expect(Get4bitEquivalentColorAttribute(DeepPink)).To(BeEquivalentTo(Attribute(95)))             // DeepPink matches BrightMagenta (#95)
			Expect(Get4bitEquivalentColorAttribute(Fuchsia)).To(BeEquivalentTo(Attribute(95)))              // Fuchsia matches BrightMagenta (#95)
			Expect(Get4bitEquivalentColorAttribute(HotPink)).To(BeEquivalentTo(Attribute(95)))              // HotPink matches BrightMagenta (#95)
			Expect(Get4bitEquivalentColorAttribute(Magenta)).To(BeEquivalentTo(Attribute(95)))              // Magenta matches BrightMagenta (#95)
			Expect(Get4bitEquivalentColorAttribute(MediumOrchid)).To(BeEquivalentTo(Attribute(95)))         // MediumOrchid matches BrightMagenta (#95)
			Expect(Get4bitEquivalentColorAttribute(Orchid)).To(BeEquivalentTo(Attribute(95)))               // Orchid matches BrightMagenta (#95)
			Expect(Get4bitEquivalentColorAttribute(Plum)).To(BeEquivalentTo(Attribute(95)))                 // Plum matches BrightMagenta (#95)
			Expect(Get4bitEquivalentColorAttribute(Violet)).To(BeEquivalentTo(Attribute(95)))               // Violet matches BrightMagenta (#95)
			Expect(Get4bitEquivalentColorAttribute(Aqua)).To(BeEquivalentTo(Attribute(96)))                 // Aqua matches BrightCyan (#96)
			Expect(Get4bitEquivalentColorAttribute(Aquamarine)).To(BeEquivalentTo(Attribute(96)))           // Aquamarine matches BrightCyan (#96)
			Expect(Get4bitEquivalentColorAttribute(Cyan)).To(BeEquivalentTo(Attribute(96)))                 // Cyan matches BrightCyan (#96)
			Expect(Get4bitEquivalentColorAttribute(LightBlue)).To(BeEquivalentTo(Attribute(96)))            // LightBlue matches BrightCyan (#96)
			Expect(Get4bitEquivalentColorAttribute(MediumTurquoise)).To(BeEquivalentTo(Attribute(96)))      // MediumTurquoise matches BrightCyan (#96)
			Expect(Get4bitEquivalentColorAttribute(PaleTurquoise)).To(BeEquivalentTo(Attribute(96)))        // PaleTurquoise matches BrightCyan (#96)
			Expect(Get4bitEquivalentColorAttribute(PowderBlue)).To(BeEquivalentTo(Attribute(96)))           // PowderBlue matches BrightCyan (#96)
			Expect(Get4bitEquivalentColorAttribute(SkyBlue)).To(BeEquivalentTo(Attribute(96)))              // SkyBlue matches BrightCyan (#96)
			Expect(Get4bitEquivalentColorAttribute(Turquoise)).To(BeEquivalentTo(Attribute(96)))            // Turquoise matches BrightCyan (#96)
			Expect(Get4bitEquivalentColorAttribute(AliceBlue)).To(BeEquivalentTo(Attribute(97)))            // AliceBlue matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(AntiqueWhite)).To(BeEquivalentTo(Attribute(97)))         // AntiqueWhite matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(Azure)).To(BeEquivalentTo(Attribute(97)))                // Azure matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(Beige)).To(BeEquivalentTo(Attribute(97)))                // Beige matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(Bisque)).To(BeEquivalentTo(Attribute(97)))               // Bisque matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(BlanchedAlmond)).To(BeEquivalentTo(Attribute(97)))       // BlanchedAlmond matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(Cornsilk)).To(BeEquivalentTo(Attribute(97)))             // Cornsilk matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(FloralWhite)).To(BeEquivalentTo(Attribute(97)))          // FloralWhite matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(Gainsboro)).To(BeEquivalentTo(Attribute(97)))            // Gainsboro matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(GhostWhite)).To(BeEquivalentTo(Attribute(97)))           // GhostWhite matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(Honeydew)).To(BeEquivalentTo(Attribute(97)))             // Honeydew matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(Ivory)).To(BeEquivalentTo(Attribute(97)))                // Ivory matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(Lavender)).To(BeEquivalentTo(Attribute(97)))             // Lavender matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(LavenderBlush)).To(BeEquivalentTo(Attribute(97)))        // LavenderBlush matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(LemonChiffon)).To(BeEquivalentTo(Attribute(97)))         // LemonChiffon matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(LightCyan)).To(BeEquivalentTo(Attribute(97)))            // LightCyan matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(LightGoldenrodYellow)).To(BeEquivalentTo(Attribute(97))) // LightGoldenrodYellow matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(LightGray)).To(BeEquivalentTo(Attribute(97)))            // LightGray matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(LightYellow)).To(BeEquivalentTo(Attribute(97)))          // LightYellow matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(Linen)).To(BeEquivalentTo(Attribute(97)))                // Linen matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(MintCream)).To(BeEquivalentTo(Attribute(97)))            // MintCream matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(MistyRose)).To(BeEquivalentTo(Attribute(97)))            // MistyRose matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(Moccasin)).To(BeEquivalentTo(Attribute(97)))             // Moccasin matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(NavajoWhite)).To(BeEquivalentTo(Attribute(97)))          // NavajoWhite matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(OldLace)).To(BeEquivalentTo(Attribute(97)))              // OldLace matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(PapayaWhip)).To(BeEquivalentTo(Attribute(97)))           // PapayaWhip matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(PeachPuff)).To(BeEquivalentTo(Attribute(97)))            // PeachPuff matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(Seashell)).To(BeEquivalentTo(Attribute(97)))             // Seashell matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(Snow)).To(BeEquivalentTo(Attribute(97)))                 // Snow matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(Wheat)).To(BeEquivalentTo(Attribute(97)))                // Wheat matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(White)).To(BeEquivalentTo(Attribute(97)))                // White matches BrightWhite (#97)
			Expect(Get4bitEquivalentColorAttribute(WhiteSmoke)).To(BeEquivalentTo(Attribute(97)))           // WhiteSmoke matches BrightWhite (#97)
		})
	})
})
