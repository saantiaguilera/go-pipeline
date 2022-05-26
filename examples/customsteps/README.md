The pipeline API provides the barebones and founding steps that allows us to create workflows. However, in more complex scenarios it's useful to create custom steps that satisfy specific business needs, hence this package provides a quick overview through examples on how to create them. 

The API should be flexible enough to create whatever you like or need, as long as you are compliant with the Step contract which is highly flexible.

This package provides custom steps that showcase particular use cases someone may need, feel free to copy them / modify them however you please / create your own / etc.

It's worth mentioning none of the structs inside are exposed so if you plan on using them, make sure to expose the copied versions + create them a constructor function.