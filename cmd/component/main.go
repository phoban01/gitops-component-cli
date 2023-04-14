package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

type app struct {
	Build    Build    `cmd:"" help:"Build"`
	Push     Push     `cmd:"" help:"Push a component to a registry"`
	Render   Render   `cmd:"" help:"Render a component from a registry"`
	Sign     Sign     `cmd:""`
	Verify   Verify   `cmd:""`
	Describe Describe `cmd:"" help:"Show the component descriptor manifest for a component"`
	Get      Get      `cmd:""`
}

func main() {
	ctx := kong.Parse(new(app), kong.UsageOnError())
	if err := ctx.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
