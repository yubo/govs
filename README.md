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
core          kni_deq kni_deq_err
2          37625569804 37625368197

id          rx_ring_c  rx_ring_i rx_nic_q_c rx_nic_q_i tx_nic_p_c tx_nic_p_i   kni_port kni_rx_pkt kni_rx_drop kni_tx_pkt kni_tx_drop
0                   0          0     120146     119922          0          0          0      99129          0         70          0
1                   0          0     123571     123284          0          0          1     102556          0      10318          0

```

- core: core id
- kni_deq:     calls number of rte_ring_sc_dequeue_burst(..)
- kni_deq_err: calls number of failed rte_ring_sc_dequeue_burst(..)
- rx_ring_c:   calls number of rte_ring_sp_enqueue_bulk()
- rx_ring_i:   calls number of failed rte_ring_sp_enqueue_bulk()
- rx_nic_q_c:  number of mbufs get from rte_eth_rx_burst(port, queue, ...)
- rx_nic_q_i:  calls number of function calls at rte_eth_rx_burst(port, queue, ...)
- tx_nic_p_c:  number of mbufs get from rte_eth_tx_burst(port, ...)
- tx_nic_p_i:  calls number of rte_eth_tx_burst(port, ...)
- kni_port:    port number of kni
- kni_rx_pkt:  number of pkts received from NIC, and sent to KNI
- kni_rx_drop: number of pkts received from NIC, but failed to send to KNI
- kni_tx_pkt:  number of pkts received from KNI, and sent to NIC
- kni_tx_drop: number of pkts received from KNI, but failed to send to NIC


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
