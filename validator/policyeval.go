package validator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
)

func PolicyEval() {

	ctx := context.Background()

	// Define a simple policy.
	module := `package example
default deny = false
deny {
	some i
	input.request.kind.kind == "Pod"
	image := input.request.object.spec.containers[i].image
	not startswith(image, "hooli.com/")
	msg := sprintf("Image '%v' comes from untrusted registry", [image])
}
	`

	// Compile the module. The keys are used as identifiers in error messages.
	compiler, err := ast.CompileModules(map[string]string{
		"example.rego": module,
	})
	data := `{
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
                        "image": "hooli.com/nginx",
                        "name": "nginx-frontend"
                    },
                    {
                        "image": "mysql",
                        "name": "mysql-backend"
                    }
                ]
            }
        }
    }
}
`
	var raw interface{}
	json.NewEncoder(bytes.NewBufferString(data)).Encode(raw)
	// Create a new query that uses the compiled policy from above.
	rego := rego.New(
		rego.Query("data.example.deny"),
		rego.Compiler(compiler),
		rego.Input(
			raw,
		),
	)

	// Run evaluation.
	rs, err := rego.Eval(ctx)

	if err != nil {
		// Handle error.
	}

	// Inspect results.
	fmt.Println("len:", len(rs))
	fmt.Println("value:", rs[0].Expressions[0].Value)

	// Output:
	//
	// len: 1
	// value: true
}
