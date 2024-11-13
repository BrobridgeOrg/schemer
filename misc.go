package schemer

import (
	"strconv"
)

func parsePathEntry(entry string) (string, int) {

	if len(entry) == 0 {
		return entry, -1
	}

	if entry[len(entry)-1] != ']' {
		return entry, -1
	}

	index := -1
	//	split := strings.Split(entry, "[")

	for i := len(entry) - 1; i >= 0; i-- {
		if entry[i] == '[' {
			index = i
			break
		}
	}

	if index == -1 {
		return entry, -1
	}

	key := entry[:index]
	indexStr := entry[index+1 : len(entry)-1]
	idx, err := strconv.Atoi(indexStr)
	if err != nil {
		return key, -1
	}

	index = idx

	return key, index
}
