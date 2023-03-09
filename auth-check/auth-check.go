package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	client "github.com/jackstockley89/github-actions/github-api/client"
	get "github.com/jackstockley89/github-actions/github-api/get"
	"github.com/sethvargo/go-githubactions"
)

var (
	token      = flag.String("token", os.Getenv("GITHUB_TOKEN"), "GihHub Personel token string")
	githubrepo = flag.String("githubrepo", os.Getenv("GITHUB_REPOSITORY"), "Github Repository string")
	githubref  = flag.String("githubref", os.Getenv("GITHUB_REF"), "Github Respository PR ref string")
	c          = client.ClientConnect(*token)
	g          = get.GetPullRequestData(*githubrepo, *githubref, *token)
)

func collaboratorCheck() (bool, error) {
	// get pull request
	pr := g.User
	if pr == "" {
		return false, fmt.Errorf("no pull request user found")
	}
	fmt.Println("Pull Request User:", pr)
	// compare pr user with the repo collaborators
	collab, _, err := c.Repositories.IsCollaborator(context.Background(), g.Owner, g.Repository, pr)
	if err != nil {
		return false, err
	}
	fmt.Println("Collaborator Status:", collab)
	return collab, nil
}

func main() {
	var body string
	flag.Parse()
	collab, err := collaboratorCheck()
	if err != nil {
		log.Fatal(err)
	}
	if collab {
		// create review on pull request
		body = fmt.Sprintf("Known collaborator %s", g.User)
	} else {
		// create comment on pull request
		body = fmt.Sprintf("Unknown collaborator %s", g.User)
	}
	githubactions.New().SetOutput("review", body)
}
