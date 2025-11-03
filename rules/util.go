package rules

import (
	"fmt"
	"regexp"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
    snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
    snake  = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
    return strings.ToLower(snake)
}

func ReferenceLink(ruleName string) string {
	return fmt.Sprintf("https://github.com/utilitywarehouse/tflint-ruleset-uw/blob/main/rules/%s.md", toSnakeCase(ruleName))
}
