package main

import "github.com/krobus00/product-service/cmd"

var (
	name    string
	version string
)

func main() {
	cmd.Init(name, version)
	cmd.Execute()
}
