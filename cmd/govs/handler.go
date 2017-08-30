/*
 * Copyright 2016 Xiaomi Corporation. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 *
 * Authors:    Yu Bo <yubo@xiaomi.com>
 */
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yubo/gotool/flags"
	"github.com/yubo/govs"
)

func init() {
	flags.CommandLine.Usage = fmt.Sprintf("Usage: %s COMMAND [OPTIONS] host[:port]\n\n",
		os.Args[0])

	// version
	flags.NewCommand("version", "show dpvs version information", version_handle, flag.ExitOnError)

	// status
	cmd := flags.NewCommand("stats", "get dpvs stats io stats", stats_handle, flag.ExitOnError)
	cmd.StringVar(&govs.CmdOpt.Typ, "t", "io", "type of the stats name(io/w/we/dev/ctl/mem)")
	cmd.IntVar(&govs.CmdOpt.Id, "i", -1, "id of the stats object")

	// flush
	flags.NewCommand("flush", "Flush the virtual service", flush_handle, flag.ExitOnError)

	// zero
	cmd = flags.NewCommand("zero", "zero conters in Service/all", zero_handle, flag.ExitOnError)
	cmd.StringVar(&govs.CmdOpt.TCP, "t", "", "tcp service")
	cmd.StringVar(&govs.CmdOpt.UDP, "u", "", "udp service")

	// timeout
	cmd = flags.NewCommand("timeout", "show/set timeout", timeout_handle, flag.ExitOnError)
	cmd.StringVar(&govs.CmdOpt.Timeout_s, "set", "", "set <tcp,tcp_fin,udp>")

	// list
	cmd = flags.NewCommand("list", "list -t|u host:[port]", list_handle, flag.ExitOnError)
	cmd.StringVar(&govs.CmdOpt.TCP, "t", "", "tcp service")
	cmd.StringVar(&govs.CmdOpt.UDP, "u", "", "udp service")
	cmd.BoolVar(&govs.CmdOpt.L, "G", false, "get local address")

	// add
	cmd = flags.NewCommand("add", "add vs/rs/laddr", add_handle, flag.ExitOnError)
	cmd.StringVar(&govs.CmdOpt.TCP, "t", "", "tcp service")
	cmd.StringVar(&govs.CmdOpt.UDP, "u", "", "udp service")
	cmd.Var(&govs.CmdOpt.Netmask, "m", "netmask default 0.0.0.0")
	cmd.StringVar(&govs.CmdOpt.Sched_name, "sched", "rr", "the service sched name rr/wrr")
	cmd.UintVar(&govs.CmdOpt.Flags, "flags", 0, "the service flags")

	// adddest
	cmd.Var(&govs.CmdOpt.Daddr, "dest", "service-address is host[:port]")
	cmd.UintVar(&govs.CmdOpt.Conn_flags, "conn_flags", 0, "the conn flags")
	cmd.IntVar(&govs.CmdOpt.Weight, "weight", 0, "capacity of real server")
	cmd.UintVar(&govs.CmdOpt.U_threshold, "x", 0, "upper threshold of connections")
	cmd.UintVar(&govs.CmdOpt.L_threshold, "y", 0, "lower threshold of connections")
	// addladdr
	cmd.Var(&govs.CmdOpt.Lip, "laddr", "local-address is host")

	// edit
	cmd = flags.NewCommand("edit", "edit vs/rs/laddr", edit_handle, flag.ExitOnError)
	cmd.StringVar(&govs.CmdOpt.TCP, "t", "", "tcp service")
	cmd.StringVar(&govs.CmdOpt.UDP, "u", "", "udp service")
	cmd.StringVar(&govs.CmdOpt.Sched_name, "sched", "rr", "the service sched name")
	cmd.UintVar(&govs.CmdOpt.Flags, "flags", 0, "the service flags")

	// editdest
	cmd.Var(&govs.CmdOpt.Daddr, "dest", "service-address is host[:port]")
	cmd.UintVar(&govs.CmdOpt.Conn_flags, "conn_flags", 0, "the conn flags")
	cmd.IntVar(&govs.CmdOpt.Weight, "weight", 0, "capacity of real server")
	cmd.UintVar(&govs.CmdOpt.U_threshold, "x", 0, "upper threshold of connections")
	cmd.UintVar(&govs.CmdOpt.L_threshold, "y", 0, "lower threshold of connections")
	// editladdr
	cmd.Var(&govs.CmdOpt.Lip, "laddr", "local-address is host")

	// del
	cmd = flags.NewCommand("del", "del vs/rs/laddr", del_handle, flag.ExitOnError)
	cmd.StringVar(&govs.CmdOpt.TCP, "t", "", "tcp service")
	cmd.StringVar(&govs.CmdOpt.UDP, "u", "", "udp service")
	// deldest
	cmd.Var(&govs.CmdOpt.Daddr, "dest", "service-address is host[:port]")
	// delladdr
	cmd.Var(&govs.CmdOpt.Lip, "laddr", "local-address is host")
}

func version_handle(arg interface{}) {
	if err, version := govs.Get_version(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(version)
	}
}

