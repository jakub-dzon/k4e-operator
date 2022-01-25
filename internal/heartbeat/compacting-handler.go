package heartbeat

import (
	"context"
	"github.com/jakub-dzon/k4e-operator/internal/repository/edgedevice"
	"k8s.io/client-go/tools/record"
)

type CompactingHandler struct {
	deviceRepository edgedevice.Repository
	recorder         record.EventRecorder
	notifications    chan Notification
}

func NewCompactingHandler(deviceRepository edgedevice.Repository, recorder record.EventRecorder) *CompactingHandler {
	return &CompactingHandler{
		deviceRepository: deviceRepository,
		recorder:         recorder,
		notifications:    make(chan Notification),
	}
}

// https://medium.com/capital-one-tech/building-an-unbounded-channel-in-go-789e175cd2cd

func (h *CompactingHandler) Start() {

}

func (h *CompactingHandler) Process(_ context.Context, notification Notification) error {
	h.notifications <- notification
	return nil
}
