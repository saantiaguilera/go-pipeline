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
You can import a version with a guaranteed stable API via http://gopkg.in/saantiaguilera/go-pipeline.v0

## Example

_The following graph creation, execution and representation can be found under the [examples](examples/usages/cooking_a_recipe_pipeline) directory._

Imagine we are making a dish, we need to:
1. Put the eggs to boil and cut them.
2. Wash the carrots and cut them.
3. Make a salad with the cut eggs and carrots.
4. Start the oven. 
5. If the meat is too big, cut it. 
6. Put the meat in the oven.
7. Serve when the meat and the salad are done.

This workflow is represented as such (with this same API, no need to draw it on your own)

![](examples/usages/cooking_a_recipe_pipeline/template.svg)

To build this, we simply need to create a step / unit of work for each given task and then "link" them however we want them to be traversed later in the graph. The graph creation can be seen [here](https://github.com/saantiaguilera/go-pipeline/blob/master/examples/usages/cooking_a_recipe_pipeline/main.go#L18)
