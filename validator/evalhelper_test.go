package validator

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func Test_ParseJson(t *testing.T) {
	tests := []struct {
		name string
		data string
		want error
	}{
		{name: "test parse json format good", data: "./fixture/pod.json", want: nil},
		{name: "test parse json format bad", data: "./fixture/pod_policy_deny", want: fmt.Errorf("invalid character 'p' looking for beginning of value")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := ioutil.ReadFile(tt.data)
			if err != nil {
				t.Fatal(err)
			}
			_, err = ParseJSON(string(data))
			if err != tt.want {
				if err.Error() != tt.want.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}

func Test_YamlToJSON(t *testing.T) {
	tests := []struct {
		name string
		data string
		want error
	}{
		{name: "test parse json format good", data: "./fixture/pod.yaml", want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := ioutil.ReadFile(tt.data)
			if err != nil {
				t.Fatal(err)
			}
			yToJ, err := YamlToJSON(string(data))
			if err != nil {
				t.Fatal(err)
			}
			input, err := ParseJSON(string(yToJ))
			fmt.Println(input)
			if err != tt.want {
				t.Fatal(err)
			}
		})
	}
}
