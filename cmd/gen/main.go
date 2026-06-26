package main

import (
	"os"

	"github.com/dave/jennifer/jen"
)

func main() {
	file := jen.NewFilePath("nuconv/test")

	reused := jen.Qual("fmt", "Println").Call(jen.Lit("hello"))

	file.Func().Id("init").Params().Block(
		jen.Id("x").Op(":=").Qual("nuconv/test", "Test").Block(),
		reused,
		reused,
	)

	err := file.Render(os.Stdout)
	if err != nil {
		panic(err)
	}
}
