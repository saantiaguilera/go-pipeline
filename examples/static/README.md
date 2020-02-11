# Static samples

This directory contains samples of statically immutable graphs.

This type of graphs are stateless by nature (they can't contain information about a specific transaction) and because of this they can be reused as many times as someone wants to.

They are commonly used across web-services since they usually have a much better performance than dynamic graphs (because of the non-existent overhead of creating the whole graph on each transaction), but have a much weaker type safety because of the static communication between elements of the graph lacking more explicit function declarations (eg. instead of `func Run(userId int, codeId int)` we have `func Run(ctx Context)` and assert if the ctx contains `userId` and `codeId`)