package util

import (
	"regexp"

	"github.com/TwiN/go-away"
)

func CensorProfanity(input string) string {
	asciiOnly := filterNonASCII(input)
	return goaway.Censor(asciiOnly)
}

func ContainsProfanity(input ...string) bool {
	for _, s := range input {
		asciiOnly := filterNonASCII(s)
		if goaway.IsProfane(asciiOnly) {
			return true
		}
	}
	return false
}

func filterNonASCII(input string) string {
	reg := regexp.MustCompile("[[:^ascii:]]")
	return reg.ReplaceAllString(input, "")
}
