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
	"context"
	"flag"

	"github.com/qqiao/cli"
)

var (
	tmplFlag   string
	outputFlag string
)

// NewComponent returns the cli component for all yaml related commands.
func NewComponent() *cli.Component {
	return &cli.Component{
		UsageLine: "yaml command",
		Short:     "tools for generating Google App Engine .yaml files",

		Components: []*cli.Component{
			newGenerateComponent(),
		},

		Run: func(ctx context.Context, comp *cli.Component, args []string) {
			if flag.ErrHelp == comp.Flag.Parse(args) {
				return
			}

			if comp.Flag.NArg() < 1 {
				comp.Flag.Usage()
				return
			}

			a := comp.Flag.Args()
			for _, c := range comp.Components {
				if a[0] == c.Name() {
					if c.Runnable() {
						c.Flag.Usage = c.Usage
						c.Run(ctx, c, a[1:])
						return
					}
				}
			}
			comp.Flag.Usage()
		},
	}
}
