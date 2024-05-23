package Url

import (
	NetUrl "net/url"
	"regexp"
	"strings"
)

const args = "args"

type Url struct {
	origin   NetUrl.URL
	query    map[string]interface{}
	Dirs     []string
	dirIndex int
}

type Query map[string]interface{}

func New(origin *NetUrl.URL) Url {
	url := Url{}
	url.origin = *origin
	url.Dirs = strings.Split(url.origin.Path[1:], "/")
	url.dirIndex = -1
	url.queryInit()
	return url
}

func (url *Url) Shift() string {
	index := url.dirIndex + 1
	if index >= len(url.Dirs) {
		return ""
	}
	url.dirIndex = index
	return url.Dirs[index]
}

func (url *Url) queryInit() {
	url.query = make(map[string]interface{})

	query, _ := NetUrl.QueryUnescape(url.origin.RawQuery)

	if query != "" {
		for _, q := range strings.Split(query, "&") {
			equal := strings.LastIndex(q, "=")
			key := args
			value := q

			if equal > -1 {
				key = q[:equal]
				value = q[equal+1:]
			}

			queryValue(&url.query, key, value)
		}

	}
}

func queryValue(parent *map[string]interface{}, key string, value string) {
	parentRegexp := regexp.MustCompile(`.*\[.*\]`)
	//判斷是否有括弧
	if matched := parentRegexp.MatchString(key); matched {

		leftBracket := strings.Index(key, "[")
		rightBracket := strings.LastIndex(key, "]")
		parentKey := key[:leftBracket]
		subKey := key[leftBracket+1 : rightBracket]

		if container, exist := (*parent)[parentKey]; exist {
			if containerSlice, isArray := container.([]string); isArray {
				newContainer := make(map[string]interface{})
				newContainer[args] = containerSlice
				(*parent)[parentKey] = newContainer
				queryValue(&newContainer, subKey, value)
			}
			if containerMap, isArray := container.(map[string]interface{}); isArray {
				queryValue(&containerMap, subKey, value)
			}
		} else {
			newContainer := make(map[string]interface{})
			(*parent)[parentKey] = newContainer
			queryValue(&newContainer, subKey, value)
		}

	} else {
		if container, exist := (*parent)[key]; exist {
			if containerSlice, isArray := container.([]string); isArray {
				(*parent)[key] = append(containerSlice, value)
			}
			if containerMap, isArray := container.(map[string]interface{}); isArray {
				queryValue(&containerMap, args, value)
			}
		} else {
			(*parent)[key] = []string{value}
		}

	}

}

func (url *Url) GetQuery() map[string]interface{} {
	return url.query
}
