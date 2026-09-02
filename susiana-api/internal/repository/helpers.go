package repository

import "encoding/json"

func nullString(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func nullJSON(raw json.RawMessage) any {
	if len(raw) == 0 {
		return nil
	}
	return []byte(raw)
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
