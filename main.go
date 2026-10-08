package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	spec, err := ebpf.LoadCollectionSpec("xdp_counter.o")
	if err != nil {
		fmt.Println("Error:", err)
		log.Fatal(err)
	}
	c, err := ebpf.NewCollection(spec)
	if err != nil {
		fmt.Println("Error:", err)
		log.Fatal(err)
	}
	defer c.Close()
	intrface, err := net.InterfaceByName("veth-a")
	if err != nil {
		fmt.Println("Error:", err)
		log.Fatal(err)
	}
	isAttached, err := link.AttachXDP(link.XDPOptions{
		Program:   c.Programs["xdp_counter"],
		Interface: intrface.Index,
		Flags:     link.XDPGenericMode,
	})
	if err != nil {
		fmt.Println("Error:", err)
		log.Fatal(err)
	}
	defer isAttached.Close()
	fmt.Println("Xdp Attached to", intrface.Name)

	counters := c.Maps["counters"]
	registerCounter := func(name string, key uint32) {
		prometheus.MustRegister(prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name: name,
				Help: "Cumulative packets counted by XDP",
			},
			func() float64 {
				var value uint64

				if err := counters.Lookup(key, &value); err != nil {
					log.Printf("BPF lookup error: %v", err)
					return 0
				}

				return float64(value)
			},
		))
	}
	registerCounter("xdp_total_packet", 0)
	registerCounter("xdp_syn_packet", 1)
	http.Handle("/Metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe("127.0.0.1:9091", nil))
	for {
		var total, syn uint64
		counters.Lookup(uint32(0), &total)
		counters.Lookup(uint32(1), &syn)
		fmt.Println(total, syn)
		time.Sleep(time.Second)

	}

}
