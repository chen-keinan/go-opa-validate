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
		{name: "test validate policy deny pod name json format", data: "./fixture/pod.json", policyRule: []string{"example.deny"}, policy: "./fixture/pod_policy_deny", want: true, wantError: nil},
		{name: "test validate policy deny pod name yaml format", data: "./fixture/pod.yaml", policyRule: []string{"example.deny"}, policy: "./fixture/pod_policy_deny", want: true, wantError: nil},
		{name: "test validate policy allow pod name", data: "./fixture/allow_pod.json", policyRule: []string{"example.deny"}, policy: "./fixture/pod_policy_deny", want: false, wantError: nil},
		{name: "test validate policy bad data", data: "./fixture/badJson.json", policyRule: []string{"example.deny"}, policy: "./fixture/pod_policy_deny", want: false, wantError: nil},
		{name: "test validate policy bad policy", data: "./fixture/badJson.json", policyRule: []string{"example.deny"}, policy: "./fixture/pod_policy_deny_bad", want: false, wantError: fmt.Errorf("1 error occurred: eval.rego:5: rego_parse_error: unexpected } token\n\t}\n\t^")},
		{name: "test validate policy bad policy", data: "./fixture/strict_policy.json", policyRule: []string{"itsio.allow"}, policy: "./fixture/deny_strict.policy", want: true, wantError: nil},
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
			got, err := NewPolicyEval().EvaluatePolicy(tt.policyRule, string(policy), string(data))
			if err != nil {
				goErr := err.Error()
				if goErr != tt.wantError.Error() {
					t.Fatal(err)
				}
			}
			if err == nil {
				if got[0].Value != tt.want {
					t.Errorf("Test_PolicyEval() = %v, want %v", got[0], tt.want)
				}
			}
		})
	}
}
