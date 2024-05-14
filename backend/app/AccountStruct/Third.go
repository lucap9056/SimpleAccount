package AccountStruct

type Third struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func (t *Third) Empty() bool {
	if t.Id != 0 {
		return false
	}
	if t.Name != "" {
		return false
	}
	return true
}

func (third *Third) MappingTable(args []string) []interface{} {
	tableMap := map[string]interface{}{
		"third.id":   &third.Id,
		"third.name": &third.Name,
	}

	var columns []interface{}
	for _, arg := range args {
		refer, exist := tableMap[arg]
		if !exist {
			var i interface{}
			refer = &i
		}
		columns = append(columns, refer)
	}
	return columns
}
