package logging

type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

func (f Format) String() string {
	return string(f)
}

func (f Format) Valid() bool {
	switch f {
	case FormatText, FormatJSON:
		return true
	default:
		return false
	}
}
