package heartbeat

import (
	"context"
	"fmt"
	"github.com/jakub-dzon/k4e-operator/internal/repository/edgedevice"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

type ProcessAllHandler struct {
	deviceRepository edgedevice.Repository
	updater          Updater
	notifications    chan Notification
}

func NewProcessAllHandler(deviceRepository edgedevice.Repository, recorder record.EventRecorder) *ProcessAllHandler {
	return &ProcessAllHandler{
		deviceRepository: deviceRepository,
		notifications:    make(chan Notification),
		updater: Updater{
			deviceRepository: deviceRepository,
			recorder:         recorder,
		},
	}
}

func (h *ProcessAllHandler) Start() {
	for i := 0; i < 5; i++ {
		go func() {
			for j := range h.notifications {
				err := h.process(context.Background(), j)
				if err != nil {
					j.Retry++
					h.notifications <- j
				}
			}
		}()
	}
}

func (h *ProcessAllHandler) Process(_ context.Context, notification Notification) error {
	h.notifications <- notification
	return nil
}

func (h *ProcessAllHandler) process(ctx context.Context, notification Notification) error {
	logger := log.FromContext(ctx, "DeviceID", notification.DeviceID, "Namespace", notification.Namespace)
	heartbeat := notification.Heartbeat
	logger.V(1).Info("processing heartbeat", "content", heartbeat)
	edgeDevice, err := h.deviceRepository.Read(ctx, notification.DeviceID, notification.Namespace)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.V(1).Info("Device not found")
			return nil
		}
		return fmt.Errorf("can't read device %s/%s", notification.Namespace, notification.DeviceID)
	}

	if edgeDevice.Status.LastSeenTime.After(time.Time(heartbeat.Time)) {
		logger.V(1).Info("heartbeat outdated")
	} else {
		err := h.updater.updateStatus(ctx, edgeDevice, heartbeat)
		if err != nil {
			return err
		}
	}

	err = h.updater.updateLabels(ctx, edgeDevice, heartbeat)
	if err != nil {
		return err
	}

	// Produce k8s events based on the device-worker events:
	if notification.Retry == 0 {
		h.updater.processEvents(edgeDevice, heartbeat.Events)
	}
	return nil
}
