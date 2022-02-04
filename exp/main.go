package main

import (
	"html/template"
	"os"
)

type Dog struct {
	Name  string
	Breed string
	Age   int
}

type User struct {
	Name string
	Dogs []Dog
}

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	dogs := []Dog{{
		Name:  "Fido",
		Age:   4,
		Breed: "Mix",
	},
		{
			Name:  "Fluf",
			Age:   7,
			Breed: "Husky",
		}}

	data := User{
		Name: "John Smith",
		Dogs: dogs,
	}

	err = t.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}
