package chanctx

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var Optional = []validation.Rule{
	validation.Length(1, 10),
	validation.Each(validation.Required),
}
