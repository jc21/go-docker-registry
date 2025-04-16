# Docker Registry Library for Golang

::NOTE::
This is WIP and not ready for use yet

This package can be included in your golang apps to talk to a docker registry.

## Example

```bash
go get github.com/jc21/go-docker-registry@latest
```

```go
package main

import (
	"fmt"
	"github.com/jc21/go-docker-registry"
)

// Setup your custom logger
type myLogger struct {}
func (l *myLogger) Debug(format string, args ...any) {
	fmt.Printf("[DEBUG] " + format + "\n", ...args)
}
func (l *myLogger) Info(format string, args ...any) {
	fmt.Printf("[INFO] " + format + "\n", ...args)
}
func (l *myLogger) Error(format string, args ...any) {
	fmt.Printf("[ERROR] " + format + "\n", ...args)
}

// Run
func main() {
	reg, _ = registry.NewInsecureServer("http://registry.local:5000", "", "", &myLogger{})
	images, _ := reg.GetCatalog()

	for _, image := range images {
		tags, err := reg.GetTags(image)
		fmt.Printf("%s tags: %v\n", image, tags)
	}
}
```
