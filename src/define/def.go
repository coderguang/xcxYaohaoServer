package yaohaoDef

import (
	"sync"

	"github.com/coderguang/GameEngine_go/sgtime"
)

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
