/*
 * Copyright 2016 Xiaomi Corporation. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 *
 * Authors:    Yu Bo <yubo@xiaomi.com>
 */
package govs

func Get_version() (error, *Vs_version_r) {
	var reply Vs_version_r
	args := Vs_cmd_q{VS_CMD_GET_INFO}

	err := client.Call("api", args, &reply)
	return err, &reply
}

func Get_info() (error, *Vs_info_r) {
	var reply Vs_info_r
	args := Vs_cmd_q{VS_CMD_GET_INFO}

	err := client.Call("api", args, &reply)
	return err, &reply
}

func Get_services(o *CmdOptions) (error, *Vs_list_services_r) {
	var reply Vs_list_services_r

	args := Vs_list_q{
		Cmd:          VS_CMD_GET_SERVICES,
		Num_services: o.Num_services,
	}

	err := client.Call("api", args, &reply)
	return err, &reply
}

func Get_service(o *CmdOptions) (error, *Vs_list_service_r) {
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
	return err, &reply
}

func Get_dests(o *CmdOptions) (error, *Vs_list_dests_r) {
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
	return err, &reply
}

func Get_laddrs(o *CmdOptions) (error, *Vs_list_laddrs_r) {
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
	return err, &reply
}

func Get_timeout(o *CmdOptions) (error, *Vs_timeout_r) {
	args := Vs_timeout_q{Cmd: VS_CMD_GET_CONFIG}
	reply := &Vs_timeout_r{}

	err := client.Call("api", args, reply)
	return err, reply
}

func Get_stats_io(id int) (error, *Vs_stats_io_r) {
	args := Vs_stats_q{Type: VS_STATS_IO, Id: id}
	reply := &Vs_stats_io_r{}

	err := client.Call("stats", args, reply)
	return err, reply
}

func Get_stats_worker(id int) (error, *Vs_stats_worker_r) {
	args := Vs_stats_q{Type: VS_STATS_WORKER, Id: id}
	reply := &Vs_stats_worker_r{}

	err := client.Call("stats", args, reply)
	return err, reply
}

func Get_stats_dev(id int) (error, *Vs_stats_dev_r) {
	args := Vs_stats_q{Type: VS_STATS_DEV, Id: id}
	reply := &Vs_stats_dev_r{}

	err := client.Call("stats", args, reply)
	return err, reply
}

func Get_stats_ctl() (error, *Vs_stats_ctl_r) {
	args := Vs_stats_q{Type: VS_STATS_CTL}
	reply := &Vs_stats_ctl_r{}

	err := client.Call("stats", args, reply)
	return err, reply
}
