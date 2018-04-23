// Copyright 2018 Qian Qiao
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

package execute

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/qqiao/cli"
)

var (
	outputFlag string
	tmplFlag   string
	dataFlag   string
)

// NewComponent creates a new CLI component for executing a template
func NewComponent() *cli.Component {
	comp := &cli.Component{
		UsageLine: "execute",
	}

	flagSet := comp.FlagSet()

	comp.Run = func(ctx context.Context, comp *cli.Component, args []string) {
		if flag.ErrHelp == flagSet.Parse(args) {
			return
		}

		tmplFlag = sanitizePath(tmplFlag,
			fmt.Sprintln("Please specify template file..."), comp)
		outputFlag = sanitizePath(outputFlag,
			fmt.Sprintln("Please specify output file..."), comp)

		var data map[string]interface{}
		if err := json.Unmarshal([]byte(dataFlag), &data); nil != err {
			fmt.Fprintf(os.Stderr, "Unable to parse data. Error: %s",
				err.Error())
			os.Exit(1)
		}

		ch := Execute(tmplFlag, outputFlag, data)
		if !<-ch {
			os.Exit(1)
		}
	}

	flagSet.StringVar(&tmplFlag, "t", "",
		"path to the template file")
	flagSet.StringVar(&outputFlag, "o", "",
		"path to the output file")
	flagSet.StringVar(&dataFlag, "data", "",
		"data for the template in JSON format")

	return comp
}

// Execute runs template based on the data.
// writes generate file to the specified output file
func Execute(tmplPath, outputPath string, data interface{}) <-chan bool {
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
