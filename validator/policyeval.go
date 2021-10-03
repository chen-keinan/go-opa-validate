package validator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
)

type Evaluator interface {
	EvaluatePolicy(pkgName string, policyRule []string, policy string, data string) ([]bool,error)
}

type policyEval struct {
}

func NewPolicyEval() Evaluator {
	return &policyEval{}
}

func (pe policyEval) EvaluatePolicy(pkgName string, policyRule []string, policy string, data string) ([]bool,error) {
	ctx := context.Background()
	/*err:=json.NewEncoder(bytes.NewBufferString(data)).Encode(raw)
	if err != nil{
		return nil,err
	}
*/
	d := json.NewDecoder(bytes.NewBufferString(data))

	// Numeric values must be represented using json.Number.
	d.UseNumber()

	var input interface{}

	if err := d.Decode(&input); err != nil {
		panic(err)
	}


	// Compile the module. The keys are used as identifiers in error messages.
	compiler, err := ast.CompileModules(map[string]string{
		fmt.Sprintf("%s.rego",pkgName): policy,
	})
	if err != nil {
		return nil,err
	}
	regoFunc:=make([]func(r *rego.Rego),0)
	for _,pr:=range policyRule{
		regoFunc = append(regoFunc,rego.Query(fmt.Sprintf("data.%s.%s",pkgName,pr)))
	}
	regoFunc = append(regoFunc,rego.Compiler(compiler))
	regoFunc = append(regoFunc,rego.Input(input))
	rego := rego.New(regoFunc...)
	res,err:=rego.Eval(ctx)
	if err != nil {
		return nil,err
	}
	return []bool{res[0].Expressions[0].Value.(bool)},nil
}

func Eval() []string {

	ctx := context.Background()

	// Define a simple policy.
	module := `package example
default deny = true
allow {
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
	fmt.Println("value:", rs.Allowed())

	// Output:
	//
	// len: 1
	// value: true
	return []string{}
}
