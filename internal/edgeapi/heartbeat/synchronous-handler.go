package heartbeat

import (
	"context"
	mtrcs "github.com/project-flotta/flotta-operator/internal/common/metrics"
	"github.com/project-flotta/flotta-operator/internal/common/repository/edgedevice"
	"time"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
)

type SynchronousHandler struct {
	deviceRepository edgedevice.Repository
	updater          Updater
	logger           *zap.SugaredLogger
}

func NewSynchronousHandler(deviceRepository edgedevice.Repository, recorder record.EventRecorder, metrics mtrcs.Metrics,
	logger *zap.SugaredLogger) *SynchronousHandler {
	return &SynchronousHandler{
		logger:           logger,
		deviceRepository: deviceRepository,
		updater: Updater{
			deviceRepository: deviceRepository,
			recorder:         recorder,
			metrics:          metrics,
		},
	}
}

func (h *SynchronousHandler) Start() {
	// noop
}

func (h *SynchronousHandler) Process(ctx context.Context, notification Notification) error {
	// retry patching the edge device status
	var err error
	var retry bool
	for i := 1; i < 5; i++ {
		err, retry = h.process(ctx, notification)
		if err == nil {
			return nil
		}
		if !retry {
			break
		}

		notification.Retry++
		time.Sleep(time.Duration(i*50) * time.Millisecond)
	}
	return err
}

func (h *SynchronousHandler) process(ctx context.Context, notification Notification) (error, bool) {
	logger := h.logger.With("DeviceID", notification.DeviceID, "Namespace", notification.Namespace)
	heartbeat := notification.Heartbeat
	logger.Debug("processing heartbeat", "content", heartbeat, "retry", notification.Retry)
	edgeDevice, err := h.deviceRepository.Read(ctx, notification.DeviceID, notification.Namespace)
	if err != nil {
		if errors.IsNotFound(err) {
			return err, false
		}
		return err, true
	}

	// Produce k8s events based on the device-worker events:
	if notification.Retry == 0 {
		h.updater.processEvents(edgeDevice, heartbeat.Events)
	}

	err = h.updater.updateStatus(ctx, edgeDevice, heartbeat)
	if err != nil {
		return err, true
	}
	err = h.updater.updateLabels(ctx, edgeDevice, heartbeat)
	if err != nil {
		return err, true
	}

	return nil, false
}
