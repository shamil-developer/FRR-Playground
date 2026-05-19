package handlers

import (
	"frr-playground/frrpb"
	"frr-playground/internal/models"
)

func parseCapabilityModules(
	step models.Step,
) []string {

	rawModules := step.Params["modules"].([]any)

	modules := make([]string, 0, len(rawModules))

	for _, item := range rawModules {
		modules = append(
			modules,
			item.(string),
		)
	}

	return modules
}

func hasCapabilityModule(
	response *frrpb.GetCapabilitiesResponse,
	module string,
) bool {

	for _, item := range response.GetSupportedModules() {

		if item.GetName() == module {
			return true
		}
	}

	return false
}
