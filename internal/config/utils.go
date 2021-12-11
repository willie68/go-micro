package config

import (
	"fmt"
)

func GetConfigValueAsString(properties map[string]interface{}, key string) (string, error) {
	if _, ok := properties[key]; !ok {
		return "", fmt.Errorf("missing config value for %s", key)
	}
	value, ok := properties[key].(string)
	if !ok {
		return "", fmt.Errorf("config value for %s is not a string", key)
	}
	return value, nil
}

func GetConfigValueAsBool(properties map[string]interface{}, key string) (bool, error) {
	if _, ok := properties[key]; !ok {
		return false, fmt.Errorf("missing config value for %s", key)
	}
	value, ok := properties[key].(bool)
	if !ok {
		return false, fmt.Errorf("config value for %s is not a string", key)
	}
	return value, nil
}

func GetConfigValueAsInt(properties map[string]interface{}, key string) (int64, error) {
	if _, ok := properties[key]; !ok {
		return 0, fmt.Errorf("missing config value for %s", key)
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
