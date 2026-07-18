package main

import (
	"fmt"
	"net/url"
)


type PageData struct{
	URL string
	Heading string
	FirstParagraph string
	OutgoingLinks []string
	ImageURLs []string
}


func extractPageData(html, pageURL string) PageData {
	parseURL,err:=url.Parse(pageURL)
	if err!=nil{
		return PageData{}
	}
	heading,err:= getHeadingFromHTML(html)
	if err!=nil{
		fmt.Println(err.Error())
	}
	fp,err:=getFirstParagraphFromHTML(html)
	if err!=nil{
		fmt.Println(err.Error())
	}
	url.Parse(pageURL)
	outlinks,err:= getURLsFromHTML(html,parseURL)
	if err!=nil{
		fmt.Println(err.Error())
	}
	images,err:= getImagesFromHTML(html,parseURL)
	if err!=nil{
		fmt.Println(err.Error())
	}

	return PageData{
		URL: pageURL,
		Heading: heading,
		FirstParagraph: fp,
		OutgoingLinks: outlinks,
		ImageURLs: images,
	}
}