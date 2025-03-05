# liberator [![Go Reference](https://pkg.go.dev/badge/github.com/nais/liberator.svg)](https://pkg.go.dev/github.com/nais/liberator)

Collection of small Go libraries used by Nais operators.

## Contents

* Kubernetes _Custom Resource Definitions_
* Name generation functions
* ...more to come

## Usage

Install dependency:
```
go get github.com/nais/liberator
```

To update to the latest version, run:
```
go get -u github.com/nais/liberator@HEAD
```

## Developing

### Using a locally checked out copy of Liberator

When developing applications using new Liberator features, you can point the application to link
against your local Liberator by running the following command in the application repository:

```
go mod edit -replace github.com/nais/liberator=../liberator
```

This will add a line in your `go.mod` file. Make sure this change isn't checked in when committing,
otherwise the code won't build in the CI tool. To revert your changes, run:

```
go mod edit -dropreplace github.com/nais/liberator
go get -u github.com/nais/liberator@HEAD
```

### Immutable parts of spec

You can define fields as immutable using the field tag `nais:"immutable"`:
```go
type MySpec struct {
	MutableString string `json:"mutableString"`
	ImmutableString string `json:"immutableString" nais:"immutable"`
}
```

You can also use this on slices. When defined on a field with type slice, the entire slice is immutable and cannot be changed. If the slice contains structs, individual elements of the slice will be tested against each other. This requires another field flag, `nais:"key"`, to identify which of the elements in the slice to match against each other. You can specify more than one in a "one of" matching process.

```go
type Apps struct {
	Name string `json:"name" nais:"immutable,key"`
	Awesomeness int `json:"awesomeness" nais:"immutable"`
}

type AppList struct {
	Items []Apps `json:"name"`
}
```

#### Documentation

Documenting immutable fields must be done using the comment `// +nais:doc:Immutable=true`:

```go
type Apps struct {
	// +nais:doc:Immutable=true
	Name string `json:"name" nais:"immutable,key"`
	// +nais:doc:Immutable=true
	Awesomeness int `json:"awesomeness" nais:"immutable"`
}

type AppList struct {
	Items []Apps `json:"name"`
}
```

### Kubernetes dependencies

The `controller-tools` dependency in `go.mod` locks the versions of
upstream libraries `api`, `apiextensions`, and `apimachinery`,
effectively determining which Kubernetes version is compatible.

Check for compatibility here: https://github.com/kubernetes-sigs/controller-tools/commits/master/go.mod

| Kubernetes version | controller-gen version |
|--------------------|------------------------|
| 1.17               | v0.2.5                 |
| 1.18               | v0.4.1                 |
| 1.19               | ??                     |
| 1.20               | ??                     |
| 1.21               | v0.6.2                 |

### Code generation

Make sure `controller-gen` is of a compatible version by doing `go get sigs.k8s.io/controller-tools@VERSION`

Run `make generate` to generate deep copy functions and CRD files.
