# liberator [![Go Reference](https://pkg.go.dev/badge/github.com/nais/liberator.svg)](https://pkg.go.dev/github.com/nais/liberator)

Collection of small Go libraries used by NAIS operators.

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

The library is semantically versioned. Minor and patch level releases contain
only new features and bugfixes, respectively. API breaking changes occurs only
in major version releases.

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

### Kubernetes dependencies

The `controller-tools` dependency in `go.mod` locks the versions of
upstream libraries `api`, `apiextensions`, and `apimachinery`,
effectively determining which Kubernetes version is compatible.

Check for compatibility here: https://github.com/kubernetes-sigs/controller-tools/commits/master/go.mod

| Kubernetes version | controller-gen version |
|--------------------|------------------------|
| 1.17               | v0.2.5                 |
| 1.18               | v0.4.1                 |
| 1.19               | master                 |

### Code generation

Make sure `controller-gen` is of a compatible version by modifying `Makefile` and running `make controller-gen`.
 
Run `make generate` to generate deep copy functions and CRD files.
