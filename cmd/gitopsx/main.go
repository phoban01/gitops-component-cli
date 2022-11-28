package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

type app struct {
	Component Component `cmd:"" help:""`
}

func main() {
	ctx := kong.Parse(new(app), kong.UsageOnError())
	if err := ctx.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
