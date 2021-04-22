package validators

import (
	"github.com/gin-gonic/gin/binding"

	"github.com/go-playground/validator/v10"
)

func Providers() []interface{} {
	return []interface{}{
		func() *validator.Validate {
			t, ok := binding.Validator.Engine().(*validator.Validate)
			if !ok {
				panic("cannot init engine validator")
			}

			return t
		},
		func() Interface {
			if v, ok := binding.Validator.Engine().(Interface); ok {
				return v
			}
			panic("cannot init engine validator")
		},
	}
}
