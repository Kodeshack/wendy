# Wendy
### Simple Scaffolding Library for Go

![example workflow](https://github.com/Kodeshack/wendy/actions/workflows/tests.yaml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/Kodeshack/wendy.svg)](https://pkg.go.dev/github.com/Kodeshack/wendy)


## Usage Example

```go
g := &FSGenerator{
	RootDir: "./generate",
}

err := g.Generate(
	PlainFile("README.md", "# Wendy"),
	Dir("bin",
		Dir("cli",
			PlainFile("main.go", "package main"),
		),
	),
	Dir("pkg",
		PlainFile("README.md", "how to use this thing"),
		Dir("cli",
			PlainFile("cli.go", "package cli..."),
			PlainFile("run.go", "package cli...run..."),
		),
	),
)
````
