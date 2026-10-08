package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var packets = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "demo_packets_total",
		Help: "Total demo packets",
	},
)

// func main() {
// 	prometheus.MustRegister(packets)

// 	go func() {
// 		for {
// 			packets.Add(100)
// 			time.Sleep(time.Second)
// 		}
// 	}()

// 	http.Handle("/metrics", promhttp.Handler())

// 	fmt.Println("Metrics at http://localhost:9091/metrics")
// 	http.ListenAndServe(":9091", nil)
// }
