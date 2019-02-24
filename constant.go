package hwpush

import (
	"encoding/json"
)

/*
 * 接口文档:
 * https://developer.huawei.com/consumer/cn/service/hms/catalog/huaweipush_agent.html?page=hmssdk_huaweipush_api_reference_agent_s2
 */

var accessToken = AccessToken{AccessToken: "", Expires: 0}

// url
const (
	TOKEN_URL = "https://login.cloud.huawei.com/oauth2/v2/token"
	PUSH_URL  = "https://api.push.hicloud.com/pushsend.do"
)

// config
const (
	GRANTTYPE = "client_credentials"
	NSP_SVC   = "openpush.message.api.send"
)

/**
 **************************************** 结构体
 */

type HuaweiPushClient struct {
	ClientId     string
	ClientSecret string
	AppPkgName   string
	NspCtx       string
}

type Vers struct {
	Ver   string `json:"ver"`
	AppID string `json:"appId"`
}

type AccessToken struct {
	AccessToken string
	Expires     int64
}

type TokenResStruct struct {
	AccessToken string `json:"access_token"`
	Expires     int64  `json:"expires_in"`
	Scope       string `json:"scope,omitempty"`
	ErrorCode   string `json:"error,omitempty"`
	ErrorMsg    string `json:"error_description,omitempty"`
}

type PushResStruct struct {
	StatusCode int         `json:"statusCode"`
	NspStatus  string      `json:"nspStatus,omitempty"`
	PushCode   string      `json:"code"`
	Msg        string      `json:"msg"`
	RequestID  string      `json:"requestId"`
	Ext        interface{} `json:"ext,omitempty"`
}

/**
 **************************************** 消息体
 */

type Message struct {
	Hps Hps `json:"hps"`
}

type Hps struct {
	Msg Msg `json:"msg"`
	Ext Ext `json:"ext"`
}
type Msg struct {
	Type   int    `json:"type"`
	Body   Body   `json:"body"`
	Action Action `json:"action"`
}
type Body struct {
	Content string `json:"content"`
	Title   string `json:"title"`
}
type Action struct {
	Type  int   `json:"type"`
	Param Param `json:"param"`
}
type Param struct {
	Intent     string `json:"intent,omitempty"`
	AppPkgName string `json:"appPkgName"`
}

type ExtObj struct {
	Name string
}
type Ext struct {
	BiTag     string                   `json:"biTag,omitempty"`
	Icon      string                   `json:"icon,omitempty"`
	Customize []map[string]interface{} `json:"customize,omitempty"`
	Action    string                   `json:"action,omitempty"`
	Func      string                   `json:"func,omitempty"`
	Collect   string                   `json:"collect,omitempty"`
	Title     string                   `json:"title,omitempty"`
	Content   string                   `json:"content,omitempty"`
	Url       string                   `json:"url,omitempty"`
}

/**
 **************************************** 封装
 */

func (this *Message) SetBiTag(tag string) *Message {
	this.Hps.Ext.BiTag = tag
	return this
}

func (this *Message) SetIcon(icon string) *Message {
	this.Hps.Ext.Icon = icon
	return this
}

func (this *Message) SetCustomize(data []map[string]interface{}) *Message {
	this.Hps.Ext.Customize = data
	return this
}

func (this *Message) SetContent(content string) *Message {
	if content == "" {
		content = " "
	}
	this.Hps.Msg.Body.Content = content
	return this
}

func (this *Message) SetTitle(title string) *Message {
	if title == "" {
		title = " "
	}
	this.Hps.Msg.Body.Title = title
	return this
}

func (this *Message) SetIntent(intent string) *Message {
	this.Hps.Msg.Action.Param.Intent = intent
	return this
}

func (this *Message) SetAppPkgName(appPkgName string) *Message {
	this.Hps.Msg.Action.Param.AppPkgName = appPkgName
	return this
}

func (this *Message) SetExtAction(Action string) *Message {
	this.Hps.Ext.Action = Action
	return this
}
func (this *Message) SetExtFunc(Func string) *Message {
	this.Hps.Ext.Func = Func
	return this
}
func (this *Message) SetExtCollect(collect string) *Message {
	this.Hps.Ext.Collect = collect
	return this
}
func (this *Message) SetExtTitle(title string) *Message {
	this.Hps.Ext.Title = title
	return this
}
func (this *Message) SetExtContent(content string) *Message {
	this.Hps.Ext.Collect = content
	return this
}

func (this *Message) SetExtUrl(url string) *Message {
	this.Hps.Ext.Url = url
	return this
}

func (this *Message) Json() string {
	bytes, err := json.Marshal(this)
	if err != nil {
		return ""
	}
	return string(bytes)
}
