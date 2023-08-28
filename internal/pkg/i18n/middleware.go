package i18n

import (
	"context"

	"golang.org/x/text/language"

	"github.com/gin-gonic/gin"
)

// I18nHandlerGin 设置语言
func I18nHandlerGin() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.GetHeader(AcceptLanguage)
		c.Request = c.Request.WithContext(context.WithValue(c, AcceptLanguage, lang))
		c.Set(AcceptLanguage, lang)
		c.Next()
	}
}

func GetLang(ctx context.Context) string {
	lang, ok := ctx.Value(AcceptLanguage).(string)
	if !ok {
		lang = language.English.String()
	}
	return lang
}
