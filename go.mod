module github.com/citihub/probr-core

go 1.14

require (
	github.com/citihub/probr-sdk v0.0.11
	github.com/hashicorp/go-hclog v0.14.1
	github.com/hashicorp/go-plugin v1.4.0
)

// replace github.com/citihub/probr-sdk => ../probr-sdk

//Line above is intended to be used during dev only when editing modules locally.
