package assembler

import (
	"strings"
)

const (
	DEFAULT_SEGMENT = ""
)

func getSegment(str string) (string, bool) {
	if strings.HasPrefix(str, ".") {
		return str[1:], true
	}
	return str, false
}

func trimLine(str string) string {
	ind := strings.Index(str, "#")
	if ind != -1 {
		str = str[0:ind]
	}
	return strings.Trim(str, " \t")
}

func TrimSplitSegment(content []string) map[string][]string {
	var result map[string][]string = make(map[string][]string)
	currentSeg := DEFAULT_SEGMENT
	result[currentSeg] = make([]string, 0)
	for _, str := range content {
		str = trimLine(str)
		if len(str) == 0 {
			continue
		}

		if seg, ok := getSegment(str); ok {
			_, ok := result[seg]
			if !ok {
				result[seg] = make([]string, 0)
			}
			currentSeg = seg
			continue
		}
		val, _ := result[currentSeg]
		result[currentSeg] = append(val, str)
	}
	return result
}
