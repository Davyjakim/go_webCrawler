package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)



func main(){
	cliArg:= os.Args[1:]
	maxConcurency:=20
	maxPage:=100
	if len(cliArg)<1{
		log.Fatal("no website provided")
	}
	if len(cliArg)>3{
		log.Fatal("too many arguments provided")
	}
	if len(cliArg)==2{
		num, err := strconv.Atoi(cliArg[1])
		if err!=nil{
			log.Fatal(err)
		}
		maxConcurency =num
	}
	if len(cliArg)==3{
		num, err := strconv.Atoi(cliArg[2])
		if err!=nil{
			log.Fatal(err)
		}
		maxPage =num
	}

	rawBaseURL := cliArg[0]
	cfg,err:= configure(rawBaseURL,maxPage,maxConcurency)
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Printf("starting crawl of: %s...\n", rawBaseURL)
	cfg.wg.Add(1)
	go cfg.crawlPage(rawBaseURL)
	cfg.wg.Wait()

	
	writeJSONReport(cfg.pages, "report.json")

}