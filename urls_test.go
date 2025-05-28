package main

import (
	"reflect"
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
		wantErr  bool
	}{
		{
			name:     "remove scheme https",
			inputURL: "https://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove trailing forward slash",
			inputURL: "https://blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "lowercase capital letters",
			inputURL: "https://BLOG.boot.dev/PATH",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove scheme and capitals and trailing slash",
			inputURL: "http://BLOG.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "handle invalid URL",
			inputURL: `:\\invalidURL`,
			expected: "",
			wantErr:  true,
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if (err != nil) != tc.wantErr {
				t.Errorf("Test %v - FAIL: unexpected error: %v", i, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - FAIL: expected URL: %v, actual: %v", i, tc.expected, actual)
			}
		})
	}
}

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
		wantErr   bool
	}{
		{
			name:     "absolute and relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
	<body>
		<a href="/path/one">
			<span>Boot.dev</span>
		</a>
		<a href="https://other.com/path/one">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
		{
			name:     "absolute URL",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
	<body>
		<a href="https://blog.boot.dev">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`,
			expected: []string{"https://blog.boot.dev"},
		},
		{
			name:     "relative URL",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
	<body>
		<a href="/path/one">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`,
			expected: []string{"https://blog.boot.dev/path/one"},
		},
		{
			name:     "no href",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
	<body>
		<a>
			<span>Boot.dev></span>
		</a>
	</body>
</html>
`,
			expected: make([]string, 0),
		},
		{
			name:     "bad HTML",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html body>
	<a href="path/one">
		<span>Boot.dev></span>
	</a>
</html body>
`,
			expected: []string{"https://blog.boot.dev/path/one"},
		},
		{
			name:     "invalid href URL",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
	<body>
		<a href=":\\invalidURL">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`,
			expected: make([]string, 0),
		},
		{
			name:     "handle invalid base URL",
			inputURL: `:\\invalidBaseURL`,
			inputBody: `
<html>
	<body>
		<a href="/path">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`,
			expected: nil,
			wantErr:  true,
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(tc.inputBody, tc.inputURL)
			if (err != nil) != tc.wantErr {
				t.Errorf("Test %v - FAIL: unexpected error: %v", i, err)
				return
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - FAIL: expected: %v, actual: %v", i, tc.expected, actual)
			}
		})
	}
}
