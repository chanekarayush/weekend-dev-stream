package main

import (
	"log"

	"github.com/go-git/go-git/v6"
)

const daysInLastSixMonths = 183
const outOfRange = 99999 

func stats(email string)  {
	commits := processRepos(email)
	printCommits(commits)
}

/**
* @param email of the user for which the commits are going to be scanned
* @return map[int]int A map with number of commits per day for a 
* given length of time
*/
func processRepos(email string) map[int]int {
	filePath := getDotFilePath()
	repos, err := parseLinesToSlice(filePath)
	if err != nil{
		log.Fatal(err)
	}
	daysInMap := daysInLastSixMonths
	commits := make(map[int]int)
	for i:=1; i<=daysInMap; i++{
		commits[i] = 0
	}

	for _, path := range repos{
		commits = fillCommits(email, path, commits)
	}
	return commits
}

func fillCommits(email string, path string, commits map[int]int) map[int]int {
	repo, err := git.PlainOpen(path)
	if err != nil{
		log.Panic(err)
	}
	ref, err := repo.Head()
	if err != nil{
		log.Panic(err)
	}
	iterator, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil{
		log.Panic(err)
	}

	offset := calcOffset()
	err = iterator.ForEach(func(c *object.Commit) error {
		daysAgo := counDaysSinceDate(c.Author.When) + offset
		if c.Author.Email != email{
			return nil
		}
		if daysAgo != outOfRange{
			commits[daysAgo]++
		}
		return nil
	})
	if err != nil{
		log.Panic(err)
	}
	
	return commits
}

func printCommits(commits map[int]int)  {
	
}
