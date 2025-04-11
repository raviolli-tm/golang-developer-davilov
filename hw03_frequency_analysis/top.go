package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type Pair struct {
	Key   string
	Value int
}

func Top(x string, taskWithAsteriskIsCompleted bool) []string {
	if !taskWithAsteriskIsCompleted {
		return Top10(x)
	}
	return Top10Asterisk(x)
}

func Top10(x string) []string {
	stringsArr := strings.Fields(x)
	if len(stringsArr) == 0 {
		return []string{}
	}
	sort.Strings(stringsArr)

	topFrequencies := make([]Pair, 0)
	start := 0
	for i := 0; i < len(stringsArr); i++ {
		if stringsArr[start] != stringsArr[i] {
			topFrequencies = append(topFrequencies, Pair{stringsArr[start], i - start})
			start = i
		}
	}
	topFrequencies = append(topFrequencies, Pair{stringsArr[start], len(stringsArr) - start})

	sort.Slice(topFrequencies, func(i, j int) bool {
		if topFrequencies[i].Value == topFrequencies[j].Value {
			return topFrequencies[i].Key < topFrequencies[j].Key
		}
		return topFrequencies[i].Value > topFrequencies[j].Value
	})
	var result []string

	for i := 0; i < len(topFrequencies); i++ {
		result = append(result, topFrequencies[i].Key)
	}
	return result[:10]
}

func Top10Asterisk(_ string) []string {
	return []string{}
}
