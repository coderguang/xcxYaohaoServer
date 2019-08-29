package yaohaoDb

import (
	"database/sql"
	"strconv"
	yaohaoData "xcxYaohaoServer/src/data"
	yaohaoDef "xcxYaohaoServer/src/define"

	"github.com/coderguang/GameEngine_go/sgstring"

	"github.com/coderguang/GameEngine_go/sglog"
	"github.com/coderguang/GameEngine_go/sgmysql"
)

var globalmysqldb *sql.DB
var globalmysqlstmtHistory *sql.Stmt
var globalmysqlInsertData *sql.Stmt

func InitDbConnection(user string, pwd string, url string, port string, dbname string) {
	conn, err := sgmysql.Open(user, pwd, url, port, dbname, "utf8")
	if err != nil {
		sglog.Fatal("connection to db %s error,%s", url, err)
	}
	globalmysqldb = conn
	historySqlstr := "replace into " + yaohaoData.GetHistoryTableName() + " (url,status,title,tips) values (?,?,?,?);"
	globalmysqlstmtHistory, err = globalmysqldb.Prepare(historySqlstr)
	if err != nil {
		sglog.Fatal("db stmt prepare history sql error,err=%s", err)
	}
	insertSqlstr := "insert into " + yaohaoData.GetDataTableName() + " set type=?,card_type=?,code=?,name=?,time=?,tips=?;"
	globalmysqlInsertData, err = globalmysqldb.Prepare(insertSqlstr)
	if err != nil {
		sglog.Fatal("db stmt prepare insert db data sql error,err=%s", err)
	}
	sglog.Info("InitDbConnection complete")
}

func InitDbURLData() {
	sqlStr := "select * from " + yaohaoData.GetHistoryTableName()
	rows, rowErr := globalmysqldb.Query(sqlStr)
	if rowErr != nil {
		sglog.Error("connect db error,err:%s", rowErr)
		return
	}
	defer rows.Close()
	downloadSize := 0
	for rows.Next() {
		columns, _ := rows.Columns()
		scanArgs := make([]interface{}, len(columns))
		values := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}
		_ = rows.Scan(scanArgs...)

		data := new(yaohaoDef.HistoryUrlData)
		for i, col := range values {
			if col != nil {
				if "url" == columns[i] {
					data.Url = string(col.([]byte))
				}
				if "status" == columns[i] {
					data.Status = string(col.([]byte))
				}
				if "title" == columns[i] {
					data.Title = string(col.([]byte))
				}
				if "tips" == columns[i] {
					data.Desc = string(col.([]byte))
				}
			}
		}
		downloadSize++
		yaohaoData.InitDbURLData(data)
	}

	sglog.Info("init download data from db size:%d", downloadSize)
}

func UpdateDbURLData(url string, status string, title string, desc string) {
	globalmysqlstmtHistory.Exec(url, status, title, desc)
}

func UpdateCardData(data *yaohaoDef.SData) bool {
	sqlStr := "select * from " + yaohaoData.GetDataTableName() + " where code=\"" + data.Code + "\" and name=\"" + data.Name + "\" and time=\"" + data.Time + "\""
	rowsEx, rowErr := globalmysqldb.Query(sqlStr)
	if nil != rowErr {
		rowsEx.Close()
		sglog.Fatal("find data error,code=%s,name=%s,time=%s,err=%s", data.Code, data.Name, data.Time, rowErr)
		return false
	}
	if rowsEx.Next() {
		rowsEx.Close()
		return false
	}

	if rowErr != nil {
		sglog.Error("find doquery error,err:%s", rowErr)
	}
	rowsEx.Close()
	_, err := globalmysqlInsertData.Exec(data.Type, data.CardType, data.Code, data.Name, data.Time, data.Desc)
	if err != nil {
		sglog.Error("insert error,code=%s,name=%s,time=%s,err=%s", data.Code, data.Name, data.Time, rowErr)
		return false
	}
	return true
}

func InitRequireServerDbData() {
	sglog.Info("InitRequireServerDbData:start load data from db")
	sqlStr := "select * from " + yaohaoData.GetDataTableName()
	rows, rowErr := globalmysqldb.Query(sqlStr)
	if rowErr != nil {
		sglog.Fatal("init require server data from db error,connect db error,err:%s", rowErr)
	}

	yaohaoData.SetUpdateFlag(true)

	defer yaohaoData.SetUpdateFlag(false)

	initNum := 0
	for rows.Next() {
		columns, _ := rows.Columns()
		scanArgs := make([]interface{}, len(columns))
		values := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}
		_ = rows.Scan(scanArgs...)

		data := new(yaohaoDef.SData)
		for i, col := range values {
			if col != nil {
				if "name" == columns[i] {
					data.Name = string(col.([]byte))
				}
				if "code" == columns[i] {
					data.Code = string(col.([]byte))
				}
				if "time" == columns[i] {
					data.Time = string(col.([]byte))
				}
				if "type" == columns[i] {
					data.Type, _ = strconv.Atoi(string(col.([]byte)))
				}
				if "card_type" == columns[i] {
					data.CardType, _ = strconv.Atoi(string(col.([]byte)))
				}
				if "desc" == columns[i] {
					data.Desc = string(col.([]byte))
				}
			}
		}

		data.Name = sgstring.RemoveSpaceAndLineEnd(data.Name)
		data.Code = sgstring.RemoveSpaceAndLineEnd(data.Code)

		yaohaoData.AddCardData(data)

		initNum++

	}

	sglog.Info("init data from db ok,total data size %d", initNum)
}
