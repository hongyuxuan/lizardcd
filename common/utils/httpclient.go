package utils

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hongyuxuan/lizardcd/common/errorx"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"

	"github.com/imroc/req/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type HttpClient struct {
	*req.Client
}

func NewHttpClient(tracer trace.Tracer) *HttpClient {
	client := req.C().
		SetTimeout(60 * time.Second). // client超时设置为60s
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
						err = errorx.NewError(res.StatusCode, res.String(), nil)
					}
				}()
				ress := make(map[string]interface{})
				res.UnmarshalJson(&ress)
				err = errorx.NewError(res.StatusCode, ress["message"].(string), ress["data"])
			}
			if res.Err != nil {
				err = errorx.NewError(http.StatusInternalServerError, res.Err.Error(), nil)
			}
			return
		}).
		WrapRoundTripFunc(func(rt req.RoundTripper) req.RoundTripFunc {
			return func(req *req.Request) (resp *req.Response, err error) {
				spanName, ok := req.Context().Value(commontypes.TraceIDKey{}).(string)
				if !ok {
					spanName = req.URL.Path
				}
				spanCtx, span := tracer.Start(req.Context(), spanName)
				otel.GetTextMapPropagator().Inject(spanCtx, propagation.HeaderCarrier(req.Headers))
				defer span.End()
				span.SetAttributes(
					attribute.String("http.url", req.URL.String()),
					attribute.String("http.method", req.Method),
					attribute.String("http.request.header", req.HeaderToString()),
				)
				if len(req.Body) > 0 {
					span.SetAttributes(
						attribute.String("http.request.body", string(req.Body)),
					)
				}
				resp, err = rt.RoundTrip(req)
				if err != nil {
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())
					return
				}
				span.SetAttributes(
					attribute.Int("http.status_code", resp.StatusCode),
					attribute.String("http.response.header", resp.HeaderToString()),
					attribute.String("http.response.body", resp.String()),
				)
				return
			}
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
