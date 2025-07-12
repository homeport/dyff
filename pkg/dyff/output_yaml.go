package dyff

import (
	"bufio"
	"github.com/gonvenience/neat"
	"github.com/gonvenience/ytbx"
	"io"
)

type YAMLReport struct {
	Report
}

type YAMLReportDiff struct {
	Details map[string]string
	Path    string
}

type YAMLReportOutput struct {
	APIVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   map[string]string `yaml:"metadata"`
	Diffs      []YAMLReportDiff  `yaml:"diffs"`
}

// TODO: Support non-Kubernetes yaml documents
func (report *YAMLReport) WriteReport(out io.Writer) error {
	writer := bufio.NewWriter(out)
	defer writer.Flush()
	consolidatedDiff, err := report.consolidateDiff()
	if err != nil {
		return err
	}
	for file, diffs := range consolidatedDiff {
		meta, err := K8sMetaFromName(file)
		if err != nil {
			return err
		}
		var d []YAMLReportDiff
		for _, diff := range diffs {
			d = append(d, YAMLReportDiff{
				Path:    diff.Path,
				Details: diff.Details,
			})
		}
		data := YAMLReportOutput{
			APIVersion: meta.APIVersion,
			Kind:       meta.Kind,
			Metadata:   meta.Metadata,
			Diffs:      d,
		}

		// Use neat to format the YAML output
		yamlData, err := neat.NewOutputProcessor(false, true, nil).ToYAML(data)
		if err != nil {
			return err
		}

		if _, err := writer.WriteString(yamlData); err != nil {
			return err
		}
	}

	_, _ = writer.WriteString("\n") // Ensure a newline at the end of the report
	return nil
}

func (report *YAMLReport) consolidateDiff() (map[string][]YAMLReportDiff, error) {
	fileDiffs := make(map[string][]YAMLReportDiff)

	for _, diff := range report.Diffs {
		deet := make(map[string]string)
		switch len(diff.Details) {
		case 1:
			switch diff.Details[0].Kind {
			case ADDITION:
				ytbx.RestructureObject(diff.Details[0].To)
				output, err := neat.NewOutputProcessor(false, true, nil).ToYAML(diff.Details[0].To)
				if err != nil {
					return nil, err
				}
				deet["to"] = output
				deet["from"] = ""
				deet["kind"] = "addition"
			case REMOVAL:
				ytbx.RestructureObject(diff.Details[0].From)
				output, err := neat.NewOutputProcessor(false, true, nil).ToYAML(diff.Details[0].From)
				if err != nil {
					return nil, err
				}
				deet["to"] = ""
				deet["from"] = output
				deet["kind"] = "removal"
			case MODIFICATION:
				ytbx.RestructureObject(diff.Details[0].To)
				outputTo, err := neat.NewOutputProcessor(false, true, nil).ToYAML(diff.Details[0].To)
				ytbx.RestructureObject(diff.Details[0].From)
				outputFrom, err := neat.NewOutputProcessor(false, true, nil).ToYAML(diff.Details[0].From)
				if err != nil {
					return nil, err
				}
				deet["to"] = outputTo
				deet["from"] = outputFrom
				deet["kind"] = "modification"
			case ORDERCHANGE:
				ytbx.RestructureObject(diff.Details[0].To)
				outputTo, err := neat.NewOutputProcessor(false, true, nil).ToYAML(diff.Details[0].To)
				ytbx.RestructureObject(diff.Details[0].From)
				outputFrom, err := neat.NewOutputProcessor(false, true, nil).ToYAML(diff.Details[0].From)
				if err != nil {
					return nil, err
				}
				deet["to"] = outputTo
				deet["from"] = outputFrom
				deet["kind"] = "orderchange"
			}
		case 2:
			for _, detail := range diff.Details {
				switch detail.Kind {
				case ADDITION:
					ytbx.RestructureObject(detail.To)
					output, err := neat.NewOutputProcessor(false, true, nil).ToYAML(detail.To)
					if err != nil {
						return nil, err
					}
					deet["to"] = output
				case REMOVAL:
					ytbx.RestructureObject(detail.From)
					output, err := neat.NewOutputProcessor(false, true, nil).ToYAML(detail.From)
					if err != nil {
						return nil, err
					}
					deet["from"] = output
				}
			}
			deet["kind"] = "modification"
		}

		fileDiffs[diff.Path.RootDescription()] = append(fileDiffs[diff.Path.RootDescription()], YAMLReportDiff{
			Path:    diff.Path.String(),
			Details: deet,
		})
	}

	return fileDiffs, nil
}
