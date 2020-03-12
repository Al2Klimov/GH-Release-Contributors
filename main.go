package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/go-github/v29/github"
	"net/http"
	"os"
	"sort"
)

func main() {
	user := flag.String("user", "", "GitHub user")
	owner := flag.String("owner", "", "GitHub repo owner")
	repo := flag.String("repo", "", "GitHub repo")
	milestone := flag.String("milestone", "", "GitHub repo milestone")

	flag.Parse()

	if *user == "" {
		fmt.Fprintln(os.Stderr, "-user missing")
		os.Exit(2)
	}

	if *owner == "" {
		fmt.Fprintln(os.Stderr, "-owner missing")
		os.Exit(2)
	}

	if *repo == "" {
		fmt.Fprintln(os.Stderr, "-repo missing")
		os.Exit(2)
	}

	if *milestone == "" {
		fmt.Fprintln(os.Stderr, "-milestone missing")
		os.Exit(2)
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "GITHUB_TOKEN env var missing")
		os.Exit(2)
	}

	contributors := map[string]struct{}{}

	{
		gh := github.NewClient(&http.Client{
			Transport: &github.BasicAuthTransport{
				Username: *user,
				Password: token,
			},
		})

		prlo := github.PullRequestListOptions{
			State: "closed",
			ListOptions: github.ListOptions{
				Page: 1,
			},
		}

		for {
			prs, _, errLP := gh.PullRequests.List(context.Background(), *owner, *repo, &prlo)
			if errLP != nil {
				fmt.Fprintln(os.Stderr, errLP.Error())
				os.Exit(1)
			}

			if len(prs) < 1 {
				break
			}

			for _, pr := range prs {
				if pr != nil && pr.Milestone != nil && pr.Milestone.Title != nil && *pr.Milestone.Title == *milestone {
					if pr.User != nil && pr.User.Login != nil {
						contributors[*pr.User.Login] = struct{}{}
					}
				}
			}

			prlo.Page++
		}
	}

	sortedContributors := make([]string, 0, len(contributors))
	for contributor := range contributors {
		sortedContributors = append(sortedContributors, contributor)
	}

	sort.Strings(sortedContributors)

	for _, contributor := range sortedContributors {
		fmt.Println(contributor)
	}
}
