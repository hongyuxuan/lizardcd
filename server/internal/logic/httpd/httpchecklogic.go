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
	"go.opentelemetry.io/otel"

	"github.com/zeromicro/go-zero/core/logx"
)

type HttpcheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHttpcheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HttpcheckLogic {
	return &HttpcheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HttpcheckLogic) Httpcheck(req *types.HttpCheckReq) (resp *types.Response, err error) {
	client := utils.NewHttpClient(otel.Tracer("imroc/req"))
	if l.svcCtx.Config.Log.Level == "debug" {
		client.EnableDebug(true)
	}
	request := client.SetBaseURL(req.HttpUrl).SetCommonHeaders(req.HttpHeader).R()
	var res *reqv3.Response
	if req.Method == "get" {
		res, err = request.Get(req.HttpPath)
	} else if req.Method == "post" {
		var httpBody map[string]string
		if err = json.Unmarshal([]byte(req.HttpBody), &httpBody); err != nil {
			l.Logger.Error(err)
			return
		}
		res, err = request.SetBody(httpBody).Post(req.HttpPath)
	} else { // no health check
		return &types.Response{
			Code: http.StatusOK,
		}, nil
	}
	if err != nil {
		l.Logger.Error(err)
		return
	}
	if res.IsError() {
		e := errorx.NewDefaultError("Http check return statusCode: %d, data: %v", res.StatusCode, res.String())
		l.Logger.Error(e)
		return nil, e
	}

	var resInterface interface{}
	res.Unmarshal(&resInterface)
	l.Logger.Debugf("Http check response: %s", res.String())
	finished := true
	success := true
	var lookup interface{}
	var message string
	if req.FinishJsonpath != "" {
		if lookup, err = jsonpath.JsonPathLookup(resInterface, req.FinishJsonpath); err != nil {
			return
		}
		lookupstr := utils.AnyToString(lookup)
		re, _ := regexp.Compile(req.FinishKeyword)
		if !re.MatchString(lookupstr) {
			finished = false
		}
	}
	if req.SuccessJsonpath != "" {
		if lookup, err = jsonpath.JsonPathLookup(resInterface, req.SuccessJsonpath); err != nil {
			return
		}
		lookupstr := utils.AnyToString(lookup)
		re, _ := regexp.Compile(req.SuccessKeyword)
		if !re.MatchString(lookupstr) {
			success = false
		}
	}
	if req.MsgJsonpath != "" {
		re, _ := regexp.Compile(`.*(\{\{(\$.*)\}\}).*`)
		matches := re.FindStringSubmatch(req.MsgJsonpath)
		if len(matches) >= 3 {
			var lookup interface{}
			if lookup, err = jsonpath.JsonPathLookup(resInterface, matches[2]); err != nil {
				return
			}
			lookupstr := utils.AnyToString(lookup)
			message = strings.ReplaceAll(req.MsgJsonpath, matches[1], lookupstr)
		}
	}

	resp = &types.Response{
		Code: http.StatusOK,
		Data: types.HttpcheckResponse{
			Finished: finished,
			Success:  success,
			Message:  message,
		},
	}
	return
}
