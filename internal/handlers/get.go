package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"google.golang.org/grpc"

	"frr-playground/frrpb"
	"frr-playground/internal/models"
	"frr-playground/internal/utils"
)

type GetHandler struct{}

func NewGetHandler() *GetHandler {
	return &GetHandler{}
}

func (h *GetHandler) Execute(
	ctx context.Context,
	conn *grpc.ClientConn,
	flow map[string]any,
	step models.Step,
) error {

	client := frrpb.NewNorthboundClient(conn)

	requestType := frrpb.GetRequest_DataType(
		step.Params["requestType"].(int),
	)

	rawPaths := step.Params["path"].([]any)

	paths := make([]string, 0)

	for _, item := range rawPaths {

		paths = append(
			paths,
			item.(string),
		)
	}

	withDefaults := true

	if rawWithDefaults, ok := step.Params["withDefaults"]; ok {
		withDefaults = rawWithDefaults.(bool)
	}

	stream, err := client.Get(
		ctx,
		&frrpb.GetRequest{
			Type: requestType,

			Encoding: frrpb.Encoding_JSON,

			WithDefaults: withDefaults,

			Path: paths,
		},
	)

	if err != nil {
		return err
	}

	for {

		response, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		var parsed any

		err = json.Unmarshal(
			[]byte(response.GetData().GetData()),
			&parsed,
		)

		if err != nil {

			parsed = response.GetData().GetData()
		}

		fmt.Println()

		fmt.Println(
			"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━",
		)

		err = utils.PrintAny(
			map[string]any{
				"timestamp": response.GetTimestamp(),

				"encoding": response.GetData().
					GetEncoding().
					String(),

				"data": parsed,
			},
		)

		if err != nil {
			return err
		}
	}

	return nil
}
