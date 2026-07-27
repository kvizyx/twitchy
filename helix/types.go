package helix

import "encoding/json"

// StringEnum is the common base for named string enums. Unknown wire values
// remain representable because JSON decoding does not validate the value.
type StringEnum string

func (p Pagination) Cursor() string { return p.cursor }

func (p *Pagination) UnmarshalJSON(data []byte) error {
	var wire struct {
		Cursor string `json:"cursor"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	p.cursor = wire.Cursor
	return nil
}

func (p Pagination) MarshalJSON() ([]byte, error) {
	if p.cursor == "" {
		return []byte("{}"), nil
	}
	return json.Marshal(struct {
		Cursor string `json:"cursor"`
	}{Cursor: p.cursor})
}
