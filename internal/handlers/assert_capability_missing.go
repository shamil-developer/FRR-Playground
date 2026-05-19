package handlers

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"frr-playground/internal/models"
	"frr-playground/internal/utils"
)

type AssertCapabilityMissingHandler struct{}

func NewAssertCapabilityMissingHandler() *AssertCapabilityMissingHandler {
	return &AssertCapabilityMissingHandler{}
}

func (h *AssertCapabilityMissingHandler) Execute(
	ctx context.Context,
	conn *grpc.ClientConn,
	flow map[string]any,
	step models.Step,
) error {

	response, err := getCapabilities(
		ctx,
		conn,
	)

	if err != nil {
		return err
	}

	modules := parseCapabilityModules(
		step,
	)

	for _, module := range modules {

		if hasCapabilityModule(response, module) {
			return fmt.Errorf("assert capability missing failed: %s exists", module)
		}

		err = utils.PrintAny(
			map[string]any{
				"assert": "capability-missing",
				"module": module,
				"status": "ok",
			},
		)

		if err != nil {
			return err
		}
	}

	return nil
}
