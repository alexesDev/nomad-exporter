package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	taskStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "nomad_task_status",
	}, []string{
		"job",
		"taskgroup",
		"task",
	})
)

func main() {
	log.SetFlags(0)

	prometheus.MustRegister(taskStatus)

	go record()

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":5577"
	}

	log.Println("Start listening", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func record() {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal("failed to create nomad client:", err)
	}

	statuses := map[string]float64{
		"running":  1,
		"pending":  2,
		"complete": 3,
		"failed":   4,
		"lost":     5,
	}

	recordCompleted := os.Getenv("RECORD_COMPLETED") != ""

	for {
		allocs, _, err := client.Allocations().List(nil)
		if err != nil {
			log.Println("failed to get allocations", err)
			continue
		}

		for _, alloc := range allocs {
			if !recordCompleted && alloc.ClientStatus == "complete" {
				continue
			}

			statusValue, ok := statuses[alloc.ClientStatus]
			if !ok {
				statusValue = -1
			}

			taskStatus.WithLabelValues(alloc.JobID, alloc.TaskGroup, alloc.Name).Set(statusValue)
		}

		time.Sleep(30 * time.Second)
	}
}
