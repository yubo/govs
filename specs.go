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
	"strconv"

	"errors"
	"fmt"
	"strings"
)

const (
	VS_CTL_S_SYNC = iota
	VS_CTL_S_PENDING
	VS_CTL_S_LAST
)

const (
	VS_STATS_IO = iota
	VS_STATS_WORKER
	VS_STATS_DEV
	VS_STATS_CTL
	VS_STATS_MEM
	__VS_STATS_MAX
)

const (
	VS_CMD_UNSPEC = iota

	VS_CMD_NEW_SERVICE  /* add service */
	VS_CMD_SET_SERVICE  /* modify service */
	VS_CMD_DEL_SERVICE  /* delete service */
	VS_CMD_GET_SERVICE  /* get service info */
	VS_CMD_GET_SERVICES /* get services info */

	VS_CMD_NEW_DEST /* add destination */
	VS_CMD_SET_DEST /* modify destination */
	VS_CMD_DEL_DEST /* delete destination */
	VS_CMD_GET_DEST /* get destination info */

	VS_CMD_NEW_DAEMON /* start sync daemon */
	VS_CMD_DEL_DAEMON /* stop sync daemon */
	VS_CMD_GET_DAEMON /* get sync daemon status */

	VS_CMD_SET_CONFIG /* set config settings */
	VS_CMD_GET_CONFIG /* get config settings */

	VS_CMD_SET_INFO /* only used in GET_INFO reply */
	VS_CMD_GET_INFO /* get general IPVS info */

	VS_CMD_ZERO  /* zero all counters and stats */
	VS_CMD_FLUSH /* flush services and dests */

	VS_CMD_NEW_LADDR /* add local address */
	VS_CMD_DEL_LADDR /* del local address */
	VS_CMD_GET_LADDR /* dump local address */
)

const (
	IPPROTO_IP   = 0  /* Dummy protocol for TCP		*/
	IPPROTO_ICMP = 1  /* Internet Control Message Protocol	*/
	IPPROTO_IGMP = 2  /* Internet Group Management Protocol	*/
	IPPROTO_IPIP = 4  /* IPIP tunnels (older KA9Q tunnels use 94) */
	IPPROTO_TCP  = 6  /* Transmission Control Protocol	*/
	IPPROTO_EGP  = 8  /* Exterior Gateway Protocol		*/
	IPPROTO_PUP  = 12 /* PUP protocol				*/
	IPPROTO_UDP  = 17 /* User Datagram Protocol		*/
	IPPROTO_IDP  = 22 /* XNS IDP protocol			*/
	IPPROTO_DCCP = 33 /* Datagram Congestion Control Protocol */
	IPPROTO_RSVP = 46 /* RSVP protocol			*/
	IPPROTO_GRE  = 47 /* Cisco GRE tunnels (rfc 1701,1702)	*/

	IPPROTO_IPV6 = 41 /* IPv6-in-IPv4 tunnelling		*/

	IPPROTO_ESP    = 50  /* Encapsulation Security Payload protocol */
	IPPROTO_AH     = 51  /* Authentication Header protocol       */
	IPPROTO_BEETPH = 94  /* IP option pseudo header for BEET */
	IPPROTO_PIM    = 103 /* Protocol Independent Multicast	*/

	IPPROTO_COMP    = 108 /* Compression Header protocol */
	IPPROTO_SCTP    = 132 /* Stream Control Transport Protocol	*/
	IPPROTO_UDPLITE = 136 /* UDP-Lite (RFC 3828)			*/

	IPPROTO_RAW = 255 /* Raw IP packets			*/
)

type Addr4 struct {
	Ip   uint32
	Port uint16
}

var errIpv4 = errors.New("syntax error: expect 192.168.0.1")
var errIpv4Addr = errors.New("syntax error: expect 192.168.0.1 or 192.168.0.1:80")
var errProtocol = errors.New("syntax error: expect tcp or udp")
var errTimeout = errors.New("syntax error: expect '1,3,5'  (second)")

