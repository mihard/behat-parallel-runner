package main

import (
	"github.com/mihard/behat-parallel-runner/runner"
	"log"
	"os"
	"github.com/jessevdk/go-flags"
)

func main() {
	log.Printf("Run behat tests")
	type Options struct {
    	BehatPath string `short:"p" long:"behat-path" description:"Path to behat executable" required:"true"`
    	ProcessNumber int `short:"n" long:"process-number" description:"Process number" default:"4"`
    }

    var opts Options

    var parser = flags.NewParser(&opts, flags.IgnoreUnknown)
    args, err := parser.Parse()
    if err != nil {
    	panic(err)
    }

	index, err := runner.GetIndexOfScenarios(opts.BehatPath, args)

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

	for w := 1; w <= opts.ProcessNumber; w++ {
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
			log.Printf("%s [Done in %.2f sec] %s... OK \n", res.Scenario.File, res.ExecutionTime.Seconds(), res.Scenario.Scenario)
		} else {
			log.Printf("%s [Done in %.2f sec] %s... FAILED\n", res.Scenario.File, res.ExecutionTime.Seconds(), res.Scenario.Scenario)

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
