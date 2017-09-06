/*
 * Copyright 2016 Xiaomi Corporation. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 *
 * Authors:    Yu Bo <yubo@xiaomi.com>
 */
package govs

import (
	"encoding/json"
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
	"strings"
)

const (
	URL = "/tmp/dpvs.sock"

	/*
	 *      IPVS Connection Flags
	 */
	VS_CONN_F_FWD_MASK   = 0x000f /* mask for the fwd methods */
	VS_CONN_F_MASQ       = 0x0000 /* masquerading/NAT */
	VS_CONN_F_LOCALNODE  = 0x0001 /* local node */
	VS_CONN_F_TUNNEL     = 0x0002 /* tunneling */
	VS_CONN_F_DROUTE     = 0x0003 /* direct routing */
	VS_CONN_F_BYPASS     = 0x0004 /* cache bypass */
	VS_CONN_F_FULLNAT    = 0x0005 /* full nat */
	VS_CONN_F_DSNAT      = 0x0008 /* dsnat flag */
	VS_CONN_F_SYNC       = 0x0020 /* entry created by sync */
	VS_CONN_F_HASHED     = 0x0040 /* hashed entry */
	VS_CONN_F_NOOUTPUT   = 0x0080 /* no output packets */
	VS_CONN_F_INACTIVE   = 0x0100 /* not established */
	VS_CONN_F_OUT_SEQ    = 0x0200 /* must do output seq adjust */
	VS_CONN_F_IN_SEQ     = 0x0400 /* must do input seq adjust */
	VS_CONN_F_SEQ_MASK   = 0x0600 /* in/out sequence mask */
	VS_CONN_F_NO_CPORT   = 0x0800 /* no client port set yet */
	VS_CONN_F_TEMPLATE   = 0x1000 /* template, not connection */
	VS_CONN_F_ONE_PACKET = 0x2000 /* forward only one packet */
	VS_CONN_F_SYNPROXY   = 0x8000 /* syn proxy flag */

	VS_SVC_F_PERSISTENT = 0x0001             /* persistent port */
	VS_SVC_F_HASHED     = 0x0002             /* hashed entry */
	VS_SVC_F_ONEPACKET  = 0x0004             /* one-packet scheduling */
	VS_SVC_F_SYNPROXY   = VS_CONN_F_SYNPROXY /* synproxy flag */
	VS_SVC_F_DSNAT      = VS_CONN_F_DSNAT    /* dsnat flag */
	VS_SVC_F_MASK       = (VS_SVC_F_PERSISTENT |
		VS_SVC_F_ONEPACKET |
		VS_SVC_F_SYNPROXY |
		VS_SVC_F_DSNAT)
)

var (
	client *rpc.Client
	CmdOpt CmdOptions
)

type CallOptions struct {
	Opt  CmdOptions
	Args []string
}

type CmdOptions struct {
	/* status */
	Typ string
	Id  int
	/* service */
	Addr       Addr4
	Nic        uint
	Protocol   Protocol
	TCP        string
	UDP        string
	Sched_name string
	Flags      uint
	Number     int
	Timeout    uint
	Timeout_s  string
	Netmask    Be32

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
	Lip  Be32

	/* timeout */
	Tcp_timeout     int
	Tcp_fin_timeout int
	Udp_timeout     int
}

type Be32 uint32

func (p *Be32) Set(value string) error {
	if value == "" {
		return nil
	}
	if ip := net.ParseIP(value).To4(); ip != nil {
		*p = Htonl(ipToU32(ip))
	}
	return nil
}

func (p Be32) String() string {
	return be32_to_addr(p)
}

func (p *Be32) UnmarshalJSON(data []byte) error {
	var i int32
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	*p = Be32(uint32(i))
	return nil
}

type Be16 uint16

func (p *Be16) Set(value string) error {
	if value == "" {
		return nil
	}

	port, err := strconv.Atoi(value)
	if err != nil {
		return err
	}

	*p = Htons(uint16(port))
	return nil
}

func (p Be16) String() string {
	return fmt.Sprintf("%d", Ntohs(p))
}

type Addr4 struct {
	Ip   Be32
	Port Be16
}

