/*
 * Copyright 2016 Xiaomi Corporation. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 *
 * Authors:    Yu Bo <yubo@xiaomi.com>
 */
package govs

import (
	"fmt"
	"strings"
)

type Vs_laddr_user struct {
	Nic  uint8
	Addr Be32
}

type Vs_laddr_user_r struct {
	Addr          Be32
	Conn_counts   uint32
	Port_conflict uint64
}

func Laddr_title() string {
	return fmt.Sprintf("    %15s %8s %8s",
		"Addr", "Conn_counts", "Port_conflict")
}

func (l Vs_laddr_user_r) String() string {
	return fmt.Sprintf("    %15s %8d %8d",
		l.Addr.String(),
		l.Conn_counts, l.Port_conflict)
}

type Vs_list_laddrs_r struct {
	Code   int
	Msg    string
	Laddrs []Vs_laddr_user_r
}

func (r Vs_list_laddrs_r) String() string {
	var s string
	if r.Code != 0 {
		return fmt.Sprintf("%s:%s", Ecode(r.Code), r.Msg)
	}

	for _, laddr := range r.Laddrs {
		s += fmt.Sprintf("%s\n", laddr)
	}
	return strings.TrimRight(s, "\n")
}

func Get_laddrs(o *CmdOptions) (*Vs_list_laddrs_r, error) {
	var reply Vs_list_laddrs_r
	args := Vs_list_q{
		Cmd: VS_CMD_GET_LADDR,
		Service: Vs_service_user{
			Addr:     o.Addr.Ip,
			Port:     o.Addr.Port,
			Protocol: uint8(o.Protocol),
			Number:   o.Number,
		},
	}

	err := client.Call("api", args, &reply)
	return &reply, err
}

type Vs_laddr_q struct {
	Cmd     int
	Service Vs_service_user
	Laddr   Vs_laddr_user
}

func Set_addladdr(o *CmdOptions) (*Vs_cmd_r, error) {
	var reply Vs_cmd_r
	args := Vs_laddr_q{
		Cmd: VS_CMD_NEW_LADDR,
		Service: Vs_service_user{
			Protocol: uint8(o.Protocol),
			Addr:     o.Addr.Ip,
			Port:     o.Addr.Port,
		},
		Laddr: Vs_laddr_user{
			Nic:  uint8(o.Lnic),
			Addr: o.Lip,
		},
	}

	err := client.Call("api", args, &reply)
	return &reply, err
}

func Set_delladdr(o *CmdOptions) (*Vs_cmd_r, error) {
	var reply Vs_cmd_r
	args := Vs_laddr_q{
		Cmd: VS_CMD_DEL_LADDR,
		Service: Vs_service_user{
			Protocol: uint8(o.Protocol),
			Addr:     o.Addr.Ip,
			Port:     o.Addr.Port,
		},
		Laddr: Vs_laddr_user{
			Nic:  uint8(o.Lnic),
			Addr: o.Lip,
		},
	}

	err := client.Call("api", args, &reply)
	return &reply, err
}
