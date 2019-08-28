package yaohaoSpider

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	yaohaoData "xcxYaohaoServer/src/data"
	yaohaoDb "xcxYaohaoServer/src/db"
	yaohaoDef "xcxYaohaoServer/src/define"

	"github.com/coderguang/GameEngine_go/sgfile"
	"github.com/coderguang/GameEngine_go/sglog"
	"github.com/coderguang/GameEngine_go/sgregex"
	"github.com/coderguang/GameEngine_go/sgstring"
	"github.com/coderguang/GameEngine_go/sgthread"
	"github.com/coderguang/GameEngine_go/sgtime"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

func ReadTxtFileAndInsertToDb(fileDir string) (string, int, int, int) {
	contents, err := sgfile.GetFileContentAsStringLines(fileDir)
	if err != nil {
		sglog.Error("read txt file error,file:%s,err=%s", fileDir, err)
		return "", 0, 0, 0
	}
	timestr := ""
	totalNum := 0
	memberType := 0 // 1:个人  2:单位
	cardType := 0   //1:普通  2:节能
	startParseData := false
	ignoreNumMath := false
	details := map[string]*yaohaoDef.SData{}
	for _, v := range contents {
		if startParseData {
			if strings.Contains(v, yaohaoData.GetFinishTxt()) {
				ignoreNumMath = true
			}
			strlist := strings.Split(v, " ")
			strlistex := []string{}
			for _, v := range strlist {
				if v != " " && v != "" {
					strlistex = append(strlistex, v)
				}
			}
			if len(strlistex) < 3 {
				continue
			}
			if len(strlistex) > 3 {
				sglog.Error("not suport current format data,please check,fileDir=%s", fileDir)
				return "", 0, 0, 0
			}

			data := new(yaohaoDef.SData)
			data.Time = timestr
			data.CardType = cardType
			data.Type = memberType
			codemaxlen := 50
			namemaxlen := 300
			data.Code = strlistex[1]
			data.Name = strlistex[2]
			if len(data.Code) > codemaxlen {
				sglog.Error("code len more than %d,it is %d,old code=%s", codemaxlen, len(data.Code), data.Code)
				data.Code = data.Code[0 : codemaxlen-1]
				data.Desc = "code cut "
				sglog.Error("new code=%s", data.Code)
			}
			if len(data.Name) > namemaxlen {
				sglog.Error("name len more than %d,it is %d,,old name=%s", namemaxlen, len(data.Name), data.Name)
				data.Name = data.Name[0 : namemaxlen-1]
				data.Desc += "name cut "
				sglog.Error("new name=%s", data.Name)
			}
			details[data.Code] = data
		} else {
			if strings.Contains(v, "序号") {
				if "" == timestr || 0 == totalNum || 0 == memberType || 0 == cardType {
					sglog.Error("parse file params error,file:%s", fileDir)
					return "", 0, 0, 0
				}
				startParseData = true
				sglog.Info("parse data detail complete star,time:%s,num:%d,type:%d,cardType:%d", timestr, totalNum, memberType, cardType)
				continue
			} else {
				if "" == timestr {
					if strings.Contains(v, yaohaoData.GetTimeTxt()) {
						strlist := strings.Split(v, "：")
						if len(strlist) >= 2 {
							timestr = strlist[len(strlist)-1]
						}
					}
				}
				if 0 == totalNum {
					if strings.Contains(v, yaohaoData.GetTotalNumTxt()) {
						strlist := strings.Split(v, "：")
						if len(strlist) >= 2 {
							totalNum, _ = strconv.Atoi(strlist[len(strlist)-1])
						}
					}
				}
				if 0 == memberType {
					if strings.Contains(v, yaohaoData.GetPersonTxt()) {
						memberType = 1
					}
					if strings.Contains(v, yaohaoData.GetCompanyTxt()) {
						memberType = 2
					}
				}
				if 0 == cardType {
					if strings.Contains(v, yaohaoData.GetNormalTxt()) {
						cardType = 1
					}
					if sgstring.ContainsWithOr(v, yaohaoData.GetNewEngineTxt()) {
						cardType = 2
					}
				}
			}
		}
	}

	allRight := false
	if !ignoreNumMath && len(details) != totalNum {
		sglog.Error("parse data total num error,need %d,but only get %d", totalNum, len(details))
	} else {
		allRight = true
	}

	sglog.Info("parse data detail complete from txt,num:%d,start insert data to db,it make a little seconds", totalNum)

	if err != nil {
		sglog.Error("read txt file:%s error ex,err:%e", fileDir, err)
		return "", 0, 0, 0
	}
	updateNum := 0

	noticeDatas := []string{}

	yaohaoData.SetUpdateFlag(true)

	for _, v := range details {
		if yaohaoDb.UpdateCardData(v) {
			updateNum++
			yaohaoData.AddCardData(v)
			if 1 == memberType && 1 == cardType {
				noticeDatas = append(noticeDatas, v.Code)
			}
		}
	}

	yaohaoData.SetUpdateFlag(false)

	//only personal check
	if allRight && len(noticeDatas) > 0 && memberType == 1 {
		sglog.Info("all right ,now try to send datas to notices server,timestr=%s,size=%d", timestr, len(noticeDatas))
		now := sgtime.New()
		currentYearMonth := now.YearString() + now.MonthString()
		if currentYearMonth == timestr {
			noticeUrl := yaohaoData.GetNoticeUrl()
			if noticeUrl == "" {
				sglog.Info("all right ,notices server url is empty,timestr=%s,size=%d", timestr, len(noticeDatas))
			} else {
				lenstr := strconv.Itoa(len(noticeDatas))
				cardTypeStr := strconv.Itoa(cardType)
				// data,title,time,cardType,len,detail
				params := "key=data," + yaohaoData.GetTitle() + "," + timestr + "," + cardTypeStr + "," + lenstr
				for _, v := range noticeDatas {
					params += "," + v
				}
				sglog.Info("all right ,time not match,timestr=%s,size=%d,pos msg is %s", timestr, len(noticeDatas), params)

				response, err := http.Post(noticeUrl, "application/x-www-form-urlencoded", strings.NewReader(params))
				if err != nil {
					sglog.Error("all right,post data error,err:%e", err)
				} else {

					defer response.Body.Close()

					body, err := ioutil.ReadAll(response.Body)

					if err != nil {
						sglog.Error("all right,post data error when parse response,err:%e", err)
					} else {
						sglog.Info("all right response is %s", string(body))
					}
				}

			}
		} else {
			sglog.Info("all right ,time not match,timestr=%s,size=%d", timestr, len(noticeDatas))
		}
	}

	sglog.Info("parse data detail complete to db,time:%s,num:%d,type:%d,cardType:%d,update num is %d,start add to current data", timestr, totalNum, memberType, cardType, updateNum)
	return timestr, memberType, cardType, updateNum
}

