package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

func (h *GroupHandler) Update(update *UpdateDef) error {
	if h.DbMaster.Status {
		updatestr := update.RawUpdate
		if updatestr == "" {
			colstr := strings.Join(update.UpdateCols, ",")
			var condstr string
			if len(update.Conds) > 1 {
				condstr = strings.Join(update.Conds, " and ")
			} else {
				condstr = update.Conds[0]
			}
			updatestr = fmt.Sprintf("UPDATE %s SET %s WHERE %s",
				update.Table, colstr, condstr)
		}
		err := updateDb(h.DbMaster.Db, updatestr, update.Vals...)
		return err
	}

	return errors.Errorf("%s's master is offline", h.Name)
}

//批量更新表数据
func (h *GroupHandler) Updates(updates []*UpdateDef) error {
	if h.DbMaster.Status {
		if len(updates) > 0 {
			//开启事务
			tx, err := h.DbMaster.Db.Begin()
			//提交事务
			defer func() {
				if err != nil {
					tx.Rollback()
				} else {
					tx.Commit()
				}
			}()

			if err != nil {
				return errors.Wrap(err, "Updates")
				return err
			}
			//循环执行更新表数据操作
			for _, v := range updates {
				colstr := strings.Join(v.UpdateCols, ",")
				var condstr string
				if len(v.Conds) > 1 {
					condstr = strings.Join(v.Conds, " and ")
				} else {
					condstr = v.Conds[0]
				}
				updatestr := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
					v.Table, colstr, condstr)
				_, err := tx.Exec(updatestr, v.Vals...)
				if err != nil {
					return errors.Wrap(err, "Updates")
				}
			}
		}
		return nil
	}
	return errors.Errorf("%s's master is offline", h.Name)
}

func (h *GroupHandler) Insert(insert *InsertDef) error {
	if h.DbMaster.Status {
		insertstr := insert.RawInsert
		vals := insert.InsertVals
		if insertstr == "" {
			colstr := strings.Join(insert.InsertCols, ",")
			valpls := make([]string, len(insert.InsertCols))
			for i := range insert.InsertCols {
				valpls[i] = "?"
			}
			valplstr := strings.Join(valpls, ",")
			insertstr = fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)",
				insert.Table, colstr, valplstr)
		}
		err := insertDb(h.DbMaster.Db, insertstr, vals...)
		return err
	}

	return errors.Errorf("%s's master is offline", h.Name)
}

func (h *GroupHandler) InsertOrUpdate(insert *InsertDef) error {
	if h.DbMaster.Status {
		insertstr := insert.RawInsert
		vals := insert.InsertVals
		if insertstr == "" {
			colstr := strings.Join(insert.InsertCols, ",")
			valpls := make([]string, len(insert.InsertCols))
			upcols := make([]string, len(insert.UpCols))
			j := 0
			for i, colname := range insert.InsertCols {
				valpls[i] = "?"
				if _, ok := insert.UpCols[colname]; ok {
					upcols[j] = colname + " = ?"
					vals = append(vals, insert.InsertVals[i])
					j++
				}
			}
			valplstr := strings.Join(valpls, ",")
			upcolstr := strings.Join(upcols, ",")

			insertstr = fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s) ON DUPLICATE KEY UPDATE %s",
				insert.Table, colstr, valplstr, upcolstr)
		}
		err := insertDb(h.DbMaster.Db, insertstr, vals...)
		return err

		/*
			err := insertDb(h.DbMaster.Db, "INSERT INTO device_token (username, domain, customer_id, pwd, device_token, receive_title_status, app_push, uptime ,addtime) "+
				"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE pwd = ?, uptime = ?,addtime = ?",
				paras["username"], paras["domain"], paras["customer_id"], paras["pwd"], paras["device_token"],
				paras["receivetitlestatus"], paras["apppushstatus"], paras["uptime"], time.Now(),
				paras["pwd"], paras["uptime"], time.Now())

			return err
		*/
	}

	return errors.Errorf("%s's master is offline", h.Name)
}

