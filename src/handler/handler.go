package yaohaoHandle

import (
	"encoding/json"
	"net/http"
	yaohaoData "xcxYaohaoServer/src/data"

	"github.com/coderguang/GameEngine_go/sgstring"

	"github.com/coderguang/GameEngine_go/sglog"
)

//=======================server handler==========================

type wx_xcx_require_handler struct{}

func doCheck(w http.ResponseWriter, r *http.Request, flag chan bool) {

	r.ParseForm()

	if len(r.Form["key"]) <= 0 {
		w.Write([]byte("{\"errcode\":1}")) // not param keys
		sglog.Debug("get require from client,times=%d,no key in this handle", yaohaoData.GetTotalRequireTimes())
		flag <- true
		return
	}

	openid := ""
	if len(r.Form["openid"]) <= 0 {
		sglog.Debug("get require from client,no openid")
	} else {
		openid = r.Form["openid"][0]
		openid = sgstring.RemoveSpaceAndLineEnd(openid)
	}

	yaohaoData.AddTotalRequireTimes()

	if yaohaoData.IsInUpdateCardData() {
		str := "{\"errcode\":3}"
		w.Write([]byte(str))
		sglog.Info("get require from client,times=%d,data are update,please check later,openID:%s", yaohaoData.GetTotalRequireTimes(), openid)
		flag <- true
		return
	}

	key := r.Form["key"][0]
	key = sgstring.RemoveSpaceAndLineEnd(key)

	if key == "time" {
		timeStr := yaohaoData.GetLastesTimeStr()
		w.Write([]byte(timeStr))
		sglog.Info("get require from client,times=%d,someone open,require times,openID:%s", yaohaoData.GetTotalRequireTimes(), openid)
		flag <- true
		return
	}

	if ok, v := yaohaoData.GetMatchData(key); ok {
		jsonBytes, _ := json.Marshal(v)
		str := "{\"errcode\":0,\"data\":" + string(jsonBytes) + "}"
		w.Write([]byte(str))
		sglog.Info("get require from client,times=%d,match data by key %s,openID:%s", yaohaoData.GetTotalRequireTimes(), key, openid)
		flag <- true
		return
	}

	w.Write([]byte("{\"errcode\":2}")) // not find

	sglog.Info("get require from client,times=%d,no this data,key=%s,openID:%s", yaohaoData.GetTotalRequireTimes(), key, openid)
	flag <- true
}

func (h *wx_xcx_require_handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpChan := make(chan bool)
	go doCheck(w, r, tmpChan)
	<-tmpChan
}

func HttpRequireServer(checkPort string) {
	http.Handle("/", &wx_xcx_require_handler{})
	port := "0.0.0.0:" + checkPort
	sglog.Info("start require server.listen port:%s", checkPort)
	http.ListenAndServe(port, nil)
}
