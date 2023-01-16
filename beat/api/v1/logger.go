package v1

import "github.com/gogf/gf/v2/frame/g"

type BaseResponseInfo struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
}

type GetLoggerReq struct {
	g.Meta      `path:"/log" method:"get" tags:"fortest" summary:"test summary"`
	ClusterID   string
	ClusterName string
	ProjectID   string
	TimeRange   []int
	ID          string
}

type GetLoggerRes struct {
}
