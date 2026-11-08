package uniqsym

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend ConvertFromString
// goverter:extend ConvertToString
// goverter:extend ConvertFromNullString
// goverter:extend ConvertToNullString
var (
	ConvertFromStrings func([]string) ([]ADT, error)
	ConvertToStrings   func([]ADT) []string
)
