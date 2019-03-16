package main

import (
	"flag"
	"github.com/mihard/behat-parallel-runner/runner"
	"log"
)

const BEHAT = "./vendor/bin/behat"

func main() {
	log.Printf("Run behat tests")

	var wNum int
	var behatCmd string
	var verbose bool

	flag.IntVar(&wNum, "workers", 4, "amount of workers")
	flag.StringVar(&behatCmd, "behat", BEHAT, "path to behat executable")
	flag.BoolVar(&verbose, "verbose", false, "add the flag to see outputs of all scenarios")

	flag.Parse()

	behatArgs := flag.Args()

	index, err := runner.GetIndexOfScenarios(behatCmd, behatArgs)
	if err != nil {
		log.Fatalf("Unable to fetch list of scenarios: %s", err.Error())
	}

	if len(index) < 1 {
		log.Println("No scenarios found")
		return
	}

	log.Printf("%d Scenario(s) found", len(index))

	sc := make(chan runner.Scenario, len(index))
	rc := make(chan runner.Result, len(index))

	for w := 1; w <= wNum; w++ {
		go runner.Worker(behatCmd, w, sc, rc)
	}

	for _, s := range index {
		sc <- s
	}
	close(sc)

	fCnt := 0

	for r := 0; r < len(index); r++ {
		res := <-rc

		if res.Ok {
			log.Printf("%s [Done in %.2f sec] %s... OK \n", res.Scenario.File, res.ExecutionTime.Seconds(), res.Scenario.Scenario)

			if verbose {
				dumpOutput(res.Output)
			}
		} else {
			log.Printf("%s [Done in %.2f sec] %s... FAILED\n", res.Scenario.File, res.ExecutionTime.Seconds(), res.Scenario.Scenario)

			fCnt++

			dumpOutput(res.Output)
		}
	}
	close(rc)

	if fCnt > 0 {
		log.Println("")
		log.Fatalf("%d Scenario(s) of %d failed ", fCnt, len(index))
	}
}

func dumpOutput(output []byte) {
	for _, l := range runner.AsLineArray(output) {
		log.Printf("\t\t%s", l)
	}

	log.Println("---")
}
