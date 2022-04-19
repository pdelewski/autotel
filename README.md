# autotel
Automatic manual tracing :)
The aim of this project is to show how golang can be used to automatically inject 
open telemetry tracing (https://github.com/open-telemetry/opentelemetry-go).
It's also a place where compiler meet open telemetry.

## How to use it

```
./autotel [path to your go project]
```

## How it works

autotel will search for all root functions anotated with following call

```
	rtlib.SumoAutoInstrument()
```

where rtlib is small runtime library. Then all function calls starting from this function will be 
intrumented automatically. Example, below

```
package main

import (
	"fmt"

	"sumologic.com/autotel/rtlib"
)

func main() {
	rtlib.SumoAutoInstrument()
	fmt.Println(FibonacciHelper(10))
}
```

We can imagine other methods to say what needs to be instrumented (by argument(s) passed to autotel or configuration file).
