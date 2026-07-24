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
func configure(rawBaseURL string,maxPage,maxConcurrency int) (*config, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse base URL: %v", err)
	}

	return &config{
		pages:              make(map[string]PageData),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages: maxPage,
	}, nil
}

func (cfg *config)crawlPage(rawCurrentURL string){
	cfg.concurrencyControl <-struct{}{}
	defer func() { <-cfg.concurrencyControl 
		cfg.wg.Done()}() 
	currentURL,err:=url.Parse(rawCurrentURL)
	if err!=nil{
			fmt.Println(currentURL,err)
			return
	}
	if cfg.pagesLen()>=cfg.maxPages{
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
	fmt.Printf("crawling %s\n", rawCurrentURL)

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

func (cfg *config) pagesLen()int{
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	return len(cfg.pages)
}