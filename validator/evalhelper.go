package validator

import (
	"bytes"
	"encoding/json"
	"gopkg.in/yaml.v2"
)

//ParseJSON parse json file return json data or error
func ParseJSON(data string) (interface{}, error) {
	var input interface{}
	d := json.NewDecoder(bytes.NewBufferString(data))
	// Numeric values must be represented using json.Number.
	d.UseNumber()
	if err := d.Decode(&input); err != nil {
		return nil, err
	}
	return input, nil
}

//YamlToJSON convert yaml to json , accept yaml data and return json data or error
func YamlToJSON(data string) ([]byte, error) {
	var body interface{}
	if err := yaml.Unmarshal([]byte(data), &body); err != nil {
		return nil, err
	}
	bodyJSON := convert(body)
	var b []byte
	var err error
	if b, err = json.Marshal(bodyJSON); err != nil {
		return nil, err
	}
	return b, nil
}

func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i
}
