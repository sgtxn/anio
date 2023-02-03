# anio

Nothing to see here yet, move on, stranger.

# Development
## Enabling linter
VSCode golang plugin should load `.golangci-ling.yml` configuratin automatically, just make sure you have:
- `go.lintTool` set to `golangci-lint` 
- `go.lintOnSave` set to `package` preferably

For manual runs execute `make lint`, see the Makefile for further info.
