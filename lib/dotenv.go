/*
Copyright © 2019 Nuxeo

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package lib

import (
	"os"
	"errors"
	"bufio"
	"strings"
	"text/template"
)

type EnvVar struct {
	Name string
	Value string
}

var tpl = `'use strict'
window._env_ = { {{ range . }}
    {{ .Name }}: "{{ .Value }}",{{end}}
}
`


func GenerateDotEnv(srcEnv string, dstEnvJs string ) error {

		if _, err := os.Stat(srcEnv); os.IsNotExist(err) {
			return errors.New(".env file is not present (" + srcEnv + ")")
		}

		err, vars := getVars(srcEnv)
		if(err != nil) {
			return err
		}

		return renderDotEnv(vars, dstEnvJs)

}

func renderDotEnv(vars []EnvVar, dstEnvJs string) error {
	envConfig, err := os.Create(dstEnvJs)
  if err != nil {
  	return errors.New("Unable to create destination file: " + dstEnvJs)
  }
  defer envConfig.Close()

  t := template.Must(template.New("vars").Parse(tpl))
  err = t.Execute(envConfig, vars)
  if err != nil {
  	return err
  }
  return nil;

}


func getVars(dotEnv string) (error, []EnvVar) {
		file, err := os.Open(dotEnv)
		if err != nil {
			return err, nil
    }
    defer file.Close()

    vars := make([]EnvVar,0)

    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        line := strings.Split(scanner.Text(), "=")

        name := line[0]
        value, exists := os.LookupEnv(name)
        if(!exists) {
        	value = line[1]
        }
        vars = append(vars,EnvVar{name, value})
    }

    if err := scanner.Err(); err != nil {
      return err, nil
    }
    return nil, vars
}
