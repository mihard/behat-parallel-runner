package runner

import (
	"fmt"
	"os/exec"
	"time"
)

type Result struct {
	Scenario      Scenario
	WorkerNumber  int
	Ok            bool
	Output        []byte
	ExecutionTime time.Duration
}

func Worker(wn int, sc chan Scenario, rc chan Result) {
	for s := range sc {
		rc <- run(wn, s)
	}
}

func run(wn int, s Scenario) Result {
	start := time.Now()
	cmd := exec.Command(BEHAT, s.File, fmt.Sprintf("--name=%s", s.Scenario))
	executionTime := time.Since(start)

	output, err := cmd.Output()

	if err != nil {
		return Result{
			WorkerNumber:  wn,
			Scenario:      s,
			Ok:            false,
			Output:        output,
			ExecutionTime: executionTime,
		}
	}

	return Result{WorkerNumber: wn, Scenario: s, Ok: true, ExecutionTime: executionTime}
}
