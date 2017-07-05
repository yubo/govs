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

type Vs_stats_ifa struct {
	Port       int
	Rx_packets int64
	Rx_dropped int64
	Tx_packets int64
	Tx_dropped int64
}

type Vs_stats_io_entry struct {
	Core                int
	Rx_rings_count      []int64
	Rx_rings_iters      []int64
	Rx_nic_queues_count []int64
	Rx_nic_queues_iters []int64
	Tx_nic_ports_count  []int64
	Tx_nic_ports_iters  []int64
	Kni                 []Vs_stats_ifa
	Kni_deq             int64
	Kni_deq_err         int64
}

type Vs_stats_io_r struct {
	Code int
	Msg  string
	Io   []Vs_stats_io_entry
}

func max_len(is ...[]int64) (l int) {
	for _, v := range is {
		if l < len(v) {
			l = len(v)
		}
	}
	return l
}

/*
id kni_deq kni_err  k_rx rx_d tx tx_d
*/
func (r Vs_stats_io_r) String() string {
	if r.Code != 0 {
		return fmt.Sprintf("%s:%s", Ecode(r.Code), r.Msg)
	}
	var ret string
	for _, e := range r.Io {
		ret += fmt.Sprintf("%-10s %10s %10s\n",
			"core", "kni_deq", "kni_deq_err")
		ret += fmt.Sprintf("%-10d %10d %10d\n\n",
			e.Core, e.Kni_deq, e.Kni_deq_err)
		n := max_len(e.Rx_rings_count, e.Rx_nic_queues_count, e.Tx_nic_ports_count)
		if n > len(e.Kni) {
			n = len(e.Kni)
		}

		m := make([][11]int64, n)
		for k, _ := range e.Rx_rings_count {
			m[k][0] = e.Rx_rings_count[k]
			m[k][1] = e.Rx_rings_iters[k]
		}
		for k, _ := range e.Rx_nic_queues_count {
			m[k][2] = e.Rx_nic_queues_count[k]
			m[k][3] = e.Rx_nic_queues_iters[k]
		}
		for k, _ := range e.Tx_nic_ports_count {
			m[k][4] = e.Tx_nic_ports_count[k]
			m[k][5] = e.Tx_nic_ports_iters[k]
		}
		for k, v := range e.Kni {
			m[k][6] = int64(v.Port)
			m[k][7] = v.Rx_packets
			m[k][8] = v.Rx_dropped
			m[k][9] = v.Tx_packets
			m[k][10] = v.Tx_dropped
		}

		ret += fmt.Sprintf("%-10s %10s %10s %10s %10s "+
			"%10s %10s %10s %10s %10s "+
			"%10s %10s\n",
			"id", "rx_ring_c", "rx_ring_i", "rx_nic_q_c", "rx_nic_q_i",
			"tx_nic_p_c", "tx_nic_p_i", "kni_port", "kni_rx_pkt", "kni_rx_drop",
			"kni_tx_pkt", "kni_tx_drop")
		for i := 0; i < n; i++ {
			ret += fmt.Sprintf("%-10d %10d %10d %10d %10d "+
				"%10d %10d %10d %10d %10d "+
				"%10d %10d\n",
				i, m[i][0], m[i][1], m[i][2], m[i][3],
				m[i][4], m[i][5], m[i][6], m[i][7], m[i][8],
				m[i][9], m[i][10])
		}

		ret += fmt.Sprintf("\n")

	}
	return ret
}

type Vs_stats_worker_entry struct {
	Core        int
	Ipmiss      int64
	Frag        int64
	Icmp        int64
	V4pkt       int64
	V4sctp      int64
	V4ospf      int64
	V4unknow    int64
	V4drop      int64
	Kni_enq     int64
	Kni_enq_err int64
	Arp         int64
	Ipv6        int64
	Unknow      int64
}

type Vs_stats_worker_r struct {
	Code   int
	Msg    string
	Worker []Vs_stats_worker_entry
}

func (r Vs_stats_worker_r) String() string {
	if r.Code != 0 {
		return fmt.Sprintf("%s:%s", Ecode(r.Code), r.Msg)
	}
	ret := fmt.Sprintf("%-5s %10s %10s %10s %10s "+
		"%10s %10s %10s %10s %10s "+
		"%10s %10s %10s %10s\n",
		"core", "ipmiss", "frag", "icmp", "pkt",
		"v4sctp", "ospf", "unknow(v4)", "drop", "kni_enq",
		"kni_err", "arp", "ipv6", "unknow")
	for _, e := range r.Worker {
		ret += fmt.Sprintf("%-5d %10d %10d %10d %10d "+
			"%10d %10d %10d %10d %10d "+
			"%10d %10d %10d %10d\n",
			e.Core, e.Ipmiss, e.Frag, e.Icmp, e.V4pkt,
			e.V4sctp, e.V4ospf, e.V4unknow, e.V4drop, e.Kni_enq,
			e.Kni_enq_err, e.Arp, e.Ipv6, e.Unknow)
	}
	return ret
}

type Vs_stats_dev_entry struct {
	Port      int
	Ipackets  int64
	Opackets  int64
	Ibytes    int64
	Obytes    int64
	Imissed   int64
	Ierrors   int64
	Oerrors   int64
	Rx_nombuf int64
}

type Vs_stats_dev_r struct {
	Code int
	Msg  string
	Dev  []Vs_stats_dev_entry
}

func (r Vs_stats_dev_r) String() string {
	if r.Code != 0 {
		return fmt.Sprintf("%s:%s", Ecode(r.Code), r.Msg)
	}
	ret := fmt.Sprintf("%-10s %10s %10s %10s %10s %10s %10s %10s %10s\n",
		"port", "ipackets", "Opackets", "ibytes", "obytes",
		"imissed", "ierrors", "oerrors", "rx_nombuf")
	for _, e := range r.Dev {
		ret += fmt.Sprintf("%-10d %10d %10d %10d %10d %10d %10d %10d %10d\n",
			e.Port, e.Ipackets, e.Opackets, e.Ibytes, e.Obytes,
			e.Imissed, e.Ierrors, e.Oerrors, e.Rx_nombuf)
	}
	return ret
}

type Vs_stats_ctl_r struct {
	Code         int
	Msg          string
	Num_services int
	Seq          int
	Workers      []struct {
		Worker_id    int
		Num_services int
		Seq          int
		State        int
	}
}

/*
id                seq      n_svc      state
-                   0          0          -

0                   0          0          s
*/
func (r Vs_stats_ctl_r) String() string {
	if r.Code != 0 {
		return fmt.Sprintf("%s:%s", Ecode(r.Code), r.Msg)
	}
	ret := fmt.Sprintf("%-10s %10s %10s %10s\n",
		"id", "seq", "n_svc", "state")
	ret += fmt.Sprintf("%-10s %10d %10d %10c\n\n",
		"-", r.Seq, r.Num_services, '-')
	for _, e := range r.Workers {
		ret += fmt.Sprintf("%-10d %10d %10d %10c\n",
			e.Worker_id, e.Seq, e.Num_services, ctl_state_name(e.State))

	}
	return ret
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
