package main

import (
	"reflect"
	"testing"
)

func TestExtractPageDate(t *testing.T){
	cases:= []struct{
		name string
		inputURL string
		inputBody string
		expected PageData
	}{
		{
			name: "Extracting page data",
			inputURL : "https://crawler-test.com",
			inputBody :`<html><body>
				<h1>Test Title</h1>
				<p>This is the first paragraph.</p>
				<a href="/link1">Link 1</a>
				<img src="/image1.jpg" alt="Image 1">
			</body></html>`,
			expected : PageData{
				URL:             "https://crawler-test.com",
				Heading:         "Test Title",
				FirstParagraph: "This is the first paragraph.",
				OutgoingLinks:  []string{"https://crawler-test.com/link1"},
				ImageURLs:      []string{"https://crawler-test.com/image1.jpg"},
			},
		},
		{
			name: "Absent Heading anchor",
			inputURL : "https://crawler-test.com",
			inputBody :`<html><body>
				
				<p>This is the first paragraph.</p>
				
				<img src="/image1.jpg" alt="Image 1">
			</body></html>`,
			expected : PageData{
				URL:             "https://crawler-test.com",
				Heading:         "",
				FirstParagraph: "This is the first paragraph.",
				OutgoingLinks:  nil,
				ImageURLs:      []string{"https://crawler-test.com/image1.jpg"},
			},
		},
		{
			name: "input url is not the baseUrl",
			inputURL : "https://crawler-test.com/testing",
			inputBody :`<html><body>
				
				<p>This is the first paragraph.</p>
				
				<img src="/image1.jpg" alt="Image 1">
			</body></html>`,
			expected : PageData{
				URL:             "https://crawler-test.com/testing",
				Heading:         "",
				FirstParagraph: "This is the first paragraph.",
				OutgoingLinks:  nil,
				ImageURLs:      []string{"https://crawler-test.com/image1.jpg"},
			},
		},
	}
	for _,test :=range cases{
		t.Run(test.name,func(t *testing.T) {
			actual:=extractPageData(test.inputBody,test.inputURL)
			if !reflect.DeepEqual(actual,test.expected){
				t.Errorf("Test: %s ,expected %+v, got %+v",test.name, test.expected, actual)
			}
		})
	}
}