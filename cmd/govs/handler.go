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
	cmd.StringVar(&govs.CmdOpt.Typ, "t", "io", "type of the stats name")
	cmd.IntVar(&govs.CmdOpt.Id, "i", 0, "id of the stats object")

	/*
		// flush
		flags.NewCommand("flush", "Flush the virtual service", flush_handle, flag.ExitOnError)

		// zero
		cmd = flags.NewCommand("zero", "zero conters in Service/all", zero_handle, flag.ExitOnError)
		cmd.BoolVar(&govs.CmdOpt.TCP, "t", false, "tcp service")
		cmd.BoolVar(&govs.CmdOpt.UDP, "u", false, "udp service")

		// info
		flags.NewCommand("info", "show dpvs information", info_handle, flag.ExitOnError)

		// timeout
		cmd = flags.NewCommand("timeout", "show/set timeout", timeout_handle, flag.ExitOnError)
		cmd.StringVar(&govs.CmdOpt.Timeout_s, "set", "", "set tcp,tcp_fin,udp timeout")

		// list
		cmd = flags.NewCommand("list", "list -t|u host:[port]", list_handle, flag.ExitOnError)
		cmd.BoolVar(&govs.CmdOpt.TCP, "t", false, "tcp service")
		cmd.BoolVar(&govs.CmdOpt.UDP, "u", false, "udp service")
		cmd.BoolVar(&govs.CmdOpt.L, "G", false, "get local address")
		cmd.IntVar(&govs.CmdOpt.Num_services, "num_services", 0, "max service entries list")
		cmd.IntVar(&govs.CmdOpt.Number, "n", 0, "max laddr/dest entries list")

		// add
		cmd = flags.NewCommand("add", "add vs/rs/laddr", add_handle, flag.ExitOnError)
		cmd.UintVar(&govs.CmdOpt.Nic, "nic", 0, "the service addr bind nic port")
		cmd.BoolVar(&govs.CmdOpt.TCP, "t", false, "tcp service")
		cmd.BoolVar(&govs.CmdOpt.UDP, "u", false, "udp service")
		cmd.StringVar(&govs.CmdOpt.Sched_name, "sched", "rr", "the service sched name")
		cmd.UintVar(&govs.CmdOpt.Flags, "flags", 0, "the service flags")
		cmd.UintVar(&govs.CmdOpt.Timeout, "persistent", 0, "persistent service -persistent [timeout]")
		// adddest
		cmd.Var(&govs.CmdOpt.Daddr, "dest", "service-address is host[:port]")
		cmd.UintVar(&govs.CmdOpt.Dnic, "dest-nic", 0, "service-address out nic port(no need)")
		cmd.UintVar(&govs.CmdOpt.Conn_flags, "conn_flags", 0, "the conn flags")
		cmd.IntVar(&govs.CmdOpt.Weight, "weight", 0, "capacity of real server")
		cmd.UintVar(&govs.CmdOpt.U_threshold, "x", 0, "upper threshold of connections")
		cmd.UintVar(&govs.CmdOpt.L_threshold, "y", 0, "lower threshold of connections")
		// addladdr
		cmd.Var(&govs.CmdOpt.Lip, "laddr", "local-address is host")
		cmd.UintVar(&govs.CmdOpt.Lnic, "lnic", 0, "local addr bind nic port")

		// edit
		cmd = flags.NewCommand("edit", "edit vs/rs/laddr", edit_handle, flag.ExitOnError)
		cmd.UintVar(&govs.CmdOpt.Nic, "nic", 0, "the service addr bind nic port")
		cmd.BoolVar(&govs.CmdOpt.TCP, "t", false, "tcp service")
		cmd.BoolVar(&govs.CmdOpt.UDP, "u", false, "udp service")
		cmd.StringVar(&govs.CmdOpt.Sched_name, "sched", "rr", "the service sched name")
		cmd.UintVar(&govs.CmdOpt.Flags, "flags", 0, "the service flags")
		cmd.UintVar(&govs.CmdOpt.Timeout, "persistent", 0, "persistent service -persistent [timeout]")
		// editdest
		cmd.Var(&govs.CmdOpt.Daddr, "dest", "service-address is host[:port]")
		cmd.UintVar(&govs.CmdOpt.Dnic, "dest-nic", 0, "service-address out nic port(no need)")
		cmd.UintVar(&govs.CmdOpt.Conn_flags, "conn_flags", 0, "the conn flags")
		cmd.IntVar(&govs.CmdOpt.Weight, "weight", 0, "capacity of real server")
		cmd.UintVar(&govs.CmdOpt.U_threshold, "x", 0, "upper threshold of connections")
		cmd.UintVar(&govs.CmdOpt.L_threshold, "y", 0, "lower threshold of connections")
		// editladdr
		cmd.Var(&govs.CmdOpt.Lip, "laddr", "local-address is host")
		cmd.UintVar(&govs.CmdOpt.Lnic, "lnic", 0, "local addr bind nic port")

		// del
		cmd = flags.NewCommand("del", "del vs/rs/laddr", del_handle, flag.ExitOnError)
		cmd.BoolVar(&govs.CmdOpt.TCP, "t", false, "tcp service")
		cmd.BoolVar(&govs.CmdOpt.UDP, "u", false, "udp service")
		// deldest
		cmd.Var(&govs.CmdOpt.Daddr, "dest", "service-address is host[:port]")
		// delladdr
		cmd.Var(&govs.CmdOpt.Lip, "laddr", "local-address is host")
	*/
}

func version_handle(arg interface{}) {
	if err, version := govs.Get_version(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(version)
	}
}

func info_handle(arg interface{}) {
	if err, info := govs.Get_info(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(info)
	}
}

func timeout_handle(arg interface{}) {
	opt := arg.(*govs.CallOptions)
	if len(opt.Args) > 0 {
		opt.Opt.Addr.Set(opt.Args[0])
	}
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

func list_handle(arg interface{}) {
	opt := arg.(*govs.CallOptions)
	govs.Parse_service(opt)
	o := &opt.Opt

	if o.Addr.Ip != 0 {

		svc, err := govs.Get_service(o)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(govs.Svc_title())
		if !o.L {
			fmt.Println(govs.Dest_title())
		} else {
			fmt.Println(govs.Laddr_title())
		}
		fmt.Println(svc)
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

	svcs, err := govs.Get_services(o)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, svc := range svcs.Services {
		fmt.Println(govs.Svc_title())
		if !o.L {
			fmt.Println(govs.Dest_title())
		} else {
			fmt.Println(govs.Laddr_title())
		}
		fmt.Println(svc)
		o.Addr.Ip = svc.Addr
		o.Addr.Port = svc.Port
		o.Protocol = govs.Protocol(svc.Protocol)

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
	}
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
	switch govs.CmdOpt.Typ {
	case "io":
		relay, err := govs.Get_stats_io(govs.CmdOpt.Id)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(relay)
	case "worker":
		relay, err := govs.Get_stats_worker(govs.CmdOpt.Id)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(relay)
	case "dev":
		relay, err := govs.Get_stats_dev(govs.CmdOpt.Id)
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
