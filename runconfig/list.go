package runconfig

import "encoding/json"

// List is a base struct created for backward compatibility reasons to be able
// to Unmarshall old Config/Hostconfig []string fields when they're just strings
type List struct {
	parts []string
}

func (l *List) MarshalJSON() ([]byte, error) {
	if l == nil {
		return []byte{}, nil
	}
	return json.Marshal(l.Slice())
}

func (l *List) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return nil
	}

	var parts []string
	if err := json.Unmarshal(b, &parts); err != nil {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		parts = append(parts, s)
	}
	l.parts = parts

	return nil
}

func (l *List) Len() int {
	if l == nil {
		return 0
	}
	return len(l.parts)
}

func (l *List) Slice() []string {
	if l == nil {
		return nil
	}
	return l.parts
}
