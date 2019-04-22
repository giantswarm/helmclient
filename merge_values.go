package helmclient

import (
	"fmt"

	"github.com/giantswarm/microerror"
	yaml "gopkg.in/yaml.v2"
)

// yamlToStringMap unmarshals the YAML input into a map[string]interface{}
// with string keys. This is necessary because the default behaviour of the
// YAML parser is to return map[interface{}]interface{} types.
// See https://github.com/go-yaml/yaml/issues/139.
//
func yamlToStringMap(input []byte) (map[string]interface{}, error) {
	var raw interface{}
	var result map[string]interface{}

	err := yaml.Unmarshal(input, &raw)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	output := processMapValue(raw)
	result = output.(map[string]interface{})

	return result, nil
}

func processInterfaceArray(in []interface{}) []interface{} {
	res := make([]interface{}, len(in))
	for i, v := range in {
		res[i] = processMapValue(v)
	}
	return res
}

func processInterfaceMap(in map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range in {
		res[fmt.Sprintf("%v", k)] = processMapValue(v)
	}
	return res
}

func processMapValue(v interface{}) interface{} {
	switch v := v.(type) {
	case bool:
		return v
	case float64:
		return v
	case int:
		return v
	case string:
		return v
	case []interface{}:
		return processInterfaceArray(v)
	case map[interface{}]interface{}:
		return processInterfaceMap(v)
	default:
		return microerror.Maskf(executionFailedError, "%#v with type %T not supported")
	}
}
