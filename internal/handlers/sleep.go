package handlers

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"frr-playground/internal/models"
)

type SleepHandler struct{}

func NewSleepHandler() *SleepHandler {
	return &SleepHandler{}
}

func (h *SleepHandler) Execute(
	ctx context.Context,
	conn *grpc.ClientConn,
	flow map[string]any,
	step models.Step,
) error {

	seconds := step.Params["seconds"].(int)

	select {
	case <-time.After(time.Duration(seconds) * time.Second):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
