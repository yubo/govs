## govs

This set of packages provide the API for communication with [DPVS](http://git.n.xiaomi.com/dpvs/vs/) from Go.

#### install

```
# install golang
go get github.com/yubo/govs/cmd/govs
```

#### howto

```
govs -h
govs stats -t io/worker/dev/ctl [-i id]
```

io core information

```
#govs stats -t io
core_id                                   2
rx_ring_0                                 0          0
rx_ring_1                                 0          0
Rx_nic_port0_queue0                   20537       5539
Rx_nic_port1_queue1                     495        495
Rx_nic_port2_queue2                       5          5
Rx_nic_port3_queue3                       0          0
tx_nic_port0                              0          0
tx_nic_port1                           1192       1192
veth0                            Rx_packets         45 Rx_dropped          0 Tx_packets          8 Tx_dropped          0
veth1                            Rx_packets          0 Rx_dropped          0 Tx_packets       4000 Tx_dropped          0
kni_deq                          128726359132
kni_deq_err                      128726359087
```

- core_id: core id
- rx_ring_worker?:  calls number of rte_ring_sp_enqueue_bulk()
- rx_nic_port?_queue?:  number of mbufs get from rte_eth_rx_burst()
- tx_nic_port?:  number of mbufs get from rte_eth_tx_burst()
- veth?: kni dev counter
- kni_deq:     calls number of rte_ring_sc_dequeue_burst(..)
- kni_deq_err: calls number of failed rte_ring_sc_dequeue_burst(..)


```
#govs stats -t worker
core      ipmiss       frag       icmp        pkt     v4sctp       ospf unknow(v4)       drop    kni_enq    kni_err        arp       ipv6     unknow
3              0          0          0          0          0          0          0          0     201277          0          0          0          0
```

- core: number of worker core
- ipmiss: number of pkts ip miss 
- frag: number of pkt is fragmented
- icmp: number of icmp pkt
- pkt: number of pkt
- v4sctp: number of ipv4 sctp
- ospf: number of ipv4 ospf
- unknow(v4): number of unknow ipv4 protocol
- drop: number of drop pkt
- kni_enq:   calls number of rte_ring_sp_enqueue_bulk(lp->kni_rings_out[skb->port], ...)
- kni_err:   calls number of failed rte_ring_sp_enqueue_bulk(lp->kni_rings_out[skb->port], ...)
- arp: number of arp pkt
- ipv6: number of ipv6 pkt
- unknow: number of unknow L3 protocol

```
#govs stats -t dev
port         ipackets   opackets     ibytes     obytes    imissed    ierrors    oerrors  rx_nombuf
0              120092         70   30942493       8118          0          0          0          0
1              123517      10318   31196462     711370          0          0          0          0
```

- port: number of port
- ipackets: Total number of successfully received packets
- opackets: Total number of successfully transmitted packets
- ibytes: Total number of successfully received bytes.
- obytes: Total number of successfully transmitted bytes
- imissed: Total of RX packets dropped by the HW, * because there are no available buffer (i.e. RX queues are full)
- ierrors: Total number of erroneous received packets
- oerrors: Total number of failed transmitted packets
- rx_nombuf: Total number of RX mbuf allocation failures



```
#govs stats -t ctl
id                seq      n_svc      state
-                   0          0          -

0                   0          0          s
```

- id: number of core
- seq: sequence number of config
- n_svc: total number of virtual service
- state: state of the worker, s(sync), p(pending)


#### AUTHOR

Written by Yu Bo.

#### REPORTING BUGS

Report bugs to yubo@xiaomi.com
