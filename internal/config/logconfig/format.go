package logconfig

type LogFormat int

const (
	JSONFormat = 1 << iota
	TextFormat
)

func (f LogFormat) String() string {
	switch f {
	case JSONFormat:
		return "json"
	case TextFormat:
		return "text"
	default:
		return "unidentified"
	}
}

func (f LogFormat) MarshalText() (text []byte, err error) {
	return []byte(f.String()), nil
}
