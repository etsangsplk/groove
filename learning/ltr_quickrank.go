package learning

import (
	"io"
	"os/exec"
	"bufio"
	"log"
	"fmt"
	"os"
	"math"
)

type QuickRankQueryCandidateSelector struct {
	// The path to the binary file for execution.
	binary string
	// Maximum depth allowed to generate queries.
	depth int
	// Command-line arguments for configuration.
	arguments map[string]interface{}
}

func makeArguments(a map[string]interface{}) []string {
	// Load the arguments from the map.
	args := make([]string, len(a)*2)
	i := 0
	for k, v := range a {
		args[i] = fmt.Sprintf("--%s", k)
		args[i+1] = fmt.Sprintf("%v", v)
		i += 2
	}
	return args
}

func (qr QuickRankQueryCandidateSelector) Select(query CandidateQuery, transformations []CandidateQuery) (CandidateQuery, QueryChainCandidateSelector, error) {
	args := makeArguments(qr.arguments)
	args = append(args, "--test", "tmp.features")

	// Create a temporary file to contain the features for testing.
	f, err := os.OpenFile("tmp.features", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Write the features of the variation to temporary file.
	for _, applied := range transformations {
		_, err := f.WriteString(fmt.Sprintf("0 qid:%s %s\n", query.Topic, applied.Features.String()))
		if err != nil {
			return query, qr, err
		}
	}

	// Configure the command.
	cmd := exec.Command(qr.binary, args...)

	// Open channels to stdout and stderr.
	r, err := cmd.StdoutPipe()
	if err != nil {
		return query, qr, err
	}
	defer r.Close()

	e, err := cmd.StderrPipe()
	if err != nil {
		return query, qr, err
	}
	defer e.Close()

	// Start the command.
	cmd.Start()

	// Output the stdout pipe.
	go func() {
		s := bufio.NewScanner(r)
		for s.Scan() {
			log.Println(s.Text())
		}
		return
	}()

	// Output the stderr pipe.
	go func() {
		s := bufio.NewScanner(e)
		for s.Scan() {
			log.Println(s.Text())
		}
		return
	}()

	// Wait for the command to finish.
	if err := cmd.Wait(); err != nil {
		return query, qr, err
	}

	// Grab the top-most ranked query from the candidates.
	candidate, err := getRanking(qr.arguments["scores"].(string), transformations)
	if err != nil {
		return query, qr, err
	}

	// Totally remove the file.
	f.Truncate(0)
	f.Seek(0, 0)
	err2 := os.Remove("tmp.features")
	if err2 != nil {
		return CandidateQuery{}, nil, err2
	}

	qr.depth++

	if query.Query.String() == candidate.String() {
		qr.depth = math.MaxInt32
	}

	q := candidate.Append(query)
	return q, qr, nil
}

func (qr QuickRankQueryCandidateSelector) Train(lfs []LearntFeature) ([]byte, error) {
	args := makeArguments(qr.arguments)

	// Configure the command.
	cmd := exec.Command(qr.binary, args...)

	// Open channels to stdout and stderr.
	r, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer r.Close()

	e, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	defer e.Close()

	// Start the command.
	cmd.Start()

	// Output the stdout pipe.
	go func() {
		s := bufio.NewScanner(r)
		for s.Scan() {
			log.Println(s.Text())
		}
		return
	}()

	// Output the stderr pipe.
	go func() {
		s := bufio.NewScanner(e)
		for s.Scan() {
			log.Println(s.Text())
		}
		return
	}()

	// Wait for the command to finish.
	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (QuickRankQueryCandidateSelector) Output(lf LearntFeature, w io.Writer) error {
	_, err := lf.WriteLibSVMRank(w)
	return err
}

func (qr QuickRankQueryCandidateSelector) StoppingCriteria() bool {
	return qr.depth >= 5
}

func NewQuickRankQueryCandidateSelector(binary string, arguments map[string]interface{}) QuickRankQueryCandidateSelector {
	return QuickRankQueryCandidateSelector{
		binary:    binary,
		arguments: arguments,
		depth:     0,
	}
}