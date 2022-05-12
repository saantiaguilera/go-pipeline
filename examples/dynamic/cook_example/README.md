# Dynamic Cook Example

This example is a simple demonstration on how to serve a dish (which in this case consists of meat + salad).

The approach used for this example was a dynamic graph, meaning:
- We New the graph for each time we want to execute it. Each graph is "unique" since it carries state of the "transaction" (dish sizes and stuff for this example)
- The context doesn't serve much since we are constructing the graph with an already defined state. Yet, it can be powerful to hold external data (not bound to any part of the graph)

For demonstration purposes, the graph is built in a single method. Note that you can easily decouple this into N more meaningfull methods (one for creating a stage of the meat, another for the salad, etc)

The Newd graph looks exactly like this

![](template.svg)