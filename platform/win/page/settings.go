package page

import (
	"anto/cfg"
	"anto/lib/log"
	"anto/lib/util"
	"anto/platform/win/ui/msg"
	"anto/platform/win/ui/pack"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"go.uber.org/zap"
)

var (
	apiSettings        *Settings
	onceSettings       sync.Once
	stdLineEditSize    = Size{Width: 100}
	stdNumLineEditSize = Size{Width: 40}
)

func GetSettings() *Settings {
	onceSettings.Do(func() {
		apiSettings = new(Settings)
		apiSettings.id = util.Uid()
		apiSettings.name = "设置"
	})
	return apiSettings
}

type Settings struct {
	id         string
	name       string
	mainWindow *walk.MainWindow
	rootWidget *walk.Composite

	ptrNiutransAppKey *walk.LineEdit

	ptrLingVADataId              *walk.LineEdit
	ptrLingVAMaxSingleTextLength *walk.LineEdit

	ptrBaiduAppId               *walk.LineEdit
	ptrBaiduAppKey              *walk.LineEdit
	ptrBaiduMaxSingleTextLength *walk.LineEdit

	ptrTencentCloudMTSecretId     *walk.LineEdit
	ptrTencentCloudMTSecretKey    *walk.LineEdit
	ptrTencentMaxSingleTextLength *walk.LineEdit

	ptrOpenAPIYouDaoAppKey        *walk.LineEdit
	ptrOpenAPIYouDaoAppSecret     *walk.LineEdit
	ptrOpenAPIMaxSingleTextLength *walk.LineEdit

	ptrAliCloudMTAkId                *walk.LineEdit
	ptrAliCloudMTAkSecret            *walk.LineEdit
	ptrAliCloudMTMaxSingleTextLength *walk.LineEdit

	ptrCaiYunAIToken               *walk.LineEdit
	ptrCaiYunAIMaxSingleTextLength *walk.LineEdit

	ptrVolcEngineAccessKey *walk.LineEdit
	ptrVolcEngineSecretKey *walk.LineEdit

	ptrHuaweiCloudAKId                *walk.LineEdit
	ptrHuaweiCloudSKKey               *walk.LineEdit
	ptrHuaweiCloudAKProjectId         *walk.LineEdit
	ptrHuaweiCloudMaxSingleTextLength *walk.LineEdit
}

func (customPage *Settings) GetId() string {
	return customPage.id
}

func (customPage *Settings) GetName() string {
	return customPage.name
}

func (customPage *Settings) BindWindow(win *walk.MainWindow) {
	customPage.mainWindow = win
}

func (customPage *Settings) SetVisible(isVisible bool) {
	if customPage.rootWidget != nil {
		customPage.rootWidget.SetVisible(isVisible)
	}
}

func (customPage *Settings) GetWidget() Widget {
	defer customPage.Reset()
	return StdRootWidget(&customPage.rootWidget,
		pack.UIScrollView(pack.NewUIScrollViewArgs(nil).SetChildren(
			pack.NewWidgetGroup().Append(
				pack.UIGroupBox(pack.NewUIGroupBoxArgs(nil).SetTitle("翻译引擎").SetLayoutVBox(false).SetWidgets(
					pack.NewWidgetGroup().Append(
						customPage.getNiuWidget(),
						customPage.getLingVAWidget(),
						customPage.getCaiYunAIWidget(),
						customPage.getVolcWidget(),
						customPage.getBaiduWidget(),
						customPage.getTencentCloudMTWidget(),
						customPage.getOpenAPIYouDaoWidget(),
						customPage.getAliCloudMTWidget(),
						customPage.getHuaweiCloudNlpWidget(),

						pack.UIComposite(pack.NewUICompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
							pack.NewWidgetGroup().Append(
								pack.UIPushBtn(pack.NewUIPushBtnArgs(nil).SetText("保存").SetOnClicked(customPage.eventSync)),
								pack.UIPushBtn(pack.NewUIPushBtnArgs(nil).SetText("一键还原").SetOnClicked(customPage.eventRestore)),
							).AppendZeroHSpacer().GetWidgets(),
						)),
					).AppendZeroHSpacer().AppendZeroVSpacer().GetWidgets(),
				)),
			).AppendZeroVSpacer().GetWidgets())),
	)
}

