package main

import (
	"html/template"
	"os"
)

type User struct {
	Name   string
	Dog    string
	DogAge int
}

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	data := User{
		Name:   "John Smith",
		Dog:    "Fido",
		DogAge: 4,
	}

	err = t.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}
