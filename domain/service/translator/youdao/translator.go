package youdao

import (
	"anto/domain/service/translator"
	"anto/lib/log"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-module/carbon"
	"go.uber.org/zap"
	"net/url"
	"strings"
	"sync"
)

var api = "https://fanyi.youdao.com/translate?&doctype=json"

var (
	apiTranslator  *Translator
	onceTranslator sync.Once
)

func Singleton() *Translator {
	onceTranslator.Do(func() {
		apiTranslator = New()
	})
	return apiTranslator
}

func New() *Translator {
	return &Translator{
		id:            "youdao",
		name:          "有道翻译",
		qps:           50,
		procMax:       20,
		textMaxLen:    2000,
		sep:           "\n",
		langSupported: langSupported,
	}
}

type Translator struct {
	id            string
	name          string
	qps           int
	procMax       int
	textMaxLen    int
	langSupported []translator.LangPair
	sep           string
}

func (customT *Translator) Init(_ interface{}) {}

func (customT *Translator) GetId() string                           { return customT.id }
func (customT *Translator) GetShortId() string                      { return "yd" }
func (customT *Translator) GetName() string                         { return customT.name }
func (customT *Translator) GetCfg() interface{}                     { return nil }
func (customT *Translator) GetQPS() int                             { return customT.qps }
func (customT *Translator) GetProcMax() int                         { return customT.procMax }
func (customT *Translator) GetTextMaxLen() int                      { return customT.textMaxLen }
func (customT *Translator) GetLangSupported() []translator.LangPair { return customT.langSupported }
func (customT *Translator) GetSep() string                          { return customT.sep }
func (customT *Translator) IsValid() bool                           { return true }

func (customT *Translator) Translate(ctx context.Context, args *translator.TranslateArgs) (*translator.TranslateRes, error) {
	timeStart := carbon.Now()
	urlQueried := fmt.Sprintf(
		"%s&type=%s2%s&i=%s", api,
		strings.ToUpper(args.FromLang), strings.ToUpper(args.ToLang),
		url.QueryEscape(args.TextContent),
	)

	respBytes, err := translator.RequestSimpleGet(ctx, customT, urlQueried)
	if err != nil {
		return nil, err
	}
	youDaoResp := new(youDaoMTResp)
	if err = json.Unmarshal(respBytes, youDaoResp); err != nil {
		log.Singleton().ErrorF("解析报文异常, 引擎: %s, 错误: %s", customT.GetName(), err)
		return nil, fmt.Errorf("解析报文出现异常, 错误: %s", err.Error())
	}
	if youDaoResp.ErrorCode != 0 {
		log.Singleton().ErrorF("接口响应异常, 引擎: %s, 错误: %s", customT.GetName(), err, zap.String("result", string(respBytes)))
		return nil, fmt.Errorf("翻译异常, 代码: %d", youDaoResp.ErrorCode)
	}
	srcArrSplit := strings.Split(args.TextContent, customT.sep)
	if len(srcArrSplit) != len(youDaoResp.TransResult) {
		return nil, translator.ErrSrcAndTgtNotMatched
	}
	ret := new(translator.TranslateRes)
	for idx, transBlockArray := range youDaoResp.TransResult {
		var tgtArrEvaluated []string
		// ?+空格 会导致意外分行, 搞不清楚这个服务的换行标识是什么, 多标准的么
		for _, block := range transBlockArray {
			tgtArrEvaluated = append(tgtArrEvaluated, block.Tgt)
		}
		ret.Results = append(ret.Results, &translator.TranslateResBlock{
			Id:             srcArrSplit[idx],
			TextTranslated: strings.Join(tgtArrEvaluated, " "),
		})

	}

	ret.TimeUsed = int(carbon.Now().DiffInSeconds(timeStart))
	return ret, nil
}

type youDaoMTResp struct {
	Type        string `json:"type"`
	ErrorCode   int    `json:"errorCode"`
	ElapsedTime int    `json:"elapsedTime"`
	TransResult [][]struct {
		Src string `json:"src,omitempty"` // 原文
		Tgt string `json:"tgt,omitempty"` // 译文
	} `json:"translateResult"`
}
