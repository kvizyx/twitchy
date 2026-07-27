package helix

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	DefaultResponseBodyLimit = 1 << 20
	DefaultErrorExcerptLimit = 64 << 10
)

var ErrInvalidBodyLimit = errors.New("helix: invalid body limit")

type JSONCodec struct {
	Marshal   func(any) ([]byte, error)
	Unmarshal func([]byte, any) error
}

func standardJSONCodec() JSONCodec {
	return JSONCodec{Marshal: json.Marshal, Unmarshal: json.Unmarshal}
}

type BodyLimits struct {
	Response     int64
	ErrorExcerpt int64
}

func (limits BodyLimits) responseLimit() int64 {
	if limits.Response == 0 {
		return DefaultResponseBodyLimit
	}
	return limits.Response
}

func (limits BodyLimits) errorExcerptLimit() int64 {
	if limits.ErrorExcerpt == 0 {
		return DefaultErrorExcerptLimit
	}
	return limits.ErrorExcerpt
}

type DecodeOptions struct {
	Codec  JSONCodec
	Limits BodyLimits
}

type RequestEncodingError struct {
	Reason string
}

func (e *RequestEncodingError) Error() string { return "helix: request encoding: " + e.Reason }

type UnsupportedTagError struct {
	Field string
	Tag   string
}

func (e *UnsupportedTagError) Error() string {
	return fmt.Sprintf("helix: unsupported %s tag option on %s", e.Tag, e.Field)
}

type ExclusiveParametersError struct {
	First  string
	Second string
}

func (e *ExclusiveParametersError) Error() string {
	return fmt.Sprintf("helix: parameters %q and %q are mutually exclusive", e.First, e.Second)
}

type BodyLimitError struct {
	Limit int64
}

func (e *BodyLimitError) Error() string {
	return fmt.Sprintf("helix: response body exceeds %d bytes", e.Limit)
}

type JSONDecodeError struct {
	Err error
}

func (e *JSONDecodeError) Error() string { return "helix: decode JSON: " + e.Err.Error() }
func (e *JSONDecodeError) Unwrap() error { return e.Err }

func encodeJSONBody(value any) ([]byte, error) {
	return encodeJSONBodyWithCodec(value, standardJSONCodec())
}

func encodeJSONBodyWithCodec(value any, codec JSONCodec) ([]byte, error) {
	codec = withDefaultCodec(codec)
	encoded, err := codec.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("marshal JSON: %w", err)
	}
	return encoded, nil
}

func decodeResponse[T any](status int, body io.Reader, options DecodeOptions) (*Response[T], error) {
	var response Response[T]
	if status == 204 {
		return &response, nil
	}
	data, err := readBounded(body, options.Limits.responseLimit())
	if err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return &response, nil
	}
	codec := withDefaultCodec(options.Codec)
	if err := codec.Unmarshal(data, &response); err != nil {
		return nil, &JSONDecodeError{Err: err}
	}
	return &response, nil
}

func readBounded(body io.Reader, limit int64) ([]byte, error) {
	if limit < 1 {
		return nil, ErrInvalidBodyLimit
	}
	data, err := io.ReadAll(io.LimitReader(body, limit+1))
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	if int64(len(data)) > limit {
		return nil, &BodyLimitError{Limit: limit}
	}
	return data, nil
}

func readErrorExcerpt(body io.Reader, limits BodyLimits) ([]byte, bool, error) {
	limit := limits.errorExcerptLimit()
	if limit < 1 {
		return nil, false, ErrInvalidBodyLimit
	}
	data, err := io.ReadAll(io.LimitReader(body, limit+1))
	if err != nil {
		return nil, false, fmt.Errorf("read error excerpt: %w", err)
	}
	truncated := int64(len(data)) > limit
	if truncated {
		data = data[:limit]
	}
	return data, truncated, nil
}

func withDefaultCodec(codec JSONCodec) JSONCodec {
	standard := standardJSONCodec()
	if codec.Marshal == nil {
		codec.Marshal = standard.Marshal
	}
	if codec.Unmarshal == nil {
		codec.Unmarshal = standard.Unmarshal
	}
	return codec
}
