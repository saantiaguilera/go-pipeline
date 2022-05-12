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
- No specific mocks are needed for testing, every element is completely decoupled and atomic. You can New your own ones however you deem fit.

## API stability

Pipeline follows semantic versioning and provides API stability via the gopkg.in service.
You can import a version with a guaranteed stable API via http://gopkg.in/saantiaguilera/go-pipeline.v0

## Example

_The following graph creation, execution and representation can be found under the [examples](examples/static/cook_example/) directory._

Imagine we are making a dish, we need to:
1. Put the eggs to boil and cut them.
2. Wash the carrots and cut them.
3. Make a salad with the cut eggs and carrots.
4. Start the oven. 
5. If the meat is too big, cut it. 
6. Put the meat in the oven.
7. Serve when the meat and the salad are done.

This workflow is represented as such

![](examples/static/cook_example/template.svg)

This workflow can be built and executed as such.
```go
// Complete stage. Its sequential because we can't serve
// before all the others are done. 
graph := pipeline.NewSequentialGroup(
    // Concurrent stage, given we can do the salad / meat separately.
    pipeline.NewConcurrentGroup(
        // This will be the salad flow.
        pipeline.NewSequentialGroup( 
            // Eggs and carrots can be operated concurrently too.
            pipeline.NewConcurrentGroup(
                // Sequential stage for the eggs flow.
                pipeline.NewSequentialStage(
                    NewBoilEggsStep(),
                    NewCutEggsStep(),
                ),
                // Another sequential stage for the carrots (eggs and carrots will be concurrent though!)
                pipeline.NewSequentialStage(
                    NewWashCarrotsStep(),
                    NewCutCarrotsStep(),
                ),
            ),
            // This is sequential. When carrots and eggs are done, this will run.
            pipeline.NewSequentialStage(
                NewMakeSaladStep(),
            ),
        ),
        // Another sequential stage for the meat (concurrently with salad)
        pipeline.NewSequentialGroup(
            // If we end up cutting the meat, we can optimize it with the oven operation
            pipeline.NewConcurrentGroup(
                // Conditional stage, the meat might be too big
                pipeline.NewConditionalStage(
                    pipeline.NewStatement("is_meat_too_big", IsMeatTooBigForTheOven),
                    // True:
                    NewCutMeatStep(),
                    // False:
                    nil,
                ),
                pipeline.NewSequentialStage(
                    NewTurnOvenOnStep(),
                ),
            ),
            pipeline.NewSequentialStage(
                NewPutMeatInOvenStep(),
            ),
        ),
    ),
    // When everything is done. Serve.
    pipeline.NewSequentialStage(
        NewServeStep(),
    ),
)

pipe := pipeline.NewPipeline(NewYourExecutor())
pipe.Run(graph, NewYourInputContext())
```
_Note that, for showing purposes, this is all in a single function. You can easily decouple this into more atomic ones that take care of specific responsibilities (eg. making the salad)._