func (h *GroupHandler) InsertOrUpdateParam(insert *InsertDef) error {
	if h.DbMaster.Status {
		insertstr := insert.RawInsert
		vals := insert.InsertVals
		if insertstr == "" {
			colstr := strings.Join(insert.InsertCols, ",")
			valpls := make([]string, len(insert.InsertCols))
			upcols := make([]string, len(insert.UpCols))
			j := 0
			for i, colname := range insert.InsertCols {
				valpls[i] = "?"
				if value, ok := insert.UpCols[colname]; ok {
					//为true，进行两值相加，否则不进行
					if value.(bool) {
						upcols[j] = colname + " = " + colname + " + ?"
						vals = append(vals, insert.InsertVals[i])
					} else {
						upcols[j] = colname + " = ?"
						vals = append(vals, insert.InsertVals[i])
					}
					j++
				}
			}
			valplstr := strings.Join(valpls, ",")
			upcolstr := strings.Join(upcols, ",")

			insertstr = fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s) ON DUPLICATE KEY UPDATE %s",
				insert.Table, colstr, valplstr, upcolstr)
		}
		err := insertDb(h.DbMaster.Db, insertstr, vals...)
		return err

		/*
			err := insertDb(h.DbMaster.Db, "INSERT INTO device_token (username, domain, customer_id, pwd, device_token, receive_title_status, app_push, uptime ,addtime) "+
				"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE pwd = ?, uptime = ?,addtime = ?",
				paras["username"], paras["domain"], paras["customer_id"], paras["pwd"], paras["device_token"],
				paras["receivetitlestatus"], paras["apppushstatus"], paras["uptime"], time.Now(),
				paras["pwd"], paras["uptime"], time.Now())

			return err
		*/
	}

	return errors.Errorf("%s's master is offline", h.Name)
}

//批量插入或更新
func (h *GroupHandler) BatchInsertOrUpdate(insert *BatchInsertOrUpdateDef) error {
	if !h.DbMaster.Status {
		return errors.Errorf("%s's master is offline", h.Name)
	}

	vals := insert.InsertVals
	if len(vals) > 0 {
		//开启事务
		tx, err := h.DbMaster.Db.Begin()
		//提交事务
		defer func() {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}()

		if err != nil {
			return errors.Wrap(err, "BatchInsertOrUpdate")
			return err
		}

		colstr := strings.Join(insert.InsertCols, ",")

		for _, val := range vals {
			valpls := make([]string, len(insert.InsertCols))
			upcols := make([]string, len(insert.UpCols))

			j := 0
			v := val
			for i, colname := range insert.InsertCols {
				valpls[i] = "?"
				if _, ok := insert.UpCols[colname]; ok {
					upcols[j] = colname + " = ?"
					v = append(v, val[i])
					j++
				}
			}
			valplstr := strings.Join(valpls, ",")
			upcolstr := strings.Join(upcols, ",")

			insertstr := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s) ON DUPLICATE KEY UPDATE %s",
				insert.Table, colstr, valplstr, upcolstr)

			_, err := tx.Exec(insertstr, v...)
			if err != nil {
				return errors.Wrap(err, "BatchInsertOrUpdate")
			}

		}

	}
	return nil
}

func (h *GroupHandler) InsertOrUpdateAndUpdate(insert *InsertAndUpdateTable2Def) error {
	if h.DbMaster.Status {
		tx, err := h.DbMaster.Db.Begin()
		//提交事务
		defer func() {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}()

		//插入或者更新数据
		insertstr := insert.RawInsert
		vals := insert.InsertVals
		if insertstr == "" {
			colstr := strings.Join(insert.InsertCols, ",")
			valpls := make([]string, len(insert.InsertCols))
			upcols := make([]string, len(insert.UpCols))
			j := 0
			for i, colname := range insert.InsertCols {
				valpls[i] = "?"
				if _, ok := insert.UpCols[colname]; ok {
					upcols[j] = colname + " = ?"
					vals = append(vals, insert.InsertVals[i])
					j++
				}
			}
			valplstr := strings.Join(valpls, ",")
			upcolstr := strings.Join(upcols, ",")

			insertstr = fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s) ON DUPLICATE KEY UPDATE %s",
				insert.Table, colstr, valplstr, upcolstr)
		}
		_, err = h.DbMaster.Db.Exec(insertstr, vals...)
		if err != nil {
			return errors.Wrap(err, "insertDB")
		}
		//更新表2
		if insert.Table2 != "" {
			colstr := strings.Join(insert.UpCols2, ",")
			var condstr string
			if len(insert.UpConds2) > 1 {
				condstr = strings.Join(insert.UpConds2, " and ")
			} else {
				condstr = insert.UpConds2[0]
			}
			updatestr := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
				insert.Table2, colstr, condstr)
			_, err := h.DbMaster.Db.Exec(updatestr, insert.UpColsVals2...)
			if err != nil {
				return errors.Wrap(err, "updateDb")
			}
		}

		return err
	}

	return errors.Errorf("%s's master is offline", h.Name)
}

