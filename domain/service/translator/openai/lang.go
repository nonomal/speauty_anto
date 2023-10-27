package openai

import (
	"anto/domain/service/translator"
)

var langSupported = []translator.LangPair{
	{"Chinese", "中文"},
	{"English", "英语"},
}