func info_handle(arg interface{}) {
	if err, info := govs.Get_version(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(info)
	}
}

func timeout_handle(arg interface{}) {
	opt := arg.(*govs.CallOptions)
	o := &opt.Opt

	if o.Timeout_s != "" {
		if timeout, err := govs.Set_timeout(o); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(timeout)
		}
	} else {
		if timeout, err := govs.Get_timeout(o); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(timeout)
		}
	}
}

func list_svc_handle(o *govs.CmdOptions) {

	ret, err := govs.Get_service(o)
	if err != nil {
		fmt.Println(err)
		return
	}

	if ret.Code != 0 {
		fmt.Println(ret.Msg)
		return
	}

	fmt.Println(govs.Svc_title())
	if !o.L {
		fmt.Println(govs.Dest_title())
	} else {
		fmt.Println(govs.Laddr_title())
	}

	fmt.Println(ret)
	if !o.L {
		dests, err := govs.Get_dests(o)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(dests)
	} else {
		laddrs, err := govs.Get_laddrs(o)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(laddrs)
	}
	return
}

func list_svcs_handle(o *govs.CmdOptions) {

	ret, err := govs.Get_services(o)

	if err != nil {
		fmt.Println(err)
		return
	}

	if ret.Code != 0 {
		fmt.Println(ret.Msg)
		return
	}

	fmt.Println(govs.Svc_title())
	if !o.L {
		fmt.Println(govs.Dest_title())
	} else {
		fmt.Println(govs.Laddr_title())
	}

	for _, svc := range ret.Services {
		fmt.Println(svc)
		o.Addr.Ip = svc.Addr
		o.Addr.Port = svc.Port
		o.Protocol = govs.Protocol(svc.Protocol)

		if !o.L {
			dests, err := govs.Get_dests(o)
			if err != nil || dests.Code != 0 ||
				len(dests.Dests) == 0 {
				//fmt.Println(err)
				continue
			}
			fmt.Println(dests)
		} else {
			laddrs, err := govs.Get_laddrs(o)
			if err != nil || laddrs.Code != 0 ||
				len(laddrs.Laddrs) == 0 {
				//fmt.Println(err)
				return
			}
			fmt.Println(laddrs)
		}
	}

}

func list_handle(arg interface{}) {
	opt := arg.(*govs.CallOptions)
	govs.Parse_service(opt)
	o := &opt.Opt

	if o.Addr.Ip != 0 {
		list_svc_handle(o)
		return
	}

	list_svcs_handle(o)
}

func flush_handle(arg interface{}) {
	if reply, err := govs.Set_flush(nil); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(reply)
	}
}

func zero_handle(arg interface{}) {
	opt := arg.(*govs.CallOptions)
	govs.Parse_service(opt)

	if reply, err := govs.Set_zero(&opt.Opt); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(reply)
	}

}

func add_handle(arg interface{}) {
	var err error
	var reply *govs.Vs_cmd_r

	opt := arg.(*govs.CallOptions)
	if err := govs.Parse_service(opt); err != nil {
		fmt.Println(err)
		return
	}
	o := &opt.Opt

	if o.Lip != 0 {
		reply, err = govs.Set_addladdr(o)
	} else if o.Daddr.Ip != 0 {
		reply, err = govs.Set_adddest(o)
	} else {
		reply, err = govs.Set_add(o)
	}

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(reply)
	}
}

func edit_handle(arg interface{}) {
	var err error
	var reply *govs.Vs_cmd_r

	opt := arg.(*govs.CallOptions)
	if err := govs.Parse_service(opt); err != nil {
		fmt.Println(err)
		return
	}
	o := &opt.Opt

	if o.Daddr.Ip != 0 {
		reply, err = govs.Set_editdest(o)
	} else {
		reply, err = govs.Set_edit(o)
	}

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(reply)
	}
}

func del_handle(arg interface{}) {
	var err error
	var reply *govs.Vs_cmd_r

	opt := arg.(*govs.CallOptions)
	if err := govs.Parse_service(opt); err != nil {
		fmt.Println(err)
		return
	}
	o := &opt.Opt

	if o.Lip != 0 {
		reply, err = govs.Set_delladdr(o)
	} else if o.Daddr.Ip != 0 {
		reply, err = govs.Set_deldest(o)
	} else {
		reply, err = govs.Set_del(o)
	}

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(reply)
	}
}

func stats_handle(arg interface{}) {
	id := govs.CmdOpt.Id

	switch govs.CmdOpt.Typ {
	case "io":
		relay, err := govs.Get_stats_io(id)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%+v\n", relay)
	case "w":
		relay, err := govs.Get_stats_worker(id)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(relay)
	case "we":
		relay, err := govs.Get_estats_worker(id)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(relay)
	case "dev":
		relay, err := govs.Get_stats_dev(id)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(relay)
	case "ctl":
		relay, err := govs.Get_stats_ctl()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(relay)
	case "mem":
		relay, err := govs.Get_stats_mem()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(relay)
	default:
		fmt.Println("govs stats -t io/worker/dev/ctl")
	}
}
