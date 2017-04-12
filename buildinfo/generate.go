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

package buildinfo

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/qqiao/buildinfo"
	"github.com/qqiao/cli"
)

// Generate generates build information.
func Generate(path string) <-chan bool {
	done := make(chan bool, 1)

	go func() {
		defer close(done)

		rCh := gitRevision()

		bytes, err := json.Marshal(buildinfo.BuildInfo{
			BuildTime: time.Now().UnixNano(),
			Revision:  <-rCh,
		})
		if nil != err {
			fmt.Fprintf(os.Stderr,
				"Unable to marshal build information. Error: %s\n",
				err)
			done <- false
			return
		}

		if err = ioutil.WriteFile(path, bytes, 0644); nil != err {
			fmt.Fprintf(os.Stderr,
				"Unable to write to build_info.json. Error: %s",
				err)
			done <- false
			return
		}

		done <- true
	}()

	return done
}

func gitRevision() <-chan string {
	out := make(chan string, 1)

	go func() {
		defer close(out)

		// check if we have GIT
		if _, err := exec.LookPath("git"); nil != err {
			fmt.Fprintln(os.Stderr, "Cannot find git in your PATH environmental variable")
			out <- "UNKNOWN"
			return
		}

		r, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
		if nil != err {
			fmt.Fprintf(os.Stderr, "Unable to figuring out repository revision: %s\n", err)
			out <- "UNKNOWN"
			return
		}
		out <- strings.TrimSpace(string(r))
	}()

	return out
}

// newGenerateComponent returns a cli component that is for generating a new
// build info.
func newGenerateComponent() *cli.Component {
	comp := &cli.Component{
		UsageLine: "generate [-f file_name]",
		Short:     "generates a new build information file",
		Run: func(ctx context.Context, comp *cli.Component, args []string) {
			comp.Flag.Parse(args)

			outputFile := fileFlag
			if !filepath.IsAbs(outputFile) {
				cwd, err := os.Getwd()
				if nil != err {
					fmt.Fprintf(os.Stderr,
						"Unable to determine current working directory. Error: %s\n",
						err.Error())
					return
				}
				outputFile = filepath.Clean(filepath.Join(cwd, outputFile))
			}

			if !<-Generate(outputFile) {
				os.Exit(1)
			}
		},
	}

	comp.Flag.StringVar(&fileFlag, "f", "build_info.json", "file to write to")

	return comp
}
