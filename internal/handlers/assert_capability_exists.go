package handlers

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"frr-playground/frrpb"
	"frr-playground/internal/models"
	"frr-playground/internal/utils"
)

type AssertCapabilityExistsHandler struct{}

func NewAssertCapabilityExistsHandler() *AssertCapabilityExistsHandler {
	return &AssertCapabilityExistsHandler{}
}

func (h *AssertCapabilityExistsHandler) Execute(
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

		if !hasCapabilityModule(response, module) {
			return fmt.Errorf("assert capability exists failed: %s does not exist", module)
		}

		err = utils.PrintAny(
			map[string]any{
				"assert": "capability-exists",
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

func getCapabilities(
	ctx context.Context,
	conn *grpc.ClientConn,
) (*frrpb.GetCapabilitiesResponse, error) {

	client := frrpb.NewNorthboundClient(conn)

	return client.GetCapabilities(
		ctx,
		&frrpb.GetCapabilitiesRequest{},
	)
}
