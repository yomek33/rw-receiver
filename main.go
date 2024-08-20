package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/storage/remote"
)

func main() {
	log.SetOutput(os.Stdout)
	log.Println("Starting server on :1234")

	http.HandleFunc("/receive", func(w http.ResponseWriter, r *http.Request) {
		req, err := remote.DecodeWriteRequest(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, ts := range req.Timeseries {
			m := make(model.Metric, len(ts.Labels))
			for _, l := range ts.Labels {
				m[model.LabelName(l.Name)] = model.LabelValue(l.Value)
			}
			fmt.Println(m)
			for _, s := range ts.Samples {
				fmt.Printf("\tSample:  %f %d\n", s.Value, s.Timestamp)
			}
		}
	})

	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
