package dyff_test

import (
	"os"
	"strings"
	"testing"

	"github.com/gonvenience/ytbx"
	"github.com/homeport/dyff/pkg/dyff"
)

func TestChangedEntriesReport_Issue525(t *testing.T) {
	from, to, err := ytbx.LoadFiles(assets("issues/issue-525/from.yaml"), assets("issues/issue-525/to.yaml"))
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	report, err := dyff.CompareInputFiles(from, to)
	if err != nil {
		t.Fatalf("compare: %v", err)
	}

	writer := &dyff.ChangedEntriesReport{Report: report}
	var b strings.Builder
	if err := writer.WriteReport(&b); err != nil {
		t.Fatalf("write: %v", err)
	}

	got := b.String()
	expectedBytes, err := os.ReadFile(assets("issues/issue-525/expected.changed-entries"))
	if err != nil {
		t.Fatalf("expected: %v", err)
	}
	expected := string(expectedBytes)

	if strings.TrimSpace(got) != strings.TrimSpace(expected) {
		// show diff-like output
		linesGot := strings.Split(got, "\n")
		linesExp := strings.Split(expected, "\n")
		max := len(linesGot)
		if len(linesExp) > max {
			max = len(linesExp)
		}
		var sb strings.Builder
		for i := 0; i < max; i++ {
			var g, e string
			if i < len(linesExp) {
				e = linesExp[i]
			}
			if i < len(linesGot) {
				g = linesGot[i]
			}
			if e != g {
				sb.WriteString("-" + e + "\n")
				sb.WriteString("+" + g + "\n")
			}
		}
		if sb.Len() == 0 {
			sb.WriteString("whitespace mismatch\n")
		}
		// Fail with details
		if len(got) > 4000 {
			got = got[:4000] + "..."
		}
		if len(expected) > 4000 {
			expected = expected[:4000] + "..."
		}
		// Additional helpful context
		if !strings.Contains(got, "additional/image") {
			t.Logf("missing additional/image in output")
		}
		if !strings.Contains(got, "oh-look/another-flaky") {
			t.Logf("missing flaky replacement item")
		}
		t.Fatalf("changed-entries output mismatch:\nExpected:\n%s\nGot:\n%s\nDiff:\n%s", expected, got, sb.String())
	}
}
