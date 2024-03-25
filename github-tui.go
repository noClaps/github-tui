package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
)

type User struct {
	Login   string `json:"login"`
	HtmlUrl string `json:"html_url"`
}

type Gist struct {
	HtmlUrl     string `json:"html_url"`
	Description string `json:"description"`
}

type Repository struct {
	Name        string `json:"name"`
	HtmlUrl     string `json:"html_url"`
	Description string `json:"description"`
}

var (
	name     string
	userData User
	option   string
)

func form() {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What's the username?").
				Value(&name).
				Validate(func(name string) error {
					res, err := http.Get("https://api.github.com/users/" + name)
					if err != nil {
						return err
					}

					if res.StatusCode == 404 {
						return errors.New("User not found")
					}

					return nil
				}),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What would you like to see?").
				Options(
					huh.NewOption("About", "about"),
					huh.NewOption("Followers", "followers"),
					huh.NewOption("Following", "following"),
					huh.NewOption("Gists", "gists"),
					huh.NewOption("Starred repos", "starred"),
					huh.NewOption("Subscriptions", "subscriptions"),
					huh.NewOption("Organizations", "organizations"),
					huh.NewOption("Repos", "repos"),
					huh.NewOption("Quit", "quit"),
				).
				Validate(func(option string) error {
					if option == "quit" || option == "about" {
						return nil
					}

					var optionData []any
					getData("https://api.github.com/users/"+name+"/"+option, &optionData)

					if len(optionData) == 0 {
						return errors.New(option + ": Not found")
					}

					return nil
				}).
				Value(&option),
		),
	)

	err := form.Run()
	if err != nil {
		log.Error(err)
	}
}

func getData(url string, store any) {
	res, err := http.Get(url)
	if err != nil {
		log.Error(err)
	}
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(store)
}

func main() {
	logger := log.New(os.Stderr)
	form()
	getData("https://api.github.com/users/"+name, &userData)

	if option == "about" {
		logger.Printf(
			`Username: %s
GitHub URL: %s`, userData.Login, userData.HtmlUrl)
	}

	if option == "followers" || option == "following" || option == "organizations" {
		var users []User
		var url string

		if option == "followers" {
			url = "https://api.github.com/users/" + name + "/followers"
		} else if option == "following" {
			url = "https://api.github.com/users/" + name + "/following"
		} else if option == "organizations" {
			url = "https://api.github.com/users/" + name + "/orgs"
		}

		getData(url, &users)

		for _, user := range users {
			logger.Printf("%s - %s", user.Login, user.HtmlUrl)
		}
	}
	if option == "gists" {
		var gists []Gist
		getData("https://api.github.com/users/"+name+"/gists", &gists)

		for _, gist := range gists {
			logger.Printf("%s - %s", gist.Description, gist.HtmlUrl)
		}
	}
	if option == "starred" || option == "subscriptions" || option == "repos" {
		var repos []Repository
		var url string

		if option == "starred" {
			url = "https://api.github.com/users/" + name + "/starred"
		} else if option == "subscriptions" {
			url = "https://api.github.com/users/" + name + "/subscriptions"
		} else if option == "repos" {
			url = "https://api.github.com/users/" + name + "/repos"
		}

		getData(url, &repos)

		for _, repo := range repos {
			logger.Printf(`
Name: %s
Description: %s
URL: %s`, repo.Name, repo.Description, repo.HtmlUrl)
		}
	}
}
