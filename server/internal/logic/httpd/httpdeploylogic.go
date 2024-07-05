package httpd

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/hongyuxuan/lizardcd/common/errorx"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	reqv3 "github.com/imroc/req/v3"
	"github.com/oliveagle/jsonpath"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel"
)

type HttpdeployLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHttpdeployLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HttpdeployLogic {
	return &HttpdeployLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HttpdeployLogic) Httpdeploy(req *types.HttpDeployReq) (resp *types.Response, err error) {
	req.HttpBody = strings.ReplaceAll(req.HttpBody, "{{artifact_url}}", req.ArtifactUrl)
	client := utils.NewHttpClient(otel.Tracer("imroc/req"))
	if l.svcCtx.Config.Log.Level == "debug" {
		client.EnableDebug(true)
	}
	request := client.SetBaseURL(req.HttpUrl).SetCommonHeaders(req.HttpHeader).R()
	if req.HttpContentType == "json" {
		var httpBody map[string]interface{}
		if err = json.Unmarshal([]byte(req.HttpBody), &httpBody); err != nil {
			l.Logger.Error(err)
			return
		}
		request.SetBody(httpBody)
	} else if req.HttpContentType == "x-www-form-urlencoded" {
		var httpBody map[string]string
		if err = json.Unmarshal([]byte(req.HttpBody), &httpBody); err != nil {
			l.Logger.Error(err)
			return
		}
		request.SetFormData(httpBody)
	} else {
		return nil, errorx.NewDefaultError("Unsupported content-type: application/%s", req.HttpContentType)
	}
	var res *reqv3.Response
	if req.HttpMethod == "post" {
		res, err = request.Post(req.HttpPath)
	} else if req.HttpMethod == "put" {
		res, err = request.Put(req.HttpPath)
	} else {
		return nil, errorx.NewDefaultError("Unsupported http method: %s", req.HttpMethod)
	}
	if err != nil {
		l.Logger.Error(err)
		return
	}
	if res.IsError() {
		e := errorx.NewDefaultError("Http deploy return statusCode: %d, data: %v", res.StatusCode, res.String())
		l.Logger.Error(e)
		return nil, e
	}
	var resInterface interface{}
	res.Unmarshal(&resInterface)
	if req.ResJsonpath != "" {
		var lookup interface{}
		if lookup, err = jsonpath.JsonPathLookup(resInterface, req.ResJsonpath); err != nil {
			return
		}
		lookupstr := utils.AnyToString(lookup)
		re, _ := regexp.Compile(req.ResKeyword)
		if !re.MatchString(lookupstr) {
			return nil, errorx.NewDefaultError("Http deploy failed with response: %s", res.String())
		}
	}
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "任务提交成功",
		Data:    resInterface,
	}
	return
}
