package svc

import (
	"context"
	"fmt"
	"net/url"

	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel"
)

type RepoService struct {
	logx.Logger
	ctx    context.Context
	svcCtx *ServiceContext
	// esUtil *utils.EsUtil
}

func NewRepoService(ctx context.Context, svcCtx *ServiceContext) *RepoService {
	return &RepoService{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (r *RepoService) GetJrogArtifactList(repo commontypes.ImageRepository, repoName, imageName string) (files []commontypes.JfrogFileItem, err error) {
	type response struct {
		Files []commontypes.JfrogFileItem `json:"files"`
	}
	resp := new(response)
	client := utils.NewHttpClient(otel.Tracer("imroc/req"))
	if r.svcCtx.Config.Log.Level == "debug" {
		client.EnableDebug(true)
	}
	if err = client.SetBaseURL(repo.RepoUrl).Get(fmt.Sprintf("/artifactory/api/storage/%s/%s?list&deep=0&listFolders=1", repoName, imageName)).SetHeader("X-JFrog-Art-Api", repo.RepoPassword).SetResult(&resp).Do(context.WithValue(r.ctx, "SpanName", "http.GetJrogArtifactList")).Err; err != nil {
		return
	}
	return resp.Files, nil
}

func (r *RepoService) GetHarborArtifactList(repo commontypes.ImageRepository, repoName, imageName string) (files []commontypes.HarborFileItem, err error) {
	client := utils.NewHttpClient(otel.Tracer("imroc/req"))
	if r.svcCtx.Config.Log.Level == "debug" {
		client.EnableDebug(true)
	}
	imageName = url.QueryEscape(imageName)
	if err = client.SetBaseURL(repo.RepoUrl).Get(fmt.Sprintf("/api/v2.0/projects/%s/repositories/%s/artifacts", repoName, imageName)).SetBasicAuth(repo.RepoAccount, repo.RepoPassword).SetResult(&files).Do(context.WithValue(r.ctx, "SpanName", "http.GetHarborArtifactList")).Err; err != nil {
		return
	}
	return
}
