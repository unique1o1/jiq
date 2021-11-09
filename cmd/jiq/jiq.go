package main

import (
	"fmt"
	"os"

	"github.com/fiatjaf/jiq"
)

var version string

// Both runner and mockRunner will use the jiqRunner interface
type jiqRunner interface {
	run() int
}

type runner struct {
	engine        *jiq.Engine
	outputQuery   bool
	outputResults bool
}

func (r runner) run() int {
	result := r.engine.Run()
	if result.Err != nil {
		return 2
	}
	printResults(result, r.outputQuery, r.outputResults)
	return 0
}

func main() {
	content := os.Stdin

	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println(version)
		os.Exit(0)
	}
	if len(os.Args) > 1 && os.Args[1] == "--help" {
		fmt.Print(`jiq - interactive commandline JSON processor
Usage: <json string> | jiq [options] [initial filter]

    jiq is a tool that allows you to play with jq filters interactively
    acting direcly on a JSON source of your choice, given through STDIN.
    For all the details about which filters you can use to transform your
    JSON string, see jq(1) manpage or https://stedolan.github.io/jq

    jiq supports all command line arguments jq supports, plus
     -b         will print the ending filter and filtered results.
     -q         will print the ending filter to STDOUT, instead of
                printing the resulting filtered JSON, the default.
     --help     prints this help message.
     --version  prints version
      `)
		os.Exit(0)
	}

	initialquery := "."
	outputquery := false
	outputresults := true
	jqargs := os.Args[1:]
	for i, arg := range os.Args[1:] {
		i = i + 1
		if arg == "-b" {
			outputquery = true
			jqargs = os.Args[1:i]
			jqargs = append(jqargs, os.Args[i+1:]...)
			break
		} else if arg == "-q" {
			outputquery = true
			outputresults = false
			jqargs = os.Args[1:i]
			jqargs = append(jqargs, os.Args[i+1:]...)
			break
		} else if arg[0] != '-' {
			initialquery = arg
			jqargs = os.Args[1:i]
			jqargs = append(jqargs, os.Args[i+1:]...)
			break
		}
	}

	e := jiq.NewEngine(content, jqargs, initialquery)
	r := runner{engine: e, outputQuery: outputquery, outputResults: outputresults}
	os.Exit(doRun(r))
}

// TODO: All of the runs should return an EngineResult or OutputResult for proper mocking/unit testing
func doRun(jr jiqRunner) int {
	return jr.run()
}

func printResults(res *jiq.EngineResult, outputquery bool, outputresults bool) {
	if outputquery && outputresults {
		fmt.Printf("%s\njq '%s'\n", res.Content, res.Qs)
	} else if outputquery {
		fmt.Printf("%s\n", res.Qs)
	} else {
		fmt.Printf("%s\n", res.Content)
	}
}
