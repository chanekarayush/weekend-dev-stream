package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func scanForGitDir(folderPath string){
	gitDirs := make([]string, 1)
	err := filepath.WalkDir(folderPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err;
		}
		if (d.Name() == "node_modules" || d.Name() == "env" || d.Name() == "venv"){
			return filepath.SkipDir
		}
		if (d.IsDir()){
			println(path)
		}
		if (d.IsDir() && d.Name()==".git"){
			gitDirs = append(gitDirs, path)
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil{
		log.Fatal(err)
	}

	for i:=range gitDirs{
		fmt.Printf("Scanned: %s\n", gitDirs[i])
	}
}

/**
* Scans all the subdirs of the given valid directory for `.git` folders
* @param folderPath is the path of the parent folder for the subdirs with `.git` directories
*/
func scan (folderPath string)  {
	fileinfo, err := os.Stat(folderPath)
	if (err != nil || !fileinfo.IsDir()){
		log.Fatalln("The given folder does not exist.\nPlease enter a valid Directory Path!")
	}
	scanForGitDir(folderPath)
}

func computeStats(stats int)  {
	print(stats)	
}

func main()  {
	var folder string
	var email string
	flag.StringVar(&folder,"add-dir", "", "Directory to be added for include in stats")
	flag.StringVar(&email,"email", "", "User Email to be searched for in the commits")
	flag.Parse();

	if (folder != ""){
		scan(folder)
		// if (email != ""){
		// 	fmt.Printf("[Email]: %s\n",email)
		// }
		return
	}
	fmt.Println(folder)
}
