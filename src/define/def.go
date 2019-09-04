package yaohaoDef

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/coderguang/GameEngine_go/sglog"
	"github.com/coderguang/GameEngine_go/sgtime"
)

const wx_access_token_error_code string = "errcode"
const wx_access_token_error_msg string = "errmsg"
const wx_open_id string = "openid"

type Config struct {
	Title        string   `json:"title"`
	IndexUrl     string   `json:"indexUrl"`
	AllowUrls    []string `json:"allowUrls"`
	IgnoreUrls   []string `json:"ignoreUrls"`
	DbUrl        string   `json:"dbUrl"`
	DbPort       string   `json:"dbPort"`
	DbUser       string   `json:"dbUser"`
	DbPwd        string   `json:"dbPwd"`
	DbName       string   `json:"dbName"`
	DbTable      string   `json:"dbTable"`
	HistoryTable string   `json:"historyTable"`
	ListenPort   string   `json:"listenPort"`
	FinishTxt    string   `json:"finishTxt"`
	TimeTxt      string   `json:"timeTxt"`
	TotalNumTxt  string   `json:"totalNumTxt"`
	PersonTxt    string   `json:"personTxt"`
	CompanyTxt   string   `json:"companyTxt"`
	NormalTxt    string   `json:"normalTxt"`
	NewEngineTxt []string `json:"newEngineTxt"`
	PageTxt      []string `json:"pageTxt"`
	ResultDate   int      `json:"resultDate"`
	Http         string   `json:"http"`
	NoticeUrl    string   `json:"noticeUrl"`
	Appid        string   `json:"appid"`
	Secret       string   `json:"secret"`
}

const (
	URL_status_visting  = "0"
	URL_status_complete = "1"
	URL_status_error    = "2"
)

type HistoryUrlData struct {
	Url    string
	Status string
	Title  string
	Desc   string
}

type SData struct {
	Type     int    `json:"type"`
	CardType int    `json:"cardtype"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Time     string `json:"time"`
	Desc     string `json:"desc"`
}

type SecureSData struct {
	Data map[string][]*SData
	Lock sync.RWMutex
}

type SLastestCardData struct {
	TimeStr         string
	PersonalNormal  int
	PersonalJieNeng int
	CompanyNormal   int
	CompanyJieNeng  int
}

func (data *SLastestCardData) Reset() {
	now := sgtime.New()
	data.TimeStr = now.YearString() + now.MonthString()
	data.PersonalNormal = 0
	data.PersonalJieNeng = 0
	data.CompanyNormal = 0
	data.CompanyJieNeng = 0
}

type SWxOpenid struct {
	Code   string
	Openid string
	Time   *sgtime.DateTime
}

type SecureWxOpenid struct {
	Data map[string]*SWxOpenid
	Lock sync.RWMutex
}

func (data *SWxOpenid) GetOpenIdFromWx(appid string, secret string, noticeUrl string) {
	//GET https://api.weixin.qq.com/sns/jscode2session?appid=APPID&secret=SECRET&js_code=JSCODE&grant_type=authorization_code
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=" + appid + "&secret=" + secret + "&js_code=" + data.Code + "&grant_type=authorization_code"
	resp, err := http.Get(url)

	if nil != err {
		sglog.Error("get wx openid from %s error,err=%s", url, err)
	} else {

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if nil != err {
			sglog.Error("get wx openid error,read resp body error,err=%s", err)
			return
		}
		str := string(body)
		sglog.Info("openid:%s", str)
		decoder := json.NewDecoder(bytes.NewBufferString(str))
		decoder.UseNumber()
		var result map[string]interface{}
		if err := decoder.Decode(&result); err != nil {
			sglog.Error("json parse failed,str=%s,err=%s", str, err)
			return
		}
		sglog.Info("parse %s json", str)

		if _, ok := result[wx_access_token_error_code]; ok {
			sglog.Error("error openid,code=%s", result[wx_access_token_error_code])
			sglog.Error("errmsg=%s", result[wx_access_token_error_msg])
			return
		}

		tmp_openid := result[wx_open_id]
		tmp_openid_value, ok := tmp_openid.(string)
		if !ok {
			sglog.Error("parse tmp_openid failed,tmp_openid=%s", tmp_openid)

			return
		}
		sglog.Info("tmp_openid_value:%s", tmp_openid_value)

		data.Time = sgtime.New()
		data.Openid = tmp_openid_value

		if noticeUrl == "" {
			sglog.Info("get openid notices server url is empty,code=%s,openid=%d", data.Code, data.Openid)
		} else {
			params := "/?key=openid," + data.Code + "," + data.Openid
			sendUrl := noticeUrl + params
			_, err := http.Get(sendUrl)
			if err != nil {
				sglog.Error("send openid to notice server,post data error,err:%s", err)
			}
		}

	}

}

func (data *SWxOpenid) String() string {
	str := "code:" + data.Code +
		"openid:" + data.Openid
	return str
}
