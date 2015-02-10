# go-yrapi
A toy project for learning Go in a useful way


```bash
$ go get github.com/runeb/go-yrapi/yrapi
```

```Go
package main

import(
  "github.com/runeb/go-yrapi/yrapi"
  "github.com/davecgh/go-spew/spew"
)

func main() {
  weatherData,err := yrapi.LocationforecastLTS(59.95, 10.75)
  if(err == nil) {
    spew.Dump(weatherData)
  }
}
```

