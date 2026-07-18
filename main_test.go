package main

import (
	"net/url"
	"reflect"
	"strings"
	"testing"
)


func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name          string
		inputURL      string
		expected      string
		errorContains string
	}{
		{
			name:     "remove scheme",
			inputURL: "https://crawler-test.com/path",
			expected: "crawler-test.com/path",
		},
		{
			name:     "remove trailing slash",
			inputURL: "https://crawler-test.com/path/",
			expected: "crawler-test.com/path",
		},
		{
			name:     "lowercase capital letters",
			inputURL: "https://CRAWLER-TEST.com/PATH",
			expected: "crawler-test.com/path",
		},
		{
			name:     "remove scheme and capitals and trailing slash",
			inputURL: "http://CRAWLER-TEST.com/path/",
			expected: "crawler-test.com/path",
		},
		{
			name:          "handle invalid URL",
			inputURL:      `:\\invalidURL`,
			expected:      "",
			errorContains: "couldn't parse URL",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil && !strings.Contains(err.Error(), tc.errorContains) {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			} else if err != nil && tc.errorContains == "" {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			} else if err == nil && tc.errorContains != "" {
				t.Errorf("Test %v - '%s' FAIL: expected error containing '%v', got none.", i, tc.name, tc.errorContains)
				return
			}

			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetHeadingFromHTML(t *testing.T){
	tests := []struct {
		name          string
		inputHTML      string
		expected      string
		errorContains string
	}{
		{
			name :"Extraction h1 content",
			inputHTML: `<html>
			<body>
				<h1>Welcome to Boot.dev</h1>
				<main>
				<p>Learn to code by building real projects.</p>
				<p>This is the second paragraph.</p>
				</main>
			</body>
			</html>`,
			expected: "Welcome to Boot.dev",
		},
		{
			name :"h2 content as fallback",
			inputHTML: `<html>
			<body>
				<h2>Welcome to Boot.dev</h2>
				<main>
				<p>Learn to code by building real projects.</p>
				<p>This is the second paragraph.</p>
				</main>
			</body>
			</html>`,
			expected: "Welcome to Boot.dev",
		},
		{
			name :"No H1 present",
			inputHTML: `<html>
			<body>
				<p>Welcome to Boot.dev</p>
				<main>
				<p>Learn to code by building real projects.</p>
				<p>This is the second paragraph.</p>
				</main>
			</body>
			</html>`,
			expected: "",
		},
	}
	for i,test:=range tests{
		t.Run(test.name, func(t *testing.T) {
			heading,err:=getHeadingFromHTML(test.inputHTML)
			if err!=nil && !strings.Contains(err.Error(),test.errorContains){
				t.Errorf("Test %v - %s FAIL: expected error: %v, actual: %v", i, test.name, test.errorContains, heading)
			} else if heading!=test.expected{
				t.Errorf("Test %v - %s FAIL: expected: %v, actual: %v", i, test.name, test.expected, heading)
			}
		})
	}
}

func TestGetFirstParagraphFromHTMLMainPriority(t *testing.T){
	tests := []struct {
		name          string
		inputHTML      string
		expected      string
		errorContains string
	}{
		{
			name :"Extraction p in the main content",
			inputHTML: `<html>
			<body>
				<h1>Welcome to Boot.dev</h1>
				<main>
				<p>Learn to code by building real projects.</p>
				<p>This is the second paragraph.</p>
				</main>
			</body>
			</html>`,
			expected: "Learn to code by building real projects.",
		},
		{
			name :"main section is absent",
			inputHTML: `<html>
			<body>
				<h2>Welcome to Boot.dev</h2>
				<div>
				<p>Learn to code by building real projects.</p>
				<p>This is the second paragraph.</p>
				</div>
			</body>
			</html>`,
			expected: "Learn to code by building real projects.",
		},
		{
			name :"p tag is outside main",
			inputHTML: `<html>
			<body>
				<p>Welcome to Boot.dev</p>
				<main>
				<p>Learn to code by building real projects.</p>
				<p>This is the second paragraph.</p>
				</main>
			</body>
			</html>`,
			expected: "Learn to code by building real projects.",
		},
	}
	for i,test:=range tests{
		t.Run(test.name, func(t *testing.T) {
			heading,err:=getFirstParagraphFromHTML(test.inputHTML)
			if err!=nil && !strings.Contains(err.Error(),test.errorContains){
				t.Errorf("Test %v - %s FAIL: expected error: %v, actual: %v", i, test.name, test.errorContains, heading)
			} else if heading!=test.expected{
				t.Errorf("Test %v - %s FAIL: expected: %v, actual: %v", i, test.name, test.expected, heading)
			}
		})
	}
}

func TestGetUrlsFromHTML(t *testing.T){
	tests := []struct {
		name          string
		input      struct{
			html string
			baseUrl *url.URL
		}
		expected     [] string
		errorContains string
	}{
		{
			name:"getting urls",
			input: struct{html string; baseUrl *url.URL}{
				html: `<html>
			<body>
				<a href="https://crawler-test.com">Welcome to Boot.dev</a>
		
				<a href="https://crawler-test.com/home">Learn to code by building real projects.</a>
				<a href="https://crawler-test.com/settings">This is the second paragraph.</a>
				
			</body>
			</html>`,
			baseUrl: &url.URL{Scheme: "https", Host: "crawler-test.com"},
			},
			expected: []string{
				"https://crawler-test.com",
				"https://crawler-test.com/home",
				"https://crawler-test.com/settings",
			},
		},
		{
			name:"converting relative urls",
			input: struct{html string; baseUrl *url.URL}{
				html: `<html>
			<body>
				<a href="/test"><span>Boot.dev</span></a></body>
		
				<a href="/test1">Learn to code by building real projects.</a>
				<a href="/test2">This is the second paragraph.</a>
				
			</body>
			</html>`,
			baseUrl: &url.URL{Scheme: "https", Host: "crawler-test.com"},
			},
			expected: []string{
				"https://crawler-test.com/test",
				"https://crawler-test.com/test1",
				"https://crawler-test.com/test2",
			},
		},
	}

	for i,test:=range tests{
		t.Run(test.name, func(t *testing.T) {
			urls,err:=getURLsFromHTML(test.input.html,test.input.baseUrl)
			if err!=nil && !strings.Contains(err.Error(),test.errorContains){
				t.Errorf("Test %v - %s FAIL: expected error: %v, actual: %v", i, test.name, test.errorContains,urls)
			}else if !reflect.DeepEqual(test.expected,urls){
				t.Errorf("name: %v, expected %v, got %v", test.name, test.expected,urls)
			}
		})
	}
}

func TestGetImagesFromHTML(t *testing.T){
	tests := []struct {
		name          string
		input      struct{
			html string
			baseUrl *url.URL
		}
		expected     [] string
		errorContains string
	}{
		{
			name:"getting image",
			input: struct{html string; baseUrl *url.URL}{
				html: `<html>
			<body>
				<img src="/logo.png" alt="Logo">
		
				<a href="https://crawler-test.com/home">Learn to code by building real projects.</a>
				<a href="https://crawler-test.com/settings">This is the second paragraph.</a>
				
			</body>
			</html>`,
			baseUrl: &url.URL{Scheme: "https", Host: "crawler-test.com"},
			},
			expected: []string{
				"https://crawler-test.com/logo.png",
			},
		},
		{
			name:"attribute is missing",
			input: struct{html string; baseUrl *url.URL}{
				html: `<html>
			<body>
				<img alt="Logo">
		
				<a href="/test1">Learn to code by building real projects.</a>
				<a href="/test2">This is the second paragraph.</a>
				
			</body>
			</html>`,
			baseUrl: &url.URL{Scheme: "https", Host: "crawler-test.com"},
			},
			expected: nil,
		},
		{
			name:"not images present",
			input: struct{html string; baseUrl *url.URL}{
				html: `<html>
			<body>
				
		
				<a href="/test1">Learn to code by building real projects.</a>
				<a href="/test2">This is the second paragraph.</a>
				
			</body>
			</html>`,
			baseUrl: &url.URL{Scheme: "https", Host: "crawler-test.com"},
			},
			expected: nil,
		},
	}

	for i,test:=range tests{
		t.Run(test.name, func(t *testing.T) {
			images,err:=getImagesFromHTML(test.input.html,test.input.baseUrl)
			if err!=nil && !strings.Contains(err.Error(),test.errorContains){
				t.Errorf("Test %v - %s FAIL: expected error: %v, actual: %v", i, test.name, test.errorContains,images)
			}else if !reflect.DeepEqual(test.expected,images){
				t.Errorf("name: %v, expected %v, got %v", test.name, test.expected,images)
			}
		})
	}
}