func (customPage *Settings) Reset() {}

func (customPage *Settings) getNiuWidget() Widget {
	return pack.UIGroupBox(pack.NewUIGroupBoxArgs(nil).SetTitle("小牛翻译").SetLayoutVBox(false).SetWidgets(
		pack.NewWidgetGroup().Append(
			pack.UIComposite(pack.NewUICompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
				pack.NewWidgetGroup().Append(
					pack.UILabel(pack.NewUILabelArgs(nil).SetText("密钥")),
					pack.UILineEdit(pack.NewUILineEditArgs(&customPage.ptrNiutransAppKey).
						SetText(cfg.Singleton().Niutrans.AppKey).SetCustomSize(stdLineEditSize)),
				).AppendZeroHSpacer().GetWidgets(),
			)),
		).AppendZeroHSpacer().GetWidgets(),
	))
}

func (customPage *Settings) getLingVAWidget() Widget {
	return pack.UIGroupBox(pack.NewUIGroupBoxArgs(nil).SetTitle("LingVA").SetLayoutVBox(false).SetWidgets(
		pack.NewWidgetGroup().Append(
			pack.UIComposite(pack.NewUICompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
				pack.NewWidgetGroup().Append(
					pack.UILabel(pack.NewUILabelArgs(nil).SetText("应用")),
					pack.UILineEdit(pack.NewUILineEditArgs(&customPage.ptrLingVADataId).
						SetText(cfg.Singleton().LingVA.DataId).SetCustomSize(stdLineEditSize)),
					pack.UILabel(pack.NewUILabelArgs(nil).SetText("最长请求")),
					pack.UILineEdit(pack.NewUILineEditArgs(&customPage.ptrLingVAMaxSingleTextLength).
						SetText(fmt.Sprintf("%d", cfg.Singleton().LingVA.GetMaxCharNum())).SetCustomSize(stdNumLineEditSize)),
				).AppendZeroHSpacer().GetWidgets(),
			)),
		).AppendZeroHSpacer().GetWidgets(),
	))
}

func (customPage *Settings) getCaiYunAIWidget() Widget {
	return pack.UIGroupBox(pack.NewUIGroupBoxArgs(nil).SetTitle("彩云小译").SetLayoutVBox(false).SetWidgets(
		pack.NewWidgetGroup().Append(
			pack.UIComposite(pack.NewUICompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
				pack.NewWidgetGroup().Append(
					pack.UILabel(pack.NewUILabelArgs(nil).SetText("密钥")),
					pack.UILineEdit(pack.NewUILineEditArgs(&customPage.ptrCaiYunAIToken).
						SetText(cfg.Singleton().CaiYunAI.Token).SetCustomSize(stdLineEditSize)),
					pack.UILabel(pack.NewUILabelArgs(nil).SetText("最长请求")),
					pack.UILineEdit(pack.NewUILineEditArgs(&customPage.ptrCaiYunAIMaxSingleTextLength).
						SetText(fmt.Sprintf("%d", cfg.Singleton().CaiYunAI.GetMaxCharNum())).SetCustomSize(stdNumLineEditSize)),
				).AppendZeroHSpacer().GetWidgets(),
			)),
		).AppendZeroHSpacer().GetWidgets(),
	))
}

func (customPage *Settings) getVolcWidget() Widget {
	return pack.UIGroupBox(pack.NewUIGroupBoxArgs(nil).SetTitle("火山翻译").SetLayoutVBox(false).SetWidgets(
		pack.NewWidgetGroup().Append(
			pack.UIComposite(pack.NewUICompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
				pack.NewWidgetGroup().Append(
					pack.UILabel(pack.NewUILabelArgs(nil).SetText("应用")),
					pack.UILineEdit(pack.NewUILineEditArgs(&customPage.ptrVolcEngineAccessKey).
						SetText(cfg.Singleton().VolcEngine.AccessKey).SetCustomSize(stdLineEditSize)),
					pack.UILabel(pack.NewUILabelArgs(nil).SetText("密钥")),
					pack.UILineEdit(pack.NewUILineEditArgs(&customPage.ptrVolcEngineSecretKey).
						SetText(cfg.Singleton().VolcEngine.SecretKey).SetCustomSize(stdLineEditSize)),
				).AppendZeroHSpacer().GetWidgets(),
			)),
		).AppendZeroHSpacer().GetWidgets(),
	))
}

