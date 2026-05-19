package handlers

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"frr-playground/frrpb"
	"frr-playground/internal/models"
	"frr-playground/internal/utils"
)

type LoadToCandidateHandler struct{}

func NewLoadToCandidateHandler() *LoadToCandidateHandler {
	return &LoadToCandidateHandler{}
}

func (h *LoadToCandidateHandler) Execute(
	ctx context.Context,
	conn *grpc.ClientConn,
	flow map[string]any,
	step models.Step,
) error {
	client := frrpb.NewNorthboundClient(conn)
	candidateId := flow["candidateId"].(uint32)

	loadType := frrpb.LoadToCandidateRequest_MERGE
	if rawType, ok := step.Params["type"]; ok && rawType != nil {
		switch fmt.Sprint(rawType) {
		case "MERGE":
			loadType = frrpb.LoadToCandidateRequest_MERGE
		case "REPLACE":
			loadType = frrpb.LoadToCandidateRequest_REPLACE
		default:
			return fmt.Errorf("unknown load type %q", rawType)
		}
	}

	data := "{}"
	if rawData, ok := step.Params["data"]; ok && rawData != nil {
		data = fmt.Sprint(rawData)
	}

	response, err := client.LoadToCandidate(
		ctx,
		&frrpb.LoadToCandidateRequest{
			CandidateId: candidateId,
			Type:        loadType,
			Config: &frrpb.DataTree{
				Encoding: frrpb.Encoding_JSON,
				Data:     data,
			},
		},
	)
	if err != nil {
		return err
	}

	return utils.PrintAny(response)
}
