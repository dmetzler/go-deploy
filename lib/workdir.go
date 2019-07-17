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
	"github.com/otiai10/copy"
	"io/ioutil"
)

func BuildWorkDir(srcDir string, dotEnv string, configName string) (error, string) {
    // Create temporary workdir
		workdir, err := ioutil.TempDir("/tmp", "go-deploy")
		if err != nil {
		    return err, ""
		}

    // Copy the Source directory into our workdir
    err = copy.Copy(srcDir, workdir)
    if err != nil {
		    return err, ""
		}

		// Generate env-config.js in workdir
		err = GenerateDotEnv(srcDir + "/" + dotEnv, workdir + "/" + configName )
		if err != nil {
		    return err, ""
		}

		return nil, workdir

}