func insertDb(dbMaster *sql.DB, query string, args ...interface{}) error {
	_, err := dbMaster.Exec(query, args...)
	if err != nil {
		//logging.Logger.Error("db insertDb fail:" + err.Error())
		return errors.Wrap(err, "insertDB")
	}

	//rowAffected, err := rs.RowsAffected()
	//if err != nil {
	//	logging.Logger.Error("db insertDb get rowsAffected fail:" + err.Error())
	//	return errors.Wrap(err, "insertDB")
	//}
	//
	//if rowAffected < 1 {
	//	logging.Logger.Error("db insertDb fail rowAffected = " + strconv.Itoa(int(rowAffected)))
	//	return errors.Errorf("insertDb: record is absent")
	//}

	//logging.Logger.Info("db insertDb ", args)

	return nil
}

func updateDb(dbMaster *sql.DB, query string, args ...interface{}) error {
	_, err := dbMaster.Exec(query, args...)
	if err != nil {
		//logging.Logger.Error("db updateDb fail:" + err.Error())
		return errors.Wrap(err, "updateDb")
	}
	//
	//rowAffected, err := rs.RowsAffected()
	//if err != nil {
	//	logging.Logger.Error("db updateDb get rowsAffected fail:" + err.Error())
	//	return errors.Wrap(err, "updateDb")
	//}
	//
	//if rowAffected < 1 {
	//	logging.Logger.Error("db updateDb fail rowAffected = " + strconv.Itoa(int(rowAffected)))
	//	return errors.Errorf("updateDb: record is absent")
	//}

	//logging.Logger.Info("db updateDb ", args)

	return nil
}

func BatchInsertDb(dbMaster *sql.DB, insert *BatchInsertDef) error {
	tx, err := dbMaster.Begin()
	//提交事务
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if err != nil {
		return errors.Wrap(err, "batchInsertDb")
		return err
	}

	errs := make([]error, len(insert.InsertVals))
	colstr := strings.Join(insert.InsertCols, ",")
	for i := 0; i < len(insert.InsertVals); i++ {
		valpls := make([]string, len(insert.InsertCols))
		for i := range insert.InsertCols {
			valpls[i] = "?"
		}
		valplstr := strings.Join(valpls, ",")
		insertstr := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", insert.Table, colstr, valplstr)
		_, err := tx.Exec(insertstr, insert.InsertVals[i]...)
		errs[i] = err
	}
	for _, v := range errs {
		if v != nil {
			err = v
			return err
		}
	}
	return nil
}

func BatchInsertAndDelDb(dbMaster *sql.DB, obj *BatchInsertAndDelDef) error {
	tx, err := dbMaster.Begin()
	//提交事务
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if err != nil {
		return errors.Wrap(err, "batchInsertDb")
		return err
	}

	errs := make([]error, len(obj.InsertVals))
	colstr := strings.Join(obj.InsertCols, ",")
	for i := 0; i < len(obj.InsertVals); i++ {
		valpls := make([]string, len(obj.InsertCols))
		for i := range obj.InsertCols {
			valpls[i] = "?"
		}
		valplstr := strings.Join(valpls, ",")
		insertstr := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", obj.InsertTable, colstr, valplstr)
		_, err := tx.Exec(insertstr, obj.InsertVals[i]...)
		errs[i] = err
	}
	for _, v := range errs {
		if v != nil {
			err = v
			return err
		}
	}
	var querykeysstr string
	if len(obj.Conds) > 1 {
		querykeysstr = strings.Join(obj.Conds, " and ")
	} else {
		querykeysstr = obj.Conds[0]
	}
	querystr := "delete  from " + obj.DelTable + " where " + querykeysstr
	_, err = tx.Exec(querystr, obj.CondsVals...)

	return nil
}

func BatchInsertOrUpdateAndDelDb(dbMaster *sql.DB, obj *BatchInsertAndDelDef) error {
	tx, err := dbMaster.Begin()
	//提交事务
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if err != nil {
		return errors.Wrap(err, "batchInsertDb")
		return err
	}

	errs := make([]error, len(obj.InsertVals))
	colstr := strings.Join(obj.InsertCols, ",")

	for i := 0; i < len(obj.InsertVals); i++ {
		upcols := make([]string, len(obj.UpCols))
		valpls := make([]string, len(obj.InsertCols))
		j := 0
		for idx, colname := range obj.InsertCols {
			valpls[idx] = "?"
			if _, ok := obj.UpCols[colname]; ok {
				upcols[j] = colname + "=?"
				obj.InsertVals[i] = append(obj.InsertVals[i], obj.InsertVals[i][idx])
				j++
			}
		}
		valplstr := strings.Join(valpls, ",")
		upcolstr := strings.Join(upcols, ",")
		insertstr := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s) ON DUPLICATE KEY UPDATE %s", obj.InsertTable, colstr, valplstr, upcolstr)
		_, err := tx.Exec(insertstr, obj.InsertVals[i]...)
		errs[i] = err
	}
	for _, v := range errs {
		if v != nil {
			err = v
			return err
		}
	}
	if len(obj.DelTable) > 0 && len(obj.Conds) > 0 {
		var querykeysstr string
		if len(obj.Conds) > 1 {
			querykeysstr = strings.Join(obj.Conds, " and ")
		} else {
			querykeysstr = obj.Conds[0]
		}
		querystr := "delete  from " + obj.DelTable + " where " + querykeysstr
		_, err = tx.Exec(querystr, obj.CondsVals...)
		return err
	}
	return nil
}

