package validator

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func Test_PolicyEval(t *testing.T) {
	tests := []struct {
		name       string
		data       string
		policy     string
		pkgName    string
		policyRule []string
		want       bool
		wantError  error
	}{
		{name: "test validate policy deny pod name json format", data: "./fixture/pod.json", pkgName: "example", policyRule: []string{"deny"}, policy: "./fixture/pod_policy_deny", want: true, wantError: nil},
		{name: "test validate policy deny pod name yaml format", data: "./fixture/pod.yaml", pkgName: "example", policyRule: []string{"deny"}, policy: "./fixture/pod_policy_deny", want: true, wantError: nil},
		{name: "test validate policy allow pod name", data: "./fixture/allow_pod.json", pkgName: "example", policyRule: []string{"deny"}, policy: "./fixture/pod_policy_deny", want: false, wantError: nil},
		{name: "test validate policy bad data", data: "./fixture/badJson.json", pkgName: "example", policyRule: []string{"deny"}, policy: "./fixture/pod_policy_deny", want: false, wantError: nil},
		{name: "test validate policy bad policy", data: "./fixture/badJson.json", pkgName: "example", policyRule: []string{"deny"}, policy: "./fixture/pod_policy_deny_bad", want: false, wantError: fmt.Errorf("1 error occurred: example.rego:5: rego_parse_error: unexpected } token\n\t}\n\t^")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := ioutil.ReadFile(tt.data)
			if err != nil {
				t.Fatal(err)
			}
			policy, err := ioutil.ReadFile(tt.policy)
			if err != nil {
				t.Fatal(err)
			}
			got, err := NewPolicyEval().EvaluatePolicy("example", []string{"deny"}, string(policy), string(data))
			if err != nil {
				goErr := err.Error()
				if goErr != tt.wantError.Error() {
					t.Fatal(err)
				}
			}
			if len(got) > 0 {
				if got[0].Value != tt.want {
					t.Errorf("Test_PolicyEval() = %v, want %v", got[0], tt.want)
				}
			}
		})
	}
}
