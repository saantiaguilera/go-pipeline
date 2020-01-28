# Go-Pipeline

Go module for building pipelines. For more information visit the [GoDoc API](https://godoc.org/github.com/saantiaguilera/go-pipeline)

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
graph := stage.CreateSequentialGroup(
    // Concurrent stage, given we are 3, we can do the salad / meat separately
    stage.CreateConcurrentGroup(
        // This will be the salad flow. It can be done concurrently with the meat
        stage.CreateSequentialGroup( 
            // Eggs and carrots can be operated concurrently too
            stage.CreateConcurrentGroup(
                // Sequential stage for the eggs flow
                stage.CreateSequentialStage(
                    // Use a mean of communication. Channels could be one.
                    your_step.CreateBoilEggsStep(eggsChan),
                    your_step.CreateCutEggsStep(eggsChan),
                ),
                // Another sequential stage for the carrots (eggs and carrots will be concurrent though!)
                stage.CreateSequentialStage(
                    // Use a mean of communication. Channels could be one.
                    your_step.CreateWashCarrotsStep(carrotsChan),
                    your_step.CreateCutCarrotsStep(carrotsChan),
                ),
            ),
            // This is sequential. When carrots and eggs are done, this will run
            stage.CreateSequentialStage(
                your_step.MakeSaladStep(carrotsChan, eggsChan, saladChan)
            )
        ),
        // Another sequential stage for the meat (concurrently with salad)
        stage.CreateSequentialGroup(
            // If we end up cutting the meat, we can optimize it with the oven operation
            stage.CreateConcurrentGroup(
                // Conditional stage, the meat might be too big
                stage.CreateConditionalStage(
                    func() bool {
                        return isMeatTooBigForOven()
                    },
                    // True:
                    your_step.CutMeatStep(meatChan),
                    // False:
                    nil,
                ),
                stage.CreateSequentialStage(
                    your_step.TurnOvenOnStep(),
                ),
            ),
            stage.CreateSequentialStage(
                your_step.PutMeatInOvenStep(meatChan),
            ),
        ),
    ),
    // When everything is done. Serve
    stage.CreateSequentialStage(
        your_step.ServeStep(meatChan, saladChan)
    ),
)

pipe := pipeline.CreatePipeline(CreateYourExecutor())
pipe.Run(graph)
```