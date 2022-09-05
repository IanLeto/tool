package apis

import (
	"github.com/gogf/gf/net/ghttp"
	"trace2/model/job"
	"trace2/utils"
)

func QueryJob(r *ghttp.Request) {
	var (
		req = &job.GetJobRequest{}
		err error
	)
	err = r.Parse(req)
	if err != nil {
		r.Response.WriteJson("err")
	}
	client := utils.NewGoOriginClient()
	res, err := client.QueryJob(req.ID)
	if err != nil {
		r.Response.WriteJson("err")
	}
	r.Response.WriteJson(res)
}
