## gitops component CLI

** A tool to manage components for GitOps **

## Getting started

```
# make
go install./cmd/gitopsx

# build
gitopsx component build github.com/acme/mycomponent:v1.0.0

# push
gitopsx component push github.com/acme/mycomponent:v1.0.0 ghcr.io/$GITHUB_USER

# render (currently very slow...)
gitopsx component render -f Application.cue -oyaml
```

