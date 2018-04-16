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
	"fmt"

	"github.com/HeavyWombat/yaml"
)

var knownKeyOrders = [][]string{
	{"name", "director_uuid", "releases", "stemcells", "update", "instance_groups", "addons"},             // https://bosh.io/docs/manifest-v2.html
	{"name", "director_uuid", "releases", "instance_groups", "networks", "resource_pools", "compilation"}, // Random actual example ...
	{"jobs", "resources", "resource_types"},                                                               // https://concourse-ci.org/pipelines.html
	{"name", "type", "source"},                                                                            // https://concourse-ci.org/resources.html
	{"get"},                                                                                               // https://concourse-ci.org/steps.html
	{"put"},                                                                                               // https://concourse-ci.org/steps.html
	{"task"},                                                                                              // https://concourse-ci.org/steps.html
	{"name"},                                                                                              // Universal default #1 ... name should always be first
	{"key"},                                                                                               // Universal default #2 ... key should always be first
	{"id"},                                                                                                // Universal default #3 ... id should always be first
}

func lookupMap(list []string) map[string]int {
	result := make(map[string]int, len(list))
	for idx, entry := range list {
		result[entry] = idx
	}

	return result
}

func hasAll(keys, list []string) bool {
	counter := 0
	target := len(list)

	lookup := lookupMap(keys)
	for _, key := range list {
		if _, ok := lookup[key]; ok {
			counter++

			if counter == target {
				return true
			}
		}
	}

	return false
}

func neword(input yaml.MapSlice, keys []string) yaml.MapSlice {
	// Add all keys from the input MapSlice that are not part of the ordered keys list
	lookup := lookupMap(keys)
	for _, mapitem := range input {
		key := mapitem.Key.(string)
		if _, ok := lookup[key]; !ok {
			keys = append(keys, key)
		}
	}

	// Rebuild a new YAML MapSlice key by key using in provided keys list for the order
	result := yaml.MapSlice{}
	for _, key := range keys {
		result = append(result, yaml.MapItem{
			Key:   key,
			Value: GetKeyValueOrPanic(input, key),
		})
	}

	return result
}

func foobar(keys []string) func(yaml.MapSlice) yaml.MapSlice {
	for _, candidate := range knownKeyOrders {
		if hasAll(keys, candidate) {
			return func(input yaml.MapSlice) yaml.MapSlice {
				return neword(input, candidate)
			}
		}
	}

	return nil
}

func ListStringKeys(mapslice yaml.MapSlice) ([]string, error) {
	keys := make([]string, len(mapslice))
	for i, mapitem := range mapslice {
		switch mapitem.Key.(type) {
		case string:
			keys[i] = mapitem.Key.(string)

		default:
			return nil, fmt.Errorf("Provided mapslice mapitem contains non-string key: %#v", mapitem.Key)
		}
	}

	return keys, nil
}

func RestructureMapSlice(mapslice yaml.MapSlice) yaml.MapSlice {
	// Restructure the YAML MapSlice keys
	if keys, err := ListStringKeys(mapslice); err == nil {
		if fn := foobar(keys); fn != nil {
			mapslice = fn(mapslice)
		}
	}

	// Restructure the values of the respective keys of this YAML MapSlice
	for _, mapitem := range mapslice {
		mapitem.Value = RestructureObject(mapitem.Value)
	}

	return mapslice
}

func RestructureObject(obj interface{}) interface{} {
	switch obj.(type) {
	case yaml.MapSlice:
		return RestructureMapSlice(obj.(yaml.MapSlice))

	case []interface{}:
		list := obj.([]interface{})
		for i := range list {
			list[i] = RestructureObject(list[i])
		}
		return list

	default:
		return obj
	}
}
