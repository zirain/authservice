// Copyright 2024 Tetrate
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"encoding/base64"
	"strings"
)

// GetPathQueryFragment splits the given path into path, query, and fragment.
// See https://tools.ietf.org/html/rfc3986#section-3.4 and https://tools.ietf.org/html/rfc3986#section-3.5 for more information.
func GetPathQueryFragment(fullPath string) (path string, query string, fragment string) {
	// inter and hash hold the index of the first `?` and `#` respectively
	// `?` must be present before `#` if both are present to consider the query
	var inter, hash int

	hash = strings.Index(fullPath, "#")
	if hash != -1 {
		inter = strings.Index(fullPath[:hash], "?")
	} else {
		inter = strings.Index(fullPath, "?")
	}

	switch {
	case inter != -1 && hash != -1:
		// both query and fragment defined
		path = fullPath[:inter]
		query = fullPath[inter+1 : hash]
		fragment = fullPath[hash+1:]
	case inter != -1:
		// only query defined
		path = fullPath[:inter]
		query = fullPath[inter+1:]
	case hash != -1:
		// only fragment defined
		path = fullPath[:hash]
		fragment = fullPath[hash+1:]
	default:
		// neither query nor fragment defined
		path = fullPath
	}

	return
}

// DecodeCookiesHeader parses the value of the Cookie header to find all the cookies set.
// It returns a map of name->value for all the found valid cookies.
func DecodeCookiesHeader(headerValue string) map[string]string {
	cookies := make(map[string]string, 0)
	for _, c := range strings.Split(headerValue, ";") {
		parts := strings.Split(strings.TrimSpace(c), "=")
		if len(parts) != 2 {
			// invalid cookie it must be Name=Value
			continue
		}
		cookies[parts[0]] = parts[1]
	}
	return cookies
}

// EncodeCookieHeader builds the value of the Set-Cookie header from the given cookie name, value and directives.
func EncodeCookieHeader(name string, value string, directives []string) string {
	b := strings.Builder{}
	_, _ = b.WriteString(name + "=" + value)
	for _, directive := range directives {
		_, _ = b.WriteString("; " + directive)
	}
	return b.String()
}

func BasicAuthHeader(id string, secret string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(id+":"+secret))
}
