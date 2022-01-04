package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

// When adding metric names, see https://prometheus.io/docs/practices/naming/#metric-names
const (
	EdgeDeviceSuccessfulRegistrationQuery = "flotta_operator_edge_devices_successful_registration"
	EdgeDeviceFailedRegistrationQuery     = "flotta_operator_edge_devices_failed_registration"
	EdgeDeviceUnregistrationQuery         = "flotta_operator_edge_devices_unregistration"
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
	SetPatchEdgeDeviceStatusTime(duration int64)
	SetPatchEdgeDeviceTime(duration int64)
	SetProcessHeartbeatTime(duration int64)
}

func New() Metrics {
	return &metricsImpl{}
}

type metricsImpl struct{}

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