func TransportPDFToTxt(title string, rawFileName string, pdfFileName string) {
	sglog.Info("start transform pdf to txt")

	cmd := exec.Command("python3", "index.py", title, rawFileName)
	sglog.Info("command is python3 index.py %s, %s ", title, rawFileName)
	out, err := cmd.Output()
	if err != nil {
		sglog.Error("exec parse pdf to txt by python error,file=%s,err=%e", pdfFileName, err)
		return
	}
	sglog.Info("output is :\n%s", string(out))
}

func TransportTxtFileToDb(txtFileName string) {
	timestr, memberType, cardType, updateNum := ReadTxtFileAndInsertToDb(txtFileName)

	memberStr := ""
	if memberType == 1 {
		memberStr = "personal"
	} else {
		memberStr = "Company"
	}
	cardStr := ""
	if cardType == 1 {
		cardStr = "normal"
	} else {
		cardStr = "conservation"
	}
	rename := timestr + "_" + memberStr + "_" + cardStr
	sglog.Info("by cmd get data success:%s,update total:%d", rename, updateNum)
}

func DownloadFile(src string, title string) {

	if !yaohaoData.NeedDownloadFile(src, title) {
		return
	}

	sglog.Info("get pdf data src:%s", src)
	yaohaoData.ChangeDownloadStatus(src, yaohaoDef.URL_status_visting, title, "")
	yaohaoDb.UpdateDbURLData(src, yaohaoDef.URL_status_visting, title, "")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Get(src)
	if err != nil {
		sglog.Error("download file error,all try failed,src:%s ,error:%e", src, err)
		if res != nil {
			res.Body.Close()
		}
		return
	}
	defer res.Body.Close()
	str := strings.Split(src, "/")

	if len(str) <= 0 {
		return
	}

	rawFileName := str[len(str)-1]
	pdfFileName := yaohaoData.GetPDFDir() + rawFileName

	f, err := os.Create(pdfFileName)
	if err != err {
		sglog.Error("create file error,src:%s,error:%e", src, err)
		return
	}
	io.Copy(f, res.Body)
	sglog.Info("download file %s ,complete,file in %s", src, pdfFileName)

	TransportPDFToTxt(yaohaoData.GetTitle(), rawFileName, pdfFileName)

	txtFileName := strings.Replace(rawFileName, "pdf", "txt", -1)
	txtFileName = yaohaoData.GetTxtDir() + txtFileName

	timestr, memberType, cardType, updateNum := ReadTxtFileAndInsertToDb(txtFileName)

	memberStr := ""
	if memberType == 1 {
		memberStr = "personal"
	} else {
		memberStr = "Company"
	}
	cardStr := ""
	if cardType == 1 {
		cardStr = "normal"
	} else {
		cardStr = "conservation"
	}

	rename := timestr + "_" + memberStr + "_" + cardStr
	sglog.Info("get data success:%s,update total:%d", rename, updateNum)

	//rename and move file
	datePDFDir := yaohaoData.GetFinishPDFDir() + timestr + "/"
	sgfile.AutoMkDir(datePDFDir)
	newPDFFile := datePDFDir + rename + ".pdf"
	sgfile.Rename(pdfFileName, newPDFFile)

	dateTxtFDir := yaohaoData.GetFinishTxtDir() + timestr + "/"
	sgfile.AutoMkDir(dateTxtFDir)
	newTxtile := dateTxtFDir + rename + ".txt"
	sgfile.Rename(txtFileName, newTxtile)

	sgthread.SleepByMillSecond(1000)

	yaohaoData.ChangeDownloadStatus(src, yaohaoDef.URL_status_complete, title, "")
	yaohaoDb.UpdateDbURLData(src, yaohaoDef.URL_status_complete, title, "")

}

