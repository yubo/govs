/*
 * Copyright 2016 Xiaomi Corporation. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 *
 * Authors:    Yu Bo <yubo@xiaomi.com>
 */
package govs

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

const (
	URL = "/tmp/dpvs.sock"
)

var (
	client *rpc.Client
	CmdOpt CmdOptions
)

func Vs_dial() error {
	conn, err := net.Dial("unix", URL)
	if err != nil {
		return err
	}
	client = jsonrpc.NewClient(conn)
	return nil
}

func Vs_close() {
	if client != nil {
		client.Close()
	}
}
