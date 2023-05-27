# anio

Nothing to see here yet, move on, stranger.

# Development
## Enabling linter
VSCode golang plugin should load `.golangci-ling.yml` configuration from the repository root automatically, just make sure you have:
- `go.lintTool` set to `golangci-lint` 
- `go.lintOnSave` set to `package` preferably

`make lint` will execute the linter manually in terminal.  
`make lint_fix` will attempt to fix some of the issues.  
See the Makefile for further info.  
