package test

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"testing"

	"github.com/flanksource/canary-checker/api/context"
	"github.com/flanksource/canary-checker/checks"
	dutyContext "github.com/flanksource/duty/context"
	"github.com/flanksource/duty/tests/setup"

	"github.com/flanksource/canary-checker/cmd"
	"github.com/flanksource/canary-checker/pkg"
	"github.com/flanksource/commons/logger"
)

var testFolder string
var DefaultContext dutyContext.Context
var verbosity = 0

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func init() {
	flag.IntVar(&verbosity, "verbose", 0, "Add verbose logging")
	flag.StringVar(&testFolder, "test-folder", "fixtures/minimal", "The folder containing test fixtures to run")
}

func TestRunChecks(t *testing.T) {
	kommonsClient, k8s, err := pkg.NewKommonsClient()
	if err != nil {
		logger.Warnf("Failed to get kommons client, features that read kubernetes configs will fail: %v", err)
	}

	DefaultContext = setup.BeforeSuiteFn().WithDBLogLevel("trace").WithTrace().WithKubernetes(k8s).WithKommons(kommonsClient)

	logger.StandardLogger().SetLogLevel(verbosity)
	logger.Infof("Testing %s", testFolder)
	files, _ := os.ReadDir(fmt.Sprintf("../%s", testFolder))
	t.Logf("Folder: %s", testFolder)
	wg := sync.WaitGroup{}
	for _, fixture := range files {
		name := path.Base(fixture.Name())
		if strings.HasPrefix(name, "_") || !strings.HasSuffix(name, ".yaml") || name == "kustomization.yaml" {
			continue
		}
		wg.Add(1)
		go func() {
			runFixture(t, name)
			wg.Done()
		}()
	}
	wg.Wait()
}

func runFixture(t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		canaries, err := pkg.ParseConfig(fmt.Sprintf("../%s/%s", testFolder, name), "")
		if err != nil {
			t.Error(err)
			return
		}

		for _, canary := range canaries {
			if canary.Namespace == "" {
				canary.Namespace = "default"
			}
			if canary.Name == "" {
				canary.Name = cmd.CleanupFilename(name)
			}
			context := context.New(DefaultContext, canary)

			checkResults, err := checks.RunChecks(context)
			if err != nil {
				t.Error(err)
				return
			}

			for _, res := range checkResults {
				if res == nil {
					t.Errorf("Result in %v returned nil:\n", name)
				} else {
					if strings.Contains(name, "_mix") {
						t.Logf("%v: %v", name, res.String())
					} else if strings.Contains(name, "fail") && res.Pass {
						t.Errorf("Expected test to fail, but it passed: %s", res)
					} else if !strings.Contains(name, "fail") && !res.Pass {
						t.Errorf("Expected test to pass but it failed %s", res)
					} else {
						t.Logf("%v: %v", name, res.String())
					}
				}
			}
		}
	})
}
