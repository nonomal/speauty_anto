package caiyunai

import (
	"anto/domain/service/translator"
	"anto/lib/log"
	"anto/lib/util"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-module/carbon"
	"net/http"
	"strings"
	"sync"
)

var (
	apiTranslator  *Translator
	onceTranslator sync.Once
	api            = "http://api.interpreter.caiyunai.com/v1/translator"
)

func Singleton() *Translator {
	onceTranslator.Do(func() {
		apiTranslator = New()
	})
	return apiTranslator
}

func New() *Translator {
	return &Translator{
		id:            "caiyun_ai",
		name:          "彩云小译",
		qps:           10,
		procMax:       20,
		textMaxLen:    5000,
		sep:           "\n",
		langSupported: langSupported,
	}
}

type Translator struct {
	id            string
	name          string
	cfg           *Cfg
	qps           int
	procMax       int
	textMaxLen    int
	langSupported []translator.LangPair
	sep           string
}

func (customT *Translator) Init(cfg interface{}) { customT.cfg = cfg.(*Cfg) }

func (customT *Translator) GetId() string       { return customT.id }
func (customT *Translator) GetShortId() string  { return "cy" }
func (customT *Translator) GetName() string     { return customT.name }
func (customT *Translator) GetCfg() interface{} { return nil }
func (customT *Translator) GetQPS() int         { return customT.qps }
func (customT *Translator) GetProcMax() int     { return customT.procMax }
func (customT *Translator) GetTextMaxLen() int {
	if customT.cfg.MaxSingleTextLength > 0 {
		return customT.cfg.MaxSingleTextLength
	}
	return customT.textMaxLen
}
func (customT *Translator) GetLangSupported() []translator.LangPair { return customT.langSupported }
func (customT *Translator) GetSep() string                          { return customT.sep }
func (customT *Translator) IsValid() bool                           { return customT.cfg != nil && customT.cfg.Token != "" }

func (customT *Translator) Translate(ctx context.Context, args *translator.TranslateArgs) (*translator.TranslateRes, error) {
	timeStart := carbon.Now()
	texts := strings.Split(args.TextContent, customT.GetSep())
	req := new(translateReq)
	req.Source = texts
	req.RequestId = util.Uid()
	req.TransType = fmt.Sprintf("%s2%s", args.FromLang, args.ToLang)

	reqBytes, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest(http.MethodPost, api, bytes.NewReader(reqBytes))
	httpReq.Header.Set("content-type", "application/json")
	httpReq.Header.Set("x-authorization", fmt.Sprintf("token %s", customT.cfg.Token))
	respBytes, err := translator.RequestSimpleHttp(ctx, customT, httpReq)
	if err != nil {
		return nil, err
	}

	resp := new(translateResp)
	if err = json.Unmarshal(respBytes, resp); err != nil {
		log.Singleton().ErrorF("解析报文异常, 引擎: %s, 错误: %s", customT.GetName(), err)
		return nil, fmt.Errorf("解析报文出现异常, 错误: %s", err.Error())
	}

	if resp.Msg != "" {
		log.Singleton().ErrorF("接口响应异常, 引擎: %s, 错误: %s", customT.GetName(), resp.Msg)
		return nil, fmt.Errorf("接口响应异常, 引擎: %s, 错误: %s", customT.GetName(), resp.Msg)
	}

	if len(texts) != len(resp.Target) {
		return nil, translator.ErrSrcAndTgtNotMatched
	}

	ret := new(translator.TranslateRes)

	for textIdx, textTarget := range resp.Target {
		ret.Results = append(ret.Results, &translator.TranslateResBlock{
			Id:             texts[textIdx],
			TextTranslated: textTarget,
		})
	}

	ret.TimeUsed = int(carbon.Now().DiffAbsInSeconds(timeStart))
	return ret, nil

}

type translateReq struct {
	Source    []string `json:"source"`
	TransType string   `json:"trans_type"`
	RequestId string   `json:"request_id"`
}

type translateResp struct {
	Msg        string   `json:"message,omitempty"`
	SrcTgt     []string `json:"src_tgt,omitempty"`
	Target     []string `json:"target,omitempty"`
	Confidence float64  `json:"confidence,omitempty"` // 可信度?
	Rc         int      `json:"rc,omitempty"`
}
