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
	"github.com/dmetzler/go-deploy/lib"
	"strings"
	"net/http"
	"github.com/spf13/cobra"
)

// containsDotFile reports whether name contains a path element starting with a period.
// The name is assumed to be a delimited by forward slashes, as guaranteed
// by the http.FileSystem interface.
func containsDotFile(name string) bool {
	parts := strings.Split(name, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, ".") {
			return true
		}
	}
	return false
}

// dotFileHidingFile is the http.File use in dotFileHidingFileSystem.
// It is used to wrap the Readdir method of http.File so that we can
// remove files and directories that start with a period from its output.
type dotFileHidingFile struct {
	http.File
}

// Readdir is a wrapper around the Readdir method of the embedded File
// that filters out all files that start with a period in their name.
func (f dotFileHidingFile) Readdir(n int) (fis []os.FileInfo, err error) {
	files, err := f.File.Readdir(n)
	for _, file := range files { // Filters out the dot files
		if !strings.HasPrefix(file.Name(), ".") {
			fis = append(fis, file)
		}
	}
	return
}

// dotFileHidingFileSystem is an http.FileSystem that hides
// hidden "dot files" from being served.
type dotFileHidingFileSystem struct {
	http.FileSystem
}

// Open is a wrapper around the Open method of the embedded FileSystem
// that serves a 403 permission error when name has a file or directory
// with whose name starts with a period in its path.
func (fs dotFileHidingFileSystem) Open(name string) (http.File, error) {
	if containsDotFile(name) { // If dot file, return 403 response
		return nil, os.ErrPermission
	}

	file, err := fs.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}
	return dotFileHidingFile{file}, err
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the web app (for development use only)",
	Long: `.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check that SRC_DIR exist
		srcDir, exists := os.LookupEnv("SRC_DIR")
    if(!exists) {
    	log.Fatal("SRC_DIR env variable does not exist")
    }

    if _, err := os.Stat(srcDir); os.IsNotExist(err) {
			log.Fatal("Source directory does not exist (SRC_DIR: " + srcDir + ")")
		}

		configName, _:= cmd.Flags().GetString("configname")
		dotenv, _:= cmd.Flags().GetString("env")
		port, _:= cmd.Flags().GetString("port")

		err, workdir := lib.BuildWorkDir(srcDir, dotenv, configName )
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("[WARN] This is a development server, don't use for production")
		log.Printf("[INFO] Listening for connection on port :%s",port)
		fs := dotFileHidingFileSystem{http.Dir(workdir)}
		log.Fatal(http.ListenAndServe(":"+port, http.FileServer(fs)))

	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringP("port", "p", "8080", "Listening port")
	serveCmd.Flags().StringP("env", "e", ".env", "Source dotenv file")
	serveCmd.Flags().StringP("configname", "c", "env-config.js", "Name of the generated config file")
}
