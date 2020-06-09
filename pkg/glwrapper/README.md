# GL wrapper

This application is the wrapper for the gl lib. It was written to support the change between gl versions. With the previous solution, the gl lib was included in several files. Now it's included only in the wrapper package, and the wrapper is included in the apps instead of the gl lib.
The advantage of this solution, that you only need to update the included gl lib in the wrapper, and then change the version constants.
The disadvantage is that you need to write a new wrapper function if you want to use a new gl function.
