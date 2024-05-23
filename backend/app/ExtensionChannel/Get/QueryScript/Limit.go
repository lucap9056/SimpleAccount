package QueryScript

import (
	"fmt"
	"simple_account/app/Http/Message"
	"strconv"
)

func Limit(context *Message.Context) string {
	query := context.Url.GetQuery()
	offset := 0
	limit := 20

	if rawOffsets, ok := query["offset"]; ok {
		offsetsStr, isStrArray := rawOffsets.([]string)
		if isStrArray {
			offset64, err := strconv.ParseInt(offsetsStr[0], 10, 0)
			if err == nil {
				offset = int(offset64)
			}
		}
	}

	if rawLimits, ok := query["limit"]; ok {
		limitsStr, isStrArray := rawLimits.([]string)
		if isStrArray {
			limit64, err := strconv.ParseInt(limitsStr[0], 10, 0)
			if err == nil {
				limit = int(limit64)
			}
		}
	}

	if limit > 40 {
		limit = 40
	}

	return fmt.Sprintf(" LIMIT %d,%d", offset, limit)
}
