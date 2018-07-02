package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
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
	exlcudedCategoryGroups := map[string]bool{"Income": true, "Internal Master Category": true, "Credit Card Payments": true, "Hidden Categories": true}

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

	hostName := "smtp.gmail.com"
	auth := smtp.PlainAuth("", "mail.mulroy@gmail.com", os.Getenv("MAIL_PASSWD"), hostName)

	var categories string

	for _, categoryGroup := range data.Data.CategoryGroups {
		if !exlcudedCategoryGroups[categoryGroup.Name] {
			categories = categories + fmt.Sprintf("---------- %s ----------\n\n", categoryGroup.Name)
			for _, category := range categoryGroup.Categories {
				categories = categories + fmt.Sprintf("%s: $%.2f\n\n", category.Name, category.Balance/1000)
			}
		}
	}

	msg := []byte("To: dillon.mulroy@gmail.com, cbechdel@2u.com\r\n" +
		"Subject: Daily Budget Report!\r\n" +
		"\r\n" + categories)

	err = smtp.SendMail(hostName+":587", auth, "mail.mulroy@gmail.com", []string{"dillon.mulroy@gmail.com", "ccbechdel@2u.com"}, msg)

	if err != nil {
		fmt.Printf("Something went wrong sending mail: %s", err)
	}
}
