package runconfig

import "encoding/json"

type List struct {
	Parts []string
}

func NewList(parts []string) *List {
	return &List{parts}
}

func (l *List) MarshalJSON() ([]byte, error) {
	if l == nil {
		return []byte{}, nil
	}
	return json.Marshal(l.Parts)
}

// UnmarshalJSON decoded the entrypoint whether it's a string or an array of strings.
func (l *List) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return nil
	}

	p := make([]string, 0, 1)
	if err := json.Unmarshal(b, &p); err != nil {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		p = append(p, s)
	}
	l.Parts = p
	return nil
}

func (l *List) Len() int {
	if l == nil {
		return 0
	}
	return len(l.Parts)
}

func (l *List) Slice() []string {
	if l == nil {
		return nil
	}
	return l.Parts
}

func (l *List) Append(s []string) {
	if l == nil {

	}
	l.Parts = append(l.Parts, s)
}
