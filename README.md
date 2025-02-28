Command gencodec generates marshaling methods for Go struct types.

The generated methods add features which json and other marshaling packages cannot offer.

	gencodec -dir . -type MyType -formats json,yaml,toml -out mytype_json.go

See [the documentation for more details](https://godoc.org/github.com/fjl/gencodec).

To pin a specific version of gencodec in your project, you can use, for example for version v0.1.0:

- For Go >= 1.24, use `go get -tool github.com/fjl/gencodec@v0.1.0`
- For Go <= 1.23:
  1. Add the tool dependency with `go get github.com/fjl/gencodec@v0.1.0`
  1. Create a `tools.go` file at the root of your project with content:

    ```go
    package yourproject

    import (
        _ "github.com/fjl/gencodec/importable"
    )
    ```

  1. Run `go mod tidy`
