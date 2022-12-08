# Podinfo example

Build the component:

```shell
gitopsx component build -f componentfile.cue github.com/acme/podinfo:v1.0.0
```

Push to registry:

```shell
gitopsx component push github.com/acme/podinfo:v1.0.0 ghcr.io/${GITHUB_USER}
```
