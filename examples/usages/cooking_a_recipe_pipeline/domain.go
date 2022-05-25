package main

// This package contains all domain entities for this sample.
// Consider it a `domain`/`entities`/`whatever you call it` package of your project
type (
	MealMaterials struct {
		Eggs    []Egg
		Carrots []Carrot
		Meat    Meat

		Oven Oven
	}

	Egg struct {
		Boiled bool
	}

	Carrot struct {
		Washed bool
	}

	CutCarrot struct {
		// Something a cut carrot could have for you to process it
	}

	CutEgg struct {
		// Something a cut egg could have for you to process it
	}

	Meat struct {
		Size int
	}

	Vegetables struct {
		Eggs    []CutEgg
		Carrots []CutCarrot
	}

	Salad struct {
		Vegetables
		Mixed bool
	}

	Oven struct {
		Ignited bool
	}

	CookingTools struct {
		Oven Oven
		Meat Meat
	}

	CookedMeat struct {
		// Something a cooked meats could have for you to process it
	}

	DishContents struct {
		Salad Salad
		Meat  CookedMeat
	}

	Dish struct {
		Salad Salad
		Meat  CookedMeat
		// Something a dish could have for you to process it
	}
)
