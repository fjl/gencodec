Command gencodec generates marshaling methods for struct types.

When gencodec is invoked on a directory and type name, it creates a Go source file
containing JSON and YAML marshaling methods for the type. The generated methods add
features which the standard json package cannot offer.

	gencodec -dir . -type MyType -out mytype_json.go

### Struct Tags

All fields are required unless the "optional" struct tag is present. The generated
unmarshaling method return an error if a required field is missing. Other struct tags are
carried over as is. The standard "json" and "yaml" tags can be used to rename a field when
marshaling to/from JSON.

Example:

	type foo {
		Required string
		Optional string `optional:""`
		Renamed  string `json:"otherName"`
	}

### Field Type Overrides

An invocation of gencodec can specify an additional 'field override' struct from which
marshaling type replacements are taken. If the override struct contains a field whose name
matches the original type, the generated marshaling methods will use the overridden type
and convert to and from the original field type.

In this example, the specialString type implements json.Unmarshaler to enforce additional
parsing rules. When json.Unmarshal is used with type foo, the specialString unmarshaler
will be used to parse the value of SpecialField.

	//go:generate gencodec -dir . -type foo -field-override fooMarshaling -out foo_json.go

	type foo struct {
		Field        string
		SpecialField string
	}

	type fooMarshaling struct {
		SpecialField specialString // overrides type of SpecialField when marshaling/unmarshaling
	}
