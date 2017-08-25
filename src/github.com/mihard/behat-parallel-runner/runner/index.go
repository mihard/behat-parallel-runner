package runner

import (
	"bufio"
	"bytes"
	"errors"
	"os/exec"
	"regexp"
)

type Scenario struct {
	File     string
	Scenario string
}

func GetIndexOfScenarios(behatArgs []string) (index []Scenario, err error) {

	indexerArgs := append([]string{"--dry-run", "--no-colors", "--no-interaction"}, behatArgs...)

	cmd := exec.Command(BEHAT, indexerArgs...)
	output, err := cmd.Output()

	if err != nil {
		return
	}

	return readOutput(output)
}

func readOutput(output []byte) (index []Scenario, err error) {
	rx := regexp.MustCompile("Scenario(?:\\sOutline)?\\s*:\\s*(\\w.*\\w)\\s*#\\s*(.*):")

	index = []Scenario{}

	for _, l := range AsLineArray(output) {
		if s, el := readLine(l, rx); el == nil {
			index = append(index, s)
		}
	}

	return
}

func readLine(line string, rx *regexp.Regexp) (s Scenario, err error) {

	out := rx.FindStringSubmatch(line)

	if len(out) < 3 {
		return s, errors.New("no scenario at this line")
	}

	return Scenario{
		out[2],
		out[1],
	}, nil
}

func AsLineArray(output []byte) []string {
	b := bytes.NewBuffer(output)
	r := bufio.NewReader(b)

	lines := []string{}

	for {
		l, err := r.ReadString('\n')

		if err != nil {
			break
		}

		lines = append(lines, l)
	}

	return lines
}