func (p *Addr4) Set(value string) error {
	if value == "" {
		p.Ip = 0
		p.Port = 0
		return nil
	}

	fields := strings.Split(value, ":")
	if len(fields) == 1 || len(fields) == 2 {
		if ip := net.ParseIP(fields[0]).To4(); ip != nil {
			p.Ip = ipToU32(ip)
			if len(fields) == 2 {
				if port, err := strconv.Atoi(fields[1]); err != nil {
					return errIpv4Addr
				} else {
					p.Port = uint16(port)
				}
			} else {
				p.Port = 0
			}
			return nil
		}
	}
	return errIpv4Addr
}

func (p *Addr4) String() string {
	return fmt.Sprintf("%s:%d", u32_to_addr(p.Ip), p.Port)
}

type Ipv4 uint32

func (p *Ipv4) Set(value string) error {
	if ip := net.ParseIP(value).To4(); ip != nil {
		*p = Ipv4(ipToU32(ip))
	}
	return errIpv4Addr
}

func (p *Ipv4) String() string {
	return u32_to_addr(uint32(*p))
}

type Protocol uint8

func (p *Protocol) Set(value string) error {
	if value == "tcp" {
		*p = IPPROTO_TCP
	} else if value == "udp" {
		*p = IPPROTO_UDP
	} else {
		return errProtocol
	}
	return nil
}

func (p *Protocol) String() string {
	return get_protocol_name(uint8(*p))
}

type CallOptions struct {
	Opt  CmdOptions
	Args []string
}

type CmdOptions struct {
	/* status */
	Typ string
	Id  int
	/* service */
	Addr         Addr4
	Nic          uint
	Protocol     Protocol
	TCP          bool
	UDP          bool
	Sched_name   string
	Flags        uint
	Timeout      uint
	Num_services int
	Number       int
	Timeout_s    string
	Netmask      Ipv4

	/* dest */
	D           bool
	Dnic        uint
	Daddr       Addr4
	Conn_flags  uint
	Weight      int
	U_threshold uint
	L_threshold uint

	/* local addr */
	L    bool
	Lnic uint
	Lip  Ipv4

	/* timeout */
	Tcp_timeout     int
	Tcp_fin_timeout int
	Udp_timeout     int
}

type Vs_service_user struct {
	Nic        uint8
	Protocol   uint8
	Addr       uint32
	Port       uint16
	Sched_name string
	Flags      uint
	Timeout    uint
	Netmask    uint32
	Number     int /* max list laddr/dests */
}

type Vs_service_user_r struct {
	Nic        uint8
	Protocol   uint8
	Addr       uint32
	Port       uint16
	Flags      uint32
	Timeout    uint32
	Netmask    uint32
	Num_dests  uint32
	Num_laddrs uint32
	Sched_name string
	Conns      string
	Inpkts     string
	Outpkts    string
	Inbytes    string
	Outbytes   string
}

type Vs_dest_user struct {
	Nic         uint8
	Addr        uint32
	Port        uint16
	Conn_flags  uint
	Weight      int
	U_threshold uint32
	L_threshold uint32
}

type Vs_dest_user_r struct {
	Nic         uint8
	Addr        uint32
	Port        uint16
	Conn_flags  uint
	Weight      int
	U_threshold uint32
	L_threshold uint32
	Activeconns uint32
	Inactconns  uint32
	Persistent  uint32
	Conns       string
	Inpkts      string
	Outpkts     string
	Inbytes     string
	Outbytes    string
}

type Vs_laddr_user struct {
	Nic  uint8
	Addr uint32
}

type Vs_laddr_user_r struct {
	Nic           uint8
	Addr          uint32
	Conn_counts   uint32
	Port_conflict string
}

type Vs_timeout_user struct {
	Tcp_timeout     int
	Tcp_fin_timeout int
	Udp_timeout     int
}

type Vs_cmd_q struct {
	Cmd int
}

type Vs_cmd_r struct {
	Code int
	Msg  string
}

func (r Vs_cmd_r) String() string {
	if r.Code == 0 {
		return fmt.Sprintf("done")
	} else {
		return fmt.Sprintf("%s:%s", Ecode(r.Code), r.Msg)
	}
}