func BatchInsertOrUpdateAndUpdateDb(dbMaster *sql.DB, obj *BatchInsertOrUpdateAndUpdateDef) error {
	tx, err := dbMaster.Begin()
	//提交事务
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if err != nil {
		return errors.Wrap(err, "batchInsertDb")
		return err
	}
	errs := make([]error, len(obj.InsertVals))
	colstr := strings.Join(obj.InsertCols, ",")

	//插入表数据
	for i := 0; i < len(obj.InsertVals); i++ {
		upcols := make([]string, len(obj.UpCols))
		valpls := make([]string, len(obj.InsertCols))
		j := 0
		for idx, colname := range obj.InsertCols {
			valpls[idx] = "?"
			if _, ok := obj.UpCols[colname]; ok {
				upcols[j] = colname + "=?"
				obj.InsertVals[i] = append(obj.InsertVals[i], obj.InsertVals[i][idx])
				j++
			}
		}
		valplstr := strings.Join(valpls, ",")
		upcolstr := strings.Join(upcols, ",")
		insertstr := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s) ON DUPLICATE KEY UPDATE %s", obj.InsertTable, colstr, valplstr, upcolstr)
		_, err := tx.Exec(insertstr, obj.InsertVals[i]...)
		errs[i] = err
	}
	for _, v := range errs {
		if v != nil {
			err = v
			return err
		}
	}
	//更新表数据
	if len(obj.UpdateCols) > 0 && len(obj.UpdateTable) > 0 {
		colstr := strings.Join(obj.UpdateCols, ",")
		var condstr string
		if len(obj.UpdateConds) > 1 {
			condstr = strings.Join(obj.UpdateConds, " and ")
		} else {
			condstr = obj.UpdateConds[0]
		}
		updatestr := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
			obj.UpdateTable, colstr, condstr)

		errs2 := make([]error, len(obj.InsertVals))
		for idx, v := range obj.UpdateVals {
			_, errs2[idx] = tx.Exec(updatestr, v...)
		}
		for _, v := range errs2 {
			if v != nil {
				err = v
				return err
			}
		}
	}

	//删除表记录
	if len(obj.DelTable) > 0 && len(obj.DelConds) > 0 {
		var querykeysstr string
		if len(obj.DelConds) > 1 {
			querykeysstr = strings.Join(obj.DelConds, " and ")
		} else {
			querykeysstr = obj.DelConds[0]
		}
		querystr := "delete  from " + obj.DelTable + " where " + querykeysstr

		//循环执行删除语句
		errs2 := make([]error, len(obj.DelVals))
		for idx, v := range obj.DelVals {
			_, errs2[idx] = tx.Exec(querystr, v...)
		}
		for _, v := range errs2 {
			if v != nil {
				err = v
				return err
			}
		}
	}
	return nil
}

func (h *GroupHandler) BatchInsertOrUpdateNoTransaction(obj *BatchInsertOrUpdate) error {
	if h.DbMaster.Status {
		db := h.DbMaster.Db
		colstr := strings.Join(obj.InsertCols, ",")
		//插入表数据
		for i := 0; i < len(obj.InsertVals); i++ {
			upcols := make([]string, len(obj.UpCols))
			valpls := make([]string, len(obj.InsertCols))
			j := 0
			for idx, colname := range obj.InsertCols {
				valpls[idx] = "?"
				if _, ok := obj.UpCols[colname]; ok {
					upcols[j] = colname + "=?"
					obj.InsertVals[i] = append(obj.InsertVals[i], obj.InsertVals[i][idx])
					j++
				}
			}
			valplstr := strings.Join(valpls, ",")
			upcolstr := strings.Join(upcols, ",")
			insertstr := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s) ON DUPLICATE KEY UPDATE %s", obj.InsertTable, colstr, valplstr, upcolstr)
			_, err := db.Exec(insertstr, obj.InsertVals[i]...)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return errors.Errorf("%s's master is offline", h.Name)
}
