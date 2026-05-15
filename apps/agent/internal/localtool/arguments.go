package localtool

import "fmt"

func stringArg(arguments map[string]any, name string) (string, error) {
	value, exists := arguments[name]
	if !exists {
		return "", fmt.Errorf("missing argument %q", name)
	}

	text, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("argument %q must be a string", name)
	}
	return text, nil
}

func numberArg(arguments map[string]any, name string) (float64, error) {
	value, exists := arguments[name]
	if !exists {
		return 0, fmt.Errorf("missing argument %q", name)
	}

	switch number := value.(type) {
	case float64:
		return number, nil
	case float32:
		return float64(number), nil
	case int:
		return float64(number), nil
	case int64:
		return float64(number), nil
	default:
		return 0, fmt.Errorf("argument %q must be a number", name)
	}
}
