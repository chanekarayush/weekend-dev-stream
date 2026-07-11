package main

import (
	"flag"
)

func main() {
	var folder string
	var email string
	flag.StringVar(&folder, "add-dir", "", "Directory to be added for include in stats")
	flag.StringVar(&email, "email", "", "User Email to be searched for in the commits")
	flag.Parse()

	if folder != "" {
		Scan(folder)
	}
	stats(email)
}
