package handlers

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"

	"frr-playground/frrpb"
	"frr-playground/internal/models"
)

type assertJSONCheck struct {
	Path   string
	Equals *string
}

func parseAssertJSONChecks(
	step models.Step,
) []assertJSONCheck {

	rawChecks := step.Params["checks"].([]any)

	checks := make([]assertJSONCheck, 0, len(rawChecks))

	for _, item := range rawChecks {

		rawCheck := item.(map[string]any)

		check := assertJSONCheck{
			Path: rawCheck["path"].(string),
		}

		if rawEquals, ok := rawCheck["equals"]; ok {

			equals := fmt.Sprint(rawEquals)

			check.Equals = &equals
		}

		checks = append(
			checks,
			check,
		)
	}

	return checks
}

func getJSONDocument(
	ctx context.Context,
	conn *grpc.ClientConn,
	step models.Step,
) (string, error) {

	client := frrpb.NewNorthboundClient(conn)

	requestType := frrpb.GetRequest_DataType(
		step.Params["requestType"].(int),
	)

	rawPaths := step.Params["path"].([]any)

	paths := make([]string, 0, len(rawPaths))

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
			Type:         requestType,
			Encoding:     frrpb.Encoding_JSON,
			WithDefaults: withDefaults,
			Path:         paths,
		},
	)

	if err != nil {
		return "", err
	}

	document := "{}"

	for {

		response, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			return "", err
		}

		if data := response.GetData().GetData(); data != "" {
			document = data
		}
	}

	return document, nil
}
