package runner

import (
	"fmt"
	"os/exec"
)

type Result struct {
	Scenario
	WorkerNumber int
	Ok           bool
	Output       []byte
}

func Worker(wn int, sc chan Scenario, rc chan Result) {
	for s := range sc {
		rc <- run(wn, s)
	}
}

func run(wn int, s Scenario) Result {
	cmd := exec.Command(BEHAT, s.File, fmt.Sprintf("--name=%s", s.Scenario))
	output, err := cmd.Output()

	if err != nil {
		return Result{
			WorkerNumber: wn,
			Scenario: s,
			Ok:       false,
			Output:   output,
		}
	}

	return Result{WorkerNumber: wn, Scenario: s, Ok: true}
}
