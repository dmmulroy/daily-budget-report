package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type categoryListResponse struct {
	Data struct {
		CategoryGroups []categoryList `json:"category_groups"`
	}
}

type categoryList struct {
	ID         string
	Name       string
	Hidden     bool
	Categories []category
}

type category struct {
	ID              string
	CategoryGroupID string `json:"category_group_id"`
	Name            string
	Hidden          bool
	Note            string
	Budgeted        float32
	Activity        float32
	Balance         float32
}

func main() {
	APIToken := os.Getenv("YNAB_API_KEY")

	req, err := http.NewRequest("GET", "https://api.youneedabudget.com/v1/budgets/da7e4433-9470-4b70-b9fd-ef10722eedea/categories", nil)

	if err != nil {
		fmt.Printf("Something went wrong creating the request: %s", err)
		return
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", APIToken))

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Printf("Something went wrong with the request: %s", err)
		return
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	var data categoryListResponse

	err = decoder.Decode(&data)

	if err != nil {
		fmt.Printf("Something went wrong parsing JSON: %s", err)
		return
	}

	for _, categoryGroup := range data.Data.CategoryGroups {
		for _, category := range categoryGroup.Categories {
			fmt.Println(category.Name)
		}
	}

}
