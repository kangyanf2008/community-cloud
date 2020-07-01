package model

//单元
type UnitsDef struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

//楼号
type BuildingDef struct {
	Id    int64      `json:"id"`
	Name  string     `json:"name"`
	Units []UnitsDef `json:"units"`
}

//小区
type ParkDef struct {
	Id     int64         `json:"id"`
	Name   string        `json:"name"`
	Builds []BuildingDef `json:"building"`
}

/*type ResponseParkSearch struct {
	ResponseCode
	Parks []ParkDef `json:"data"`
}*/
