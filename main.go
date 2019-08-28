package main

import (
	"fmt"
	"log"
	"os"
	yaohaoData "xcxYaohaoServer/src/data"
	yaohaoDb "xcxYaohaoServer/src/db"
	yaohaoHandle "xcxYaohaoServer/src/handler"
	yaohaoSpider "xcxYaohaoServer/src/spider"

	"github.com/coderguang/GameEngine_go/sgcmd"

	"github.com/coderguang/GameEngine_go/sglog"
	"github.com/coderguang/GameEngine_go/sgserver"
)

// cmd func

func RegistCmd() {

	sgcmd.RegistCmd("ShowRequireTimes", "[\"ShowRequireTimes\"] :show current", yaohaoData.ShowRequireTimes)
	sgcmd.RegistCmd("PDF", "[\"PDF\",\"1459152795388.pdf\"]:transform a pdf file to txt file", yaohaoSpider.TransPDF)
	sgcmd.RegistCmd("TXT", "[\"TXT\",\"1459152795388.txt\"]:get txt file and insert it to db", yaohaoSpider.TransTxt)
	sgcmd.RegistCmd("ShowHadVisit", "[\"ShowHadVisit\"]:show all had visit url", yaohaoData.ShowHadVistUrl)
	sgcmd.RegistCmd("ShowHadDownload", "[\"ShowHadDownload\"]:show all had download url", yaohaoData.ShowHadDownVistUrl)
	sgcmd.RegistCmd("ShowNotCompleteUrl", "[\"ShowNotCompleteUrl\"]:show all not complete url", yaohaoData.ShowNotCompleteUrl)
	sgcmd.RegistCmd("ShowIgnoreUrl", "[\"ShowIgnoreUrl\"]:show all ignore url", yaohaoData.ShowIgnoreUrl)
	sgcmd.RegistCmd("AddIgnoreUrl", "[\"AddIgnoreUrl\",\"url1\",\"url2\"]:add ignore url", yaohaoData.AddIgnoreUrlByCmd)
	sgcmd.RegistCmd("RemoveIgnoreUrl", "[\"RemoveIgnoreUrl\",\"url1\",\"url2\"]:remove ignore url", yaohaoData.RemoveIgnoreUrl)
	sgcmd.RegistCmd("ShowLasteTime", "[\"ShowLasteTime\"] :show current", yaohaoData.ShowLasteTime)
}

//==============main=============
func main() {

	arg_num := len(os.Args) - 1

	if arg_num < 1 {
		log.Println("please input config file")
		return
	}
	configfile := os.Args[1]
	log.Println("read global config ...")
	yaohaoData.InitConfig(configfile)
	logPath := "./log/" + yaohaoData.GetTitle()
	sgserver.StartLogServer("debug", logPath, log.LstdFlags, true)
	sglog.Info("start yaohao program...")
	sglog.Info("start connect to db ")

	go yaohaoHandle.HttpRequireServer(yaohaoData.GetListenPort())

	yaohaoDb.InitDbConnection(yaohaoData.GetDbConnectionData())
	yaohaoDb.InitDbURLData()
	yaohaoDb.InitRequireServerDbData()

	go yaohaoSpider.AutoVisitUrl(yaohaoData.GetIndexUrl())

	RegistCmd()

	sgcmd.StartCmdWaitInputLoop()

	fmt.Println("exit by cmd")
}
