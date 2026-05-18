package handlers

import (
	"context"
	"fmt"

	"github.com/tidwall/gjson"
	"google.golang.org/grpc"

	"frr-playground/internal/models"
	"frr-playground/internal/utils"
)

type AssertJSONMissingHandler struct{}

func NewAssertJSONMissingHandler() *AssertJSONMissingHandler {
	return &AssertJSONMissingHandler{}
}

func (h *AssertJSONMissingHandler) Execute(
	ctx context.Context,
	conn *grpc.ClientConn,
	flow map[string]any,
	step models.Step,
) error {

	document, err := getJSONDocument(
		ctx,
		conn,
		step,
	)

	if err != nil {
		return err
	}

	checks := parseAssertJSONChecks(
		step,
	)

	for _, check := range checks {

		result := gjson.Get(
			document,
			check.Path,
		)

		if result.Exists() {
			return fmt.Errorf("assert json missing failed: %s exists with value %q", check.Path, result.String())
		}

		err = utils.PrintAny(
			map[string]any{
				"assert": "missing",
				"path":   check.Path,
				"status": "ok",
			},
		)

		if err != nil {
			return err
		}
	}

	return nil
}
