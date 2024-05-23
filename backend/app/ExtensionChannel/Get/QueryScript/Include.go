package QueryScript

import (
	"regexp"
	"simple_account/app/Http/Url"
	"strings"
)

func Includes(url *Url.Url) string {
	query := url.GetQuery()
	rawInclude, exist := query["include"]
	if !exist {
		return ""
	}
	includesStr := ""
	if includes, isArray := rawInclude.([]string); isArray {
		includesStr = strings.Join(upperToDot(includes), ",")
	}

	return includesStr
}

func upperToDot(values []string) []string {
	re := regexp.MustCompile(`[A-Z]`)

	for i, value := range values {

		values[i] = re.ReplaceAllStringFunc(value, func(char string) string {
			return "." + strings.ToLower(char)
		})

	}

	return values
}
