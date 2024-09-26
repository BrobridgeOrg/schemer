package schemer

import (
	"strconv"
	"strings"
)

func parsePathEntry(entry string) (string, int) {

	split := strings.Split(entry, "[")

	key := split[0]
	index := -1

	// Check if we have an array index
	if len(split) > 1 {
		indexStr := strings.TrimRight(split[1], "]")
		idx, err := strconv.Atoi(indexStr)
		if err != nil {
			return key, -1
		}

		if idx < 0 {
			return key, -1
		}

		index = idx
	}

	return key, index
}