func (customPage *Settings) getBaiduWidget() Widget {
	return pack.UIGroupBox(pack.NewUIGroupBoxArgs(nil).SetTitle("百度翻译").SetLayoutVBox(false).SetWidgets(
		pack.NewWidgetGroup().Append(
			pack.UIComposite(pack.NewUICompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
				pack.NewWidgetGroup().Append(
					pack.UILabel(pack.NewUILabelArgs(nil).SetText("应用")),
					pack.UILineEdit(pack.NewUILineEditArgs(&customPage.ptrBaiduAppId).
						SetText(cfg.Singleton().Baidu.AppId).SetCustomSize(stdLineEditSize)),
					pack.UILabel(pack.NewUILabelArgs(nil).SetText("密钥")),
					pack.UILineEdit(pack.NewUILineEditArgs(&customPage.ptrBaiduAppKey).
						SetText(cfg.Singleton().Baidu.AppKey).SetCustomSize(stdLineEditSize)),
					pack.UILabel(pack.NewUILabelArgs(nil).SetText("最长请求")),
					pack.UILineEdit(pack.NewUILineEditArgs(&customPage.ptrBaiduMaxSingleTextLength).
						SetText(fmt.Sprintf("%d", cfg.Singleton().Baidu.GetMaxCharNum())).SetCustomSize(stdNumLineEditSize)),
				).AppendZeroHSpacer().GetWidgets(),
			)),
		).AppendZeroHSpacer().GetWidgets(),
	))
}

func (customPage *Settings) getTencentCloudMTWidget() Widget {
	return pack.UIGroupBox(pack.NewUIGroupBoxArgs(nil).SetTitle("腾讯云翻译").SetLayoutVBox(false).SetWidgets(
		pack.NewWidgetGroup().Append(
			pack.UIComposite(pack.NewUICompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
				pack.NewWidgetGroup().Append(
					pack.UILabel(pack.NewUILabelArgs(nil).SetText("应用")),
					pack.UILineEdit(pack.NewUILineEditArgs(&customPage.ptrTencentCloudMTSecretId).
						SetText(cfg.Singleton().TencentCloudMT.SecretId).SetCustomSize(stdLineEditSize)),
					pack.UILabel(pack.NewUILabelArgs(nil).SetText("密钥")),
					pack.UILineEdit(pack.NewUILineEditArgs(&customPage.ptrTencentCloudMTSecretKey).
						SetText(cfg.Singleton().TencentCloudMT.SecretKey).SetCustomSize(stdLineEditSize)),
					pack.UILabel(pack.NewUILabelArgs(nil).SetText("最长请求")),
					pack.UILineEdit(pack.NewUILineEditArgs(&customPage.ptrTencentMaxSingleTextLength).
						SetText(fmt.Sprintf("%d", cfg.Singleton().TencentCloudMT.GetMaxCharNum())).SetCustomSize(stdNumLineEditSize)),
				).AppendZeroHSpacer().GetWidgets(),
			)),
		).AppendZeroHSpacer().GetWidgets(),
	))
}

func (customPage *Settings) getOpenAPIYouDaoWidget() Widget {
	return pack.UIGroupBox(pack.NewUIGroupBoxArgs(nil).SetTitle("有道智云翻译").SetLayoutVBox(false).SetWidgets(
		pack.NewWidgetGroup().Append(
			pack.UIComposite(pack.NewUICompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
				pack.NewWidgetGroup().Append(
					pack.UILabel(pack.NewUILabelArgs(nil).SetText("应用")),
					pack.UILineEdit(pack.NewUILineEditArgs(&customPage.ptrOpenAPIYouDaoAppKey).
						SetText(cfg.Singleton().OpenAPIYouDao.AppKey).SetCustomSize(stdLineEditSize)),
					pack.UILabel(pack.NewUILabelArgs(nil).SetText("密钥")),
					pack.UILineEdit(pack.NewUILineEditArgs(&customPage.ptrOpenAPIYouDaoAppSecret).
						SetText(cfg.Singleton().OpenAPIYouDao.AppSecret).SetCustomSize(stdLineEditSize)),
					pack.UILabel(pack.NewUILabelArgs(nil).SetText("最长请求")),
					pack.UILineEdit(pack.NewUILineEditArgs(&customPage.ptrOpenAPIMaxSingleTextLength).
						SetText(fmt.Sprintf("%d", cfg.Singleton().OpenAPIYouDao.GetMaxCharNum())).SetCustomSize(stdNumLineEditSize)),
				).AppendZeroHSpacer().GetWidgets(),
			)),
		).AppendZeroHSpacer().GetWidgets(),
	))
}

