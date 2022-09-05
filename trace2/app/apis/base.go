package apis

import "github.com/gogf/gf/net/ghttp"

func Healthz(r *ghttp.Request) {
	_ = r.Response.WriteJson("ok")
}
