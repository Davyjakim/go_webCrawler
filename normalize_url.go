package main

import (
	"fmt"
	"net/url"
	"strings"
)
func normalizeURL(inputURL string) (string, error){
	urlObj, err:=url.Parse(inputURL)
	if err!=nil{
		return "",fmt.Errorf("couldn't parse URL")
	}
	normalized := fmt.Sprintf("%v%v",urlObj.Host,urlObj.Path)
	normalized =strings.TrimSuffix(normalized, "/")
	normalized=strings.ToLower(normalized)

return normalized, nil
}
