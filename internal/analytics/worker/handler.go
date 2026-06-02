package worker

import (
	"context"

	"github.com/bellapacx/kids-utopia/internal/analytics/service"
)
type Worker struct {
	svc *service.Service
}
func New(svc *service.Service) *Worker {
	return &Worker{svc: svc}
}
func (w *Worker) Handle(msg string) error {
	return w.svc.ProcessMessage(context.Background(), msg)
}