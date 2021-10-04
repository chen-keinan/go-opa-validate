package validator

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"io/ioutil"
	"strings"
)

//Evaluator OPA evaluate interface
type Evaluator interface {
	EvaluatePolicy(evalProperty []string, policy string, data string) ([]*ValidateResult, error)
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
func (pe policyEval) EvaluatePolicy(evalProperty []string, policy string, data string) ([]*ValidateResult, error) {
	var pkgName string
	const policyPackage = "package"
	reader := ioutil.NopCloser(bytes.NewReader([]byte(policy)))
	defer reader.Close()
	scanner := bufio.NewScanner(reader)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, policyPackage) {
			pkgName = strings.TrimSpace(strings.Replace(line, policyPackage, "", -1))
			break
		}
	}
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
	for _, pr := range evalProperty {
		if len(pkgName) > 0 {
			regoFunc = append(regoFunc, rego.Query(fmt.Sprintf("data.%s.%s", pkgName, pr)))
		} else {
			regoFunc = append(regoFunc, rego.Query(policy))
		}
	}
	regoFunc = append(regoFunc, rego.Compiler(compiler))
	regoFunc = append(regoFunc, rego.Input(inputObject))
	rego := rego.New(regoFunc...)
	res, err := rego.Eval(ctx)
	if err != nil {
		return nil, err
	}
	validateResult := make([]*ValidateResult, 0)
	if len(res) > 0 {
		validateResult = append(validateResult, &ValidateResult{Value: res[0].Expressions[0].Value.(bool), ValidateProperty: res[0].Expressions[0].Text})
	}
	return validateResult, nil
}

//ValidateResult opa validation results
type ValidateResult struct {
	Value            bool
	ValidateProperty string
}
