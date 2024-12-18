package utils

func NormalizeToSlice(input any) []string {
	switch v := input.(type) {
	case string:
		return []string{v}
	case []string:
		return v
	default:
		return []string{}
	}
}
