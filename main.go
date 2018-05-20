package main

import "fmt"

type categoryListResponse struct {
	Data struct {
		CategoryGroups []categoryList
	}
}

type categoryList struct {
	ID         string
	Name       string
	Hidden     bool
	Categories []category
}

type category struct {
	ID               string
	CategoryStringID string
	Name             string
	Hidden           bool
	Note             string
	Budgeted         float32
	Activity         float32
	Balance          float32
}

func main() {
	clr := new(categoryListResponse)

	fmt.Print(clr)
}
