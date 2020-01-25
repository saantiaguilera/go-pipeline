# Go-Pipeline

Toy module for building decoupled pipelines in Golang

### Supported structure

- Stage: Collection of steps
- Stage group: Also a stage. Represents a collection of stages
- Step: Single unit of work

A pipeline contains a single `Stage`, from which a graph/template can be formed and be used as a template for running a workload.

There are different type of stages, eg. 

- `ConcurrentStage` which runs steps concurrently
- `SequentialStage` which runs them sequentially.
- `ConditionalStage` which given a statement, runs one stage or another.

### Executor

A pipeline can have an executor, which will define the flow for executing a Step. This is useful if we want to add hooks / circuit-breakers / tracers / etc to the pipeline.

### Example

Imagine we are 3 persons making a dish. We have to:
1. Put the eggs to boil. When done, cut them.
2. Wash the carrots. Cut them.
3. Start the oven. If the meat is too big, cut it. Put the meat in the oven.
4. Make a salad with the cut eggs and carrots
5. Serve

You might realize there are a lot of things that don't depend on each other. Eg. you can do the salad and meat separately.

This flow can be achieved as such:
```go
// Complete stage. Its sequential because we can't serve
// before all the others are done. 
graph := pipeline_stage.CreateSequentialGroup(
    // Concurrent stage, given we are 3, we can do the salad / meat separately
    pipeline_stage.CreateConcurrentGroup(
        // This will be the salad flow. It can be done concurrently with the meat
        pipeline_stage.CreateSequentialGroup(
        	// Eggs and carrots can be operated concurrently too
            pipeline_stage.CreateConcurrentGroup(
                // Sequential stage for the eggs flow
                pipeline_stage.CreateSequentialStage(
                    // Use a mean of communication. Channels could be one.
                    your_step.CreateBoilEggsStep(eggsChan),
                    your_step.CreateCutEggsStep(eggsChan),
                ),
                // Another sequential stage for the carrots (eggs and carrots will be concurrent though!)
                pipeline_stage.CreateSequentialStage(
                    // Use a mean of communication. Channels could be one.
                    your_step.CreateWashCarrotsStep(carrotsChan),
                    your_step.CreateCutCarrotsStep(carrotsChan),
                ),
            ),
            // This is sequential. When carrots and eggs are done, this will run
            pipeline_stage.CreateSequentialStage(
                your_step.MakeSaladStep(carrotsChan, eggsChan, saladChan)
            )
        ),
        // Another sequential stage for the meat (concurrently with salad)
        pipeline_stage.CreateSequentialStage(
            your_step.TurnOvenOnStep(),
            // Conditional step, the meat might be too big
            pipeline_stage.CreateConditionalStep(
                func() bool {
            	    return isMeatTooBigForOven()
            	},
                // True:
                your_step.CutMeat(meatChan),
                // False:
                nil,
            ),
            your_step.PutMeatInOvenStep(meatChan),
        ),
    ),
    // When everything is done. Serve
    pipeline_stage.CreateSequentialStage(
        your_step.ServeStep(meatChan, saladChan)
    ),
)

pipeline := CreatePipeline(CreateYourExecutor())
pipeline.Run(graph)
```