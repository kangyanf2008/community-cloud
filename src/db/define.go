package db

import (
	"database/sql"
	"sync"
)

type Handler struct {
	Db          *sql.DB
	Status      bool
	Conf        string
	Name        string
	MaxOpenConn int
	MaxIdleConn int
	MaxLifetime int
	rwlock      sync.RWMutex
}

type GroupHandler struct {
	DbMaster *Handler
	DbSlave  *Handler
	Name     string
}

type QueryDef struct {
	Table    string
	Cols     []string
	Conds    []string
	Vals     []interface{}
	Limit    []int
	RawQuery string
}

type QueryInDef struct {
	Table string
	Cols  []string
	Cond  string
	Vals  []interface{}
}

type InsertDef struct {
	Table      string
	UpCols     map[string]interface{}
	InsertCols []string
	InsertVals []interface{}
	RawInsert  string
}

type BatchInsertOrUpdateDef struct {
	Table      string
	UpCols     map[string]interface{}
	InsertCols []string
	InsertVals [][]interface{}
}

type InsertAndUpdateTable2Def struct {
	Table       string
	UpCols      map[string]interface{}
	InsertCols  []string
	InsertVals  []interface{}
	Table2      string
	UpCols2     []string
	UpColsVals2 []interface{}
	UpConds2    []string
	RawInsert   string
}

type UpdateDef struct {
	Table      string
	Conds      []string
	UpdateCols []string
	Vals       []interface{}
	RawUpdate  string
}

type BatchInsertDef struct {
	Table      string
	InsertCols []string
	InsertVals [][]interface{}
}
type BatchInsertAndDelDef struct {
	InsertTable string
	DelTable    string
	InsertCols  []string
	InsertVals  [][]interface{}
	UpCols      map[string]interface{}
	Conds       []string
	CondsVals   []interface{}
}

type BatchInsertOrUpdateAndUpdateDef struct {
	//插入表或更新
	InsertTable string
	InsertCols  []string
	InsertVals  [][]interface{}
	UpCols      map[string]interface{}
	//更新表
	UpdateTable string
	UpdateCols  []string
	UpdateVals  [][]interface{}
	UpdateConds []string
	//删除表
	DelTable string
	DelVals  [][]interface{}
	DelConds []string
}

type BatchInsertOrUpdate struct {
	//插入表或更新
	InsertTable string
	InsertCols  []string
	InsertVals  [][]interface{}
	UpCols      map[string]interface{}
}

type DelDef struct {
	Table string
	Conds []string
	Vals  []interface{}
}

type DelAndDelDef struct {
	Table string
	Conds []string
	Vals  []interface{}

	Table2 string
	Conds2 []string
	Vals2  []interface{}
}

type DelAndUpdateDef struct {
	//删除表
	Table string
	Conds []string
	Vals  []interface{}

	//更新表1
	UpdateTable1 string
	UpdateCols1  []string
	UpdateVals1  []interface{}
	UpdateConds1 []string

	//更新表2
	UpdateTable2 string
	UpdateCols2  []string
	UpdateVals2  []interface{}
	UpdateConds2 []string
}
