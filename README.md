# evoli

Color palettes: https://www.lihaoyi.com/post/BuildyourownCommandLinewithANSIescapecodes.html 

## Getting started

Create a main.go file with following context:
```GO
package main

import (
	"github.com/evolidev/evoli"
)

//go:generate go run main.go generate

type MyApp struct {
	*evoli.Application
}

func main() {
	app := &MyApp{}

	app.Start()
}
```

Next run following command `go run main.go init`. 
This command will initialize the project and creates all folders for you. 

To get an command overview just run `go run main.go`. 