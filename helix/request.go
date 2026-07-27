package helix

import (
	"bytes"
	"context"
	"encoding"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

type requestSpec struct {
	Context context.Context
	Method  string
	URL     string
	Query   any
	Body    any
	Form    bool
	Codec   JSONCodec
}

func buildRequest(spec requestSpec) (*http.Request, error) {
	if spec.Context == nil || spec.Method == "" || spec.URL == "" {
		return nil, &RequestEncodingError{Reason: "method and URL are required"}
	}
	target, err := url.Parse(spec.URL)
	if err != nil {
		return nil, fmt.Errorf("parse request URL: %w", err)
	}
	values, err := encodeQuery(spec.Query)
	if err != nil {
		return nil, fmt.Errorf("encode query: %w", err)
	}
	query := target.Query()
	for key, items := range values {
		for _, item := range items {
			query.Add(key, item)
		}
	}
	target.RawQuery = query.Encode()

	body, contentType, err := encodeRequestBody(spec.Method, spec.Body, spec.Form, spec.Codec)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(spec.Context, spec.Method, target.String(), body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}
	return request, nil
}

func encodeRequestBody(method string, value any, form bool, codec JSONCodec) (io.Reader, string, error) {
	if isNilValue(value) {
		return nil, "", nil
	}
	if method != http.MethodPost && method != http.MethodPatch && method != http.MethodPut {
		return nil, "", &RequestEncodingError{Reason: fmt.Sprintf("%s cannot have a request body", method)}
	}
	if form {
		body, err := encodeFormBody(value)
		if err != nil {
			return nil, "", fmt.Errorf("encode form body: %w", err)
		}
		return bytes.NewReader(body), "application/x-www-form-urlencoded", nil
	}
	body, err := encodeJSONBodyWithCodec(value, codec)
	if err != nil {
		return nil, "", fmt.Errorf("encode JSON body: %w", err)
	}
	return bytes.NewReader(body), "application/json", nil
}

func isNilValue(value any) bool {
	if value == nil {
		return true
	}
	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return reflected.IsNil()
	default:
		return false
	}
}

func encodeQuery(value any) (url.Values, error) {
	return encodeValues(value, "query")
}

func encodeFormBody(value any) ([]byte, error) {
	values, err := encodeValues(value, "form")
	if err != nil {
		return nil, err
	}
	return []byte(values.Encode()), nil
}

func encodeValues(value any, kind string) (url.Values, error) {
	values := make(url.Values)
	if isNilValue(value) {
		return values, nil
	}
	current := reflect.ValueOf(value)
	for current.Kind() == reflect.Pointer || current.Kind() == reflect.Interface {
		if current.IsNil() {
			return values, nil
		}
		current = current.Elem()
	}
	if current.Kind() != reflect.Struct {
		return nil, &RequestEncodingError{Reason: fmt.Sprintf("%s values must be a struct", kind)}
	}
	seen := make(map[string]bool)
	typeOfValue := current.Type()
	for index := 0; index < current.NumField(); index++ {
		field := typeOfValue.Field(index)
		if field.PkgPath != "" {
			continue
		}
		name, options, err := fieldName(field, kind)
		if err != nil {
			return nil, err
		}
		if name == "-" {
			continue
		}
		set, err := appendValues(values, name, current.Field(index), options)
		if err != nil {
			return nil, fmt.Errorf("encode %s field %s: %w", kind, field.Name, err)
		}
		if set {
			seen[name] = true
		}
	}
	if seen["after"] && seen["before"] {
		return nil, &ExclusiveParametersError{First: "after", Second: "before"}
	}
	return values, nil
}

func fieldName(field reflect.StructField, kind string) (string, map[string]bool, error) {
	tag := field.Tag.Get(kind)
	if tag == "" && kind == "query" {
		tag = field.Tag.Get("url")
	}
	if tag == "" {
		tag = field.Tag.Get("json")
	}
	if tag == "" {
		return field.Name, map[string]bool{}, nil
	}
	parts := strings.Split(tag, ",")
	options := make(map[string]bool, len(parts)-1)
	for _, option := range parts[1:] {
		if option != "omitempty" && option != "" {
			return "", nil, &UnsupportedTagError{Field: field.Name, Tag: option}
		}
		options[option] = true
	}
	return parts[0], options, nil
}

func appendValues(values url.Values, name string, item reflect.Value, options map[string]bool) (bool, error) {
	for item.Kind() == reflect.Interface {
		if item.IsNil() {
			return false, nil
		}
		item = item.Elem()
	}
	if item.Kind() == reflect.Pointer {
		if item.IsNil() {
			return false, nil
		}
		item = item.Elem()
	} else if options["omitempty"] && item.IsZero() {
		return false, nil
	}
	if item.Kind() == reflect.Slice || item.Kind() == reflect.Array {
		set := false
		for index := 0; index < item.Len(); index++ {
			if item.Index(index).Kind() == reflect.Pointer && item.Index(index).IsNil() {
				continue
			}
			encoded, err := scalarValue(item.Index(index))
			if err != nil {
				return false, err
			}
			values.Add(name, encoded)
			set = true
		}
		return set, nil
	}
	encoded, err := scalarValue(item)
	if err != nil {
		return false, err
	}
	values.Add(name, encoded)
	return true, nil
}

func scalarValue(value reflect.Value) (string, error) {
	for value.Kind() == reflect.Pointer || value.Kind() == reflect.Interface {
		if value.IsNil() {
			return "", &RequestEncodingError{Reason: "nil scalar value"}
		}
		value = value.Elem()
	}
	if value.CanInterface() {
		if marshaler, ok := value.Interface().(encoding.TextMarshaler); ok {
			encoded, err := marshaler.MarshalText()
			if err != nil {
				return "", fmt.Errorf("marshal text: %w", err)
			}
			return string(encoded), nil
		}
	}
	switch value.Kind() {
	case reflect.String:
		return value.String(), nil
	case reflect.Bool:
		return strconv.FormatBool(value.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(value.Uint(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'g', -1, value.Type().Bits()), nil
	default:
		return "", &RequestEncodingError{Reason: fmt.Sprintf("unsupported scalar type %s", value.Type())}
	}
}