func (customPage *Settings) getAliCloudMTWidget() Widget {
	return pack.UIGroupBox(pack.NewUIGroupBoxArgs(nil).SetTitle("阿里云翻译").SetLayoutVBox(false).SetWidgets(
		pack.NewWidgetGroup().Append(
			pack.UIComposite(pack.NewUICompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
				pack.NewWidgetGroup().Append(
					pack.UILabel(pack.NewUILabelArgs(nil).SetText("应用")),
					pack.UILineEdit(pack.NewUILineEditArgs(&customPage.ptrAliCloudMTAkId).
						SetText(cfg.Singleton().AliCloudMT.AKId).SetCustomSize(stdLineEditSize)),
					pack.UILabel(pack.NewUILabelArgs(nil).SetText("密钥")),
					pack.UILineEdit(pack.NewUILineEditArgs(&customPage.ptrAliCloudMTAkSecret).
						SetText(cfg.Singleton().AliCloudMT.AKSecret).SetCustomSize(stdLineEditSize)),
					pack.UILabel(pack.NewUILabelArgs(nil).SetText("最长请求")),
					pack.UILineEdit(pack.NewUILineEditArgs(&customPage.ptrAliCloudMTMaxSingleTextLength).
						SetText(fmt.Sprintf("%d", cfg.Singleton().AliCloudMT.GetMaxCharNum())).SetCustomSize(stdNumLineEditSize)),
				).AppendZeroHSpacer().GetWidgets(),
			)),
		).AppendZeroHSpacer().GetWidgets(),
	))
}

func (customPage *Settings) getHuaweiCloudNlpWidget() Widget {
	return pack.UIGroupBox(pack.NewUIGroupBoxArgs(nil).SetTitle("华为云翻译").SetLayoutVBox(false).SetWidgets(
		pack.NewWidgetGroup().Append(
			pack.UIComposite(pack.NewUICompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
				pack.NewWidgetGroup().Append(
					pack.UILabel(pack.NewUILabelArgs(nil).SetText("应用")),
					pack.UILineEdit(pack.NewUILineEditArgs(&customPage.ptrHuaweiCloudAKId).SetText(cfg.Singleton().HuaweiCloudNlp.AKId).SetCustomSize(stdLineEditSize)),
					pack.UILabel(pack.NewUILabelArgs(nil).SetText("密钥")),
					pack.UILineEdit(pack.NewUILineEditArgs(&customPage.ptrHuaweiCloudSKKey).SetText(cfg.Singleton().HuaweiCloudNlp.SkKey).SetCustomSize(stdLineEditSize)),
					pack.UILabel(pack.NewUILabelArgs(nil).SetText("项目")),
					pack.UILineEdit(pack.NewUILineEditArgs(&customPage.ptrHuaweiCloudAKProjectId).SetText(cfg.Singleton().HuaweiCloudNlp.ProjectId).SetCustomSize(stdLineEditSize)),
					pack.UILabel(pack.NewUILabelArgs(nil).SetText("最长请求")),
					pack.UILineEdit(pack.NewUILineEditArgs(&customPage.ptrHuaweiCloudMaxSingleTextLength).
						SetText(fmt.Sprintf("%d", cfg.Singleton().HuaweiCloudNlp.GetMaxCharNum())).SetCustomSize(stdNumLineEditSize)),
				).AppendZeroHSpacer().GetWidgets(),
			)),
		).AppendZeroHSpacer().GetWidgets(),
	))
}

