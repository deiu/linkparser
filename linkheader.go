package lh

import (
	"regexp"
	"strings"
)

var (
	linkexp  = regexp.MustCompile(`<[^>]*>\s*(\s*;\s*[^\(\)<>@,;:"\/\[\]\?={} \t]+=(([^\(\)<>@,;:"\/\[\]\?={} \t]+)|("[^"]*")))*(,|$)`)
	paramexp = regexp.MustCompile(`[^\(\)<>@,;:"\/\[\]\?={} \t]+=(([^\(\)<>@,;:"\/\[\]\?={} \t]+)|("[^"]*"))`)
)

func unquote(s string) string {
	if len(s) < 2 {
		return s
	}
	if s[0] != '"' {
		return s
	}
	if s[len(s)-1] != '"' {
		return s
	}
	return s[1 : len(s)-1]
}

func brack(s string) string {
	if len(s) > 0 && s[0] == '<' {
		return s
	}
	if len(s) > 0 && s[len(s)-1] == '>' {
		return s
	}
	return "<" + s + ">"
}

func debrack(s string) string {
	if s[0] != '<' {
		return s
	}
	if len(s) == 1 {
		return s[1:]
	}
	s = s[1:len(s)]

	if s[len(s)-1] != '>' {
		return s
	}
	return s[:len(s)-1]
}

// ParseHeader takes a the value of a Link header (e.g. usually using req.Header.Get("Link")) and returns a map where link values are keyed based on `rel` values.
func ParseHeader(header string) map[string]string {
	links := map[string]string{}
	if len(header) == 0 {
		return links
	}
	matches := linkexp.FindAllString(header, -1)
	for i := range matches {
		match := matches[i]
		parts := strings.Split(match, ">")
		if len(parts) == 2 {
			href := debrack(parts[0])
			rels := strings.Split(paramexp.FindString(parts[1]), "=")
			if len(rels) > 1 {
				rel := unquote(rels[1])
				if len(rel) > 0 {
					links[rel] = href
				}
			}
		}
	}
	return links
}

// AddLink returns a Link header with multiple values by adding new links to an existing Link header (which can also be empty). It only supports two parameters, the <link> and the rel="" value.
func AddLink(oldHeader, link, rel string) string {
	if len(link) == 0 || len(rel) == 0 {
		return oldHeader
	}
	links := ParseHeader(oldHeader)
	links[rel] = link
	newHeader := ""
	for k, v := range links {
		if len(newHeader) > 0 {
			newHeader += ", "
		}
		newHeader += `<` + v + `>; rel="` + k + `"`
	}
	return newHeader
}