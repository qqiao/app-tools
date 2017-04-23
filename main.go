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

package main

import (
	"context"
	"flag"

	"github.com/qqiao/app-tools/appengine"
	"github.com/qqiao/app-tools/buildinfo"
	"github.com/qqiao/cli"
)

var component = cli.Component{
	UsageLine: "app-toos command",
	Components: []*cli.Component{
		appengine.NewComponent(),
		buildinfo.NewComponent(),
	},
	Run: func(context.Context, *cli.Component, []string) {},
}

func init() {
	flag.Usage = component.Usage
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		return
	}

	// workingDir, err := os.Getwd()
	// if nil != err {
	// 	fmt.Fprintf(os.Stderr, "Unable to determine working dir. Error: %s\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Fprintf(os.Stdout, "Working Directory: %s.\n", workingDir)

	name := flag.Arg(0)
	for _, comp := range component.Components {
		if name == comp.Name() {
			if comp.Runnable() {
				args := flag.Args()[1:]
				ctx := context.Background()
				comp.Flag.Usage = comp.Usage
				comp.Run(ctx, comp, args)
				return
			}
		}
	}
	flag.Usage()
}
