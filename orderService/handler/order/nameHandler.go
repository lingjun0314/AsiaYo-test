package order

import (
	"regexp"
	"strings"
)

type NameHandler struct{}

func (n *NameHandler) IsEnglish(name string) bool {
	return regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(name)
}

func (n *NameHandler) IsCapitalized(name string) bool {
	words := strings.Fields(name)
	for _, word := range words {
		if word[0] < 'A' || word[0] > 'Z' {
			return false
		}
	}
	return true
}
