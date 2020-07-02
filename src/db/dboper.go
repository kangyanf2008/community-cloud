package db

import (
	"community-cloud/logging"
	"community-cloud/utils"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func (h *GroupHandler) GetIF(query *QueryDef) ([]map[string]interface{}, error) {
	if len(query.Conds) == 0 && query.RawQuery == "" {
		return nil, errors.New("querykeys is null")
	}

	if h.DbSlave.Status {
		querystr := query.RawQuery
		if querystr == "" {
			colstr := strings.Join(query.Cols, ",")

			var querykeysstr string
			if len(query.Conds) > 1 {
				querykeysstr = strings.Join(query.Conds, " and ")
			} else {
				querykeysstr = query.Conds[0]
			}

			querystr = "select " + colstr + " from " + query.Table + " where " + querykeysstr
		}
		result, err := queryDbIF(h.DbSlave.Db, querystr, query.Vals...)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	return nil, errors.New("GetIF db is offline")
}

func (h *GroupHandler) DelIF(query *DelDef) error {

	if h.DbMaster.Status {
		var querykeysstr string
		if len(query.Conds) > 1 {
			querykeysstr = strings.Join(query.Conds, " and ")
		} else {
			querykeysstr = query.Conds[0]
		}
		querystr := "delete  from " + query.Table + " where " + querykeysstr

		return delDb(h.DbMaster.Db, querystr, query.Vals...)
	}

	return errors.New("DelIF db is offline")
}

func (h *GroupHandler) DelAndDelDef(query *DelAndDelDef) error {

	if h.DbMaster.Status {
		tx, err := h.DbMaster.Db.Begin()
		if err != nil {
			logging.Logger.Errorf("DelAndDelDef db begin Transaction error=[%s] ", err.Error())
		}
		//提交事务
		defer func() {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}()

		//删除表1
		if len(query.Table) > 0 && len(query.Conds) > 0 {
			var querykeysstr string
			if len(query.Conds) > 1 {
				querykeysstr = strings.Join(query.Conds, " and ")
			} else {
				querykeysstr = query.Conds[0]
			}
			querystr := "delete  from " + query.Table + " where " + querykeysstr

			_, err = tx.Exec(querystr, query.Vals...)
			if err != nil {
				logging.Logger.Error("db", "DelAndDelDef fail", query, query.Vals, err.Error())
				return err
			}
		}

		//删除表2
		if len(query.Table2) > 0 && len(query.Conds2) > 0 {
			var querykeysstr string
			if len(query.Conds2) > 1 {
				querykeysstr = strings.Join(query.Conds2, " and ")
			} else {
				querykeysstr = query.Conds2[0]
			}
			querystr := "delete  from " + query.Table2 + " where " + querykeysstr

			_, err = tx.Exec(querystr, query.Vals2...)
			if err != nil {
				logging.Logger.Error("db", "DelAndDelDef fail", query, query.Vals2, err.Error())
				return err
			}
		}
		return nil
	}

	return errors.New("DelAndDelDef db is offline")
}

func (h *GroupHandler) DelDefNoTransaction(query *DelAndDelDef) error {

	if h.DbMaster.Status {
		db := h.DbMaster.Db

		//删除表1
		if len(query.Table) > 0 && len(query.Conds) > 0 {
			var querykeysstr string
			if len(query.Conds) > 1 {
				querykeysstr = strings.Join(query.Conds, " and ")
			} else {
				querykeysstr = query.Conds[0]
			}
			querystr := "delete  from " + query.Table + " where " + querykeysstr

			_, err := db.Exec(querystr, query.Vals...)
			if err != nil {
				logging.Logger.Error("db", "DelDefNoTransaction fail", query, query.Vals, err.Error())
				return err
			}
		}

		//删除表2
		if len(query.Table2) > 0 && len(query.Conds2) > 0 {
			var querykeysstr string
			if len(query.Conds2) > 1 {
				querykeysstr = strings.Join(query.Conds2, " and ")
			} else {
				querykeysstr = query.Conds2[0]
			}
			querystr := "delete  from " + query.Table2 + " where " + querykeysstr

			_, err := db.Exec(querystr, query.Vals2...)
			if err != nil {
				logging.Logger.Error("db", "DelDefNoTransaction fail", query, query.Vals2, err.Error())
				return err
			}
		}
		return nil
	}

	return errors.New("DelDefNoTransaction db is offline")
}

func (h *GroupHandler) DelAndUpdateIF(query *DelAndUpdateDef) error {

	if h.DbMaster.Status {
		tx, err := h.DbMaster.Db.Begin()
		if err != nil {
			logging.Logger.Errorf("DelAndUpdateIF db begin Transaction error=[%s] ", err.Error())
		}
		//提交事务
		defer func() {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}()

		var delCondition string
		if len(query.Conds) > 1 {
			delCondition = strings.Join(query.Conds, " and ")
		} else {
			delCondition = query.Conds[0]
		}
		//删除表
		querystr := "delete  from " + query.Table + " where " + delCondition
		_, err = tx.Exec(querystr, query.Vals...)
		if err != nil {
			logging.Logger.Error("db", "DelAndUpdateIF fail", query, query.Vals, err.Error())
			return err
		}

		//更新表1数据
		if len(query.UpdateTable1) > 0 && len(query.UpdateCols1) > 0 {
			//更新字段
			colstr := strings.Join(query.UpdateCols1, ",")
			//更新条件
			var condstr string
			if len(query.UpdateConds1) > 1 {
				condstr = strings.Join(query.UpdateConds1, " and ")
			} else {
				condstr = query.UpdateConds1[0]
			}
			//拼接更新表语句
			updatestr := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
				query.UpdateTable1, colstr, condstr)
			//执行更新数据操作
			_, err = tx.Exec(updatestr, query.UpdateVals1...)
			if err != nil {
				logging.Logger.Error("db", "DelAndUpdateIF fail", query, query.Vals, err.Error())
				return err
			}
		}
		//更新表2数据
		if len(query.UpdateTable2) > 0 && len(query.UpdateCols2) > 0 {
			//更新字段
			colstr := strings.Join(query.UpdateCols2, ",")
			//更新条件
			var condstr string
			if len(query.UpdateConds2) > 1 {
				condstr = strings.Join(query.UpdateConds2, " and ")
			} else {
				condstr = query.UpdateConds2[0]
			}
			//拼接更新表语句
			updatestr := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
				query.UpdateTable2, colstr, condstr)
			//执行更新数据操作
			_, err = tx.Exec(updatestr, query.UpdateVals2...)
			if err != nil {
				logging.Logger.Error("db", "DelAndUpdateIF fail", query, query.Vals, err.Error())
				return err
			}
		}
		return err
	}

	return errors.New("DelAndUpdateIF db is offline")
}

