package main

import (
	"bufio"
	"io/fs"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func getIgnoredDirs() (map[string]int, error) {
	parsedDirs, err := parseLinesToSlice(IGNOREFILEPATH)
	if err != nil {
		log.Panic(err)
	}
	ignoredDirs := make(map[string]int)

	for i, x := range parsedDirs {
		if strings.HasPrefix(x, "*/") || strings.HasPrefix(x, "/") {
			ignoredDirs[strings.Split(x, "/")[1]] = i
		} else if strings.HasSuffix(x, "/") || strings.HasPrefix(x, "/**") {
			ignoredDirs[strings.Split(x, "/")[0]] = i
		} else {
			ignoredDirs[x] = i
		}
	}

	return ignoredDirs, nil
}

func scanForGitDir(folderPath string) ([]string, error) {
	gitDirs := make([]string, 1)
	ignoredDirs, err := getIgnoredDirs()
	if err != nil {
		log.Panicln(err)
	}
	err = filepath.WalkDir(folderPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if _, ok := ignoredDirs[d.Name()]; ok {
			return filepath.SkipDir
		}
		if d.IsDir() && d.Name() == ".git" {
			gitDirs = append(gitDirs, path)
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return gitDirs, nil
}

func getDotFilePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	dotFilePath := filepath.Join(usr.HomeDir, "/.local-git-heatmap")
	return dotFilePath
}

func parseLinesToSlice(dotFilePath string) ([]string, error) {
	f, err := os.Open(dotFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	existingRepos := make([]string, 0)

	for sc.Scan() {
		existingRepos = append(existingRepos, sc.Text())
	}
	if err = sc.Err(); err != nil {
		return nil, err
	}

	return existingRepos, nil
}

func mergeSlices(newRepos []string, existingRepos []string) ([]string, error) {
	repoMap := make(map[string]bool)
	mergedRepos := make([]string, 0)

	for i := range existingRepos {
		if _, ok := repoMap[existingRepos[i]]; ok == false {
			repoMap[existingRepos[i]] = true
			mergedRepos = append(mergedRepos, existingRepos[i])
		}
	}

	for i := range newRepos {
		if _, ok := repoMap[newRepos[i]]; ok == false {
			repoMap[newRepos[i]] = true
			mergedRepos = append(mergedRepos, newRepos[i])
		}
	}
	return mergedRepos, nil
}

func writeToFile(dotFilePath string, dedupedRepos []string) error {
	file, err := os.Create(dotFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	for i := range dedupedRepos {
		_, err = file.WriteString(dedupedRepos[i] + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

/** deduplicate files
* write them to `dotFilePath`
 */
func addNewItemsToSlice(dotFilePath string, newRepos []string) {
	existingRepos, err := parseLinesToSlice(dotFilePath)
	if err != nil {
		log.Printf("Looks Like there was an error in reading %s\n", dotFilePath)
		file, err := os.Create(dotFilePath)
		if err != nil {
			log.Fatal(err)
		} else {
			log.Printf("File was successfully created at %s\n", dotFilePath)
		}
		defer file.Close()
	}
	dedupedRepos, err := mergeSlices(newRepos, existingRepos)
	if err != nil {
		log.Fatal(err)
	}
	err = writeToFile(dotFilePath, dedupedRepos)
	if err != nil {
		log.Fatal(err)
	}
}

/**
* Scans all the subdirs of the given valid directory for `.git` folders
* @param folderPath is the path of the parent folder for the subdirs with `.git` directories
 */
func Scan(folderPath string) {
	fileinfo, err := os.Stat(folderPath)
	if err != nil || !fileinfo.IsDir() {
		log.Fatalln("The given folder does not exist.\nPlease enter a valid Directory Path!")
	}
	folders, err := scanForGitDir(folderPath)
	if err != nil {
		log.Fatalln(err)
	}
	dotFilePath := getDotFilePath()
	addNewItemsToSlice(dotFilePath, folders)
}
