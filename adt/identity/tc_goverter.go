package identity

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend ConvertFromString
// goverter:extend ConvertToString
var (
	ConvertFromStrings func([]string) ([]ADT, error)
	ConvertToStrings   func([]ADT) []string
)
