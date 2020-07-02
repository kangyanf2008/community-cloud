package db

import (
	"community-cloud/logging"
	"time"
)

func checkAndRepair(handler *Handler) bool {
	oldStatus := handler.Status
	if handler.Db != nil {
		err := pingDB(handler.Db)
		if err != nil {
			handler.Status = false
		} else {
			handler.Status = true
		}
	} else {
		connDB(handler)
	}

	return oldStatus != handler.Status
}

func checkChange(groupHandler *GroupHandler, count int) {
	mchange := checkAndRepair(groupHandler.DbMaster)
	schange := checkAndRepair(groupHandler.DbSlave)

	if mchange || schange || count%360 == 1 {
		logging.Logger.Infof("db monitor %s %s:%v %s:%v", groupHandler.Name, "master", groupHandler.DbMaster.Status, "slave", groupHandler.DbSlave.Status)
	}
}

func Monitor(groupHandler *GroupHandler) {
	ticker := time.NewTicker(1 * time.Minute)
	count := 1
	checkChange(groupHandler, count)
	for {
		count += 1
		select {
		case <-ticker.C:
			checkChange(groupHandler, count)
		}
	}
}
