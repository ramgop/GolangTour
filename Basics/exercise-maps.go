package main

import (
	"golang.org/x/tour/wc"
	"strings"
)

func WordCount(s string) map[string]int {
	var words []string
	words = strings.Split(s, " ")
	var wordlist = make(map[string]int)
	for i := 0; i < len(words); i++ {
		wordlist[words[i]]++
	}
	return wordlist
}

func main() {
	wc.Test(WordCount)
}
