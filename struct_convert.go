package assist

import (
	"encoding/json"
	"github.com/fatih/structs"
)

func init() {
	structs.DefaultTagName = "json"
}

func Map2Struct(m map[string]interface{}, s interface{}) error {
	// TODO refactor to more smart method :)
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, s)
}

func Struct2Map(s interface{}) map[string]interface{} {
	return structs.Map(s)
}


