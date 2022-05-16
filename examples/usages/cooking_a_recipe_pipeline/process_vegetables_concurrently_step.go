package main

import (
	"context"
	"sync"

	"github.com/saantiaguilera/go-pipeline"
)

type processVegetablesConcurrently struct {
	EggStep    pipeline.Step[[]Egg, []CutEgg]
	CarrotStep pipeline.Step[[]Carrot, []CutCarrot]
}

func (p processVegetablesConcurrently) Draw(g pipeline.Graph) {
	g.AddConcurrency(p.EggStep.Draw, p.CarrotStep.Draw)
}

func (p processVegetablesConcurrently) Run(ctx context.Context, in MealMaterials) (Vegetables, error) {
	wg := new(sync.WaitGroup)
	wg.Add(2)

	errs := make(chan error, 2)
	var cutEggs *[]CutEgg
	var cutCarrots *[]CutCarrot

	go func() {
		defer wg.Done()
		v, err := p.EggStep.Run(ctx, in.Eggs)
		if err != nil {
			errs <- err
			return
		}
		cutEggs = &v
	}()

	go func() {
		defer wg.Done()
		v, err := p.CarrotStep.Run(ctx, in.Carrots)
		if err != nil {
			errs <- err
			return
		}
		cutCarrots = &v
	}()

	wg.Wait()
	close(errs)

	if len(errs) > 0 {
		return Vegetables{}, <-errs
	}
	return Vegetables{
		Eggs:    *cutEggs,
		Carrots: *cutCarrots,
	}, nil
}
