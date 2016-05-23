package util

import "encoding/json"

func JsonStringify(obj interface{}, intent bool) string {
	if intent {
		data, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			return ""
		}
		return string(data)
	} else {
		data, err := json.Marshal(obj)
		if err != nil {
			return ""
		}
		return string(data)
	}
}
