package core_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	yaml "gopkg.in/yaml.v2"
)

func TestCore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Core Suite")
}

func getYamlFromString(input string) (yaml.MapSlice, error) {
	content := yaml.MapSlice{}
	err := yaml.UnmarshalStrict([]byte(input), &content)
	if err != nil {
		return nil, err
	}

	return content, nil
}
