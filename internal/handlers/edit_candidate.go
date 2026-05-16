package handlers

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"frr-playground/frrpb"
	"frr-playground/internal/models"
	"frr-playground/internal/utils"
)

type EditCandidateHandler struct{}

func NewEditCandidateHandler() *EditCandidateHandler {
	return &EditCandidateHandler{}
}

func (h *EditCandidateHandler) Execute(
	ctx context.Context,
	conn *grpc.ClientConn,
	flow map[string]any,
	step models.Step,
) error {

	client := frrpb.NewNorthboundClient(conn)

	candidateId := flow["candidateId"].(uint32)

	var updates []*frrpb.PathValue
	var deletes []*frrpb.PathValue

	if rawUpdates, ok := step.Params["update"]; ok {

		for _, item := range rawUpdates.([]any) {

			update := item.(map[string]any)

			path := update["path"].(string)

			value := ""

			if v, ok := update["value"]; ok && v != nil {
				value = fmt.Sprint(v)
			}

			updates = append(
				updates,
				&frrpb.PathValue{
					Path:  path,
					Value: value,
				},
			)
		}
	}

	if rawDeletes, ok := step.Params["delete"]; ok {

		for _, item := range rawDeletes.([]any) {

			del := item.(map[string]any)

			deletes = append(
				deletes,
				&frrpb.PathValue{
					Path: del["path"].(string),
				},
			)
		}
	}

	response, err := client.EditCandidate(
		ctx,
		&frrpb.EditCandidateRequest{
			CandidateId: candidateId,
			Update:      updates,
			Delete:      deletes,
		},
	)

	if err != nil {
		return err
	}

	return utils.PrintAny(response)
}
