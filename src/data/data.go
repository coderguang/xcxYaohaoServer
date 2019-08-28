package yaohaoData

import (
	"sort"
	yaohaoConfig "xcxYaohaoServer/src/config"
	yaohaoDef "xcxYaohaoServer/src/define"

	"github.com/coderguang/GameEngine_go/sgfile"
	"github.com/coderguang/GameEngine_go/sglog"
)

var globalCfg *yaohaoDef.Config
var data_code *yaohaoDef.SecureSData
var data_name *yaohaoDef.SecureSData

var downloadMap map[string]*yaohaoDef.HistoryUrlData
var urlMap map[string]string //don't save to db again
var ignoreUrlMap map[string]string
var globalDataUpdateFlag bool
var global_require_times int
var global_laste_data_timestr string

func InitConfig(configfile string) {
	globalCfg = yaohaoConfig.ReadConfig(configfile)

	sgfile.AutoMkDir(GetPDFDir())
	sgfile.AutoMkDir(GetTxtDir())
	sgfile.AutoMkDir(GetFinishPDFDir())
	sgfile.AutoMkDir(GetFinishTxtDir())

	globalDataUpdateFlag = false

	urlMap = make(map[string]string)
	downloadMap = make(map[string]*yaohaoDef.HistoryUrlData)
	ignoreUrlMap = make(map[string]string)

	global_require_times = 0

	global_laste_data_timestr = "201101"

	data_code = new(yaohaoDef.SecureSData)
	data_code.Data = make(map[string][]*yaohaoDef.SData)
	data_name = new(yaohaoDef.SecureSData)
	data_name.Data = make(map[string][]*yaohaoDef.SData)

}

func GetTitle() string {
	return globalCfg.Title
}

func GetHistoryTableName() string {
	return globalCfg.HistoryTable + "_" + GetTitle()
}

func GetDataTableName() string {
	return globalCfg.DbTable + "_" + GetTitle()
}

func GetPDFDir() string {
	return "./data/" + GetTitle() + "/pdf/"
}

func GetTxtDir() string {
	return "./data/" + GetTitle() + "/txt/"
}

func GetFinishPDFDir() string {
	return "./data/" + GetTitle() + "/finish_pdf/"
}

func GetFinishTxtDir() string {
	return "./data/" + GetTitle() + "/finish_txt/"
}

func GetNoticeUrl() string {
	return globalCfg.NoticeUrl
}

func GetAllowDomains() []string {
	return globalCfg.AllowUrls
}

func GetHttpType() string {
	return globalCfg.Http
}

func GetIgnoreUrls() []string {
	return globalCfg.IgnoreUrls
}

func GetPageTxt() []string {
	return globalCfg.PageTxt
}

func GetTotalRequireTimes() int {
	return global_require_times
}
func AddTotalRequireTimes() {
	global_require_times++
}

func IsInUpdateCardData() bool {
	return globalDataUpdateFlag
}

func GetMatchData(key string) (bool, []*yaohaoDef.SData) {
	data_code.Lock.RLock()
	defer data_code.Lock.RUnlock()
	if v, ok := data_code.Data[key]; ok {
		return true, v
	}
	data_name.Lock.RLock()
	defer data_name.Lock.RUnlock()
	if v, ok := data_name.Data[key]; ok {
		return true, v
	}
	return false, nil
}

func GetDbConnectionData() (string, string, string, string, string) {
	return globalCfg.DbUser, globalCfg.DbPwd, globalCfg.DbUrl, globalCfg.DbPort, globalCfg.DbName
}

func GetIndexUrl() string {
	return globalCfg.IndexUrl
}

func GetTimeTxt() string {
	return globalCfg.TimeTxt
}

func GetTotalNumTxt() string {
	return globalCfg.TotalNumTxt
}

func GetPersonTxt() string {
	return globalCfg.PersonTxt
}

func GetCompanyTxt() string {
	return globalCfg.CompanyTxt
}

