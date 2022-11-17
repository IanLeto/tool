package controller

import (
	v1 "beat/api/v1"
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
)

var (
	Index = cIndex{}
)

type cIndex struct{}

func (c *cIndex) Ping(ctx context.Context, req *v1.PingReq) (res *v1.PongRes, err error) {
	ghttp.RequestFromCtx(ctx).Response.Writeln("pong")
	return
}
