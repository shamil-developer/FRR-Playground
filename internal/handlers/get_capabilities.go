package handlers

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc"

	"google.golang.org/protobuf/encoding/protojson"

	"frr-playground/frrpb"
	"frr-playground/internal/models"
	"frr-playground/internal/utils"
)

type GetCapabilitiesHandler struct{}

func NewGetCapabilitiesHandler() *GetCapabilitiesHandler {
	return &GetCapabilitiesHandler{}
}

func (h *GetCapabilitiesHandler) Execute(
	ctx context.Context,
	conn *grpc.ClientConn,
	flow map[string]any,
	step models.Step,
) error {

	client := frrpb.NewNorthboundClient(conn)

	response, err := client.GetCapabilities(
		ctx,
		&frrpb.GetCapabilitiesRequest{},
	)

	if err != nil {
		return err
	}

	raw, err := protojson.Marshal(
		response,
	)

	if err != nil {
		return err
	}

	var parsed any

	err = json.Unmarshal(
		raw,
		&parsed,
	)

	if err != nil {
		return err
	}

	return utils.PrintAny(parsed)
}
