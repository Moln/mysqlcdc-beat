# mysql cdcbeat

Welcome to mysql cdcbeat.

Ensure that this folder is at the following location:
`${GOPATH}/src/github.com/moln/cdcbeat`

## Getting Started with mysql cdcbeat

### Requirements

* [Golang](https://golang.org/dl/) 1.7

### Test

To test mysql cdcbeat, run the following command:

```
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`


## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make release
```

This will fetch and create all images required for the build process. The whole process to finish can take several minutes.
