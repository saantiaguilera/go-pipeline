<p align="center">
    <img width="175" align="center" src="https://github.com/saantiaguilera/go-pipeline/raw/master/logo/logo.png"/><br>
    <br>
    <b>Pipeline</b>
</p>

![Build Status](https://github.com/saantiaguilera/go-pipeline/workflows/Go/badge.svg) 
[![Coverage](https://codecov.io/gh/saantiaguilera/go-pipeline/branch/master/graph/badge.svg)](https://codecov.io/gh/saantiaguilera/go-pipeline)
[![Go Report Card](https://goreportcard.com/badge/github.com/saantiaguilera/go-pipeline)](https://goreportcard.com/report/github.com/saantiaguilera/go-pipeline)
[![GoDoc](https://godoc.org/github.com/saantiaguilera/go-pipeline?status.svg)](https://godoc.org/github.com/saantiaguilera/go-pipeline)
[![Release](https://img.shields.io/github/release/saantiaguilera/go-pipeline.svg?style=flat-square)](https://github.com/saantiaguilera/go-pipeline/releases)

Pipeline is a GPL3-licensed Go package for building, executing and representing pipelines (aka workflows / templates).

## Getting started

- API documentation and examples are available via [godoc](https://godoc.org/github.com/saantiaguilera/go-pipeline).
- The [examples](./examples) directory contains more elaborate example applications.
- No specific mocks are needed for testing, every element is completely decoupled and atomic. You can create your own ones however you deem fit.

## API stability

Pipeline follows semantic versioning and provides API stability via the gopkg.in service.
You can import a version with a guaranteed stable API via http://gopkg.in/saantiaguilera/go-pipeline.v1
