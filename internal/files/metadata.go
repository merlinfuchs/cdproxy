package files

import "net/url"

func EncodeMetadata(metadata map[string]string) string {
	if metadata == nil {
		return ""
	}

	res := ""
	for k, v := range metadata {
		// trailing & in after the last key-value pair is intended!
		res += k + "=" + url.QueryEscape(v) + "&"
	}

	return res
}

func DecodeMetadata(metadata string) (map[string]string, error) {
	values, err := url.ParseQuery(metadata)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	for k, v := range values {
		m[k] = v[0]
	}

	return m, nil
}
