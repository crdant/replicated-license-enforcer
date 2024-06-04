package main

import (
	"flag"
	"os"
  "os/signal"
	"time"

	"github.com/charmbracelet/log"
	"github.com/crdant/replicated-license-enforcer/pkg/enforce"
	"github.com/crdant/replicated-license-enforcer/pkg/version"
)

var (
	logLevel       string
	recheckInterval time.Duration
)

func init() {
	logLevel = os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	log.ParseLevel(logLevel)
}

func parseFlags() {
	flag.DurationVar(&recheckInterval, "recheck", 0, "Recheck license periodically to assure it's still valid")
	flag.Parse()
}

func main() {
	parseFlags()

	log.Infof("Version: %s, Build Time: %s, GitCommit: %s\n", version.Version, version.BuildTime, version.GitSHA)

  enforcer := enforce.DefaultEnforcer()
	err := enforcer.Validate()
	if err != nil {
		log.Error("Error checking license validity", "error", err)
		os.Exit(1)
	}

	if recheckInterval > 0 {
		enforcer.Monitor(recheckInterval)
	}

	waitForSignal()
}

func waitForSignal() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
	os.Exit(0)
}
