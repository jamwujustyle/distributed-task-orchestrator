package main

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
)

func runREPL(ctx context.Context, c pb.TaskServiceClient) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("interactive mode. ommands: upload <path>, list, submit, <key>, retrieve <id>, exit")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		parts := strings.Fields(scanner.Text())
		if len(parts) == 0 {
			continue
		}
		switch parts[0] {
		case "upload":
			data, err := os.ReadFile(os.Args[2])
			if err != nil {
				slog.Error("failed to read file", "err", err)
				continue
			}
			key, err := doUploadScript(ctx, c, data)
			if err != nil {
				slog.Error("failed 'doUploadScript' job", "err", err)
				//i wonder if the os.exit radicality justified, perhaps continue is a better fit?
				continue
			}
			slog.Info("use this key to to submit a task", "key", key)
		case "list":
			doListScripts(ctx, c)
		case "submit":
			doSubmitTask(ctx, c, os.Args[2])
		case "retrieve":
			doRetrieveTask(ctx, c, os.Args[2])
		default:
			fmt.Println("unknown command")
		}

	}
}