func (h *GroupHandler) Get(query *QueryDef) ([]map[string]string, error) {
	if len(query.Conds) == 0 && query.RawQuery == "" {
		return nil, errors.New("querykeys is null")
	}

	if h.DbSlave.Status {
		querystr := query.RawQuery
		if querystr == "" {
			colstr := strings.Join(query.Cols, ",")

			var querykeysstr string
			if len(query.Conds) > 1 {
				querykeysstr = strings.Join(query.Conds, " and ")
			} else {
				querykeysstr = query.Conds[0]
			}

			querystr = "select " + colstr + " from " + query.Table + " where " + querykeysstr
		}
		result, err := queryDb(h.DbSlave.Db, querystr, query.Vals...)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	return nil, errors.New("GetOne db is offline")
}

/**
index 开始页从1开始,默认是1开始
pageSize 每页取多少条数据，默认是50条
 */
func (h *GroupHandler) GetLimitPage(query *QueryDef, index, pageSize int32) ([]map[string]string, error) {
	if len(query.Conds) == 0 && query.RawQuery == "" {
		return nil, errors.New("querykeys is null")
	}

	if h.DbSlave.Status {
		querystr := query.RawQuery
		if index <=0 {
			index = 1
		}
		if pageSize <= 0 {
			pageSize  = 50
		}
		if querystr == "" {
			colstr := strings.Join(query.Cols, ",")

			var querykeysstr string
			if len(query.Conds) > 1 {
				querykeysstr = strings.Join(query.Conds, " and ")
			} else {
				querykeysstr = query.Conds[0]
			}

			querystr = "select " + colstr + " from " + query.Table + " where " + querykeysstr + " limit ?,? "
		}
		var values []interface{}
		values = append(values, query.Vals...)
		values = append(values, (index-1) * pageSize, pageSize)
		result, err := queryDb(h.DbSlave.Db, querystr, values...)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	return nil, errors.New("GetOne db is offline")
}

func (h *GroupHandler) GetIn(query *QueryInDef) ([]map[string]string, error) {
	if len(query.Cond) == 0 && len(query.Vals) == 0 {
		return nil, errors.New("querykeys is null")
	}

	if h.DbSlave.Status {
		colstr := strings.Join(query.Cols, ",")
		var condstr string
		if len(query.Vals) > 0 {
			condstr = query.Cond + " in ( "
			for idx, _ := range query.Vals {
				if idx != 0 {
					condstr += (", ?")
				} else {
					condstr += ("?")
				}
			}
			condstr += " )"
		}
		querystr := "select " + colstr + " from " + query.Table + " where " + condstr

		result, err := queryDb(h.DbSlave.Db, querystr, query.Vals...)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	return nil, errors.New("GetOne db is offline")
}

func (h *GroupHandler) GetOne(table, keyname, key, colname string) (string, error) {

	if h.DbSlave.Status {
		query := "select " + colname + " from " + table + " where " + keyname + "=?"
		result, err := queryDb(h.DbSlave.Db, query, key)
		if err != nil {
			return "", err
		}

		if len(result) > 0 {
			return result[0][colname], nil
		}

		return "", errors.New("notfound")
	}

	return "", errors.New("GetOne db is offline")
}

func (h *GroupHandler) GetColumns(table, keyname, key string, colname ...string) (map[string]string, error) {

	if h.DbSlave.Status {
		query := "select " + strings.Join(colname, ",") + " from " + table + " where " + keyname + "=?"
		result, err := queryDb(h.DbSlave.Db, query, key)
		if err != nil {
			return nil, err
		}

		if len(result) > 0 {
			return result[0], nil
		}

		return nil, errors.New("notfound")
	}

	return nil, errors.New("GetOne db is offline")
}

func (h *GroupHandler) GetAll(table, keyname, key string) (map[string]string, error) {

	if h.DbSlave.Status {
		query := "select * from " + table + " where " + keyname + "=?"
		result, err := queryDb(h.DbSlave.Db, query, key)
		if err != nil {
			return nil, err
		}

		if len(result) > 0 {
			return result[0], nil
		}
		return nil, errors.New("GetAll: notfound")
	}

	return nil, errors.New("GetAll: db is offline")
}

func (h *GroupHandler) GetTableColumns(table string, colname ...string) ([]map[string]string, error) {

	var err error
	if h.DbSlave.Status {
		query := "select " + strings.Join(colname, ",") + " from " + table
		var result []map[string]string
		result, err = queryDb(h.DbSlave.Db, query)
		if err != nil {
			return nil, err
		}

		if len(result) > 0 {
			return result, nil
		}

		return nil, errors.New("GetTableColumns: no content")
	}

	return nil, errors.New("GetTableColumns: db is offline")
}

func (h *GroupHandler) GetTable(table string) ([]map[string]string, error) {

	var err error
	if h.DbSlave.Status {
		query := "select * from " + table
		var result []map[string]string
		result, err = queryDb(h.DbSlave.Db, query)
		if err != nil {
			return nil, err
		}

		if len(result) > 0 {
			return result, nil
		}

		return nil, errors.New("GetTable: no content")
	}

	return nil, errors.New("GetTable: db is offline")
}

//执行sql
func (h *GroupHandler) ExeQuerySql(sql string, params []interface{}) ([]map[string]string, error) {
	var err error
	if h.DbSlave.Status {
		var result []map[string]string
		result, err = queryDb(h.DbSlave.Db, sql, params...)
		if err != nil {
			return nil, err
		}
		if len(result) > 0 {
			return result, nil
		}

		return nil, errors.New("ExeQuerySql: no content")
	}
	return nil, errors.New("ExeQuerySql: db is offline")
}

func queryDb(dbconn *sql.DB, query string, args ...interface{}) ([]map[string]string, error) {
	result := make([]map[string]string, 0, 16)

	rows, err := dbconn.Query(query, args...)
	if err != nil {
		logging.Logger.Error("db", "querydb fail", query, args, err.Error())
		return nil, err
	}

	defer rows.Close()
	cols, err := rows.Columns() // Remember to check err afterwards
	vals := make([]interface{}, len(cols))
	valstrs := make([]string, len(cols))
	for i := range cols {
		vals[i] = new(sql.RawBytes)
	}

	for rows.Next() {
		rows.Scan(vals...)
		//logging.Logger.Info("db", "querydb", vals)
		for i, val := range vals {
			valstrs[i] = string(*val.(*sql.RawBytes))
		}

		kvmap, myerr := utils.Zipmap(cols, valstrs)
		if myerr == nil {
			result = append(result, kvmap)
		}
	}

	//logging.Logger.Info("db", "querydb len", len(result))
	return result, nil
}

func queryDbIF(dbconn *sql.DB, query string, args ...interface{}) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0, 16)

	rows, err := dbconn.Query(query, args...)
	if err != nil {
		logging.Logger.Error("db", "queryDbIF fail", query, args, err.Error())
		return nil, err
	}

	defer rows.Close()
	cols, err := rows.Columns() // Remember to check err afterwards
	vals := make([]interface{}, len(cols))
	valPtrs := make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		valPtrs[i] = &vals[i]
	}

	for rows.Next() {
		err := rows.Scan(valPtrs...)
		if err != nil {
			continue
		}
		//logging.Logger.Info("db", "querydb", vals)
		kvmap := make(map[string]interface{}, len(cols))
		for i, val := range vals {
			//fmt.Println(cols[i], val, reflect.TypeOf(val).String())
			kvmap[cols[i]] = val
		}

		result = append(result, kvmap)
	}

	logging.Logger.Info("db", "queryDbIF len", len(result))
	return result, nil
}

func delDb(dbconn *sql.DB, query string, args ...interface{}) error {
	_, err := dbconn.Exec(query, args...)
	if err != nil {
		logging.Logger.Error("db", "delDb fail", query, args, err.Error())
	}
	return err
}
