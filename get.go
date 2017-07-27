/*
 * Copyright 2016 Xiaomi Corporation. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 *
 * Authors:    Yu Bo <yubo@xiaomi.com>
 */
package govs

func Get_version() (*Vs_version_r, error) {
	var reply Vs_version_r
	args := Vs_cmd_q{VS_CMD_GET_INFO}

	err := client.Call("api", args, &reply)
	return &reply, err
}

func Get_info() (*Vs_info_r, error) {
	var reply Vs_info_r
	args := Vs_cmd_q{VS_CMD_GET_INFO}

	err := client.Call("api", args, &reply)
	return &reply, err
}

func Get_services(o *CmdOptions) (*Vs_list_services_r, error) {
	var reply Vs_list_services_r

	args := Vs_list_q{
		Cmd:          VS_CMD_GET_SERVICES,
		Num_services: o.Num_services,
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

	err := client.Call("api", args, &reply)
	return &reply, err
}

func Get_dests(o *CmdOptions) (*Vs_list_dests_r, error) {
	var reply Vs_list_dests_r
	args := Vs_list_q{
		Cmd: VS_CMD_GET_DEST,
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

func Get_timeout(o *CmdOptions) (*Vs_timeout_r, error) {
	args := Vs_timeout_q{Cmd: VS_CMD_GET_CONFIG}
	reply := &Vs_timeout_r{}

	err := client.Call("api", args, reply)
	return reply, err
}
