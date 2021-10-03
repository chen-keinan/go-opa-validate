package validator

import (
	"context"
	"fmt"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
)

//Evaluator OPA evaluate interface
type Evaluator interface {
	EvaluatePolicy(pkgName string, policyRule []string, policy string, data string) ([]bool, error)
}

//policyEval opa evaluate object
type policyEval struct {
}

//NewPolicyEval instantiate new OPA eval Object
func NewPolicyEval() Evaluator {
	return &policyEval{}
}

//EvaluatePolicy evaluate opa policy against given json input , accept opa pkg name ,policy rule(deny/allow),policy and input data
// return evaluation result in a bool form
func (pe policyEval) EvaluatePolicy(pkgName string, policyRule []string, policy string, data string) ([]bool, error) {
	ctx := context.Background()
	var inputObject interface{}
	// try to read data as json format
	inputObject, err := ParseJSON(data)
	if err != nil {
		var convertedJSON []byte
		// try to read data as yaml format and convert it to json
		convertedJSON, err = YamlToJSON(data)
		if err != nil {
			return nil, err
		}
		// read data as yaml format
		inputObject, err = ParseJSON(string(convertedJSON))
		if err != nil {
			return nil, err
		}
	}
	// Compile the module. The keys are used as identifiers in error messages.
	compiler, err := ast.CompileModules(map[string]string{
		fmt.Sprintf("%s.rego", pkgName): policy,
	})
	if err != nil {
		return nil, err
	}
	regoFunc := make([]func(r *rego.Rego), 0)
	for _, pr := range policyRule {
		regoFunc = append(regoFunc, rego.Query(fmt.Sprintf("data.%s.%s", pkgName, pr)))
	}
	regoFunc = append(regoFunc, rego.Compiler(compiler))
	regoFunc = append(regoFunc, rego.Input(inputObject))
	rego := rego.New(regoFunc...)
	res, err := rego.Eval(ctx)
	if err != nil {
		return nil, err
	}
	return []bool{res[0].Expressions[0].Value.(bool)}, nil
}