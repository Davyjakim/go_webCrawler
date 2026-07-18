package main

import (
	"fmt"
	"log"
	"os"
)

func main(){
	cliArg:= os.Args[1:]
	if len(cliArg)<1{
		log.Fatal("no website provided")
	}
	if len(cliArg)>1{
		log.Fatal("too many arguments provided")
	}
	baseUrl:=cliArg[0]
	fmt.Println("starting crawl of:",baseUrl)
}