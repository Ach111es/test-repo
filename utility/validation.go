package utility

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

func ValidationErrorHandle(err error) string {
	var message string

	castedObject, ok := err.(validator.ValidationErrors)
	if ok {
		for _, v := range castedObject {
			switch v.Tag() {
			case "required":
				message = fmt.Sprintf("%s required", v.Field())
			case "min":
				message = fmt.Sprintf("%s input value must be greater than %s character", v.Field(), v.Param())
			case "max":
				message = fmt.Sprintf("%s input value must be lower than %s character", v.Field(), v.Param())
			case "lte":
				message = fmt.Sprintf("%s input value must be below %s", v.Field(), v.Param())
			case "gte":
				message = fmt.Sprintf("%s input value must be above %s", v.Field(), v.Param())
			case "numeric":
				message = fmt.Sprintf("%s input value must be numeric", v.Field())
			case "url":
				message = fmt.Sprintf("%s input value must be am url", v.Field())
			case "email":
				message = fmt.Sprintf("%s input value must be an email", v.Field())
			case "password":
				message = fmt.Sprintf("%s input value must be filled", v.Field())
			}
		}
	}

	return strings.ToLower(message)
}

func TimeParse(strCreatedAt string, config Configuration) (time.Time, error) {

	loc, err := time.LoadLocation(config.Timezone.Location)
	if err != nil {
		log.Println(err)
		return time.Time{}, errors.New("load timezone error")
	}

	parsedTime, err := time.ParseInLocation(config.Timezone.Format, strCreatedAt, loc)
	if err != nil {
		log.Println(err)
		return time.Time{}, errors.New("parsing time : invalid time format")
	}

	return parsedTime, nil

}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func NewNullInt(s int64) sql.NullInt64 {
	if s == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: s,
		Valid: true,
	}
}

func NullIfEmpty(value interface{}) string {
	switch v := value.(type) {
	case int:
		if v == 0 {
			return "NULL"
		}
		return fmt.Sprintf("%d", v)
	case string:
		if v == "" || v == "[]" {
			return "NULL"
		}
		return fmt.Sprintf("'%s'", v)
	default:
		return "NULL"
	}
}
