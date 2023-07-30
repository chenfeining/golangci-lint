`golangci-lint` is a fast Go linters runner. It runs linters in parallel, uses caching, supports `yaml` config, has integrations
with all major IDE and has dozens of linters included.

## important
- This version has a powerful ability to detect hidden nil pointers reference by enable `npecheck` linter
- Such as `golangci-lint run -E npecheck ./...` 
- Or configure it in `.golangci.yml`, similar to enable other linters, can be used with them

## Install `golangci-lint`
- Using go version >= go 1.19
- go install github.com/chenfeining/golangci-lint/cmd/golangci-lint@v1.0.0

## Documentation
- Documentation is hosted at https://golangci-lint.run.
- npecheck single repo with test case: https://github.com/chenfeining/go-npecheck
