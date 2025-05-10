package utilities

import (
	"regexp"
	"strconv"
	"strings"
)

func IsUint64Type(argument string) bool {
	_, err := strconv.ParseUint(argument, 10, 64)
	return err == nil
}

func TransformMention(argument string) string {
	input := strings.TrimSpace(argument)
	regex := regexp.MustCompile(`^<@!?(\d+)>$`)
	if matches := regex.FindStringSubmatch(input); matches != nil {
		return matches[1]
	}
	return input
}
