package validator

import (
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
	}{
		{name: "test policy deny pod name json format", data: "./fixture/pod.json", pkgName: "example", policyRule: []string{"deny"}, policy: "./fixture/pod_policy_deny", want: true},
		{name: "test policy deny pod name yaml format", data: "./fixture/pod.yaml", pkgName: "example", policyRule: []string{"deny"}, policy: "./fixture/pod_policy_deny", want: true},
		{name: "test policy allow pod name", data: "./fixture/allow_pod.json", pkgName: "example", policyRule: []string{"deny"}, policy: "./fixture/pod_policy_deny", want: false},
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
				t.Fatal(err)
			}
			if got[0] != tt.want {
				t.Errorf("Test_PolicyEval() = %v, want %v", got[0], tt.want)
			}
		})
	}
}
