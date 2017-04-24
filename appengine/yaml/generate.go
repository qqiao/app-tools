// Copyright 2017 Qian Qiao
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package yaml

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/qqiao/cli"
)

func newGenerateComponent() *cli.Component {
	comp := &cli.Component{
		UsageLine: "generate",

		Run: func(ctx context.Context, comp *cli.Component, args []string) {
			if flag.ErrHelp == comp.Flag.Parse(args) {
				return
			}

			tmplFlag = sanitizePath(tmplFlag,
				fmt.Sprintln("Please specify template file..."), comp)
			outputFlag = sanitizePath(outputFlag,
				fmt.Sprintln("Please specify output file..."), comp)

			// TODO sort out the data part
			ch := Generate(tmplFlag, outputFlag, nil)
			if !<-ch {
				os.Exit(1)
			}
		},
	}

	comp.Flag.StringVar(&tmplFlag, "tmpl", "",
		"path of the template for the .yaml file to generate")
	comp.Flag.StringVar(&outputFlag, "o", "",
		"path of the .yaml file to create")

	return comp
}

// Generate generates the .yaml file based on the template and the data.
// writes generate file to the specified output file
func Generate(tmplPath, outputPath string, data interface{}) <-chan bool {
	done := make(chan bool, 1)

	go func() {
		defer close(done)

		tmpl, err := template.ParseFiles(tmplFlag)
		if nil != err {
			fmt.Fprintf(os.Stderr, "Unable to parse template. Error: %s",
				err.Error())
			done <- false
			return
		}

		var buf bytes.Buffer

		if err = tmpl.Execute(&buf, data); nil != err {
			fmt.Fprintf(os.Stderr, "Unable to render %s. Error: %s",
				outputFlag, err.Error())
			done <- false
			return
		}

		if err = ioutil.WriteFile(outputPath, buf.Bytes(), 0644); nil != err {
			fmt.Fprintf(os.Stderr, "Unable to write to %s. Error: %s",
				outputFlag, err.Error())
			done <- false
			return
		}

		done <- true
	}()

	return done
}

func sanitizePath(path, errMsg string, comp *cli.Component) string {
	if "" == path {
		fmt.Fprintf(os.Stderr, errMsg)
		comp.Usage()
		os.Exit(2)
	}

	if !filepath.IsAbs(path) {
		cwd, err := os.Getwd()
		if nil != err {
			fmt.Fprintf(os.Stderr,
				"Unable to get current working directory. Error: %s", err.Error())
			os.Exit(1)
		}
		path = filepath.Join(cwd, path)
	}

	return filepath.Clean(path)
}
