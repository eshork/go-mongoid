package grifts

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	g "github.com/markbates/grift/grift"
)

func execCmd(cmd string) error {
	args := strings.Split(cmd, " ")
	binary, lookErr := exec.LookPath(args[0])
	if lookErr != nil {
		panic(lookErr)
	}
	env := os.Environ()
	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		return execErr
	}
	return nil
}

var _ = g.Namespace("env", func() {
	g.Desc("print", "Prints out all of the ENV variables in your environment. Pass in the name of a particular ENV variable to print just that one out. (e.g. grift env:print GOPATH)")
	g.Add("print", func(c *g.Context) error {
		if len(c.Args) >= 1 {
			for _, e := range c.Args {
				fmt.Printf("%s=%s\n", e, os.Getenv(e))
			}
		} else {
			for _, e := range os.Environ() {
				pair := strings.Split(e, "=")
				fmt.Printf("%s=%s\n", pair[0], os.Getenv(pair[0]))
			}
		}
		return nil
	})
})

var _ = g.Desc("fmt", "runs gofmt in the standard project manner")
var _ = g.Add("fmt", func(c *g.Context) error {
	execErr := execCmd("gofmt -l -w .")
	if execErr != nil {
		panic(execErr)
	}
	return nil
})

var _ = g.Desc("docs", "run a local doc server")
var _ = g.Add("docs", func(c *g.Context) error {
	fmt.Printf("Starting godoc server at: http://localhost:6060/\n")
	execErr := execCmd("go run golang.org/x/tools/cmd/godoc -index -http localhost:6060 -goroot ./")
	if execErr != nil {
		panic(execErr)
	}
	return nil
})

var _ = g.Desc("test", "run basic tests (no benchmarks)")
var _ = g.Add("test", func(c *g.Context) error {
	execErr := execCmd("go run github.com/onsi/ginkgo/ginkgo -r -skipMeasurements")
	if execErr != nil {
		panic(execErr)
	}
	return nil
})

var _ = g.Namespace("test", func() {
	g.Desc("watch", "run tests within a watch loop")
	g.Add("watch", func(c *g.Context) error {
		execErr := execCmd("go run github.com/onsi/ginkgo/ginkgo watch -r -skipMeasurements -succinct")
		if execErr != nil {
			panic(execErr)
		}
		return nil
	})

	g.Desc("force", "run all tests even if some fail")
	g.Add("force", func(c *g.Context) error {
		execErr := execCmd("go run github.com/onsi/ginkgo/ginkgo -r -keepGoing -skipMeasurements")
		if execErr != nil {
			panic(execErr)
		}
		return nil
	})

	g.Desc("bench", "run performance benchmark tests")
	g.Add("bench", func(c *g.Context) error {
		execErr := execCmd("go run github.com/onsi/ginkgo/ginkgo -r -keepGoing -focus performance")
		if execErr != nil {
			panic(execErr)
		}
		return nil
	})

	g.Desc("all", "run all tests")
	g.Add("all", func(c *g.Context) error {
		execErr := execCmd("go run github.com/onsi/ginkgo/ginkgo -r")
		if execErr != nil {
			panic(execErr)
		}
		return nil
	})

	g.Desc("ci", "run all tests as ci expects")
	g.Add("ci", func(c *g.Context) error {
		// https://onsi.github.io/ginkgo/#ginkgo-and-continuous-integration
		// -r --randomizeAllSpecs --randomizeSuites --failOnPending --cover --trace --race --progress
		execErr := execCmd("go run github.com/onsi/ginkgo/ginkgo -r --randomizeAllSpecs --randomizeSuites --failOnPending --cover --trace --race --progress")
		if execErr != nil {
			panic(execErr)
		}
		return nil
	})

})
