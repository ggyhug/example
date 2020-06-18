package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
        "encoding/json"
	"log"
	"metrics-server-exporter-go/api"
	"time"
)

// Info - Data structure of pod
type Info struct {
	Items []struct {
		Metadata struct {
			Name              string    `json:"name"`
			Namespace         string    `json:"namespace"`
			CreationTimestamp time.Time `json:"creationTimestamp"`
		} `json:"metadata"`
		Containers []struct {
			Name  string `json:"name"`
			Usage struct {
				CPU    string `json:"cpu"`
				Memory string `json:"memory"`
			} `json:"usage"`
		} `json:"containers"`
	} `json:"items"`
}


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
	
	// MetricsPodsCPU - CPU Gauge
	MetricsPodsCPU = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kube_metrics_server_pods_cpu",
			Help: "Metrics Server Pods CPU",
		},
		[]string{"pod_name", "pod_namespace", "pod_container_name"},
	)
	// MetricsPodsMEM - Memory Gauge
	MetricsPodsMEM = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kube_metrics_server_pods_mem",
			Help: "Metrics Server Pods Memory",
		},
		[]string{"pod_name", "pod_namespace", "pod_container_name"},
	)
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
	prometheus.MustRegister(MetricsPodsMeM)
	prometheus.MustRegister(MetricsPodsCPU)
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
	var pods Info
	log.Println("Starting collect POD data,")

	apiPod := api.Connect("pod")

	_ = json.NewDecoder(apiPod.Body).Decode(&pods)

	for i := range pods.Items {

		podName := pods.Items[i].Metadata.Name
		podNamespace := pods.Items[i].Metadata.Namespace

		for j := range pods.Items[i].Containers {

			MetricsPodsCPU.With(prometheus.Labels{"pod_name": podName, "pod_namespace": podNamespace, "pod_container_name": pods.Items[i].Containers[j].Name}).Add(api.ReturnFloat(pods.Items[i].Containers[j].Usage.CPU))
			MetricsPodsMEM.With(prometheus.Labels{"pod_name": podName, "pod_namespace": podNamespace, "pod_container_name": pods.Items[i].Containers[j].Name}).Add(api.ReturnFloat(pods.Items[i].Containers[j].Usage.Memory))
		}

	}
	log.Println("POD data collected.")
}
