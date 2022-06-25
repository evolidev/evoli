package use

import (
	"regexp"
	"strings"
)

type Str struct {
	value string

	strings.Replacer
	strings.Builder
	strings.Reader
}

func String(s string) *Str {
	return &Str{value: s}
}

func (s *Str) Kebab() *Str {
	s.value = s.toCase(s.value, "-")
	return s
}

func (s *Str) Snake() *Str {
	s.value = s.toCase(s.value, "_")
	return s
}

func (s *Str) Get() string {
	return s.value
}

func (s *Str) toCase(str string, delimiter string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z\\d])([A-Z])")

	snake := strings.ReplaceAll(str, " ", delimiter)
	snake = matchFirstCap.ReplaceAllString(snake, "${1}"+delimiter+"${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}"+delimiter+"${2}")

	return strings.ToLower(snake)
}
