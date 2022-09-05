package utils

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"trace2/model/job"
)

type GoOriginClient struct {
	client *resty.Client
}

// NewGoOriginClient 目前不读配置文件，写死在代码中
func NewGoOriginClient() *GoOriginClient {
	return &GoOriginClient{client: resty.New()}
}

func (q GoOriginClient) QueryJob(id string) (*job.GetJobResponse, error) {
	var res = &job.GetJobResponse{}
	resp, err := q.client.R().SetResult(res).SetQueryParams(map[string]string{
		"id": id,
	}).Get(fmt.Sprintf("http://localhost:8008/v1/job/6"))
	//resp, err := q.client.R().SetResult(res).SetQueryParams(map[string]string{
	//	"id": id,
	//}).Get("http:124.222.48.125:8008/v1/job")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if resp.StatusCode() != 200 {
		fmt.Println(resp.StatusCode())
		return nil, err
	}
	return res, err

}
