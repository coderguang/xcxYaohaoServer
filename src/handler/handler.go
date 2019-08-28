package yaohaoHandle

import (
	"encoding/json"
	"net/http"
	yaohaoData "wx/xcx_yaohao_server/src/data"

	"github.com/coderguang/GameEngine_go/sgstring"

	"github.com/coderguang/GameEngine_go/sglog"
)

//=======================server handler==========================

type wx_xcx_require_handler struct{}

func doCheck(w http.ResponseWriter, r *http.Request, flag chan bool) {

	sglog.Info("get require from client,times=%d", yaohaoData.GetTotalRequireTimes())
	r.ParseForm()

	if len(r.Form["key"]) <= 0 {
		w.Write([]byte("{\"errcode\":1}")) // not param keys
		sglog.Debug("no key in this handle")
		flag <- true
		return
	}

	yaohaoData.AddTotalRequireTimes()

	if yaohaoData.IsInUpdateCardData() {
		str := "{\"errcode\":3}"
		w.Write([]byte(str))
		sglog.Info("data are update,please check later")
		flag <- true
		return
	}

	key := r.Form["key"][0]
	key = sgstring.RemoveSpaceAndLineEnd(key)

	if key == "time" {
		timeStr := yaohaoData.GetLastesTimeStr()
		w.Write([]byte(timeStr))
		sglog.Info("someone open,require times")
		flag <- true
		return
	}

	if ok, v := yaohaoData.GetMatchData(key); ok {
		jsonBytes, _ := json.Marshal(v)
		str := "{\"errcode\":0,\"data\":" + string(jsonBytes) + "}"
		w.Write([]byte(str))
		sglog.Info("match data by key %s", key)
		flag <- true
		return
	}

	w.Write([]byte("{\"errcode\":2}")) // not find

	sglog.Info("no this data,key=%s", key)
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
