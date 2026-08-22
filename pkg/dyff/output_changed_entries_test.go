package dyff_test

import (
	"os"
	"regexp"
	"strings"

	. "github.com/gonvenience/bunt"
	"github.com/gonvenience/ytbx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/homeport/dyff/pkg/dyff"
)

// normalize output (line endings + strip ANSI + trim)
func normalizeChangedEntriesOutput(s string) string {
	// Normalize line endings
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	// Strip ANSI sequences
	ansiRE := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	s = ansiRE.ReplaceAllString(s, "")
	return strings.TrimSpace(s)
}

// diffLines returns unified like diff of two multi-line strings (without context)
func diffLines(expected, got string) string {
	if expected == got {
		return ""
	}
	eLines := strings.Split(expected, "\n")
	gLines := strings.Split(got, "\n")
	max := len(eLines)
	if len(gLines) > max {
		max = len(gLines)
	}
	var b strings.Builder
	for i := 0; i < max; i++ {
		var e, g string
		if i < len(eLines) {
			e = eLines[i]
		}
		if i < len(gLines) {
			g = gLines[i]
		}
		if e != g {
			b.WriteString("-" + e + "\n")
			b.WriteString("+" + g + "\n")
		}
	}
	if b.Len() == 0 {
		b.WriteString("whitespace mismatch\n")
	}
	return b.String()
}

var _ = Describe("changed entries report", func() {
	Context("issue-525 regression", func() {
		BeforeEach(func() { SetColorSettings(OFF, OFF) })
		AfterEach(func() { SetColorSettings(AUTO, AUTO) })

		It("should show the expected changed entries output", func() {
			from, to, err := ytbx.LoadFiles(assets("issues/issue-525/from.yaml"), assets("issues/issue-525/to.yaml"))
			Expect(err).NotTo(HaveOccurred())

			report, err := dyff.CompareInputFiles(from, to)
			Expect(err).NotTo(HaveOccurred())

			writer := &dyff.ChangedEntriesReport{Report: report}
			var sb strings.Builder
			Expect(writer.WriteReport(&sb)).To(Succeed())

			got := normalizeChangedEntriesOutput(sb.String())

			expectedBytes, err := os.ReadFile(assets("issues/issue-525/expected.changed-entries"))
			Expect(err).NotTo(HaveOccurred())
			expected := normalizeChangedEntriesOutput(string(expectedBytes))

			if got != expected {
				Fail("changed entries output mismatch:\n" + diffLines(expected, got))
			}
		})
	})

	Context("when there are no changes", func() {
		BeforeEach(func() { SetColorSettings(OFF, OFF) })
		AfterEach(func() { SetColorSettings(AUTO, AUTO) })

		It("prints a helpful message", func() {
			// Single trivial YAML document used as both from and to
			docs, err := ytbx.LoadYAMLDocuments([]byte("---\nkey: value\n"))
			Expect(err).NotTo(HaveOccurred())

			input := ytbx.InputFile{Documents: docs}
			report, err := dyff.CompareInputFiles(input, input)
			Expect(err).NotTo(HaveOccurred())

			writer := &dyff.ChangedEntriesReport{Report: report}
			var sb strings.Builder
			Expect(writer.WriteReport(&sb)).To(Succeed())
			Expect(sb.String()).To(Equal("No changed entries found.\n"))
		})
	})
})
