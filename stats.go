/*
 * Copyright 2016 Xiaomi Corporation. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 *
 * Authors:    Yu Bo <yubo@xiaomi.com>
 */
package govs

import "fmt"

const (
	VS_CTL_S_SYNC = iota
	VS_CTL_S_PENDING
	VS_CTL_S_LAST
)

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
		"port", "ipackets", "opackets", "ibytes", "obytes",
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

type Vs_stats_mem_r struct {
	Code int
	Msg  string
	Size struct {
		Mbuf  int
		Svc   int
		Rs    int
		Laddr int
		Conn  int
	}
	Available []struct {
		Socket_id int
		Mbuf      int
		Svc       int
		Rs        int
		Laddr     int
		Conn      int
	}
}

func (r Vs_stats_mem_r) String() string {
	if r.Code != 0 {
		return fmt.Sprintf("%s:%s", Ecode(r.Code), r.Msg)
	}
	ret := fmt.Sprintf("%-10s %10s %10s %10s %10s %10s\n",
		"id", "mbuf", "svc", "rs", "laddr", "conn")
	ret += fmt.Sprintf("%-10s %10d %10d %10d %10d %10d\n\n",
		"max", r.Size.Mbuf, r.Size.Svc, r.Size.Rs, r.Size.Laddr, r.Size.Conn)
	for _, e := range r.Available {
		ret += fmt.Sprintf("%-10d %10d %10d %10d %10d %10d\n",
			e.Socket_id,
			r.Size.Mbuf-e.Mbuf,
			r.Size.Svc-e.Svc,
			r.Size.Rs-e.Rs,
			r.Size.Laddr-e.Laddr,
			r.Size.Conn-e.Conn)

	}
	return ret
}

type Vs_stats_q struct {
	Type int
	Id   int
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
func Get_stats_io(id int) (*Vs_stats_io_r, error) {
	args := Vs_stats_q{Type: VS_STATS_IO, Id: id}
	reply := &Vs_stats_io_r{}

	err := client.Call("stats", args, reply)
	return reply, err
}

func Get_stats_worker(id int) (*Vs_stats_worker_r, error) {
	args := Vs_stats_q{Type: VS_STATS_WORKER, Id: id}
	reply := &Vs_stats_worker_r{}

	err := client.Call("stats", args, reply)
	return reply, err
}

func Get_stats_dev(id int) (*Vs_stats_dev_r, error) {
	args := Vs_stats_q{Type: VS_STATS_DEV, Id: id}
	reply := &Vs_stats_dev_r{}

	err := client.Call("stats", args, reply)
	return reply, err
}

func Get_stats_ctl() (*Vs_stats_ctl_r, error) {
	args := Vs_stats_q{Type: VS_STATS_CTL}
	reply := &Vs_stats_ctl_r{}

	err := client.Call("stats", args, reply)
	return reply, err
}

func Get_stats_mem() (*Vs_stats_mem_r, error) {
	args := Vs_stats_q{Type: VS_STATS_MEM}
	reply := &Vs_stats_mem_r{}

	err := client.Call("stats", args, reply)
	return reply, err
}
