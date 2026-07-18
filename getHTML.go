package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func getHTML(rawURL string) (string, error){
	
	client:=&http.Client{}
	req,err:=http.NewRequest("GET",rawURL,nil)
	if err!=nil{
		return "",err
	}
	req.Header.Set("User-Agent","Crawler/1.0")
	resp,err:=client.Do(req)
	if err!=nil{
		return "",err
	}
	defer resp.Body.Close()
	if resp.StatusCode>=400 && resp.StatusCode<500{
		return "", fmt.Errorf("Bad request")
	}
	if !strings.Contains(resp.Header.Get("Content-Type"),"text/html"){
		return "", fmt.Errorf("The content type is not text/html")
	}
	
	data,err := io.ReadAll(resp.Body)
	if err!=nil{
		return "",err
	}

	return string(data),nil
}