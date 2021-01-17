package grifts

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	. "github.com/markbates/grift/grift"
)

var _ = Namespace("env", func() {
	Desc("print", "Prints out all of the ENV variables in your environment. Pass in the name of a particular ENV variable to print just that one out. (e.g. grift env:print GOPATH)")
	Add("print", func(c *Context) error {
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

var _ = Desc("fmt", "runs gofmt in the standard project manner")
var _ = Add("fmt", func(c *Context) error {
	binary, lookErr := exec.LookPath("gofmt")
	if lookErr != nil {
		panic(lookErr)
	}
	args := []string{"gofmt", "-l", "-w", "."}
	env := os.Environ()
	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		panic(execErr)
	}
	return nil
})

var _ = Desc("docs", "run a local doc server")
var _ = Add("docs", func(c *Context) error {
	binary, lookErr := exec.LookPath("go")
	if lookErr != nil {
		panic(lookErr)
	}
	args := []string{"go", "run", "golang.org/x/tools/cmd/godoc", "-index", "-http", "localhost:6060", "-goroot", "./"}
	env := os.Environ()
	env = append(env, "GO111MODULE=on")
	fmt.Printf("Starting godoc server at: http://localhost:6060/\n")
	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		panic(execErr)
	}
	return nil
})

var _ = Desc("test", "run basic tests (no benchmarks)")
var _ = Add("test", func(c *Context) error {
	binary, lookErr := exec.LookPath("go")
	if lookErr != nil {
		panic(lookErr)
	}
	args := []string{"go", "run", "github.com/onsi/ginkgo/ginkgo", "-r", "-skipMeasurements"}
	env := os.Environ()
	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		panic(execErr)
	}
	return nil
})

var _ = Namespace("test", func() {
	Desc("watch", "run tests within a watch loop")
	Add("watch", func(c *Context) error {
		binary, lookErr := exec.LookPath("go")
		if lookErr != nil {
			panic(lookErr)
		}
		args := []string{"go", "run", "github.com/onsi/ginkgo/ginkgo", "watch", "-r", "-skipMeasurements", "-succinct"}
		env := os.Environ()
		execErr := syscall.Exec(binary, args, env)
		if execErr != nil {
			panic(execErr)
		}
		return nil
	})

	Desc("force", "run all tests even if some fail")
	Add("force", func(c *Context) error {
		binary, lookErr := exec.LookPath("go")
		if lookErr != nil {
			panic(lookErr)
		}
		args := []string{"go", "run", "github.com/onsi/ginkgo/ginkgo", "-r", "-keepGoing", "-skipMeasurements"}
		env := os.Environ()
		execErr := syscall.Exec(binary, args, env)
		if execErr != nil {
			panic(execErr)
		}
		return nil
	})

	Desc("bench", "run performance benchmark tests")
	Add("bench", func(c *Context) error {
		binary, lookErr := exec.LookPath("go")
		if lookErr != nil {
			panic(lookErr)
		}
		args := []string{"go", "run", "github.com/onsi/ginkgo/ginkgo", "-r", "-keepGoing", "-focus", "performance"}
		env := os.Environ()
		execErr := syscall.Exec(binary, args, env)
		if execErr != nil {
			panic(execErr)
		}
		return nil
	})

	Desc("all", "run all tests")
	Add("all", func(c *Context) error {
		binary, lookErr := exec.LookPath("go")
		if lookErr != nil {
			panic(lookErr)
		}
		args := []string{"go", "run", "github.com/onsi/ginkgo/ginkgo", "-r"}
		env := os.Environ()
		execErr := syscall.Exec(binary, args, env)
		if execErr != nil {
			panic(execErr)
		}
		return nil
	})

})
