package i18n

import (
	"context"
	"fmt"
	"github.com/winc-link/hummingbird/internal/pkg/i18n/locales"
	"strconv"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	//"gitlab.com/tedge/edgex/internal/pkg/i18n/locales"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

func init() {
	bundle = i18n.NewBundle(language.Chinese)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	// load message
	initMessages()
}

func initMessages() {
	err := bundle.AddMessages(language.Make(language.English.String()), locales.GetEnMessages()...)
	if err != nil {
		panic(err)
	}
	err = bundle.AddMessages(language.Make(language.Chinese.String()), locales.GetZhMessages()...)
	if err != nil {
		panic(err)
	}
}

// 常量的翻译 Trans("en", "device.name", nil)
func Trans(lang string, messageId string, params map[string]interface{}) string {
	localizer := i18n.NewLocalizer(bundle, lang)
	retStr, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageId,
		TemplateData: params,
		Funcs: template.FuncMap{
			"filter":     funcFilter,
			"trim":       funcTrim,
			"left_space": funcAddLeftSpace,
		},
	})
	if err != nil {
		return messageId
	}
	return retStr
}

// 错误码的翻译 TransCode("en",10001, nil)
func TransCode(ctx context.Context, code uint32, params map[string]interface{}) string {
	lang := GetLang(ctx)
	localizer := i18n.NewLocalizer(bundle, lang)
	strCode := strconv.Itoa(int(code))
	retStr, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    strCode,
		TemplateData: params,
		Funcs: template.FuncMap{
			"filter":     funcFilter,
			"trim":       funcTrim,
			"left_space": funcAddLeftSpace,
		},
	})
	if err != nil {
		return strCode
	}
	return retStr
}

// TransCodeDefault 错误码转换, 用于告警显示, 只显示英文错误信息
func TransCodeDefault(code uint32, params map[string]interface{}) string {
	return "[Code(" + strconv.Itoa(int(code)) + ")] " + TransCode(context.Background(), code, params)
}

func funcFilter(a interface{}) string {
	if a == nil {
		return ""
	}
	return fmt.Sprintf("%v", a)
}

func funcTrim(a interface{}) string {
	if a == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprintf("%v", a))
}

func funcAddLeftSpace(a interface{}) string {
	if a == nil || a == "" {
		return ""
	}
	return fmt.Sprintf(" %v", a)
}
