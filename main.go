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
	"os"

	"github.com/qqiao/app-tools/buildinfo"
	"github.com/qqiao/app-tools/template"

	"github.com/qqiao/cli"
)

var component = &cli.Component{
	UsageLine: "app-toos command",
	Long: `
App Tools is a set of tools designed to make the development of web
applications easier.`,
	Components: []*cli.Component{
		buildinfo.NewComponent(),
		template.NewComponent(),
	},
	Run: cli.Passthrough,
}

func main() {
	component.Run(context.Background(), component, os.Args[1:])
}
