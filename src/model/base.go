package model

type ResponseCode struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
}

type ResponseStruct struct {
	ResponseCode
	Data interface{} `json:"data"`
}