func AutoVisitUrl(indexUrl string) {

	sglog.Info("start visit url:%s", indexUrl)

	globalCollector := colly.NewCollector()
	globalCollector.IgnoreRobotsTxt = true
	globalCollector.CheckHead = true
	globalCollector.AllowedDomains = yaohaoData.GetAllowDomains()
	globalCollector.AllowURLRevisit = true
	globalCollector.WithTransport(&http.Transport{
		DisableKeepAlives: true,
	})

	globalCollector.OnHTML("a[href]", func(e *colly.HTMLElement) {

		link := e.Attr("href")
		title := e.Text

		if !strings.Contains(link, yaohaoData.GetHttpType()) {
			link = yaohaoData.GetHttpType() + ":" + link
		}

		sglog.Debug("find tilte:%s,link:%s", title, link)

		if !sgregex.URL(link) {
			sglog.Debug("%s not a valid url,title is %s", link, title)
			return
		}

		if sgstring.EqualWithOr(link, yaohaoData.GetIgnoreUrls()) {
			//sglog.Info("url:%s will be ignore by config", link)
			return
		}

		if !sgstring.ContainsWithOr(title, yaohaoData.GetPageTxt()) && !strings.Contains(title, "下一页") && !strings.Contains(title, "pdf") {
			if _, er := strconv.Atoi(title); er != nil {
				yaohaoData.ChangeUrlStatus(link, yaohaoDef.URL_status_complete)
				return
			}
		}

		if strings.Contains(title, "下一页") {
			yaohaoData.AddIgnoreUrl(link)
		}

		//sglog.Info("Link found:%s,%s", title, link)
		linkLen := len(link)
		if strings.Contains(link, "pdf") && string(link[linkLen-1]) == "f" && string(link[linkLen-2]) == "d" && string(link[linkLen-3]) == "p" {
			sgthread.SleepByMillSecond(500)
			DownloadFile(link, title)
		} else {
			sgthread.SleepByMillSecond(500)
			if yaohaoData.NeedVisitUrl(link) {
				er := e.Request.Visit(e.Request.AbsoluteURL(link))
				if er != nil {
					sglog.Error("start spider error by onHtml,url:=%s,err:=%e", e.Request.AbsoluteURL(link), er)
				}
			}
		}

	})

	globalCollector.OnError(func(r *colly.Response, err error) {
		sglog.Error("visit %s occurt error, went wrong", r.Request.URL.String())
		sgthread.SleepBySecond(60)
		er := r.Request.Visit(r.Request.URL.String())
		if er != nil {
			sglog.Error("start spider error by onError,url:=%s,err:=%e", r.Request.URL.String(), er)
		}
	})

	globalCollector.OnRequest(func(r *colly.Request) {

		sglog.Info("Visiting %s start...", r.URL.String())
		yaohaoData.ChangeUrlStatus(r.URL.String(), yaohaoDef.URL_status_visting)
	})

	globalCollector.OnResponse(func(r *colly.Response) {
		sglog.Info("Visiting %s complete", r.Request.URL.String())
		yaohaoData.ChangeUrlStatus(r.Request.URL.String(), yaohaoDef.URL_status_complete)
	})

	extensions.RandomUserAgent(globalCollector)
	extensions.Referer(globalCollector)

	sleepTime := 60

	for {
		//redownload
		downlist := yaohaoData.GetReDownloadList()
		for _, v := range downlist {
			sglog.Info("re download url:%s", v.Url)
			DownloadFile(v.Url, v.Title+" by reload")
			sgthread.SleepBySecond(1)
		}

		//revisit
		revisitlist := yaohaoData.GetRevisitList()
		for _, v := range revisitlist {
			sglog.Info("re visist url:%s", v)
			globalCollector.Visit(v)
			sgthread.SleepBySecond(1)

		}

		er := globalCollector.Visit(indexUrl)
		if er != nil {
			sglog.Error("start spider error,url:=%s,err:=%e", indexUrl, er)
		}

		nowTime := time.Now()
		timeInt := time.Duration(300) * time.Second
		if 0 == len(downlist) && 0 == len(revisitlist) {

			normalTime := time.Date(nowTime.Year(), nowTime.Month(), 26, 9, 0, 0, 0, nowTime.Location())

			if nowTime.Before(normalTime) {
				timeInt = normalTime.Sub(nowTime)
			} else {
				hour := time.Now().Hour()
				if hour < 9 {
					nextRun := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 9, 0, 0, 0, nowTime.Location())
					timeInt = nextRun.Sub(nowTime)
				} else if hour > 19 {
					nextRun := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 23, 59, 0, 0, nowTime.Location())
					timeInt = nextRun.Sub(nowTime)
				}
			}
			sleepTime = int(timeInt / time.Second)
		}
		sglog.Info("data collection now in sleep,will run after %ds,%s", sleepTime, nowTime.Add(timeInt))
		sgthread.SleepBySecond(sleepTime)
	}

}

func TransPDF(cmdstr []string) {
	sglog.Info("start transport pdf file from cmd")
	if len(cmdstr) < 2 {
		sglog.Error("transPDF at least need 2 params")
		return
	}
	TransportPDFToTxt(yaohaoData.GetTitle(), cmdstr[1], "by cmd")
}

func TransTxt(cmdstr []string) {
	sglog.Info("start transport txt file from cmd")
	if len(cmdstr) < 2 {
		sglog.Error("txt at least need 2 params")
		return
	}
	fileName := yaohaoData.GetTxtDir() + cmdstr[1]
	TransportTxtFileToDb(fileName)
}