func (p *Addr4) Set(value string) error {
	if value == "" {
		p.Ip = 0
		p.Port = 0
		return nil
	}

	fields := strings.Split(value, ":")
	if len(fields) == 1 || len(fields) == 2 {
		if ip := net.ParseIP(fields[0]).To4(); ip != nil {
			p.Ip = Htonl(ipToU32(ip))
			if len(fields) == 2 {
				if port, err := strconv.Atoi(fields[1]); err != nil {
					return errIpv4Addr
				} else {
					p.Port = Htons(uint16(port))
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
	return fmt.Sprintf("%s:%d", be32_to_addr(p.Ip), Ntohs(p.Port))
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

const (
	VS_STATS_IO = iota
	VS_STATS_WORKER
	VS_ESTATS_WORKER
	VS_STATS_DEV
	VS_STATS_CTL
	VS_STATS_MEM
	__VS_STATS_MAX
)

const (
	VS_CMD_UNSPEC       = iota
	VS_CMD_NEW_SERVICE  /* add service */
	VS_CMD_SET_SERVICE  /* modify service */
	VS_CMD_DEL_SERVICE  /* delete service */
	VS_CMD_GET_SERVICE  /* get service info */
	VS_CMD_GET_SERVICES /* get services info */
	VS_CMD_NEW_DEST     /* add destination */
	VS_CMD_SET_DEST     /* modify destination */
	VS_CMD_DEL_DEST     /* delete destination */
	VS_CMD_GET_DEST     /* get destination info */
	VS_CMD_NEW_DAEMON   /* start sync daemon */
	VS_CMD_DEL_DAEMON   /* stop sync daemon */
	VS_CMD_GET_DAEMON   /* get sync daemon status */
	VS_CMD_SET_CONFIG   /* set config settings */
	VS_CMD_GET_CONFIG   /* get config settings */
	VS_CMD_SET_INFO     /* only used in GET_INFO reply */
	VS_CMD_GET_INFO     /* get general IPVS info */
	VS_CMD_ZERO         /* zero all counters and stats */
	VS_CMD_FLUSH        /* flush services and dests */
	VS_CMD_NEW_LADDR    /* add local address */
	VS_CMD_DEL_LADDR    /* del local address */
	VS_CMD_GET_LADDR    /* dump local address */
	VS_CMD_GET_STATS    /* dump workers/ctl  stats */
	__VS_CMD_MAX
)

const (
	IPPROTO_IP      = 0   /* Dummy protocol for TCP		*/
	IPPROTO_ICMP    = 1   /* Internet Control Message Protocol	*/
	IPPROTO_IGMP    = 2   /* Internet Group Management Protocol	*/
	IPPROTO_IPIP    = 4   /* IPIP tunnels (older KA9Q tunnels use 94) */
	IPPROTO_TCP     = 6   /* Transmission Control Protocol	*/
	IPPROTO_EGP     = 8   /* Exterior Gateway Protocol		*/
	IPPROTO_PUP     = 12  /* PUP protocol				*/
	IPPROTO_UDP     = 17  /* User Datagram Protocol		*/
	IPPROTO_IDP     = 22  /* XNS IDP protocol			*/
	IPPROTO_DCCP    = 33  /* Datagram Congestion Control Protocol */
	IPPROTO_IPV6    = 41  /* IPv6-in-IPv4 tunnelling		*/
	IPPROTO_RSVP    = 46  /* RSVP protocol			*/
	IPPROTO_GRE     = 47  /* Cisco GRE tunnels (rfc 1701,1702)	*/
	IPPROTO_ESP     = 50  /* Encapsulation Security Payload protocol */
	IPPROTO_AH      = 51  /* Authentication Header protocol       */
	IPPROTO_BEETPH  = 94  /* IP option pseudo header for BEET */
	IPPROTO_PIM     = 103 /* Protocol Independent Multicast	*/
	IPPROTO_COMP    = 108 /* Compression Header protocol */
	IPPROTO_SCTP    = 132 /* Stream Control Transport Protocol	*/
	IPPROTO_UDPLITE = 136 /* UDP-Lite (RFC 3828)			*/
	IPPROTO_RAW     = 255 /* Raw IP packets			*/
)

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
func Vs_dial() error {
	conn, err := net.Dial("unix", URL)
	if err != nil {
		return err
	}
	client = jsonrpc.NewClient(conn)
	return nil
}

func Vs_close() {
	if client != nil {
		client.Close()
	}
}

func Get_version() (*Vs_version_r, error) {
	var reply Vs_version_r
	args := Vs_cmd_q{VS_CMD_GET_INFO}

	err := client.Call("api", args, &reply)
	return &reply, err
}

func Get_timeout(o *CmdOptions) (*Vs_timeout_r, error) {
	args := Vs_timeout_q{Cmd: VS_CMD_GET_CONFIG}
	reply := &Vs_timeout_r{}

	err := client.Call("api", args, reply)
	return reply, err
}

func Set_flush(o *CmdOptions) (*Vs_cmd_r, error) {
	var reply Vs_cmd_r
	args := Vs_cmd_q{VS_CMD_FLUSH}

	err := client.Call("api", args, &reply)
	return &reply, err
}

func Set_timeout(o *CmdOptions) (*Vs_cmd_r, error) {
	var reply Vs_cmd_r
	args := Vs_timeout_q{Cmd: VS_CMD_SET_CONFIG}

	if err := args.Set(o.Timeout_s); err != nil {
		return nil, err
	}

	err := client.Call("api", args, &reply)
	return &reply, err
}

func Set_zero(o *CmdOptions) (*Vs_cmd_r, error) {
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
	return &reply, err
}
