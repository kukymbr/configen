package main

import (
	"os"

	"github.com/kukymbr/configen/internal/command"
	"github.com/kukymbr/configen/internal/logger"
)

func main() {
	if err := command.Run(); err != nil {
		logger.Errorf("%s", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
