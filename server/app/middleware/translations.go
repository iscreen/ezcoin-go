package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/locales/zh_Hant_TW"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_hans_translations "github.com/go-playground/validator/v10/translations/zh"
	zh_hant_translations "github.com/go-playground/validator/v10/translations/zh_tw"
)

func Translations() gin.HandlerFunc {
	return func(c *gin.Context) {
		uni := ut.New(en.New(), zh_Hant_TW.New(), zh.New())
		locale := c.GetHeader("locale")
		trans, _ := uni.GetTranslator(locale)
		v, ok := binding.Validator.Engine().(*validator.Validate)
		if ok {
			switch locale {
			case "zh":
				_ = zh_hans_translations.RegisterDefaultTranslations(v, trans)
			case "en":
				_ = en_translations.RegisterDefaultTranslations(v, trans)
			default:
				_ = zh_hant_translations.RegisterDefaultTranslations(v, trans)
				// _ = en_translations.RegisterDefaultTranslations(v, trans)
				translateOverride(v, trans)

			}
			c.Set("trans", trans)
		}
		c.Next()
	}
}

func translateOverride(v *validator.Validate, trans ut.Translator) {
	v.RegisterTranslation("apiKey", trans, func(ut ut.Translator) error {
		return ut.Add("apiKey", "不正確", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("apiKey")
		return t
	})
}
