package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu&quot"
	"github.com/shirou/gopsutil/load&quot"
	"github.com/shirou/gopsutil/mem&quot"
	"time"
	"&quot"
)

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "request_total",
			Help:      "Number of request processed by this service.",
		}, []string{},
	)

	requestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:      "request_latency_seconds",
			Help:      "Time spent in this service.",
			Buckets:   []float64{0.01, 0.02, 0.05, 0.1, 0.2, 0.5, 1.0, 2.0, 5.0, 10.0, 20.0, 30.0, 60.0, 120.0, 300.0},
		}, []string{},
	)
	//系统cpu利用率
	cpu_usage := prometheus.NewGauge(
		prometheus.GaugeOpts{
		Name:      "cpu_usage",
		Help:      "system cpu usage.",
	})
	//系统cpu负载率
	cpu_usage := prometheus.NewGauge(
		prometheus.GaugeOpts{
		Name:      "cpu_load",
		Help:      "system cpu load.",
	})
	//系统mem使用情况
	mem := prometheus.NewGauge(
		prometheus.GaugeOpts{
		Name:      "mem",
		Help:      "system mem usage.",
	})
	
)

// AdmissionLatency measures latency / execution time of Admission Control execution
// usual usage pattern is: timer := NewAdmissionLatency() ; compute ; timer.Observe()
type RequestLatency struct {
	histo *prometheus.HistogramVec
	start time.Time
}

func Register() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestLatency)
	prometheus.MustRegister(cou_usage)
	prometheus.MustRegister(cpu_load)
	prometheus.MustRegister(mem)
}


// NewAdmissionLatency provides a timer for admission latency; call Observe() on it to measure
func NewAdmissionLatency() *RequestLatency {
	return &RequestLatency{
		histo: requestLatency,
		start: time.Now(),
	}
}

// Observe measures the execution time from when the AdmissionLatency was created
func (t *RequestLatency) Observe() {
	(*t.histo).WithLabelValues().Observe(time.Now().Sub(t.start).Seconds())
}


// RequestIncrease increases the counter of request handled by this service
func RequestIncrease() {
	requestCount.WithLabelValues().Add(1)
	cpu1 :=cpu.Percent(time.Second, false)
	cpu2 :=load.Avg()
	mem_ :=mem.VirtualMemory()
	cpu_usage.Set(cpu1)
	cpu_load.Set(cpu2)
	mem.Set(mem)
}
