package main

type Component struct {
	Build Build `cmd:"" help:"Build"`
	Push  Push  `cmd:"" help:"Push a component to a registry"`
}
