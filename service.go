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

type Vs_service_user struct {
	Nic        uint8
	Protocol   uint8
	Addr       Be32
	Port       Be16
	Sched_name string
	Flags      uint
	Timeout    uint
	Netmask    Be32
	Number     int /* max list laddr/dests */
}

type Vs_service_user_r struct {
	Protocol   uint8
	Addr       Be32
	Port       Be16
	Sched_name string
	Flags      uint32
	Timeout    uint32
	Netmask    Be32
	Conns      uint64
	Inpkts     uint64
	Outpkts    uint64
	Inbytes    uint64
	Outbytes   uint64
	Num_dests  uint32
	Num_laddrs uint32
}

const (
	fmt_svc_t = "%5s %21s %8s %8s %15s %6s %6s %5s %7s %10s %10s %10s %10s"
	fmt_svc   = "%5s %21s %08x %8d %15s %6d %6d %5s %7d %10d %10d %10d %10d"
)

func Svc_title() string {
	return fmt.Sprintf(fmt_svc_t,
		"Proto", "Addr:Port ", "Flags",
		"Timeout", "Netmask", "dests", "laddrs", "Sched",
		"Conns", "Inpkts", "Outpkts", "Inbytes", "Outbytes")
}

func (svc Vs_service_user_r) String() string {
	return fmt.Sprintf(fmt_svc,
		get_protocol_name(svc.Protocol),
		fmt.Sprintf("%s:%s", svc.Addr.String(), svc.Port.String()), svc.Flags,
		svc.Timeout, svc.Netmask.String(),
		svc.Num_dests, svc.Num_laddrs, svc.Sched_name,
		svc.Conns, svc.Inpkts, svc.Outpkts,
		svc.Inbytes, svc.Outbytes)
}

type Vs_list_q struct {
	Cmd          int
	Num_services int
	Service      Vs_service_user
	Dest         Vs_dest_user
	Laddr        Vs_laddr_user
}

type Vs_list_service_r struct {
	Code    int
	Msg     string
	Service Vs_service_user_r
}

func (r Vs_list_service_r) String() string {
	if r.Code == 0 {
		return r.Service.String()
	} else {
		return fmt.Sprintf("%s:%s", Ecode(r.Code), r.Msg)
	}
}

type Vs_list_services_r struct {
	Code         int
	Msg          string
	Num_services int
	Services     []Vs_service_user_r
}

func (r Vs_list_services_r) String() (s string) {
	if r.Code != 0 {
		return fmt.Sprintf("%s:%s", Ecode(r.Code), r.Msg)
	}
	for _, svc := range r.Services {
		s += fmt.Sprintf("%s\n", svc)
	}
	return strings.TrimRight(s, "\n")
}

type Vs_service_q struct {
	Cmd     int
	Service Vs_service_user
}

func Get_services(o *CmdOptions) (*Vs_list_services_r, error) {
	var reply Vs_list_services_r

	args := Vs_list_q{
		Cmd: VS_CMD_GET_SERVICES,
	}

	err := client.Call("api", args, &reply)
	return &reply, err
}

func Get_service(o *CmdOptions) (*Vs_list_service_r, error) {
	var reply Vs_list_service_r
	args := Vs_list_q{
		Cmd: VS_CMD_GET_SERVICE,
		Service: Vs_service_user{
			Addr:     o.Addr.Ip,
			Port:     o.Addr.Port,
			Protocol: uint8(o.Protocol),
		},
	}
	fmt.Printf("ip:%s, port:%d, protocol: %d\n",
		args.Service.Addr.String(), args.Service.Port, args.Service.Protocol)

	err := client.Call("api", args, &reply)
	return &reply, err
}

func Set_add(o *CmdOptions) (*Vs_cmd_r, error) {
	var reply Vs_cmd_r
	args := Vs_service_q{
		Cmd: VS_CMD_NEW_SERVICE,
		Service: Vs_service_user{
			Nic:        uint8(o.Nic),
			Protocol:   uint8(o.Protocol),
			Addr:       o.Addr.Ip,
			Port:       o.Addr.Port,
			Sched_name: o.Sched_name,
			Flags:      o.Flags,
			Timeout:    o.Timeout,
			Netmask:    o.Netmask,
		},
	}

	err := client.Call("api", args, &reply)
	return &reply, err
}

func Set_edit(o *CmdOptions) (*Vs_cmd_r, error) {
	var reply Vs_cmd_r
	args := Vs_service_q{
		Cmd: VS_CMD_SET_SERVICE,
		Service: Vs_service_user{
			Nic:        uint8(o.Nic),
			Protocol:   uint8(o.Protocol),
			Addr:       o.Addr.Ip,
			Port:       o.Addr.Port,
			Sched_name: o.Sched_name,
			Flags:      o.Flags,
			Timeout:    o.Timeout,
			Netmask:    o.Netmask,
		},
	}

	err := client.Call("api", args, &reply)
	return &reply, err
}

func Set_del(o *CmdOptions) (*Vs_cmd_r, error) {
	var reply Vs_cmd_r
	args := Vs_service_q{
		Cmd: VS_CMD_DEL_SERVICE,
		Service: Vs_service_user{
			Protocol: uint8(o.Protocol),
			Addr:     o.Addr.Ip,
			Port:     o.Addr.Port,
		},
	}

	err := client.Call("api", args, &reply)
	return &reply, err
}
