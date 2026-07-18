package main

import (
	"fmt"
	
	"net/url"
	"sync"
)
type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages 			int
}

func (cfg *config)crawlPage(rawCurrentURL string){
	cfg.concurrencyControl <-struct{}{}
	defer func() { <-cfg.concurrencyControl 
		cfg.wg.Done()}() 
	currentURL,err:=url.Parse(rawCurrentURL)
	if err!=nil{
		fmt.Println(err)
		return
	}
	if cfg.baseURL.Hostname()!=currentURL.Hostname(){
		return
	}
	normalCurrentURL,err := normalizeURL(rawCurrentURL)
	if err!=nil{
		fmt.Println(err)
		return
	}
	isFirst:=cfg.addPageVisit(normalCurrentURL)
	if !isFirst{
		return
	}


	html,err:=getHTML(currentURL.String())
	if err!=nil{
		fmt.Println(err)
		return
	}
	
	cfg.AddPages(normalCurrentURL,html,currentURL.String())

	urls,err:=getURLsFromHTML(html,cfg.baseURL)
	if err!=nil{
		fmt.Println(err)
		return
	}
	//fmt.Println(urls)
	for _,u:=range urls{
		cfg.wg.Add(1)
		go cfg.crawlPage(u)
	}

}
func (cfg *config)AddPages(normalURL ,html,currentUrl string){
	cfg.mu.Lock()
	cfg.pages[normalURL] = extractPageData(html,currentUrl)
	defer cfg.mu.Unlock()
}
func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool){
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	if _, ok := cfg.pages[normalizedURL]; ok {
		return false
	}
	// Mark it immediately so no other goroutine attempts to crawl it
	cfg.pages[normalizedURL] = PageData{} 
	return true
}