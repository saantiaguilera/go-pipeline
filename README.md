# Go-Pipeline

![Build](https://github.com/saantiaguilera/go-pipeline/workflows/Go/badge.svg) [![Coverage](https://codecov.io/gh/saantiaguilera/go-pipeline/branch/master/graph/badge.svg)](https://codecov.io/gh/saantiaguilera/go-pipeline)

Go module for building, executing and representing pipelines. For more information visit the [GoDoc API](https://godoc.org/github.com/saantiaguilera/go-pipeline)

### Example

_The following code and representation can be found under the [examples directory](examples/cook_example/) if you want to play with it._

Imagine we are 3 persons making a dish. We have to:
1. Put the eggs to boil. When done, cut them.
2. Wash the carrots. Cut them.
3. Start the oven. If the meat is too big, cut it. Put the meat in the oven.
4. Make a salad with the cut eggs and carrots
5. Serve

The following can be represented as such (using the `pipeline.DrawableDiagram` API)

![](examples/cook_example/template.svg)

This flow can be built and executed as such:
```go
// Complete stage. Its sequential because we can't serve
// before all the others are done. 
graph := pipeline.CreateSequentialGroup(
    // Concurrent stage, given we are 3, we can do the salad / meat separately
    pipeline.CreateConcurrentGroup(
        // This will be the salad flow. It can be done concurrently with the meat
        pipeline.CreateSequentialGroup( 
            // Eggs and carrots can be operated concurrently too
            pipeline.CreateConcurrentGroup(
                // Sequential stage for the eggs flow
                pipeline.CreateSequentialStage(
                    CreateBoilEggsStep(),
                    CreateCutEggsStep(),
                ),
                // Another sequential stage for the carrots (eggs and carrots will be concurrent though!)
                pipeline.CreateSequentialStage(
                    CreateWashCarrotsStep(),
                    CreateCutCarrotsStep(),
                ),
            ),
            // This is sequential. When carrots and eggs are done, this will run
            pipeline.CreateSequentialStage(
                CreateMakeSaladStep(),
            ),
        ),
        // Another sequential stage for the meat (concurrently with salad)
        pipeline.CreateSequentialGroup(
            // If we end up cutting the meat, we can optimize it with the oven operation
            pipeline.CreateConcurrentGroup(
                // Conditional stage, the meat might be too big
                pipeline.CreateConditionalStage(
                    pipeline.CreateSimpleStatement("is_meat_too_big", IsMeatTooBigForTheOven),
                    // True:
                    CreateCutMeatStep(),
                    // False:
                    nil,
                ),
                pipeline.CreateSequentialStage(
                    CreateTurnOvenOnStep(),
                ),
            ),
            pipeline.CreateSequentialStage(
                CreatePutMeatInOvenStep(),
            ),
        ),
    ),
    // When everything is done. Serve
    pipeline.CreateSequentialStage(
        CreateServeStep(),
    ),
)

pipe := pipeline.CreatePipeline(CreateYourExecutor())
pipe.Run(graph, CreateYourInputContext())
```
