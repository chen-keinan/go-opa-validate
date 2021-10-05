[![Go Report Card](https://goreportcard.com/badge/github.com/chen-keinan/opa-policy-validate)](https://goreportcard.com/report/github.com/chen-keinan/opa-policy-validate)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/chen-keinan/go-command-eval/blob/master/LICENSE)
<img src="./pkg/img/coverage_badge.png" alt="test coverage badge">
[![Gitter](https://badges.gitter.im/beacon-sec/community.svg)](https://gitter.im/beacon-sec/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
<br><img src="./pkg/img/opa_img_gopher.png" width="300" alt="opa_val logo"><br>

# go-opa-validate

go-opa-validate is an open-source lib that evaluates OPA (open policy agent) policy against JSON or YAML data.

* [Installation](#installation)
* [Usage](#usage)
* [Contribution](#Contribution)


## Installation

```shell
go get github.com/chen-keinan/go-opa-validate
```

## Usage 
#### (support json and yaml formats)
#### json data example: data.json
```json
{
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
            "image": "hooli.com/mysql",
            "name": "mysql-backend"
          }
        ]
      }
    }
  }
}
```
#### OPA policy example : denyPolicy

```shell
package example
default deny = false
deny {
	some i
	input.request.kind.kind == "Pod"
	image := input.request.object.spec.containers[i].image
	not startswith(image, "hooli.com/")
}
```

Full code example

```go
package main

import (
	"fmt"
	"github.com/chen-keinan/go-opa-validate/validator"
	"io/ioutil"
	"os"
)


func main() {
	data, err := ioutil.ReadFile("./example/data.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	policy, err := ioutil.ReadFile("./example/denyPolicy")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	validateResult, err := validator.NewPolicyEval().EvaluatePolicy([]string{"deny"}, string(policy), string(data))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(validateResult) > 0 {
		fmt.Println(fmt.Sprintf("eval result for property %v with value %v", validateResult[0].ValidateProperty, validateResult[0].Value))
	}
}
```


## Contribution
code contribution is welcome !!
contribution with passing tests and linter is more than welcome :)