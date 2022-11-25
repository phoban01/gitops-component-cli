## gitops component CLI

** A tool to manage components for GitOps **

## Getting started

```
# make
go build -o ./bin ./cmd/gitops

# build
./bin/gitops component build github.com/acme/mycomponent:v1.0.0

# push
./bin/gitops component push github.com/acme/mycomponent:v1.0.0 ghcr.io/$GITHUB_USER

```
