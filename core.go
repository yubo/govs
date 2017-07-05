/*
 * Copyright 2016 Xiaomi Corporation. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 *
 * Authors:    Yu Bo <yubo@xiaomi.com>
 */
package govs

import (
	"net/rpc"
	"net/rpc/jsonrpc"
)

const (
	URL = "127.0.0.1:1105"
)

var (
	client *rpc.Client
	CmdOpt CmdOptions
)

func Vs_dial() (err error) {
	client, err = jsonrpc.Dial("tcp", URL)
	return err
}

func Vs_close() {
	if client != nil {
		client.Close()
	}
}
