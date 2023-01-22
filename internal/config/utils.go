package config

import (
	"fmt"
)

// GetConfigValueAsString getting a value as a string, if possible
func GetConfigValueAsString(properties map[string]any, key string) (string, error) {
	if _, ok := properties[key]; !ok {
		return "", fmt.Errorf("missing config value for %s", key)
	}
	if properties[key] == nil {
		return "", nil
	}
	value, ok := properties[key].(string)
	if !ok {
		return "", fmt.Errorf("config value for %s is not a string", key)
	}
	return value, nil
}

// GetConfigValueAsBool getting a value as a bool, if possible
func GetConfigValueAsBool(properties map[string]any, key string) (bool, error) {
	if _, ok := properties[key]; !ok {
		return false, fmt.Errorf("missing config value for %s", key)
	}
	if properties[key] == nil {
		return false, nil
	}
	value, ok := properties[key].(bool)
	if !ok {
		return false, fmt.Errorf("config value for %s is not a string", key)
	}
	return value, nil
}

// GetConfigValueAsInt getting a value as a int64, if possible
func GetConfigValueAsInt(properties map[string]any, key string) (int64, error) {
	if _, ok := properties[key]; !ok {
		return 0, fmt.Errorf("missing config value for %s", key)
	}
	if properties[key] == nil {
		return 0, nil
	}
	var value int64
	switch v := properties[key].(type) {
	case int:
		value = int64(v)
	case int64:
		value = v
	default:
		return 0, fmt.Errorf("config value for %s is not a integer", key)
	}
	return value, nil
}
