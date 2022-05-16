package main

import (
	"context"
	"sync"

	"github.com/saantiaguilera/go-pipeline"
)

type processDishConcurrently struct {
	SaladStep pipeline.Step[MealMaterials, Salad]
	MeatStep  pipeline.Step[MealMaterials, CookedMeat]
}

func (p processDishConcurrently) Draw(g pipeline.Graph) {
	g.AddConcurrency(p.newStepGraphActivity(p.SaladStep), p.newStepGraphActivity(p.MeatStep))
}

func (p processDishConcurrently) Run(ctx context.Context, in MealMaterials) (DishContents, error) {
	wg := new(sync.WaitGroup)
	wg.Add(2)

	errs := make(chan error, 2)
	var salad *Salad
	var meat *CookedMeat

	go func() {
		defer wg.Done()
		v, err := p.SaladStep.Run(ctx, in)
		if err != nil {
			errs <- err
			return
		}
		salad = &v
	}()

	go func() {
		defer wg.Done()
		v, err := p.MeatStep.Run(ctx, in)
		if err != nil {
			errs <- err
			return
		}
		meat = &v
	}()

	wg.Wait()
	close(errs)

	if len(errs) > 0 {
		return DishContents{}, <-errs
	}
	return DishContents{
		Salad: *salad,
		Meat:  *meat,
	}, nil
}

func (p processDishConcurrently) newStepGraphActivity(drawable pipeline.DrawableGraph) pipeline.GraphDrawer {
	return func(graph pipeline.Graph) {
		drawable.Draw(graph)
	}
}