func (customPage *Settings) eventSync() {
	cntModified := 0

	{
		niutransAppKey := customPage.ptrNiutransAppKey.Text()
		if niutransAppKey != cfg.Singleton().Niutrans.AppKey {
			cfg.Singleton().Niutrans.AppKey = niutransAppKey
			cntModified++
		}
	}

	{
		lingVADataId := customPage.ptrLingVADataId.Text()
		if lingVADataId != cfg.Singleton().LingVA.DataId {
			cfg.Singleton().LingVA.DataId = lingVADataId
			cntModified++
		}
		lingVAMaxSingleTextLength := customPage.ptrLingVAMaxSingleTextLength.Text()
		lingVAMaxSingleTextLengthInt, err := strconv.Atoi(lingVAMaxSingleTextLength)
		if err != nil || lingVAMaxSingleTextLengthInt <= 0 {
			msg.Err(customPage.mainWindow, errors.New("LingVA单次最长请求无效, 请重新设置"))
			return
		}
	}

	{
		baiduAppId := customPage.ptrBaiduAppId.Text()
		baiduAppKey := customPage.ptrBaiduAppKey.Text()
		if baiduAppId != cfg.Singleton().Baidu.AppId {
			cfg.Singleton().Baidu.AppId = baiduAppId
			cntModified++
		}
		if baiduAppKey != cfg.Singleton().Baidu.AppKey {
			cfg.Singleton().Baidu.AppKey = baiduAppKey
			cntModified++
		}
		baiduSingleTextLength := customPage.ptrBaiduMaxSingleTextLength.Text()
		baiduSingleTextLengthInt, err := strconv.Atoi(baiduSingleTextLength)
		if err != nil || baiduSingleTextLengthInt <= 0 {
			msg.Err(customPage.mainWindow, errors.New("百度翻译单次最长请求无效, 请重新设置"))
			return
		}
	}

	{
		tencentCloudMTSecretId := customPage.ptrTencentCloudMTSecretId.Text()
		tencentCloudMTSecretKey := customPage.ptrTencentCloudMTSecretKey.Text()
		if tencentCloudMTSecretId != cfg.Singleton().TencentCloudMT.SecretId {
			cfg.Singleton().TencentCloudMT.SecretId = tencentCloudMTSecretId
			cntModified++
		}
		if tencentCloudMTSecretKey != cfg.Singleton().TencentCloudMT.SecretKey {
			cfg.Singleton().TencentCloudMT.SecretKey = tencentCloudMTSecretKey
			cntModified++
		}
		tencentMaxSingleTextLength := customPage.ptrTencentMaxSingleTextLength.Text()
		tencentMaxSingleTextLengthInt, err := strconv.Atoi(tencentMaxSingleTextLength)
		if err != nil || tencentMaxSingleTextLengthInt <= 0 {
			msg.Err(customPage.mainWindow, errors.New("华为云翻译单次最长请求无效, 请重新设置"))
			return
		}
	}

	{
		openAPIYouDaoAppKey := customPage.ptrOpenAPIYouDaoAppKey.Text()
		openAPIYouDaoAppSecret := customPage.ptrOpenAPIYouDaoAppSecret.Text()
		if openAPIYouDaoAppKey != cfg.Singleton().OpenAPIYouDao.AppKey {
			cfg.Singleton().OpenAPIYouDao.AppKey = openAPIYouDaoAppKey
			cntModified++
		}
		if openAPIYouDaoAppSecret != cfg.Singleton().OpenAPIYouDao.AppSecret {
			cfg.Singleton().OpenAPIYouDao.AppSecret = openAPIYouDaoAppSecret
			cntModified++
		}
		openAPIMaxSingleTextLength := customPage.ptrOpenAPIMaxSingleTextLength.Text()
		openAPIMaxSingleTextLengthInt, err := strconv.Atoi(openAPIMaxSingleTextLength)
		if err != nil || openAPIMaxSingleTextLengthInt <= 0 {
			msg.Err(customPage.mainWindow, errors.New("有道智云翻译单次最长请求无效, 请重新设置"))
			return
		}
	}

	{
		aliCloudMTAkId := customPage.ptrAliCloudMTAkId.Text()
		aliCloudMTAkSecret := customPage.ptrAliCloudMTAkSecret.Text()
		if aliCloudMTAkId != cfg.Singleton().AliCloudMT.AKId {
			cfg.Singleton().AliCloudMT.AKId = aliCloudMTAkId
			cntModified++
		}
		if aliCloudMTAkSecret != cfg.Singleton().AliCloudMT.AKSecret {
			cfg.Singleton().AliCloudMT.AKSecret = aliCloudMTAkSecret
			cntModified++
		}
		aliCloudMTMaxSingleTextLength := customPage.ptrAliCloudMTMaxSingleTextLength.Text()
		aliCloudMTMaxSingleTextLengthInt, err := strconv.Atoi(aliCloudMTMaxSingleTextLength)
		if err != nil || aliCloudMTMaxSingleTextLengthInt <= 0 {
			msg.Err(customPage.mainWindow, errors.New("阿里云翻译单次最长请求无效, 请重新设置"))
			return
		}
	}

	{
		caiYunAIToken := customPage.ptrCaiYunAIToken.Text()
		if caiYunAIToken != cfg.Singleton().CaiYunAI.Token {
			cfg.Singleton().CaiYunAI.Token = caiYunAIToken
			cntModified++
		}
		caiYunAIMaxSingleTextLength := customPage.ptrCaiYunAIMaxSingleTextLength.Text()
		caiYunAIMaxSingleTextLengthInt, err := strconv.Atoi(caiYunAIMaxSingleTextLength)
		if err != nil || caiYunAIMaxSingleTextLengthInt <= 0 {
			msg.Err(customPage.mainWindow, errors.New("彩云小译单次最长请求无效, 请重新设置"))
			return
		}
	}

	{
		volcEngineAccessKey := customPage.ptrVolcEngineAccessKey.Text()
		if volcEngineAccessKey != cfg.Singleton().VolcEngine.AccessKey {
			cfg.Singleton().VolcEngine.AccessKey = volcEngineAccessKey
			cntModified++
		}

		volcEngineSecretKey := customPage.ptrVolcEngineSecretKey.Text()
		if volcEngineSecretKey != cfg.Singleton().VolcEngine.SecretKey {
			cfg.Singleton().VolcEngine.SecretKey = volcEngineSecretKey
			cntModified++
		}
	}

	{
		huaweiCloudAKId := customPage.ptrHuaweiCloudAKId.Text()
		huaweiCloudSKKey := customPage.ptrHuaweiCloudSKKey.Text()
		huaweiCloudAKProjectId := customPage.ptrHuaweiCloudAKProjectId.Text()

		if huaweiCloudAKId != cfg.Singleton().HuaweiCloudNlp.AKId {
			cfg.Singleton().HuaweiCloudNlp.AKId = huaweiCloudAKId
			cntModified++
		}
		if huaweiCloudSKKey != cfg.Singleton().HuaweiCloudNlp.SkKey {
			cfg.Singleton().HuaweiCloudNlp.SkKey = huaweiCloudSKKey
			cntModified++
		}
		if huaweiCloudAKProjectId != cfg.Singleton().HuaweiCloudNlp.ProjectId {
			cfg.Singleton().HuaweiCloudNlp.ProjectId = huaweiCloudAKProjectId
			cntModified++
		}
		huaweiCloudMaxSingleTextLength := customPage.ptrHuaweiCloudMaxSingleTextLength.Text()
		huaweiCloudMaxSingleTextLengthInt, err := strconv.Atoi(huaweiCloudMaxSingleTextLength)
		if err != nil || huaweiCloudMaxSingleTextLengthInt <= 0 {
			msg.Err(customPage.mainWindow, errors.New("华为云翻译单次最长请求无效, 请重新设置"))
			return
		}
	}

	if cntModified == 0 {
		msg.Info(customPage.mainWindow, "暂无配置需要同步")
		return
	}
	if err := cfg.Singleton().Sync(); err != nil {
		log.Singleton().Error("同步配置到文件失败", zap.Error(err))
		msg.Err(customPage.mainWindow, errors.New("同步配置到文件失败"))
		return
	}
	msg.Info(customPage.mainWindow, "同步配置成功, 建议重启一下当前应用哦~如果没有生效的话")
}

func (customPage *Settings) eventRestore() {
	msg.Info(customPage.mainWindow, "手动删除当前目录下的[cfg.yml]文件, 重启应用即可, 配置会还原到初始默认, 请谨慎操作!!!")
}
