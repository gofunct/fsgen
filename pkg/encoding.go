package pkg

import (
	"encoding/csv"
	"encoding/json"
	"strings"
)

// toPrettyJson encodes an item into a pretty (indented) JSON string
func ToPrettyJsonString(v interface{}) string {
	output, _ := json.MarshalIndent(v, "", "  ")
	return string(output)
}

// toPrettyJson encodes an item into a pretty (indented) JSON string
func ToPrettyJson(v interface{}) []byte {
	output, _ := json.MarshalIndent(v, "", "  ")
	return output
}

func AsCSV(val string) ([]string, error) {
	if val == "" {
		return []string{}, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	return csvReader.Read()
}

func AsMap(val string) (map[string]string, error) {
	m := make(map[string]string)
	if val == "" {
		return m, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	arr, err := csvReader.Read()
	if err != nil {
		return m, err
	}
	for _, v := range arr {
		strings.TrimSpace(v)
		switch {
		case strings.Contains(v, "="):
			kv := strings.Split(v, "=")
			m[kv[0]] = kv[1]
		case strings.Contains(v, ":"):
			kv := strings.Split(v, ":")
			m[kv[0]] = kv[1]
		case strings.Contains(v, ":"):
			kv := strings.Split(v, ":")
			m[kv[0]] = kv[1]
		}
	}
	return m, nil
}
