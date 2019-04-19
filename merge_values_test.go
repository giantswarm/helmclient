package helmclient

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	deepNestedYaml = `
nested:
  deeper:
    value: "deeper"
test: test`
	mixedTypesYaml = `
bool: true
int: 1047552
float: 274877.906944
string: test
text: test with a sentence`
	nestedArrayYaml = `
nested:
  array:
  - 1: "test 1"
  - 2: "test 2"
test: test`
	nestedArrayAndMapYaml = `
nested:
  another: "test"
  array:
  - 1: "test 1"
  - 2: "test 2"
test: test`
	nestedMixedTypesYaml = `
nested:
  another: "test"
  array:
  - 1: 1
  - 2: 2
  deeper:
    bottom: true
    float: 274877.906944
test: test`
	simpleNestedYaml = `
nested:
  value: "nested"
test: test`
)

func Test_yamlToStringMap(t *testing.T) {
	testCases := []struct {
		name           string
		input          []byte
		expectedValues map[string]interface{}
		errorMatcher   func(error) bool
	}{
		{
			name:  "case 0: flat mixed types",
			input: []byte(mixedTypesYaml),
			expectedValues: map[string]interface{}{
				"bool":   true,
				"int":    1047552,
				"float":  274877.906944,
				"string": "test",
				"text":   "test with a sentence",
			},
		},
		{
			name:  "case 1: simple nested maps",
			input: []byte(simpleNestedYaml),
			expectedValues: map[string]interface{}{
				"nested": map[string]interface{}{
					"value": "nested",
				},
				"test": "test",
			},
		},
		{
			name:  "case 2: nested array",
			input: []byte(nestedArrayYaml),
			expectedValues: map[string]interface{}{
				"nested": map[string]interface{}{
					"array": []interface{}{
						map[string]interface{}{
							"1": "test 1",
						},
						map[string]interface{}{
							"2": "test 2",
						},
					},
				},
				"test": "test",
			},
		},
		{
			name:  "case 3: nested array and map",
			input: []byte(nestedArrayAndMapYaml),
			expectedValues: map[string]interface{}{
				"nested": map[string]interface{}{
					"another": "test",
					"array": []interface{}{
						map[string]interface{}{
							"1": "test 1",
						},
						map[string]interface{}{
							"2": "test 2",
						},
					},
				},
				"test": "test",
			},
		},
		{
			name:  "case 4: nested mixed types",
			input: []byte(nestedMixedTypesYaml),
			expectedValues: map[string]interface{}{
				"nested": map[string]interface{}{
					"another": "test",
					"array": []interface{}{
						map[string]interface{}{
							"1": 1,
						},
						map[string]interface{}{
							"2": 2,
						},
					},
					"deeper": map[string]interface{}{
						"bottom": true,
						"float":  274877.906944,
					},
				},
				"test": "test",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := yamlToStringMap(tc.input)

			switch {
			case err == nil && tc.errorMatcher == nil:
				// correct; carry on
			case err != nil && tc.errorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.errorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			case !tc.errorMatcher(err):
				t.Fatalf("error == %#v, want matching", err)
			}

			if !reflect.DeepEqual(result, tc.expectedValues) {
				t.Fatalf("want matching values \n %s", cmp.Diff(result, tc.expectedValues))
			}
		})
	}
}
