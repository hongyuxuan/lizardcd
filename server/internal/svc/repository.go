package svc

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/hongyuxuan/lizardcd/common/errorx"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel"
)

type RepoService struct {
	logx.Logger
	ctx    context.Context
	svcCtx *ServiceContext
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
	if err = client.SetBaseURL(repo.RepoUrl).Get(fmt.Sprintf("/artifactory/api/storage/%s/%s?list&deep=0&listFolders=1", repoName, imageName)).SetHeader("X-JFrog-Art-Api", repo.RepoPassword).SetResult(&resp).Do(context.WithValue(r.ctx, commontypes.TraceIDKey{}, "http.GetJrogArtifactList")).Err; err != nil {
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
	if err = client.SetBaseURL(repo.RepoUrl).Get(fmt.Sprintf("/api/v2.0/projects/%s/repositories/%s/artifacts", repoName, imageName)).SetBasicAuth(repo.RepoAccount, repo.RepoPassword).SetResult(&files).Do(context.WithValue(r.ctx, commontypes.TraceIDKey{}, "http.GetHarborArtifactList")).Err; err != nil {
		return
	}
	return
}

func (r *RepoService) GetDockerHubImages(repo commontypes.ImageRepository, repoName, imageName, tag string) (files []commontypes.DockerHubImageItem, err error) {
	type response struct {
		Count   int                              `json:"count"`
		Results []commontypes.DockerHubImageItem `json:"results"`
	}
	resp := new(response)
	client := utils.NewHttpClient(otel.Tracer("imroc/req"))
	if r.svcCtx.Config.Log.Level == "debug" {
		client.EnableDebug(true)
	}
	imageName = url.QueryEscape(imageName)
	url := fmt.Sprintf("/v2/namespaces/%s/repositories/%s/tags", repoName, imageName)
	if tag != "" {
		url += "?name=" + tag
	}
	cli := client.SetBaseURL(repo.RepoUrl).DisableInsecureSkipVerify()
	if httpProxy := os.Getenv("HTTP_PROXY"); httpProxy != "" {
		r.Logger.Infof("Using http_proxy=%s for fetching dockerhub images", httpProxy)
		cli = cli.SetProxyURL(httpProxy)
	}
	if err = cli.Get(url).SetBearerAuthToken(repo.RepoPassword).SetResult(&resp).Do(context.WithValue(r.ctx, commontypes.TraceIDKey{}, "http.GetDockerHubImages")).Err; err != nil {
		return
	}
	return resp.Results, nil
}

func (r *RepoService) GetS3ArtifactList(repo commontypes.ImageRepository, repoName, imageName string) (files []commontypes.ArtifactListRes, err error) {
	re, _ := regexp.Compile(`(http|https)://(.*)`)
	matches := re.FindStringSubmatch(repo.RepoUrl)
	if len(matches) != 3 {
		return nil, errorx.NewDefaultError("repo \"%s\" is not valid", repo.RepoUrl)
	}
	useSSL := false
	if matches[1] == "https" {
		useSSL = true
	}
	var client *minio.Client
	if client, err = minio.New(matches[2], &minio.Options{
		Creds:  credentials.NewStaticV4(repo.RepoAccount, repo.RepoPassword, ""),
		Secure: useSSL,
	}); err != nil {
		r.Logger.Error(err)
		return
	}
	objects := client.ListObjects(context.WithValue(r.ctx, commontypes.TraceIDKey{}, "http.GetS3ArtifactList"), repoName, minio.ListObjectsOptions{
		Prefix:    imageName,
		Recursive: true,
	})
	for o := range objects {
		if o.Err == nil {
			d, _ := time.ParseDuration("8h")
			files = append(files, commontypes.ArtifactListRes{
				ArtifactUrl:  fmt.Sprintf("%s/%s/%s", repo.RepoUrl, repoName, o.Key),
				LastModified: o.LastModified.Add(d).Format("2006-01-02 15:04:05"),
				Tag:          o.Key[len(imageName)+1:],
			})
		} else {
			r.Logger.Error(o.Err)
		}
	}
	return
}