func GetNormalTxt() string {
	return globalCfg.NormalTxt
}

func GetNewEngineTxt() []string {
	return globalCfg.NewEngineTxt
}

func GetListenPort() string {
	return globalCfg.ListenPort
}

//init data from db

func InitDbURLData(data *yaohaoDef.HistoryUrlData) {
	downloadMap[data.Url] = data
}

func ChangeUrlStatus(url string, status string) {
	if globalCfg.IndexUrl == url {
		return
	}
	if _, ok := ignoreUrlMap[url]; ok {
		return
	}
	urlMap[url] = status
}

func ChangeDownloadStatus(url string, status string, title string, desc string) {
	if _, ok := ignoreUrlMap[url]; ok {
		return
	}

	if v, ok := downloadMap[url]; ok {
		v.Status = status
		v.Title = title
		v.Desc = desc
	} else {
		data := new(yaohaoDef.HistoryUrlData)
		data.Url = url
		data.Status = status
		data.Title = title
		data.Desc = desc
		downloadMap[data.Url] = data
	}

}

func AddIgnoreUrl(url string) {
	ignoreUrlMap[url] = ""
}

func NeedVisitUrl(url string) bool {

	if v, ok := urlMap[url]; ok {
		if v == yaohaoDef.URL_status_complete {
			return false
		}
	}
	return true
}

func GetReDownloadList() []*yaohaoDef.HistoryUrlData {
	downlist := []*yaohaoDef.HistoryUrlData{}
	for _, v := range downloadMap {
		if v.Status == yaohaoDef.URL_status_visting {
			downlist = append(downlist, v)
		}
	}
	return downlist
}

func GetRevisitList() []string {
	vlist := []string{}
	for k, v := range urlMap {
		if v == yaohaoDef.URL_status_visting {
			vlist = append(vlist, k)
		}
	}
	return vlist
}

func SetUpdateFlag(flag bool) {
	globalDataUpdateFlag = flag
}

//get newData
func AddCardData(data *yaohaoDef.SData) {
	data_code.Lock.Lock()
	defer data_code.Lock.Unlock()
	if _, ok := data_code.Data[data.Code]; ok {
		// if len(data_code[data.Code]) > 5 {
		// 	if data.Time < data_code[data.Code][len(data_code[data.Code])-1].Time {
		// 		continue
		// 	}
		// 	data_code[data.Code][len(data_code[data.Code])-1] = data
		// } else {
		// 	data_code[data.Code] = append(data_code[data.Code], data)
		// }
		data_code.Data[data.Code] = append(data_code.Data[data.Code], data)
		sort.Slice(data_code.Data[data.Code], func(i, j int) bool {
			return data_code.Data[data.Code][i].Time > data_code.Data[data.Code][j].Time
		})
	} else {
		tmp := []*yaohaoDef.SData{}
		tmp = append(tmp, data)
		data_code.Data[data.Code] = tmp
	}
	//===============================================================
	data_name.Lock.Lock()
	defer data_name.Lock.Unlock()
	if _, ok := data_name.Data[data.Name]; ok {
		// if len(data_name[data.Name]) > 5 {
		// 	if data.Time < data_name[data.Name][len(data_name[data.Name])-1].Time {
		// 		//log.Println("data.Time1", data.Time)
		// 		//log.Println("2", data_name[data.Name][len(data_name[data.Name])-1].Time)
		// 		continue
		// 	}
		// 	data_name[data.Name][len(data_name[data.Name])-1] = data
		// } else {
		// 	data_name[data.Name] = append(data_name[data.Name], data)
		// }
		data_name.Data[data.Name] = append(data_name.Data[data.Name], data)
		sort.Slice(data_name.Data[data.Name], func(i, j int) bool {
			return data_name.Data[data.Name][i].Time > data_name.Data[data.Name][j].Time
		})
	} else {
		tmp := []*yaohaoDef.SData{}
		tmp = append(tmp, data)
		data_name.Data[data.Name] = tmp
	}

	if global_laste_data_timestr < data.Time {
		global_laste_data_timestr = data.Time
	}
}

