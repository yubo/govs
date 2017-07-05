/*
 * Copyright 2016 Xiaomi Corporation. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 *
 * Authors:    Yu Bo <yubo@xiaomi.com>
 */
package main

import (
	"fmt"

	"github.com/yubo/gotool/flags"
	"github.com/yubo/govs"
)

func main() {

	flags.Parse()

	cmd := flags.CommandLine.Cmd
	if cmd != nil && cmd.Action != nil {
		err := govs.Vs_dial()
		if err != nil {
			fmt.Println("cannot connection to dpvs server")
			return
		}

		defer govs.Vs_close()
		cmd.Action(&govs.CallOptions{Opt: govs.CmdOpt,
			Args: cmd.Flag.Args()})
	} else {
		flags.Usage()
	}

}
