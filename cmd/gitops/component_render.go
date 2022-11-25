package main

type Render struct {
	Repository string `arg:""`
}

func (p *Render) Help() string {
	return "This is the full help text for render"
}

func (p *Render) Validate() error {
	return nil
}

func (p *Render) Run() error {
	return nil
}
