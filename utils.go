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
	"net"
)

func max_len(is ...[]int64) (l int) {
	for _, v := range is {
		if l < len(v) {
			l = len(v)
		}
	}
	return l
}

func ipToU32(p net.IP) uint32 {
	// If IPv4, use dotted notation.
	if p4 := p.To4(); len(p4) == net.IPv4len {
		return uint32(p4[0])<<24 | uint32(p4[1])<<16 |
			uint32(p4[2])<<8 | uint32(p4[3])
	}
	return 0
}

func u32_to_addr(u uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		u&0xff000000>>24,
		u&0xff0000>>16,
		u&0xff00>>8,
		u&0xff)
}

func get_protocol_name(p uint8) string {
	if p == IPPROTO_TCP {
		return "tcp"
	} else if p == IPPROTO_UDP {
		return "udp"
	}
	return "unknown"
}

func Svc_title() string {
	return fmt.Sprintf("%3s %5s %21s %8s "+
		"%7s %15s %6s %6s %5s "+
		"%7s %10s %10s %10s %10s",
		"Nic", "Proto", "Addr:Port ", "Flags",
		"Timeout", "Netmask", "dests", "laddrs", "Sched",
		"Conns", "Inpkts", "Outpkts", "Inbytes", "Outbytes")
}

func (svc Vs_service_user_r) String() string {
	return fmt.Sprintf("%3d %5s %15s:%-5d %08x "+
		"%7d %15s %6d %6d %5s "+
		"%7s %10s %10s %10s %10s",
		svc.Nic, get_protocol_name(svc.Protocol),
		u32_to_addr(svc.Addr), svc.Port, svc.Flags,
		svc.Timeout, u32_to_addr(svc.Netmask),
		svc.Num_dests, svc.Num_laddrs, svc.Sched_name,
		svc.Conns, svc.Inpkts, svc.Outpkts,
		svc.Inbytes, svc.Outbytes)
}

func Dest_title() string {
	return fmt.Sprintf("    %3s %21s %8s "+
		"%7s %15s %10s %10s %10s "+
		"%7s %10s %10s %10s %10s",
		"Nic", "Addr:Port ", "Conn_flags",
		"Weight", "threshold", "Activeconns", "Inactconns", "Persistent",
		"Conns", "Inpkts", "Outpkts", "Inbytes", "Outbytes")
}

func (d Vs_dest_user_r) String() string {
	return fmt.Sprintf(" -> %3d %15s:%-5d %08x "+
		"%7d %7d-%-7d %10d %10d %10d "+
		"%7s %10s %10s %10s %10s",
		d.Nic, u32_to_addr(d.Addr), d.Port, d.Conn_flags,
		d.Weight, d.L_threshold, d.U_threshold,
		d.Activeconns, d.Inactconns, d.Persistent,
		d.Conns, d.Inpkts, d.Outpkts, d.Inbytes, d.Outbytes)
}

func Laddr_title() string {
	return fmt.Sprintf("    %3s %15s %8s %8s",
		"Nic", "Addr", "Conn_counts", "Port_conflict")
}

func (l Vs_laddr_user_r) String() string {
	return fmt.Sprintf("    %3d %15s %8d %8s",
		l.Nic, u32_to_addr(l.Addr),
		l.Conn_counts, l.Port_conflict)
}

func Parse_service(o *CallOptions) error {
	if len(o.Args) > 0 {
		if err := o.Opt.Addr.Set(o.Args[0]); err != nil {
			return err
		}
		if o.Opt.TCP {
			o.Opt.Protocol = IPPROTO_TCP
		} else {
			o.Opt.Protocol = IPPROTO_UDP
		}
		if o.Opt.TCP == o.Opt.UDP {
			o.Opt.Protocol = IPPROTO_TCP
		}
		return nil
	}
	return fmt.Errorf("service syntax error")
}

func ctl_state_name(s int) byte {
	switch s {
	case VS_CTL_S_SYNC:
		return 's'
	case VS_CTL_S_PENDING:
		return 'p'
	default:
		return '-'
	}
}
