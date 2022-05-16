package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/saantiaguilera/go-pipeline"
)

var render = flag.Bool("pipeline.render", false, "render pipeline")

// Graph creates static workflow for this sample. It's all in a single func completely coupled for showing purposes
// you should probably decouple this into more atomic ones (eg. a func for making the salad that returns that complete subgraph)
// This is done in a start-to-end approach. It could also be made the other way around (starting by serving the dish and going
// upwards until the start). This choice was simply done for readability purposes
func Graph() pipeline.Step[MealMaterials, Dish] {
	ovenSize := 500 // constant. We put it here to avoid having visibility on the sample package

	// Create subgraph for processing the eggs
	processEggs := pipeline.NewSequentialStep[[]Egg, []Egg, []CutEgg](
		newBoilEggsStep(),
		newCutEggsStep(5),
	)

	// Create subgraph for processing the carrots
	processCarrots := pipeline.NewSequentialStep[[]Carrot, []Carrot, []CutCarrot](
		newWashCarrotsStep(),
		pipeline.NewUnitStep("cut_carrots_step", newCarrotsCutter().Cut),
	)

	// Create subgraph for handling meat and oven preparations
	meatPreparations := pipeline.NewOptionalStep[Meat](
		pipeline.NewStatement("is_meat_too_big", func(ctx context.Context, t Meat) bool {
			return t.Size > ovenSize
		}),
		cutMeatCustomStep{
			OvenSize: ovenSize,
		},
	)

	// Concurrency cases to optimize most of the stuff above (as we can do them at the same time some of them)
	// We will showcase here a custom step of our own and the pipeline API one so we can see the power of creating our custom ones
	// to satisfy specific needs.
	processVegetables := processVegetablesConcurrently{
		EggStep:    processEggs,
		CarrotStep: processCarrots,
	}
	// With the pipeline API we need to wrap the concurrency steps in a get/put that will decouple the inner steps
	// from the real normalized input and output.
	// This is because we need in a concurrency step to have always the same input/output between steps (else the pipeline API can't know
	// which ones to provide / expect). Of course this could be avoided if inner steps had the normalized in/out already applied
	// (eg in this case both steps being Step[MealMaterials, CookingTools]), but that would couple them to this specific graph capabilities.
	processMeat := pipeline.NewConcurrentStep(
		[]pipeline.Step[MealMaterials, CookingTools]{
			pipeline.NewSequentialStep[MealMaterials, Meat, CookingTools](
				pipeline.NewUnitStep("get_meat", func(_ context.Context, in MealMaterials) (Meat, error) {
					return in.Meat, nil
				}),
				pipeline.NewSequentialStep[Meat, Meat, CookingTools](
					meatPreparations,
					pipeline.NewUnitStep("put_meat", func(_ context.Context, in Meat) (CookingTools, error) {
						return CookingTools{
							Meat: in,
						}, nil
					}),
				),
			),
			pipeline.NewSequentialStep[MealMaterials, Oven, CookingTools](
				pipeline.NewUnitStep("get_oven", func(_ context.Context, in MealMaterials) (Oven, error) {
					return in.Oven, nil
				}),
				pipeline.NewSequentialStep[Oven, Oven, CookingTools](
					newTurnOnOvenStep(),
					pipeline.NewUnitStep("put_oven", func(_ context.Context, in Oven) (CookingTools, error) {
						return CookingTools{
							Oven: in,
						}, nil
					}),
				),
			),
		},
		func(_ context.Context, a, b CookingTools) (CookingTools, error) {
			if b.Meat != (Meat{}) {
				a.Meat = b.Meat
			}
			if b.Oven != (Oven{}) {
				a.Oven = b.Oven
			}
			return a, nil
		},
	)

	// Create subgraph for making salad and cooking meat
	makeSalad := pipeline.NewSequentialStep[MealMaterials, Vegetables, Salad](
		processVegetables,
		newMakeSaladStep(),
	)
	cookMeat := pipeline.NewSequentialStep[MealMaterials, CookingTools, CookedMeat](
		processMeat,
		newCookMeatStep(),
	)

	// create complete graph that serves the meal.
	return pipeline.NewSequentialStep[MealMaterials, DishContents, Dish](
		processDishConcurrently{
			SaladStep: makeSalad,
			MeatStep:  cookMeat,
		},
		newServeStep(),
	)
}

// RunGraphRendering represents the graph in UML Activity and renders it as an SVG file (template.svg)
func RunGraphRendering() {
	if *render {
		diagram := pipeline.NewUMLGraph()
		renderer := pipeline.NewUMLRenderer(pipeline.UMLOptions{
			Type: pipeline.UMLFormatSVG,
		})
		file, _ := os.Create("template.svg")

		Graph().Draw(diagram)

		err := renderer.Render(diagram, file)

		if err != nil {
			panic(err)
		}
	}
}

// RunPipeline runs the provided pipeline.
// Output: (one of many)
//
// Turning oven on
// Washing 8 carrots
// Cutting meat of size 600 into 500
// Boiling 5 eggs
// Cutting 25 eggs
// Cooking meat
// Cutting 8 carrots into 40 pieces
// Making salad with 25 eggs and 40 carrots
// Serving dish with a salad of 25 eggs and 40 carrots (mixed: true) and a cooked meat
func RunPipeline() {
	// Create a stateless graph, so it can be evaluated as many times as we like with any input/context we want to.
	graph := Graph()

	// Initial input data
	ctx := context.Background()
	data := MealMaterials{
		Eggs:    make([]Egg, 5),
		Carrots: make([]Carrot, 8),
		Meat: Meat{
			Size: 600,
		},
		Oven: Oven{
			Ignited: false,
		},
	}

	// Run and assert.
	res, err := graph.Run(ctx, data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", res)
}

func main() {
	RunGraphRendering()
	RunPipeline()
}
