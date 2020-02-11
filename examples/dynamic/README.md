# Dynamic samples

This directory contains samples of dynamic graphs.

This type of graphs can carry the state of a transaction and because of this they can't be reused (as it may lead to unexpected behaviours).

They are commonly used in none repetitive use-cases or when the overhead of creating the graph doesn't matter (performance wise), they have a much stronger type safety because of the explicit declarations when building the graph (eg. We can have `func Run(userId int, codeId int)` instead of having `func Run(ctx Context)` and asserting if the ctx contains `userId` and `codeId` in a statical graph) 