package main

func DefaultValue(val, fallback string) string {
	if val == "" {
		return fallback
	}
	return val
}
