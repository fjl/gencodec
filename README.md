# About gencodec

Command gencodec generates marshaling methods for Go struct types.

The generated methods add features which json and other marshaling packages cannot offer.

## Build

```bash
git clone https://github.com/fjl/gencodec.git
cd gencodec
go build
./gencodec -h

# copy gencodec to any directory which is in PATH, such as: ${HOME}/go/bin
mkdir -p ${HOME}/go/bin
cp gencodec ${HOME}/go/bin
PATH=${PATH}:${HOME}/go/bin
gencodec -h
```

## Usage

Usage of gencodec:

```text
  -dir string
        input package (default ".")
  -field-override string
        type to take field type replacements from
  -formats string
        marshaling formats (e.g. "json,yaml") (default "json")
  -out string
        output file (default is stdout) (default "-")
  -type string
        type to generate methods for
```

Example:

```bash
gencodec -dir . -type MyType -formats json,yaml,toml -out mytype_json.go
```

See [the documentation for more details](https://godoc.org/github.com/fjl/gencodec).
