package main

import (
	"log/slog"
	"os"

	"github.com/jamwujustyle/logger"
)

func main() {
	logger.InitLogger(1 == 0)

	slog.Error("Hello")

	os.Stdout.Write([]byte("hello world lets put the theory to test"))
}
