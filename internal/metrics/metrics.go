package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
	"sync"
)

// When adding metric names, see https://prometheus.io/docs/practices/naming/#metric-names
const (
	EdgeDeviceSuccessfulRegistrationQuery = "flotta_operator_edge_devices_successful_registration"
	EdgeDeviceFailedRegistrationQuery     = "flotta_operator_edge_devices_failed_registration"
	EdgeDeviceUnregistrationQuery         = "flotta_operator_edge_devices_unregistration"
	EdgeDeviceHeartbeatQuery              = "flotta_operator_edge_devices_heartbeat"
	PatchEdgeDeviceStatusDurationQuery    = "flotta_operator_edge_devices_patch_status_duration_milliseconds"
	PatchEdgeDeviceDurationQuery          = "flotta_operator_edge_devices_patch_duration_milliseconds"
	ProcessHeartbeatDurationQuery         = "flotta_operator_process_heartbeat_duration_milliseconds"
)

var (
	processHeartbeatDuration = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: ProcessHeartbeatDurationQuery,
			Help: "Time in millis to process a heartbeat",
		},
	)
	patchEdgeDeviceStatusDuration = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: PatchEdgeDeviceStatusDurationQuery,
			Help: "Time in millis to patch EdgeDevices status",
		},
	)
	patchEdgeDevicesDuration = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: PatchEdgeDeviceDurationQuery,
			Help: "Time in millis to patch EdgeDevice",
		},
	)
	registeredEdgeDevices = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: EdgeDeviceSuccessfulRegistrationQuery,
			Help: "Number of successful registration EdgeDevices",
		},
	)
	failedToCompleteRegistrationEdgeDevices = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: EdgeDeviceFailedRegistrationQuery,
			Help: "Number of failed registration EdgeDevices",
		},
	)
	unregisteredEdgeDevices = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: EdgeDeviceUnregistrationQuery,
			Help: "Number of unregistered EdgeDevices",
		},
	)
)

func init() {
	// Register custom metrics with the global prometheus registry
	metrics.Registry.MustRegister(
		registeredEdgeDevices,
		failedToCompleteRegistrationEdgeDevices,
		unregisteredEdgeDevices,
		patchEdgeDeviceStatusDuration,
		patchEdgeDevicesDuration,
		processHeartbeatDuration,
	)
}

//go:generate mockgen -source=metrics.go -package=metrics -destination=mock_metrics_api.go

// Metrics is an interface representing a prometheus client for the Special Resource Operator
type Metrics interface {
	IncEdgeDeviceSuccessfulRegistration()
	IncEdgeDeviceFailedRegistration()
	IncEdgeDeviceUnregistration()
	RecordEdgeDevicePresence(namespace, name string)
	RemoveDeviceCounter(namespace, name string)
	RegisterDeviceCounter(namespace string, name string)
	SetPatchEdgeDeviceStatusTime(duration int64)
	SetPatchEdgeDeviceTime(duration int64)
	SetProcessHeartbeatTime(duration int64)
}

func New() Metrics {
	return &metricsImpl{
		devices: sync.Map{},
	}
}

type metricsImpl struct {
	devices sync.Map
}

func (m *metricsImpl) RecordEdgeDevicePresence(namespace, name string) {
	m.registerDeviceCounter(namespace, name).Inc()
}

func (m *metricsImpl) RegisterDeviceCounter(namespace string, name string) {
	m.registerDeviceCounter(namespace, name)
}

func (m *metricsImpl) SetProcessHeartbeatTime(duration int64) {
	processHeartbeatDuration.Set(float64(duration))
}

func (m *metricsImpl) SetPatchEdgeDeviceTime(duration int64) {
	patchEdgeDevicesDuration.Set(float64(duration))
}

func (m *metricsImpl) SetPatchEdgeDeviceStatusTime(duration int64) {
	patchEdgeDeviceStatusDuration.Set(float64(duration))
}

func (m *metricsImpl) IncEdgeDeviceSuccessfulRegistration() {
	registeredEdgeDevices.Inc()
}

func (m *metricsImpl) IncEdgeDeviceFailedRegistration() {
	failedToCompleteRegistrationEdgeDevices.Inc()
}
func (m *metricsImpl) IncEdgeDeviceUnregistration() {
	unregisteredEdgeDevices.Inc()
}

func (m *metricsImpl) RemoveDeviceCounter(namespace, name string) {
	if counter, ok := m.devices.LoadAndDelete(deviceKey(namespace, name)); ok {
		metrics.Registry.Unregister(counter.(prometheus.Counter))
	}
}

func (m *metricsImpl) registerDeviceCounter(namespace, name string) prometheus.Counter {
	key := deviceKey(namespace, name)
	collector, loaded := m.devices.LoadOrStore(key, prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: EdgeDeviceHeartbeatQuery,
			ConstLabels: prometheus.Labels{
				"deviceNamespace": namespace,
				"deviceID":        name,
			},
		}))

	counter := collector.(prometheus.Counter)
	if !loaded {
		metrics.Registry.MustRegister(counter)
	}
	return counter
}

func deviceKey(namespace string, name string) string {
	return fmt.Sprintf("%s/%s", namespace, name)
}
