package controllers

import (
	"encoding/json"
	"fmt"
	"gortcenter/app/models"
	"gortcenter/pkg/config"
	"gortcenter/pkg/helpers"
	"gortcenter/pkg/response"
	"gortcenter/pkg/xa"
	"log"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type ResetTransController struct {
}

type RTRequest struct {
	TransactId       string   `json:"transact_id"`
	TransactRollback []string `json:"transact_rollback"`
}

type XidMap struct {
	Xid    string
	SqlArr []models.ResetTransactSql
}

func (ctrl *ResetTransController) Commit(c *gin.Context) {
	var sql string
	req := RTRequest{}
	rtEntry := models.ResetTransact{}
	centerConf := &config.CenterConfig{}
	db, err := centerConf.GetDb("rt_center")
	if err != nil {
		log.Fatalln(err)
		return
	}

	c.BindJSON(&req)
	sql = fmt.Sprintf("select * from reset_transact where transact_id='%s'", req.TransactId)
	err = db.Get(&rtEntry, sql)
	if err != nil {
		log.Fatalln(err)
		return
	}

	var rollbackArr []string
	if err := json.Unmarshal([]byte(rtEntry.TransactRollback), &rollbackArr); err != nil {
		log.Fatalln(err)
		return
	}

	rollbackArr = append(rollbackArr, req.TransactRollback...)
	for _, chainId := range rollbackArr {
		sql = fmt.Sprintf("update reset_transact_sql set transact_status=%d where transact_id='%s' and chain_id like '%s%%'", models.STATUS_ROLLBACK, req.TransactId, chainId)
		db.MustExec(sql)
	}

	rtSqlArr := []models.ResetTransactSql{}
	sql = fmt.Sprintf("select * from reset_transact_sql where transact_id='%s' and transact_status in (%d, %d)", req.TransactId, models.STATUS_START, models.STATUS_COMMIT)

	db.Select(&rtSqlArr, sql)

	mapXidSql := make(map[string]XidMap)
	mapXid := make(map[string]string)

	for _, item := range rtSqlArr {
		value, isPresent := mapXidSql[item.Connection]
		if isPresent {
			value.SqlArr = append(value.SqlArr, item)
		} else {
			uuid := helpers.RandomString(32)
			val := XidMap{uuid, []models.ResetTransactSql{item}}
			mapXidSql[item.Connection] = val
			mapXid[item.Connection] = uuid
		}
	}
	// 开启xa
	xa := &xa.ServiceXa{make(map[string]*sqlx.DB)}
	xa.Start(mapXid)
	var wg sync.WaitGroup
	wg.Add(len(mapXid))

	for name, xidMap := range mapXidSql {
		db := xa.MapDb[name]
		go func(db *sqlx.DB, xidMap XidMap) {
			defer wg.Done()
			for _, item := range xidMap.SqlArr {
				db.MustExec(item.Sql)
			}
		}(db, xidMap)

	}
	wg.Wait()

	xa.Commit(mapXid)

	response.Success(c)
}

func (ctrl *ResetTransController) Rollback(c *gin.Context) {
	var sql string
	req := RTRequest{}
	centerConf := &config.CenterConfig{}
	db, _ := centerConf.GetDb("rt_center")
	c.BindJSON(&req)

	if strings.Contains(req.TransactId, "-") {
		chainId := req.TransactId
		transId := strings.Split(req.TransactId, "-")[0]

		for _, trId := range req.TransactRollback {
			sql = fmt.Sprintf("update reset_transact_sql set transact_status=%d where transact_id='%s' and chain_id like '%s%%'", models.STATUS_ROLLBACK, transId, trId)
			db.MustExec(sql)
		}
		sql = fmt.Sprintf("update reset_transact_sql set transact_status=%d where transact_id='%s' and chain_id like '%s%%'", models.STATUS_ROLLBACK, transId, chainId)
		db.MustExec(sql)
	} else {
		sql = fmt.Sprintf("update reset_transact_sql set transact_status=%d where transact_id='%s'", models.STATUS_ROLLBACK, req.TransactId)
		db.MustExec(sql)
		sql = fmt.Sprintf("update reset_transact set action=%d where transact_id='%s'", models.ACTION_ROLLBACK, req.TransactId)
		db.MustExec(sql)
	}

	response.Success(c)
}
