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

// App tools is a set of convenient tools that can be used to generate web
// application files.
package main

import (
	"context"
	"flag"

	"github.com/qqiao/app-tools/buildinfo"
	"github.com/qqiao/app-tools/template"

	"github.com/qqiao/cli"
)

var component = cli.Component{
	UsageLine: "app-toos command",
	Components: []*cli.Component{
		buildinfo.NewComponent(),
		template.NewComponent(),
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
