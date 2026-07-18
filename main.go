package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"sync"

)



func main(){
	cliArg:= os.Args[1:]
	if len(cliArg)<1{
		log.Fatal("no website provided")
	}
	if len(cliArg)>1{
		log.Fatal("too many arguments provided")
	}
	baseUrl,err:=url.Parse(cliArg[0])
	if err!=nil{
		log.Fatal(err)
	}
	
	cfg:= config{
		pages: make(map[string]PageData),
		baseURL: baseUrl,
		mu: &sync.Mutex{},
		concurrencyControl: make(chan struct{},25),
		wg: &sync.WaitGroup{},
		maxPages: 3,
	}

	cfg.wg.Add(1)
	go cfg.crawlPage(baseUrl.String())
	cfg.wg.Wait()

	
	for normalizedURL, data := range cfg.pages {
		fmt.Printf("%v - %v\n", data.Heading, normalizedURL)
	}

}