package validator

import (
	"fmt"
	"testing"
)

var policy = `package example
default deny = false
deny {
	some i
	input.request.kind.kind == "Pod"
	image := input.request.object.spec.containers[i].image
	not startswith(image, "hooli.com/")
}`

var data = `{
    "kind": "AdmissionReview",
    "request": {
        "kind": {
            "kind": "Pod",
            "version": "v1"
        },
        "object": {
            "metadata": {
                "name": "myapp"
            },
            "spec": {
                "containers": [
                    {
                        "image": "mysql",
                        "name": "mysql-backend"
                    }
                ]
            }
        }
    }
}`

func Test_PolicyEval(t *testing.T) {
	res, err := NewPolicyEval().EvaluatePolicy("example", []string{"deny"}, policy, data)
	if err != nil {
		t.Error(err)
	}
	if res[0] == true {
		fmt.Println("policy match")
	}
}
