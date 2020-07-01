package test

import (
	"community-cloud/model"
	"encoding/json"
	"fmt"
)

func Test() {
	unit1 :=model.UnitsDef{Id: 1,Name:"1单元"}
	unit2 :=model.UnitsDef{Id: 2,Name:"2单元"}

	Building1:=model.BuildingDef{Id:1,Name:"1号楼",Units:[]model.UnitsDef{unit1, unit2}}
	Building2:=model.BuildingDef{Id:2,Name:"2号楼",Units:[]model.UnitsDef{unit1, unit2}}

	parks1 := model.ParkDef{Id:1,Name:"碧桂园", Builds:[]model.BuildingDef{Building1, Building2}}
	parks2 := model.ParkDef{Id:2,Name:"保利", Builds:[]model.BuildingDef{Building1, Building2}}

	reponses := model.ResponseStruct{Data:[]model.ParkDef{parks1,parks2}}
	reponses.Code = 11
	reponses.Desc = "success"

	jsonData,_ := json.Marshal(reponses)
	fmt.Printf("%s \n", jsonData)

	reponses2 := &model.ResponseStruct{}
	fmt.Println(json.Unmarshal(jsonData, reponses2))
	fmt.Println(reponses2)

}
