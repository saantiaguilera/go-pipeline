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

## Creating a step

A step is a contract that allows us to represent or run a unit of work. A step requires an input and may return an output or an error depending on wether it failed or not. A step is declared as `pipeline.Step[Input, Output]`

Steps are considered the backbone of the API. The API already provides a set of steps that should suffice to create any type of pipeline, but there may be specific scenarios were the given API gets too verbose or its not enough. In these type of scenarios we can create our own custom steps to match our needs.

The steps provided by the API are:

### UnitStep

The most simple and atomic step. This step lets us run a single unit of work
```go
var step pipeline.Step[InputData, OutputData] = pipeline.NewUnitStep[InputData, OutputData](
    "name_of_the_step", 
    func(ctx context.Context, in InputData) (OutputData, error) {
        // do stuff with the InputData, returning Outputdata or error
    },
)
```

### SequentialStep

A sequential step allows us to "link" two steps together sequentially.

```go
var firstStep pipeline.Step[int, string]
var secondStep pipeline.Step[string, bool]

// in:  int
// out: bool
var sequentialStep pipeline.Step[int, bool] = pipeline.NewSequentialStep[int, string, bool](firstStep, secondStep)
```

### ConcurrentStep

A concurrent step allows us to "link" multiple steps concurrently and once they're done reduce them to a single output

```go
var concurrentSteps []pipeline.Step[int, string]
var reducer func(context.Context, a, b string) (string, error)

// in: int
// out: string
var concurrentStep pipeline.Step[int, string] = pipeline.NewConcurrentStep[int, string](concurrentSteps, reducer)
```

### ConditionalStep

A conditional step allows us to evaluate a condition and depending on its result branch to specific step.
This step allows us to branch the graph in two different branches.

```go
var trueWayStep pipeline.Step[InputData, OutputData]
var falseWayStep pipeline.Step[InputData, OutputData]

var statement pipeline.Statement[InputData] = pipeline.NewStatement(
    "name_of_the_statement",
    func(ctx context.Context, in InputData) bool {
        // evaluate statement and return branching mode
    }
)
var cond pipeline.Step[InputData, OutputData] = pipeline.NewConditionalStep(statement, trueWayStep, falseWayStep)
```

### OptionalStep

An optional step is similar to a conditional one, although it only has a single branch. It either runs the given Step or it skips it (returning the initial input), depending on the result of the statement evaluation.

```go
var step pipeline.Step[InputData, InputData]

var statement pipeline.Statement[InputData] = pipeline.NewStatement(
    "name_of_the_statement",
    func(ctx context.Context, in InputData) bool {
        // evaluate statement and return true to run / false to skip
    }
)
var opt pipeline.Step[InputData, InputData] = pipeline.NewOptionalStep(statement, step)
```

It also supports altering the output, but when doing so you need to provide how to default to it when the step is skipped

```go
var step pipeline.Step[InputData, OutputData]

var statement pipeline.Statement[InputData] = pipeline.NewStatement(
    "name_of_the_statement",
    func(ctx context.Context, in InputData) bool {
        // evaluate statement and return true to run / false to skip
    }
)
var def pipeline.Unit[InputData, OutputData] = func(ctx context.Context, in InputData) (OutputData, error) {
    // create default output data for when the step is skipped because the statement evaluation was false
}
var opt pipeline.Step[InputData, OutputData] = pipeline.NewOptionalStepWithDefault(statement, step, def)
```

## Creating a custom step

Steps need to comply to an extremely simple interface
```go
type Step[I, O any] interface {
    Draw(pipeline.Graph) // lets us represent a step in a graph
	Run(context.Context, I) (O, error) // lets us evaluate the step
}
```

Hence, we can create our own custom steps by simply creating a struct that matches the given contract. There are no restrictions besides these two so it's highly flexible when wanting to create custom behaviors or logics.

For example, a step that always succeeds and doesn't mutate the result might be:
```go
type ImmutableStepThatAlwaysSucceeds[I any] struct {
    name string
    fn   func(ctx context.Context, in I)
}

func (s ImmutableStepThatAlwaysSucceeds[I]) Draw(g pipeline.Graph) {
    g.AddActivity(s.name)
}

func (s ImmutableStepThatAlwaysSucceeds[I]) Run(ctx context.Context, in I) (I, error) {
    s.fn(ctx, in)
    return in, nil
}

func main() {
    var s pipeline.Step[int, int] = ImmutableStepThatAlwaysSucceeds[int]{
        name: "example",
        fn: func(ctx context.Context, in int) {
            // do something.
        }
    }
}
```