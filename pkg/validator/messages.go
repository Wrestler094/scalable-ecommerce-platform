package validator

import "fmt"

// messageForTag — генератор сообщений валидации по тэгу.
func messageForTag(tag string, param string) string {
	switch tag {
	case "required":
		return "field is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must be at least %s characters long", param)
	case "max":
		return fmt.Sprintf("must be at most %s characters long", param)
	case "len":
		return fmt.Sprintf("must be exactly %s characters long", param)
	case "gte":
		return fmt.Sprintf("must be greater than or equal to %s", param)
	case "lte":
		return fmt.Sprintf("must be less than or equal to %s", param)
	case "url":
		return "must be a valid URL"
	case "uuid":
		return "must be a valid UUID"
	case "eq":
		return fmt.Sprintf("must be equal to %s", param)
	case "ne":
		return fmt.Sprintf("must not be equal to %s", param)
	case "oneof":
		return fmt.Sprintf("must be one of: %s", param)
	default:
		return "invalid value"
	}
}
