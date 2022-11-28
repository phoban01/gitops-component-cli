package main

type Component struct {
	Build  Build  `cmd:"" help:"Build"`
	Push   Push   `cmd:"" help:"Push a component to a registry"`
	Render Render `cmd:"" help:"Render a component from a registry"`
}