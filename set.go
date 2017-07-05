/*
 * Copyright 2016 Xiaomi Corporation. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 *
 * Authors:    Yu Bo <yubo@xiaomi.com>
 */
package govs

func Set_flush(o *CmdOptions) (error, *Vs_cmd_r) {
	var reply Vs_cmd_r
	args := Vs_cmd_q{VS_CMD_FLUSH}

	err := client.Call("api", args, &reply)
	return err, &reply
}

func Set_timeout(o *CmdOptions) (error, *Vs_cmd_r) {
	var reply Vs_cmd_r
	args := Vs_timeout_q{Cmd: VS_CMD_SET_CONFIG}

	if err := args.Set(o.Timeout_s); err != nil {
		return err, nil
	}

	err := client.Call("api", args, &reply)
	return err, &reply
}

func Set_zero(o *CmdOptions) (error, *Vs_cmd_r) {
	var reply Vs_cmd_r
	args := Vs_service_q{
		Cmd: VS_CMD_ZERO,
		Service: Vs_service_user{
			Addr:     o.Addr.Ip,
			Port:     o.Addr.Port,
			Protocol: uint8(o.Protocol),
		},
	}

	err := client.Call("api", args, &reply)
	return err, &reply
}

func Set_add(o *CmdOptions) (error, *Vs_cmd_r) {
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
			Netmask:    uint32(o.Netmask),
		},
	}

	err := client.Call("api", args, &reply)
	return err, &reply
}

func Set_edit(o *CmdOptions) (error, *Vs_cmd_r) {
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
			Netmask:    uint32(o.Netmask),
		},
	}

	err := client.Call("api", args, &reply)
	return err, &reply
}

func Set_del(o *CmdOptions) (error, *Vs_cmd_r) {
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
	return err, &reply
}

func Set_adddest(o *CmdOptions) (error, *Vs_cmd_r) {
	var reply Vs_cmd_r
	args := Vs_dest_q{
		Cmd: VS_CMD_NEW_DEST,
		Service: Vs_service_user{
			Protocol: uint8(o.Protocol),
			Addr:     o.Addr.Ip,
			Port:     o.Addr.Port,
		},
		Dest: Vs_dest_user{
			Nic:         uint8(o.Dnic),
			Addr:        o.Daddr.Ip,
			Port:        o.Daddr.Port,
			Conn_flags:  o.Conn_flags,
			Weight:      o.Weight,
			U_threshold: uint32(o.U_threshold),
			L_threshold: uint32(o.L_threshold),
		},
	}

	err := client.Call("api", args, &reply)
	return err, &reply
}

func Set_editdest(o *CmdOptions) (error, *Vs_cmd_r) {
	var reply Vs_cmd_r
	args := Vs_dest_q{
		Cmd: VS_CMD_SET_DEST,
		Service: Vs_service_user{
			Protocol: uint8(o.Protocol),
			Addr:     o.Addr.Ip,
			Port:     o.Addr.Port,
		},
		Dest: Vs_dest_user{
			Nic:         uint8(o.Dnic),
			Addr:        o.Daddr.Ip,
			Port:        o.Daddr.Port,
			Conn_flags:  o.Conn_flags,
			Weight:      o.Weight,
			U_threshold: uint32(o.U_threshold),
			L_threshold: uint32(o.L_threshold),
		},
	}

	err := client.Call("api", args, &reply)
	return err, &reply
}

func Set_deldest(o *CmdOptions) (error, *Vs_cmd_r) {
	var reply Vs_cmd_r
	args := Vs_dest_q{
		Cmd: VS_CMD_DEL_DEST,
		Service: Vs_service_user{
			Protocol: uint8(o.Protocol),
			Addr:     o.Addr.Ip,
			Port:     o.Addr.Port,
		},
		Dest: Vs_dest_user{
			Addr: o.Daddr.Ip,
			Port: o.Daddr.Port,
		},
	}

	err := client.Call("api", args, &reply)
	return err, &reply
}

func Set_addladdr(o *CmdOptions) (error, *Vs_cmd_r) {
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
			Addr: uint32(o.Lip),
		},
	}

	err := client.Call("api", args, &reply)
	return err, &reply
}

func Set_delladdr(o *CmdOptions) (error, *Vs_cmd_r) {
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
			Addr: uint32(o.Lip),
		},
	}

	err := client.Call("api", args, &reply)
	return err, &reply
}
