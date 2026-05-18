package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"gopkg.in/yaml.v3"

	"frr-playground/internal/models"
	"frr-playground/internal/workflow"
)

func main() {

	plan := loadPlan(
		getWorkflowPath(),
	)

	options := loadOptions()

	connections := createGRPCConnections(
		options,
	)

	defer closeConnections(
		connections,
	)

	flow := map[string]any{}

	registry := workflow.CreateHandlersRegistry()

	executePlan(
		connections,
		flow,
		registry,
		plan,
	)
}

func getWorkflowPath() string {

	workflowPath := flag.String(
		"workflow",
		"configs/playground/workflows/01-check.yaml",
		"path to workflow yaml",
	)

	flag.Parse()

	if flag.NArg() > 0 {
		return flag.Arg(0)
	}

	return *workflowPath
}

func loadPlan(
	workflowPath string,
) models.Plan {

	rawPlan, err := os.ReadFile(
		workflowPath,
	)

	if err != nil {
		log.Fatal(err)
	}

	var plan models.Plan

	err = yaml.Unmarshal(
		rawPlan,
		&plan,
	)

	if err != nil {
		log.Fatal(err)
	}

	return plan
}

func loadOptions() models.Options {

	rawOptions, err := os.ReadFile(
		"configs/playground/options.yaml",
	)

	if err != nil {
		log.Fatal(err)
	}

	var options models.Options

	err = yaml.Unmarshal(
		rawOptions,
		&options,
	)

	if err != nil {
		log.Fatal(err)
	}

	return options
}

func createGRPCConnections(
	options models.Options,
) map[string]*grpc.ClientConn {

	connections := map[string]*grpc.ClientConn{}

	for daemon, daemonOptions := range options.Daemons {

		conn, err := grpc.NewClient(
			daemonOptions.Target,
			grpc.WithTransportCredentials(
				insecure.NewCredentials(),
			),
		)

		if err != nil {
			log.Fatal(err)
		}

		connections[daemon] = conn
	}

	return connections
}

func closeConnections(
	connections map[string]*grpc.ClientConn,
) {

	for _, conn := range connections {

		conn.Close()
	}
}

func executePlan(
	connections map[string]*grpc.ClientConn,
	flow map[string]any,
	registry map[string]workflow.Handler,
	plan models.Plan,
) {

	for index, step := range plan.Steps {

		printStepHeader(
			index,
			step,
		)

		handler := getHandler(
			registry,
			step,
		)

		executeStep(
			connections,
			flow,
			handler,
			step,
		)
	}
}

func printStepHeader(
	index int,
	step models.Step,
) {

	const (
		reset = "\033[0m"

		bold = "\033[1m"

		cyan   = "\033[38;5;51m"
		gray   = "\033[38;5;244m"
		yellow = "\033[38;5;226m"

		line = "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	)

	fmt.Println()

	fmt.Printf(
		"%s%s%s\n",
		gray,
		line,
		reset,
	)

	fmt.Printf(
		"%s%sSTEP %d%s %s→%s %s%s%s\n",
		bold,
		cyan,
		index+1,
		reset,
		yellow,
		reset,
		bold,
		step.Type,
		reset,
	)

	if step.Description != "" {

		fmt.Println()

		fmt.Printf(
			"%s%s%s\n",
			gray,
			step.Description,
			reset,
		)
	}

	fmt.Printf(
		"%s%s%s\n",
		gray,
		line,
		reset,
	)
}

func getHandler(
	registry map[string]workflow.Handler,
	step models.Step,
) workflow.Handler {

	handler, ok := registry[step.Type]

	if !ok {

		log.Fatalf(
			"handler not found: %s",
			step.Type,
		)
	}

	return handler
}

func executeStep(
	connections map[string]*grpc.ClientConn,
	flow map[string]any,
	handler workflow.Handler,
	step models.Step,
) {

	conn, ok := connections[step.Daemon]

	if !ok {

		log.Fatalf(
			"daemon connection not found: %s",
			step.Daemon,
		)
	}

	err := handler.Execute(
		context.Background(),
		conn,
		flow,
		step,
	)

	if err != nil {
		log.Fatal(err)
	}
}
