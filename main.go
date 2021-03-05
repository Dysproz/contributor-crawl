package main

import (
	// "flag"
	"log"
	"os"
	"fmt"
	"context"
	"time"
	"github.com/google/go-github/v33/github"
)

var ctx = context.Background()

// Contributor is a struct with information about Contributor for CSV processing
type Contributor struct {
	Repository string
	AuthorLogin string
	AuthorName string
	AuthorEmail string
	TotalContributions int
	RangeContributions int
	RangeAdditions int
	RandeDeletions int
}

func main() {
	// organization := flag.String("organization", "", "GitHub organization/profile to gather projects from.")
	// // startDate := flag.String("start-date", "", "Start date from which contributions should be counted. DD-MM-YYYY")
	// // endDate := flag.String("end-date", "", "End date to which contributions should be counted. DD-MM-YYYY")
	// flag.Parse()

	// if *organization == "" {
	// 	log.Fatal("No organization name specified.")
	// 	os.Exit(1)
	// }

	organization := "tungstenfabric"
	client := github.NewClient(nil)
	repos, _, err := client.Repositories.ListByOrg(ctx, organization, nil)
	if err != nil {
		log.Fatal("Could not list repositories for organization %s", organization)
		os.Exit(1)
	}
	fmt.Println(*repos[0].Name)
	stats, _, err := client.Repositories.ListContributorsStats(ctx, organization, *repos[0].Name)
	fmt.Println(stats[0].Weeks)
}

func countContributionsInTime(stats []*github.WeeklyStats, startDate, endDate string) (additions, deletions, commits int) {
	var additions int
	var deletions int
	var commits int
	time.Parse()
	for 
}