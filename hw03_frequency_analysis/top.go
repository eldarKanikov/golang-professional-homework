package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(input string) []string {
	wordCountMap := map[string]int{}
	words := strings.Fields(input)
	result := make([]string, 0, 10)

	for _, word := range words {
		wordCountMap[word]++
	}

	wordCountList := make([]keyValue, 0, len(wordCountMap))

	for word, count := range wordCountMap {
		wordCountList = append(wordCountList, keyValue{word, count})
	}

	sort.Slice(wordCountList, func(i, j int) bool {
		if wordCountList[i].value == wordCountList[j].value {
			return wordCountList[i].key < wordCountList[j].key
		}
		return wordCountList[i].value > wordCountList[j].value
	})

	for i, wordCount := range wordCountList {
		if i == 10 {
			break
		}
		result = append(result, wordCount.key)
	}
	return result
}

type keyValue struct {
	key   string
	value int
}