type Vs_version_r struct {
	Code    int
	Msg     string
	Version int
	Size    int
}

func (r Vs_version_r) String() string {
	if r.Code != 0 {
		return fmt.Sprintf("%s:%s", Ecode(r.Code), r.Msg)
	}

	return fmt.Sprintf("version\t\t%d.%d.%d\n"+
		"conn table size\t%d",
		(r.Version>>16)&0xff,
		(r.Version>>8)&0xff,
		(r.Version)&0xff,
		r.Size)
}

type Vs_info_r struct {
	Code         int
	Msg          string
	Version      int
	Size         int
	Num_services int
	Seq          int
	State        string
	Workers      []struct {
		Worker_id int
		Seq       int
		State     string
	}
}

func (r Vs_info_r) String() string {
	if r.Code != 0 {
		return fmt.Sprintf("%s:%s", Ecode(r.Code), r.Msg)
	}
	s := fmt.Sprintf("Version   :%d\n"+
		"Size      :%d\n"+
		"Num_svcs  :%d\n"+
		"State     :%-10s %d\n",
		r.Version, r.Size, r.Num_services, r.State, r.Seq)
	for _, w := range r.Workers {
		s += fmt.Sprintf("Worker_%02d :%-10s %d\n",
			w.Worker_id, w.State, w.Seq)
	}
	return s
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
		return fmt.Sprintf("%s\n%s\n", Svc_title(), r.Service)
	} else {
		return fmt.Sprintf("%s:%s", Ecode(r.Code), r.Msg)
	}
}

type Vs_list_services_r struct {
	Code     int
	Msg      string
	Services []Vs_service_user_r
}

func (r Vs_list_services_r) String() string {
	if r.Code != 0 {
		return fmt.Sprintf("%s:%s", Ecode(r.Code), r.Msg)
	}
	s := fmt.Sprintf("%s\n", Svc_title())
	for _, svc := range r.Services {
		s += fmt.Sprintf("%s\n", svc)
	}
	return s
}

type Vs_list_dests_r struct {
	Code  int
	Msg   string
	Dests []Vs_dest_user_r
}

func (r Vs_list_dests_r) String() string {
	var s string
	if r.Code != 0 {
		return fmt.Sprintf("%s:%s", Ecode(r.Code), r.Msg)
	}
	for _, dest := range r.Dests {
		s += fmt.Sprintf("%s\n", dest)
	}
	return s
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
	return s
}

type Vs_timeout_q struct {
	Cmd             int
	Tcp_timeout     int
	Tcp_fin_timeout int
	Udp_timeout     int
}

func (t *Vs_timeout_q) Set(value string) (err error) {
	fields := strings.Split(value, ",")
	if len(fields) != 3 {
		return errTimeout
	}
	if t.Tcp_timeout, err = strconv.Atoi(fields[0]); err != nil {
		return err
	}
	if t.Tcp_fin_timeout, err = strconv.Atoi(fields[1]); err != nil {
		return err
	}
	if t.Udp_timeout, err = strconv.Atoi(fields[2]); err != nil {
		return err
	}
	return nil
}

type Vs_timeout_r struct {
	Code            int
	Msg             string
	Tcp_timeout     int
	Tcp_fin_timeout int
	Udp_timeout     int
}

func (r Vs_timeout_r) String() string {
	if r.Code == 0 {
		return fmt.Sprintf("tcp:%d tcp_fin:%d udp:%d",
			r.Tcp_timeout, r.Tcp_fin_timeout, r.Udp_timeout)
	} else {
		return fmt.Sprintf("%s:%s", Ecode(r.Code), r.Msg)
	}
}

/* set ctl */
type Vs_service_q struct {
	Cmd     int
	Service Vs_service_user
}

type Vs_dest_q struct {
	Cmd     int
	Service Vs_service_user
	Dest    Vs_dest_user
}

type Vs_laddr_q struct {
	Cmd     int
	Service Vs_service_user
	Laddr   Vs_laddr_user
}

type Vs_stats_q struct {
	Type int
	Id   int
}
