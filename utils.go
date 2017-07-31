/*
 * Copyright 2016 Xiaomi Corporation. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 *
 * Authors:    Yu Bo <yubo@xiaomi.com>
 */
package govs

import (
	"encoding/binary"
	"fmt"
	"net"
	"unsafe"
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

func be32_to_addr(b Be32) string {
	return u32_to_addr(Ntohl(b))
}

func u32_to_addr(u uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		u&0xff000000>>24,
		u&0xff0000>>16,
		u&0xff00>>8,
		u&0xff)
}

func Ntohl(i Be32) uint32 {
	return binary.BigEndian.Uint32((*(*[4]byte)(unsafe.Pointer(&i)))[:])
}

func Htonl(i uint32) Be32 {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, i)
	return *(*Be32)(unsafe.Pointer(&b[0]))
}

func Ntohs(i Be16) uint16 {
	return binary.BigEndian.Uint16((*(*[2]byte)(unsafe.Pointer(&i)))[:])
}

func Htons(i uint16) Be16 {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, i)
	return *(*Be16)(unsafe.Pointer(&b[0]))
}

func get_protocol_name(p uint8) string {
	if p == IPPROTO_TCP {
		return "tcp"
	} else if p == IPPROTO_UDP {
		return "udp"
	}
	return "unknown"
}

func Parse_service(o *CallOptions) error {
	var addr string

	if o.Opt.UDP != "" {
		addr = o.Opt.UDP
		o.Opt.Protocol = IPPROTO_UDP
	}

	if o.Opt.TCP != "" {
		addr = o.Opt.TCP
		o.Opt.Protocol = IPPROTO_TCP
	}

	if err := o.Opt.Addr.Set(addr); err != nil {
		return err
	}

	return nil
}
