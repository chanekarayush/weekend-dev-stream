package main

import (
	"fmt"
	"flag"
)


func scan (folderPath string)  {
	fmt.Printf("Scanned %s:", folderPath)
}

func computeStats(stats int)  {
	print(stats)	
}

func main()  {
	folder := flag.String("add-dir", "", "Directory to be added for include in stats")
	flag.Parse();

	if (*folder != ""){
		scan(*folder)
	}
	fmt.Println(*folder)

}
