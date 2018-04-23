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

// Package buildinfo provides tool for generating and managing web application
// build information
package buildinfo

import (
	"github.com/qqiao/cli"

	"github.com/qqiao/app-tools/buildinfo/generate"
)

// NewComponent returns command line component that can be used to form
// complex commandline application
func NewComponent() *cli.Component {
	return &cli.Component{
		UsageLine: "buildinfo command",
		Short:     "tools for manipulating the build information",
		Run:       cli.Passthrough,
		Components: []*cli.Component{
			generate.NewComponent(),
		},
	}
}
