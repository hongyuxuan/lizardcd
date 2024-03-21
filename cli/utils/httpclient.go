package utils

import (
	"errors"
	"strconv"
	"strings"

	"github.com/imroc/req/v3"
)

type HttpClient struct {
	*req.Client
}

func NewHttpClient() *HttpClient {
	client := req.C().
		OnBeforeRequest(func(client *req.Client, req *req.Request) error {
			if req.RetryAttempt > 0 {
				return nil
			}
			req.EnableDump()
			return nil
		}).
		OnAfterResponse(func(client *req.Client, res *req.Response) (err error) {
			responseCode := strconv.Itoa(res.StatusCode)
			if !strings.HasPrefix(responseCode, "2") && !strings.HasPrefix(responseCode, "3") {
				defer func() {
					if e := recover(); e != nil {
						err = res.Err
					}
				}()
				ress := make(map[string]interface{})
				res.UnmarshalJson(&ress)
				err = errors.New(ress["message"].(string))
			}
			if res.Err != nil {
				err = res.Err
			}
			return
		})
	return &HttpClient{
		Client: client,
	}
}

func (c *HttpClient) EnableDebug(enable bool) *HttpClient {
	if enable {
		c.EnableDebugLog()
		c.EnableDumpAll()
	} else {
		c.DisableDebugLog()
		c.DisableDumpAll()
	}
	return c
}
