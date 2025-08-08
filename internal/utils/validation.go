package utils

import (
	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ptTranslations "github.com/go-playground/validator/v10/translations/pt"
)

var ptBR = pt_BR.New()
var uni = ut.New(ptBR, ptBR)
var Trans, _ = uni.GetTranslator("pt_BR")

func MakeValidator() *validator.Validate {
	Validate := validator.New()
	MustIfError(ptTranslations.RegisterDefaultTranslations(Validate, Trans))
	return Validate
}

var Validation = MakeValidator()

func MustValidateStruct(data interface{}) {
	MustIfError(Validation.Struct(data))
}
