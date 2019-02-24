package hwpush

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

// NewClient return new push client
func NewClient(clientID, clientSecret, appPkgName string) *HuaweiPushClient {
	vers := &Vers{
		Ver:   "1",
		AppID: clientID,
	}
	nspCtx, _ := json.Marshal(vers)

	return &HuaweiPushClient{
		ClientId:     clientID,
		ClientSecret: clientSecret,
		AppPkgName:   appPkgName,
		NspCtx:       string(nspCtx),
	}
}

// NewMessage return a default Message
func NewMessage() *Message {
	return &Message{
		Hps: Hps{
			Msg: Msg{
				Type: 3, // 1, 透传异步消息; 3, 系统通知栏异步消息;
				Body: Body{
					Content: " ",
					Title:   " ",
				},
				Action: Action{
					Type: 3, // 1, 自定义行为; 2, 打开URL; 3, 打开App;
					Param: Param{
						// Intent:     "#Intent;compo=com.rvr/.Activity;S.W=U;end",
						AppPkgName: "",
					},
				},
			},
			Ext: Ext{ // 扩展信息, 含 BI 消息统计, 特定展示风格, 消息折叠
			},
		},
	}
}

// FormPost do a post request
func FormPost(url string, data url.Values) (int, string, []byte, error) {
	statusCode := http.StatusBadRequest
	nspStatus := ""

	resp, err := http.Post(url, "application/x-www-form-urlencoded", ioutil.NopCloser(strings.NewReader(data.Encode())))
	defer resp.Body.Close()

	if resp != nil {
		statusCode = resp.StatusCode
	}
	if status := resp.Header.Get("NSP_STATUS"); status != "" {
		nspStatus = status
	}
	if err != nil {
		return statusCode, nspStatus, []byte(""), err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return statusCode, nspStatus, []byte(""), err
	}
	return statusCode, nspStatus, body, err
}

// GetAccessToken from TOKEN_URL
func (this HuaweiPushClient) GetAccessToken() (string, error) {
	now := time.Now().Unix()

	// log.WithFields(log.Fields{
	// 	"now": now,
	// 	"token": accessToken,
	// }).Debug("GetAccessToken")

	if accessToken.Expires > now {
		return accessToken.AccessToken, nil
	}

	param := make(url.Values)
	param["grant_type"] = []string{GRANTTYPE}
	param["client_id"] = []string{this.ClientId}
	param["client_secret"] = []string{this.ClientSecret}

	code, nsp, res, err := FormPost(TOKEN_URL, param)
	if err != nil {
		log.WithFields(log.Fields{
			"err":        err,
			"res":        string(res),
			"statusCode": code,
			"nspStatus":  nsp,
		}).Debug("Huawei Push GetAccessToken failed")
		return "", err
	}

	var tokenRes = TokenResStruct{}
	err = json.Unmarshal(res, &tokenRes)
	if err != nil {
		log.WithFields(log.Fields{
			"err":        err,
			"res":        string(res),
			"statusCode": code,
			"nspStatus":  nsp,
		}).Debug("Huawei Push Unmarshal accessToken resp failed")
		return "", err
	}

	accessToken.AccessToken = tokenRes.AccessToken
	accessToken.Expires += now + tokenRes.Expires - 10
	return tokenRes.AccessToken, nil
}

// PushMsg to PUSH_URL
func (this HuaweiPushClient) PushMsg(deviceToken, payload string) (PushResStruct, error) {
	accessToken, err := this.GetAccessToken()
	if err != nil {
		return PushResStruct{}, err
	}

	var originParam = map[string]string{
		"access_token":      accessToken,
		"nsp_svc":           NSP_SVC,
		"nsp_ts":            strconv.Itoa(int(time.Now().Unix())),
		"device_token_list": "[\"" + deviceToken + "\"]",
		"payload":           payload,
		// "expire_time":       time.Now().Format("2006-01-02T15:04"),
	}

	param := make(url.Values)
	param["access_token"] = []string{originParam["access_token"]}
	param["nsp_svc"] = []string{originParam["nsp_svc"]}
	param["nsp_ts"] = []string{originParam["nsp_ts"]}
	param["device_token_list"] = []string{originParam["device_token_list"]}
	param["payload"] = []string{originParam["payload"]}

	var result = PushResStruct{}

	reqUrl := PUSH_URL + "?nsp_ctx=" + url.QueryEscape(this.NspCtx)
	code, nsp, res, err := FormPost(reqUrl, param)
	if err != nil {
		// log.WithFields(log.Fields{
		// 	"err":        err,
		// 	"res":        string(res),
		// 	"statusCode": code,
		// 	"nspStatus":  nsp,
		// }).Debug("Huawei PushMsg failed")
		result.StatusCode = code
		result.NspStatus = nsp
		result.Msg = err.Error()
		return result, err
	}

	if len(res) > 0 && code < http.StatusInternalServerError {
		err = json.Unmarshal(res, &result)
		result.StatusCode = code
		result.NspStatus = nsp
		if err != nil {
			// log.WithFields(log.Fields{
			// 	"err":        err,
			// 	"res":        string(res),
			// 	"statusCode": code,
			// 	"nspStatus":  nsp,
			// }).Debug("Huawei Push Unmarshal pushResult failed")
			return result, err
		}

		if result.NspStatus != "" {
			// log.WithFields(log.Fields{
			// 	"res":        string(res),
			// 	"statusCode": code,
			// 	"nspStatus":  nsp,
			// }).Debug("Huawei Push Unmarshal pushResult failed")
			return result, errors.New("NSP_STATUS: " + result.NspStatus)
		} else if len(result.Msg) > 0 && result.Msg != "Success" {
			// log.WithFields(log.Fields{
			// 	"res":        string(res),
			// 	"statusCode": code,
			// 	"nspStatus":  nsp,
			// }).Debug("Huawei Push Unmarshal pushResult failed")
			return result, errors.New(result.Msg)
		}
	} else {
		// log.WithFields(log.Fields{
		// 	"res":        string(res),
		// 	"statusCode": code,
		// 	"nspStatus":  nsp,
		// }).Debug("Huawei PushMsg failed")
		return result, errors.New(fmt.Sprintf("statusCode: %d", result.StatusCode))
	}

	return result, nil
}
