# Wendy
### Simple Scaffolding Library for Go

[![MIT License](https://img.shields.io/github/license/Kodeshack/wendy?style=flat-square)](https://github.com/Kodeshack/wendy/blob/main/LICENSE)
![Test Workflow](https://github.com/Kodeshack/wendy/actions/workflows/tests.yaml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/Kodeshack/wendy.svg)](https://pkg.go.dev/github.com/Kodeshack/wendy)
[![Go Report Card](https://goreportcard.com/badge/github.com/Kodeshack/wendy)](https://goreportcard.com/report/github.com/Kodeshack/wendy)
[![Latest Release](https://img.shields.io/github/v/tag/Kodeshack/wendy?sort=semver&style=flat-square)](https://github.com/Kodeshack/wendy/releases/latest)
[![codecov](https://codecov.io/gh/Kodeshack/wendy/branch/main/graph/badge.svg?token=JMVj1pFT2l)](https://codecov.io/gh/Kodeshack/wendy)



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

## License

[MIT](https://github.com/Kodeshack/wendy/blob/main/LICENSE)
