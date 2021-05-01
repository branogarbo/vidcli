/*
Copyright © 2021 Brian Longmore brianl.ext@gmail.com

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
package main

import (
	"context"
	"os"
	"syscall"

	"github.com/branogarbo/vidcli/cmd"
	gb "github.com/thecodeteam/goodbye"
)

func main() {
	ctx := context.Background()

	defer gb.Exit(ctx, -1)
	gb.Notify(ctx, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	cmd.Execute()
}
