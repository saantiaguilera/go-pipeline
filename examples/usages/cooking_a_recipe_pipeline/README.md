# Cook Example

This example is a simple demonstration on how to serve a dish (which in this case consists of meat + salad).

The approach used for this example was a static immutable graph, meaning:
- We New the graph only once and it can be reused as many times as wanted
- The graph doesn't contain any state. The state is passed at execution time through the context

We showcase on different steps various approaches on creating them:
- `cut_eggs_step.go`: A constructor method of a single step
- `cut_carrots_step.go`: A unit of work that doesn't know about pipelines but is later injected on a step
- `cut_meat_step.go`: A custom step that allows us to do whatever complex logic we want to in a flexible and simple way. This is a custom step (similar to the ones the pipeline API provides)

For demonstration purposes, the graph is built in a single method. Note that you can easily decouple this into N more meaningfull methods (one for creating a step of the meat, another for the salad, etc)

The created graph looks exactly like this

![](template.svg)