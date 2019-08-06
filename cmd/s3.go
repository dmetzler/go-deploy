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
	"github.com/dmetzler/go-deploy/lib"
	"github.com/aws/aws-sdk-go/service/s3"
)


var (
  validStorageClasses = map[string]bool{
    "":                                      true,
    s3.ObjectStorageClassStandard:           true,
    s3.ObjectStorageClassReducedRedundancy:  true,
    s3.ObjectStorageClassGlacier:            true,
    s3.ObjectStorageClassStandardIa:         true,
    s3.ObjectStorageClassOnezoneIa:          true,
    s3.ObjectStorageClassIntelligentTiering: true,
    s3.ObjectStorageClassDeepArchive:        true,
  }
)


func init() {
	rootCmd.AddCommand(s3Cmd)
	s3Cmd.Flags().StringP("env", "e", ".env", "Source dotenv file")
	s3Cmd.Flags().StringP("configname", "c", "env-config.js", "Name of the generated config file")
  s3Cmd.Flags().StringP("access-key", "", "", "AWS Access Key")
  s3Cmd.Flags().StringP("secret-key", "", "", "AWS Secret Key")
  s3Cmd.Flags().StringP("storage-class", "", "", "S3 Storage Class")
  s3Cmd.Flags().IntP("concurrency", "", 10 , "Concurrency")
  s3Cmd.Flags().Int64P("part-size", "", 0, "Part Size in MB")
  s3Cmd.Flags().BoolP("check-md5", "", false, "Check MD5")
  s3Cmd.Flags().BoolP("dry-run", "", false, "Dry Run")
  s3Cmd.Flags().BoolP("verbose", "", false, "Verbose")
  s3Cmd.Flags().BoolP("recursive", "", true, "Recursive")
  s3Cmd.Flags().BoolP("force", "", false, "Force")
  s3Cmd.Flags().BoolP("skip-existing", "", false, "Skip existing")

}



// s3Cmd represents the s3 command
var s3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "Deploys to a S3 bucket",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
	    log.Fatal("Not enough arguments: add the destination bucket as the argument")
	  }

	  bucket := args[0]

		srcDir, exists := os.LookupEnv("SRC_DIR")
    if(!exists) {
    	log.Fatal("SRC_DIR env variable does not exist")
    }

    if _, err := os.Stat(srcDir); os.IsNotExist(err) {
			log.Fatal("Source directory does not exist (SRC_DIR: " + srcDir + ")")
		}

		configName, _:= cmd.Flags().GetString("configname")
		dotenv, _:= cmd.Flags().GetString("env")

		err, workdir := lib.BuildWorkDir(srcDir, dotenv, configName )
		if err != nil {
			log.Fatal(err)
		}

		config := &lib.Config{}
		config.AccessKey, _ = cmd.Flags().GetString("access-key")
		config.SecretKey, _  = cmd.Flags().GetString("secret-key")
		config.StorageClass, _  = cmd.Flags().GetString("storage-class")
		config.Concurrency, _  = cmd.Flags().GetInt("concurrency")
		config.PartSize, _  = cmd.Flags().GetInt64("part-size")
		config.CheckMD5, _  = cmd.Flags().GetBool("check-md5")
		config.DryRun, _  = cmd.Flags().GetBool("dry-run")
		config.Verbose, _  = cmd.Flags().GetBool("verbose")
		config.Recursive, _  = cmd.Flags().GetBool("recursive")
		config.Force, _  = cmd.Flags().GetBool("force")
		config.SkipExisting, _  = cmd.Flags().GetBool("skip-existing")

		// Some additional validation
		if _, found := validStorageClasses[config.StorageClass]; !found {
			log.Fatal("Invalid storage class provided: %s", config.StorageClass)
		}


		err = lib.S3Sync(config, workdir + "/", bucket)
		if err != nil {
			log.Fatal(err)
		}


	},
}

