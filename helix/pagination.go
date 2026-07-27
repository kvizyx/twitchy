package helix

import (
	"fmt"
	"reflect"
	"strings"
)

type paginationPlan struct {
	cursorParameter string
}

type paginationRequestError struct {
	reason string
}

func (e *paginationRequestError) Error() string {
	return "helix: pagination request: " + e.reason
}

func newPaginationPlan(request any, defaultCursorParameter string) (paginationPlan, error) {
	if defaultCursorParameter != "after" && defaultCursorParameter != "before" {
		return paginationPlan{}, &paginationRequestError{reason: fmt.Sprintf("unsupported cursor parameter %q", defaultCursorParameter)}
	}
	afterSet, err := paginationFieldSet(request, "after")
	if err != nil {
		return paginationPlan{}, err
	}
	beforeSet, err := paginationFieldSet(request, "before")
	if err != nil {
		return paginationPlan{}, err
	}
	if afterSet && beforeSet {
		return paginationPlan{}, &ExclusiveParametersError{First: "after", Second: "before"}
	}
	if beforeSet {
		defaultCursorParameter = "before"
	}
	return paginationPlan{cursorParameter: defaultCursorParameter}, nil
}

func (p paginationPlan) withCursor(request any, cursor string) (any, error) {
	value, structValue, pointer, err := cloneRequestValue(request)
	if err != nil {
		return nil, err
	}
	fieldIndex := paginationFieldIndex(structValue, p.cursorParameter)
	if fieldIndex < 0 {
		return nil, &paginationRequestError{reason: fmt.Sprintf("request has no %q cursor field", p.cursorParameter)}
	}
	if err := setCursorValue(structValue.Field(fieldIndex), cursor); err != nil {
		return nil, err
	}
	if pointer {
		return value.Interface(), nil
	}
	return structValue.Interface(), nil
}

func cloneRequestValue(request any) (reflect.Value, reflect.Value, bool, error) {
	value := reflect.ValueOf(request)
	if !value.IsValid() {
		return reflect.Value{}, reflect.Value{}, false, &paginationRequestError{reason: "request must be a struct"}
	}
	if value.Kind() == reflect.Pointer {
		if value.IsNil() || value.Elem().Kind() != reflect.Struct {
			return reflect.Value{}, reflect.Value{}, false, &paginationRequestError{reason: "request must point to a struct"}
		}
		clone := reflect.New(value.Elem().Type())
		clone.Elem().Set(value.Elem())
		return clone, clone.Elem(), true, nil
	}
	if value.Kind() != reflect.Struct {
		return reflect.Value{}, reflect.Value{}, false, &paginationRequestError{reason: "request must be a struct"}
	}
	clone := reflect.New(value.Type()).Elem()
	clone.Set(value)
	return clone, clone, false, nil
}

func paginationFieldSet(request any, parameter string) (bool, error) {
	_, value, _, err := cloneRequestValue(request)
	if err != nil {
		return false, err
	}
	index := paginationFieldIndex(value, parameter)
	if index < 0 {
		return false, nil
	}
	return !value.Field(index).IsZero(), nil
}

func paginationFieldIndex(request reflect.Value, parameter string) int {
	typeOfRequest := request.Type()
	for index := 0; index < request.NumField(); index++ {
		field := typeOfRequest.Field(index)
		if field.PkgPath != "" {
			continue
		}
		for _, tagName := range []string{"query", "url", "json"} {
			tag := strings.Split(field.Tag.Get(tagName), ",")[0]
			if tag == parameter {
				return index
			}
		}
		if strings.EqualFold(field.Name, parameter) {
			return index
		}
	}
	return -1
}

func setCursorValue(field reflect.Value, cursor string) error {
	if field.Kind() == reflect.String {
		field.SetString(cursor)
		return nil
	}
	if field.Kind() == reflect.Pointer && field.Type().Elem().Kind() == reflect.String {
		value := reflect.New(field.Type().Elem())
		value.Elem().SetString(cursor)
		field.Set(value)
		return nil
	}
	return &paginationRequestError{reason: fmt.Sprintf("cursor field %s is not string-compatible", field.Type())}
}
