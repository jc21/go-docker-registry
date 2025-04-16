# Docker Registry Library for Golang

> [!WARNING]
> This is WIP and not ready for use yet

This package can be included in your golang apps to talk to a docker registry.

## Basic Example, HTTP

```bash
go get github.com/jc21/go-docker-registry@latest
```

```go
package main

import (
	"fmt"
	registry "github.com/jc21/go-docker-registry"
)

func main() {
	server := "http://registry.local:5000"
	username := ""
	password := ""

	reg, _ := registry.NewInsecureServer(server, username, password)
	images, _ := reg.GetCatalog()

	for _, image := range images {
		tags, _ := reg.GetTags(image)
		fmt.Printf("%s -- %v\n", image, tags)
	}
}
```

## HTTPS example

Should your registry be behind a reverse proxy or other SSL termination

```go
func main() {
	// ...
	server := "https://registry.local"
	reg, err := registry.NewServer(server, username, password)
	// ...
}
```

## Custom Logger Example

```go
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

func main() {
	// ...
	reg, err := registry.NewInsecureServer(server, username, password, &myLogger)
	// ...
}
```