func GetLastesTimeStr() string {
	return global_laste_data_timestr
}

//get config

func GetFinishTxt() string {
	return globalCfg.FinishTxt
}

func NeedDownloadFile(src string, title string) bool {
	if v, ok := downloadMap[src]; ok {
		if v.Status == yaohaoDef.URL_status_complete || v.Status == yaohaoDef.URL_status_error {
			sglog.Info("ignore down file %s from %s,already download", v.Title, src)
			return false
		}
	}

	if _, ok := ignoreUrlMap[src]; ok {
		sglog.Info("ignore down file %s from %s by ignore url map", title, src)
		return false
	}
	return true
}

//for show

func ShowHadVistUrl(cmdstr []string) {
	sglog.Debug("start show had visit url maps...")

	for k, v := range urlMap {
		if v == yaohaoDef.URL_status_complete {
			sglog.Info("url:%s,status:%s", k, v)
		} else if v == yaohaoDef.URL_status_error {
			sglog.Error("url:%s,status:%s", k, v)
		} else {
			sglog.Debug("url:%s,status:%s", k, v)
		}
	}
	sglog.Debug("end show had visit url maps...")
}

func ShowHadDownVistUrl(cmdstr []string) {
	sglog.Debug("start ShowHadDownVistUrl...")
	for k, v := range downloadMap {
		if v.Status == yaohaoDef.URL_status_complete {
			sglog.Info("url:%s,status:%s", k, v.Status)
		} else if v.Status == yaohaoDef.URL_status_error {
			sglog.Error("url:%s,status:%s", k, v.Status)
		} else {
			sglog.Debug("url:%s,status:%s", k, v.Status)
		}
	}
	sglog.Debug("end ShowHadDownVistUrl maps...")
}

func ShowRequireTimes(cmdstr []string) {
	sglog.Info("current require times is %d", GetTotalRequireTimes())
}

func ShowNotCompleteUrl(cmdstr []string) {
	sglog.Debug("start ShowNotCompleteUrl maps...")
	for k, v := range urlMap {
		if v == yaohaoDef.URL_status_visting {
			sglog.Error("url:%s,status:%s", k, v)
		}
	}
	sglog.Debug("end ShowNotCompleteUrl maps...")
}

func ShowIgnoreUrl(cmdstr []string) {
	sglog.Debug("start ShowIgnoreUrl maps...")
	for k, v := range ignoreUrlMap {
		sglog.Debug("url:%s,status:%s", k, v)
	}
	sglog.Debug("end ShowIgnoreUrl maps...")
}

func ShowLasteTime(cmd []string) {
	sglog.Debug("current laste time is %s", GetLastesTimeStr())
}

func AddIgnoreUrlByCmd(cmdstr []string) {
	sglog.Info("start AddIgnoreUrl cmd")
	if len(cmdstr) < 2 {
		sglog.Error("AddIgnoreUrl at least need 2 params")
		return
	}
	for k, v := range cmdstr {
		if k == 0 {
			continue
		}
		ignoreUrlMap[v] = ""
		sglog.Info("add %s to ignore map success", v)
	}
	sglog.Info("end AddIgnoreUrl cmd,total=%d", len(cmdstr)-1)
}

func RemoveIgnoreUrl(cmdstr []string) {
	sglog.Info("start RemoveIgnoreUrl cmd")
	if len(cmdstr) < 2 {
		sglog.Error("RemoveIgnoreUrl at least need 2 params")
		return
	}
	for k, v := range cmdstr {
		if k == 0 {
			continue
		}
		delete(ignoreUrlMap, v)
		sglog.Info("add %s to ignore map success", v)
	}
	sglog.Info("end RemoveIgnoreUrl cmd,total=%d", len(cmdstr)-1)
}
