package main

import (
	"encoding/json"
	"os"
	"sort"
)

func writeJSONReport(pages map[string]PageData, filename string) error{
	var keyslice []string
	for key :=range pages{
		keyslice=append(keyslice,key)
	}
	sort.Strings(keyslice)
	sorteddata:=make(map[string]PageData)
	for _,k :=range keyslice{
		sorteddata[k]=pages[k]
	}
	data,err:=json.MarshalIndent(sorteddata,""," ")
	if err!=nil{
		return err
	}
	os.WriteFile("report.json",data,0644)
	return nil
}