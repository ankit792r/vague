package client

import (
	"context"
	"fmt"
	"os"
)

func OpenUI(ctx context.Context, files ...string) error {
	conn, err := ClientConnect()
	if err != nil {
		return err
	}

	defer func() {
		// if err := conn.UiDetach(ctx); err != nil {
		// 	slog.Error("ui detach", "error", err)
		// }
		conn.Close()
	}()

	workDir, err := clientWorkDir()
	if err != nil {
		return err
	}

	// TODO: this will be passed to the backend to open and load buffers
	fmt.Println(workDir)

	return LaunchWebView(&ctx)
}

func clientWorkDir() (string, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if workDir == "" || workDir == "/" {
		if home, err := os.UserHomeDir(); err == nil {
			return home, nil
		}
	}

	return workDir, nil
}
