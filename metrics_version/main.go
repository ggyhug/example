package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"example/metrics"
	"time"
)

func main(){
	http.HandleFunc("/abc", index)
	http.Handle("/metrics", promhttp.Handler())
	metrics.Register()
	err := http.ListenAndServe(":5565", nil) // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	timer:=metrics.NewAdmissionLatency()	
	go func() {
		for {
                       metrics.RequestIncrease()
		       time.Sleep(time.Second)
		}
	}()
	timer.Observe()
}


