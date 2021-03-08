package main

import (
	"context"
	"encoding/csv"
	"flag"
	"log"
	"os"
	"time"

	retry "github.com/avast/retry-go"
	"github.com/dnlo/struct2csv"
	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

var ctx = context.Background()

// Contributor is a struct with information about Contributor for CSV processing
type Contributor struct {
	Repository         string
	AuthorLogin        string
	TotalContributions int
	RangeContributions int
	RangeAdditions     int
	RangeDeletions     int
}

// ISOTimeLayout is a sample layout used for time date parse
var ISOTimeLayout = "2006-01-02"

func main() {
	organization := flag.String("organization", "", "GitHub organization/profile to gather projects from.")
	startDateInput := flag.String("start-date", "", "Start date from which contributions should be counted. YYYY-MM-DD")
	endDateInput := flag.String("end-date", "", "End date to which contributions should be counted. YYYY-MM-DD")
	outFile := flag.String("out-file", "contributors-crawl-out.csv", "Path to file where output data in CSV format should be written")
	repository := flag.String("repository", "", "(OPTIONAL) Specific repository to gather contributions from. If left empty then all repositories in organization will be crawled.")
	OauthToken := flag.String("oauth", "", "OAuth GitHub API token. Standard clients have 60 requests per hour while authorized clients have 15k requests per hour")
	flag.Parse()

	if *organization == "" {
		log.Fatal("No organization name specified.")
		os.Exit(1)
	}
	startDate, err := time.Parse(ISOTimeLayout, *startDateInput)
	if err != nil {
		log.Fatal("Could not parse start (", *startDateInput, ") date string; err: ", err)
		os.Exit(1)
	}
	endDate, err := time.Parse(ISOTimeLayout, *endDateInput)
	if err != nil {
		log.Fatal("Could not parse end (", *endDateInput, ") date string; err: ", err)
		os.Exit(1)
	}

	var client *github.Client
	if *OauthToken == "" {
		client = github.NewClient(nil)
	} else {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: *OauthToken},
		)
		tc := oauth2.NewClient(context.Background(), ts)
		client = github.NewClient(tc)

	}
	var contributions []Contributor
	var repos []*github.Repository
	if *repository == "" {
		opt := &github.RepositoryListByOrgOptions{ListOptions: github.ListOptions{
			PerPage: 999999,
		}}
		repos, _, err = client.Repositories.ListByOrg(ctx, *organization, opt)
		if err != nil {
			log.Fatal("Could not list repositories for organization ", *organization, "; err: ", err)
			os.Exit(1)
		}
		log.Println("Successfully gathered repositories for ", *organization, " organization.")
	} else {
		repos = []*github.Repository{
			{
				Name: repository,
			},
		}
	}
	for _, repo := range repos {
		var stats []*github.ContributorStats
		err := retry.Do(
			func() error {
				stats, _, err = client.Repositories.ListContributorsStats(ctx, *organization, *repo.Name)
				if err != nil {
					log.Println("Request to GitHub failed with err : ", err, "; Trying again in 5 seconds...")
				}
				return err
			}, retry.Delay(5*time.Second),
		)
		if err != nil {
			log.Fatal("Could not list contributors for ", *repo.Name, " repository in organization ", *organization, "; err: ", err)
			os.Exit(1)
		}
		log.Println("Successfully contributors for ", *repo.Name, " repository.")

		for _, stat := range stats {
			adds, dels, commits := countContributionsInTime(stat.Weeks, startDate, endDate)
			log.Println("Successfully gathered data for contributor ", *stat.Author.Login, " in repository ", *repo.Name)
			if adds != 0 || dels != 0 || commits != 0 {
				contributions = append(contributions, Contributor{
					Repository:         *repo.Name,
					AuthorLogin:        *stat.Author.Login,
					TotalContributions: *stat.Total,
					RangeContributions: commits,
					RangeAdditions:     adds,
					RangeDeletions:     dels,
				})
			}
		}
	}

	CSVenc := struct2csv.New()
	CSVData, err := CSVenc.Marshal(contributions)
	if err != nil {
		log.Fatal("Could not marshal gathered data to CSV with err: ", err)
		os.Exit(1)
	}

	f, err := os.Create(*outFile)
	if err != nil {
		log.Fatal("Could not create file ", outFile, " with error: ", err)
	}
	defer f.Close()

	CSVwriter := csv.NewWriter(f)
	if err := CSVwriter.WriteAll(CSVData); err != nil {
		log.Fatal("Could not write CSV data to file with err: ", err)
	}
	log.Println("Contribution crawling ended sucessfully! Find your data in file: ", *outFile)
	os.Exit(0)
}

func countContributionsInTime(stats []*github.WeeklyStats, startDate, endDate time.Time) (additions, deletions, commits int) {
	for _, stat := range stats {
		if stat.Week.Before(endDate) && stat.Week.After(startDate) {
			additions += *stat.Additions
			deletions += *stat.Deletions
			commits += *stat.Commits
		}
	}
	return additions, deletions, commits
}
