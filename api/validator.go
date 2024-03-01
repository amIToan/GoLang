package api

import (
	"github.com/go-playground/validator/v10"
	"sgithub.com/techschool/simplebank/util"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return util.IsSupportedByCurrency(currency)
	}
	return false
}
