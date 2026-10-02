package cmd

import (
	"fmt"
	"os"
)

func Start() error {
	args := os.Args[1:]
	if len(args) == 0 {
		return serve()
	}

	switch args[0] {
	case "serve":
		return serve()
	case "request":
		return request(args[1:])
	default:
		return fmt.Errorf("unknown command %q, available commands: serve, request", args[0])
	}
}
