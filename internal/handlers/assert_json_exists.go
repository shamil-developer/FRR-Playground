package handlers

import (
	"context"
	"fmt"

	"github.com/tidwall/gjson"
	"google.golang.org/grpc"

	"frr-playground/internal/models"
	"frr-playground/internal/utils"
)

type AssertJSONExistsHandler struct{}

func NewAssertJSONExistsHandler() *AssertJSONExistsHandler {
	return &AssertJSONExistsHandler{}
}

func (h *AssertJSONExistsHandler) Execute(
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

		if !result.Exists() {
			return fmt.Errorf("assert json exists failed: %s does not exist", check.Path)
		}

		if check.Equals != nil && result.String() != *check.Equals {
			return fmt.Errorf("assert json exists failed: %s expected %q, got %q", check.Path, *check.Equals, result.String())
		}

		err = utils.PrintAny(
			map[string]any{
				"assert": "exists",
				"path":   check.Path,
				"value":  result.Value(),
				"status": "ok",
			},
		)

		if err != nil {
			return err
		}
	}

	return nil
}
