package core_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/HeavyWombat/dyff/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	yaml "gopkg.in/yaml.v2"
)

func TestCore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Core Suite")
}

func yml(input string) yaml.MapSlice {
	// If input is a file loacation, load this as YAML
	if _, err := os.Open(input); err == nil {
		var content yaml.MapSlice
		var err error
		if content, err = core.LoadFile(input); err != nil {
			Fail(fmt.Sprintf("Failed to load YAML MapSlice from '%s': %v", input, err))
		}

		return content
	}

	content := yaml.MapSlice{}
	if err := yaml.UnmarshalStrict([]byte(input), &content); err != nil {
		Fail(fmt.Sprintf("Failed to create test YAML MapSlice from input string:\n%s\n\n%v", input, err))
	}

	return content
}

func path(path string) core.Path {
	// path string looks like: /additions/named-entry-list-using-id/id=new

	if path == "" {
		panic("Unable to create path using an empty string")
	}

	result := make([]core.PathElement, 0)
	for i, section := range strings.Split(path, "/") {
		if i == 0 {
			if section != "" {
				panic("Invalid Go-Patch style path, it cannot start with anything other than a slash")
			}

			continue
		}

		keyNameSplit := strings.Split(section, "=")
		switch len(keyNameSplit) {
		case 1:
			result = append(result, core.PathElement{Name: keyNameSplit[0]})

		case 2:
			result = append(result, core.PathElement{Key: keyNameSplit[0], Name: keyNameSplit[1]})

		default:
			panic(fmt.Sprintf("Invalid Go-Patch style path, path element '%s' cannot contain more than one equal sign", section))
		}
	}

	return result
}
