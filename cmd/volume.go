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
package cmd

import (
	"os"
	"github.com/spf13/cobra"
	"github.com/otiai10/copy"
	"github.com/dmetzler/go-deploy/lib"
)


// volumeCmd represents the volume command
var volumeCmd = &cobra.Command{
	Use:   "volume",
	Short: "Deploys the application in a directory, usually a Docker volume.",
	Long: `Deploys the content or $SRC_DIR to /html_dir that is usually mounted as a
volume. The command support an addional argument to specify the target directory`,
	Run: func(cmd *cobra.Command, args []string) {


		// Check that SRC_DIR exist
		srcDir, exists := os.LookupEnv("SRC_DIR")
    if(!exists) {
    	log.Fatal("SRC_DIR env variable does not exist")
			os.Exit(1)
    }

    if _, err := os.Stat(srcDir); os.IsNotExist(err) {
			log.Fatal("Source directory does not exist (SRC_DIR: " + srcDir + ")")
			os.Exit(1)
		}

		destination := "/html_dir"
		if len(args) > 1 {
	   	destination = args[0]
	  }

		configName, _:= cmd.Flags().GetString("configname")
		dotenv, _:= cmd.Flags().GetString("env")

		err, workdir := lib.BuildWorkDir(srcDir, dotenv, configName )
		if err != nil {
			log.Fatal(err)
		}

		// Copy the result into destination
		err = copy.Copy(workdir, destination)
    if err != nil {
		    log.Fatal(err)
		    os.Exit(1)
		}

		os.RemoveAll(workdir)

	},
}

func init() {
	rootCmd.AddCommand(volumeCmd)
	volumeCmd.Flags().StringP("env", "e", ".env", "Source dotenv file")
	volumeCmd.Flags().StringP("configname", "c", "env-config.js", "Name of the generated config file")
}
