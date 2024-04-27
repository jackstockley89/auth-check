package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

//Used by RepoCommits function to store output
type Repo_Commits = struct {
	Day   int
	Month time.Month
	Year  int
	Total int
}

// Used by CodeFrequency function to store output
type Code_Frequency struct {
	Day       int
	Month     time.Month
	Year      int
	Additions int
	Deletions int
}

// connects to the client that is used by other functions to call the api
func CallClient(token string) *github.Client {
	// get env token
	// Connect to giithub
	var client *github.Client
	if token == "" {
		client = github.NewClient(nil)
	} else {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)

		client = github.NewClient(tc)
	}
	return client
}

// Collect the number of commit to each repo by a user
func RepoCommits(githubrepo, githubref string) {
	client := CallClient(*token)
	//repo user and repo name
	githubrepoS := strings.Split(githubrepo, "/")
	repoUser := githubrepoS[0]
	repoName := githubrepoS[1]

	lca, _, err := client.Repositories.ListCommitActivity(context.Background(), repoUser, repoName)
	if err != nil {
		log.Fatal(err)
	}

	for i := range lca {
		p := lca[i]
		rc = Repo_Commits{Day: p.Week.Day(), Month: p.Week.Month(), Year: p.Week.Year(), Total: *p.Total}
		fmt.Println("Commit History: ", rc)
	}

}

// Collect code change frequency
func CodeFrequency(githubrepo, githubref string) {
	client := CallClient(*token)
	//repo user and repo name
	githubrepoS := strings.Split(githubrepo, "/")
	repoUser := githubrepoS[0]
	repoName := githubrepoS[1]

	wfs, _, err := client.Repositories.ListCodeFrequency(context.Background(), repoUser, repoName)
	if err != nil {
		log.Fatal(err)
	}

	for i := range wfs {
		p := wfs[i]
		cf = Code_Frequency{Day: p.Week.Day(), Month: p.Week.Month(), Year: p.Week.Year(), Additions: *p.Additions, Deletions: *p.Deletions}
		fmt.Println("Code Frequency: ", cf)
	}
}

var (
	rc Repo_Commits
	cf Code_Frequency
	t  = time.Now()
	ft = fmt.Sprintf("%d-%02d-%0d %02d:%02d:%02d",
		t.Day(), t.Month(), t.Year(),
		t.Hour(), t.Minute(), t.Second())
	token      = flag.String("token", os.Getenv("GITHUB_OAUTH_TOKEN"), "GihHub Personel token string")
	githubrepo = flag.String("githubrepo", os.Getenv("GITHUB_REPOSITORY"), "Github Repository string")
	githubref  = flag.String("githubref", os.Getenv("GITHUB_REF"), "Github Respository PR ref string")
	caColour   = "\033[1;33m%s\033[0m \033[1;36m%s\033[0m \033[1;32m[%v %v %v]\033[0m \033[1;36m%s\033[0m \033[1;32m[%v]\033[0m"
	cfColour   = "\033[1;33m%s\033[0m \033[1;36m%s\033[0m \033[1;32m[%v %v %v]\033[0m \033[1;36m%s\033[0m \033[1;32m[%v]\033[0m \033[1;36m%s\033[0m \033[1;32m[%v]\033[0m"
)

// Main function used to run all other function and Print output collected with the structs to the commandline
func main() {
	flag.Parse()
	RepoCommits(*githubrepo, *githubref)
	CodeFrequency(*githubrepo, *githubref)
	// Prints out the states collected from all above functions passed in from data saved within the structs
	// struct: Repo_Commits
	fmt.Printf(caColour, "\nCommit Activity", "\nWeek Starting:", rc.Day, rc.Month, rc.Year, "Total Commits:", rc.Total)
	// struct: Code_Frequency
	fmt.Printf(cfColour, "\n\nCode Frequency", "\nWeek Starting:", cf.Day, cf.Month, cf.Year, "Total Additions:", cf.Additions, "Total Deletions:", cf.Deletions)
	// Creating a Timestamp for when the stats ran
	fmt.Println("\n\nStats last updated on:", ft)
}
