package configtemplate

type Type string

const (
	TypeLongString              Type = "long string"
	TypeString                  Type = "string"
	TypeDurationNs              Type = "duration ns"
	TypeDurationUnsignedInteger Type = "unsigned integer"
)
