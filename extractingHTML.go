package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getHeadingFromHTML(html string) (string,error){

	doc,err:=goquery.NewDocumentFromReader(strings.NewReader(html))
	if err!=nil{
		return "",err
	}
	h1 := doc.Find("h1, h2").First().Text()
	return strings.TrimSpace(h1),nil
}

func getFirstParagraphFromHTML(html string) (string,error){
	doc,err:=goquery.NewDocumentFromReader(strings.NewReader(html))
	if err!=nil{
		return "",err
	}
	fpargraph:=""
	selector:=doc.Find("main p")
	if selector.Text()==""{
		fpargraph=doc.Find("p").First().Text()
	}else{
		fpargraph =selector.First().Text()
	}
	return strings.TrimSpace(fpargraph),nil
}

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	
	doc,err:=goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err!=nil{
		return []string{},err
	}
	var urls []string
	
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href,ok :=s.Attr("href")
		if !ok {
			return
		}
		href = strings.TrimSpace(href)
		if href == "" {
			return
		}
		
		u, err := url.Parse(href)
		if err != nil {
			fmt.Printf("couldn't parse src %q: %v\n", href, err)
			return 
		}
		absolute:=baseURL.ResolveReference(u)
		urls=append(urls,absolute.String())
    })

	return urls,nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc,err:=goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err!=nil{
		return []string{},err
	}
	var images []string
	
	doc.Find("img").Each(func(_ int, s *goquery.Selection){
		src, ok := s.Attr("src")
		if !ok || strings.TrimSpace(src) == "" {
			return
		}
		// relative
		u, err := url.Parse(src)
		if err != nil {
			fmt.Printf("couldn't parse src %q: %v\n", src, err)
			return 
		}

		absolute := baseURL.ResolveReference(u)
		images = append(images, absolute.String())
		

    })
	return images,nil
}


