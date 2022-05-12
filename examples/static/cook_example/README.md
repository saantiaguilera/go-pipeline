# Cook Example

This example is a simple demonstration on how to serve a dish (which in this case consists of meat + salad).

The approach used for this example was a static immutable graph, meaning:
- We New the graph only once and it can be reused as many times as wanted
- The graph doesn't contain any state. The state is passed at execution time through the context

For demonstration purposes, the graph is built in a single method. Note that you can easily decouple this into N more meaningfull methods (one for creating a container of the meat, another for the salad, etc)

The Newd graph looks exactly like this

![](template.svg)