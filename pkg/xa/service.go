package xa

import (
	"fmt"
	"gortcenter/pkg/config"

	"github.com/jmoiron/sqlx"
	// 自定义包名，避免与内置 viper 实例冲突
)

type ServiceXa struct {
	MapDb map[string]*sqlx.DB
}

func (c *ServiceXa) Start(mapXid map[string]string) {
	c.xaStart(mapXid)
}

func (c *ServiceXa) Commit(mapXid map[string]string) {
	c.xaEnd(mapXid)
	c.xaPrepare(mapXid)
	c.xaCommit(mapXid)
}

func (c *ServiceXa) Rollback(mapXid map[string]string) {
	c.xaEnd(mapXid)
	c.xaPrepare(mapXid)
	c.xaRollback(mapXid)
}

func (c *ServiceXa) xaStart(mapXid map[string]string) {
	serviceConf := &config.ServiceConfig{}
	for name, xid := range mapXid {
		db, _ := serviceConf.GetDb(name)
		sql := fmt.Sprintf("xa start '%s'", xid)
		db.MustExec(sql)
		c.MapDb[name] = db
	}
}

func (c *ServiceXa) xaEnd(mapXid map[string]string) {
	for name, xid := range mapXid {
		sql := fmt.Sprintf("xa end '%s'", xid)
		c.MapDb[name].MustExec(sql)
	}
}

func (c *ServiceXa) xaPrepare(mapXid map[string]string) {
	for name, xid := range mapXid {
		sql := fmt.Sprintf("xa prepare '%s'", xid)
		c.MapDb[name].MustExec(sql)
	}
}

func (c *ServiceXa) xaCommit(mapXid map[string]string) {
	for name, xid := range mapXid {
		sql := fmt.Sprintf("xa commit '%s'", xid)
		c.MapDb[name].MustExec(sql)
	}
}

func (c *ServiceXa) xaRollback(mapXid map[string]string) {
	for name, xid := range mapXid {
		sql := fmt.Sprintf("xa rollback '%s'", xid)
		c.MapDb[name].MustExec(sql)
	}
}
