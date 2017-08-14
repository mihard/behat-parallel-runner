package main

import (
	"github.com/mihard/behat-parallel-runner/runner"
	"log"
	"os"
	"strconv"
)

func main() {
	log.Printf("Run behat tests")

	wNum := 4
	behatArgsStart := 1

	if len(os.Args) > 1 {
		if parsedNum, err := strconv.Atoi(os.Args[1]); err == nil {
			wNum = parsedNum
			behatArgsStart = 2
		}
	}

	behatArgs := os.Args[behatArgsStart:]

	index, err := runner.GetIndexOfScenarios(behatArgs)

	log.Printf("%d Scenario(s) found", len(index))

	if err != nil {
		panic(err)
	}

	if len(index) < 1 {
		log.Println("No scenarios found")
		return
	}

	sc := make(chan runner.Scenario, len(index))
	rc := make(chan runner.Result, len(index))

	for w := 1; w <= wNum; w++ {
		go runner.Worker(w, sc, rc)
	}

	for _, s := range index {
		sc <- s
	}
	close(sc)

	fCnt := 0

	for r := 0; r < len(index); r++ {
		res := <-rc

		if res.Ok {
			log.Printf("%s %s... OK\n", res.File, res.Scenario.Scenario)
		} else {
			log.Printf("%s %s... FAILED\n", res.File, res.Scenario.Scenario)

			fCnt++

			for _, l := range runner.AsLineArray(res.Output) {
				log.Printf("\t\t%s", l)
			}
		}
	}
	close(rc)

	if fCnt > 0 {
		log.Println("")
		log.Printf("%d Scenario(s) of %d failed ", fCnt, len(index))

		os.Exit(1)
	}
}
