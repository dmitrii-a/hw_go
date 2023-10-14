package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var pattern = regexp.MustCompile(`([^!,.:]*)`)

func Top10(s string) []string {
	data := strings.Fields(s)
	counter := map[string]int{}
	for _, word := range data {
		normalizeWord := pattern.FindString(strings.ToLower(word))
		if normalizeWord == "-" {
			continue
		}
		counter[normalizeWord]++
	}
	keys := make([]string, 0, len(counter))
	for k := range counter {
		keys = append(keys, k)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		if counter[keys[i]] == counter[keys[j]] {
			return keys[i] < keys[j]
		}
		return counter[keys[i]] > counter[keys[j]]
	})
	if len(keys) > 10 {
		return keys[:10]
	}
	return keys
}
