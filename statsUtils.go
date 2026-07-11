package main

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
)

const (
	daysInLastSixMonths  = 183
	weeksInLastSixMonths = 26
	outOfRange           = 99999
)

type columns []int

func stats(email string) {
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
	if err != nil {
		log.Fatal(err)
	}
	daysInMap := daysInLastSixMonths
	commits := make(map[int]int)
	for i := 0; i <= daysInMap; i++ {
		commits[i] = 0
	}

	for _, path := range repos {
		commits = fillCommits(email, path, commits)
	}
	return commits
}

func fillCommits(email string, path string, commits map[int]int) map[int]int {
	repo, err := git.PlainOpen(path)
	if err != nil {
		if errors.Is(err, git.ErrRepositoryNotExists) {
			fmt.Printf("\x1b[0;1;3m\x1b[38;2;204;2;2m Repository does not exist:\x1b[0m %s\n", path)
			return commits
		}
		log.Panic(err)
	}
	ref, err := repo.Head()
	if err != nil {
		// Double-check both the sentinel error and the literal text string
		if errors.Is(err, plumbing.ErrReferenceNotFound) || err.Error() == "reference not found" {
			fmt.Printf("\x1b[0;1;3m\x1b[38;2;255;239;0m Skipping empty repository:\x1b[0m %s\n", path)
			return commits
		}
		log.Panic(err)
	}
	iterator, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		log.Panic(err)
	}

	offset := calcOffset()
	err = iterator.ForEach(func(c *object.Commit) error {
		daysAgo := countDaysSinceDate(c.Author.When) + offset
		if c.Author.Email != email {
			return nil
		}
		if daysAgo != outOfRange {
			commits[daysAgo]++
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return commits
}

// func countDaysSinceDate(commitDate time.Time) int {
//     now := getBeginningOfDay(time.Now())
//     commitDayStart := getBeginningOfDay(commitDate)
//
//     // Normalize time zones to compare apples to apples
//     days := int(now.Sub(commitDayStart).Hours() / 24)
//
//     if days < 0 || days > daysInLastSixMonths {
//         return outOfRange
//     }
//     return days
// }

func countDaysSinceDate(commitDate time.Time) int {
	now := getBeginningOfDay(time.Now())
	commitDate = getBeginningOfDay(commitDate)
	days := 0
	for commitDate.Before(now) {
		commitDate = commitDate.Add(24 * time.Hour)
		days++
		if days > daysInLastSixMonths {
			return outOfRange
		}
	}
	return days
}

func getBeginningOfDay(currentTime time.Time) time.Time {
	year, month, day := currentTime.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, currentTime.Location())
	return startOfDay
}

func calcOffset() int {
	weekday := time.Now().Weekday()
	offset := 0
	switch weekday {
	case time.Sunday:
		offset = 7
	case time.Monday:
		offset = 6
	case time.Tuesday:
		offset = 5
	case time.Wednesday:
		offset = 4
	case time.Thursday:
		offset = 3
	case time.Friday:
		offset = 2
	case time.Saturday:
		offset = 1
	}
	return offset
}

func printCommits(commits map[int]int) {
	keys := sortMapIntoSlice(commits)
	cols := generateColumns(keys, commits)
	printCells(cols)
}

func sortMapIntoSlice(commits map[int]int) []int {
	keys := make([]int, 0)
	for key := range commits {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	return keys
}

func generateColumns(keys []int, commits map[int]int) map[int]columns {
	cols := make(map[int]columns)
	col := columns{}
	for _, k := range keys {
		week := int(k / 7)
		dayInWeek := k % 7

		if dayInWeek == 0 {
			col = columns{}
		}

		col = append(col, commits[k])

		if dayInWeek == 6 {
			cols[week] = col
		}
	}
	return cols
}

// printCells prints the cells of the graph
func printCells(cols map[int]columns) {
	printMonths()
	for j := 6; j >= 0; j-- {
		for i := weeksInLastSixMonths + 1; i >= 0; i-- {
			if i == weeksInLastSixMonths+1 {
				printDayCol(j)
			}
			if col, ok := cols[i]; ok {
				// special case today
				if i == 0 && j == calcOffset()-1 {
					printCell(col[j], true)
					continue
				} else {
					if len(col) > j {
						printCell(col[j], false)
						continue
					}
				}
			}
			printCell(0, false)
		}
		fmt.Printf("\n")
	}
}

func printMonths() {
	week := getBeginningOfDay(time.Now()).Add(-(daysInLastSixMonths * time.Hour * 24))
	month := week.Month()
	fmt.Printf("         ")
	for {
		if week.Month() != month {
			fmt.Printf("%s ", week.Month().String()[:3])
			month = week.Month()
		} else {
			fmt.Printf("    ")
		}

		week = week.Add(7 * time.Hour * 24)
		if week.After(time.Now()) {
			break
		}
	}
	fmt.Printf("\n")
}

// printCell given a cell value prints it with a different format
// based on the value amount, and on the `today` flag.
func printCell(val int, today bool) {
	escape := "\033[0;37;30m"
	switch {
	case val > 0 && val < 5:
		escape = "\033[1;30;47m"
	case val >= 5 && val < 10:
		escape = "\033[1;30;43m"
	case val >= 10:
		escape = "\033[1;30;42m"
	}

	if today {
		escape = "\033[1;37;45m"
	}

	if val == 0 {
		fmt.Print(escape + "  - " + "\033[0m")
		return
	}

	str := "  %d "
	switch {
	case val >= 10:
		str = " %d "
	case val >= 100:
		str = "%d "
	}

	fmt.Printf(escape+str+"\033[0m", val)
}

func printDayCol(day int) {
	out := "     "
	switch day {
	case 0:
		out = " Sat "
	case 2:
		out = " Tue "
	case 4:
		out = " Thu "
	case 6:
		out = " Sun "
	}

	fmt.Print(out)
}
