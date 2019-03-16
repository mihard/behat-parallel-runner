package runner

import (
	"fmt"
	"os/exec"
	"time"
)

type Result struct {
	Scenario
	WorkerNumber  int
	Ok            bool
	Output        []byte
	ExecutionTime time.Duration
}

func Worker(behatCmd string, wn int, sc chan Scenario, rc chan Result) {
	for s := range sc {
		rc <- run(behatCmd, wn, s)
	}
}

func run(behatCmd string, wn int, s Scenario) Result {
	start := time.Now()
	cmd := exec.Command(behatCmd, s.File, fmt.Sprintf("--name=%s", s.Scenario))
	output, err := cmd.Output()
	executionTime := time.Since(start)

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
