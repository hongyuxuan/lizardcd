package lizardcd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/hongyuxuan/lizardcd/common/constant"
	"github.com/hongyuxuan/lizardcd/common/errorx"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListimagetagsLogic struct {
	logx.Logger
	ctx         context.Context
	svcCtx      *svc.ServiceContext
	repoService *svc.RepoService
}

func NewListimagetagsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListimagetagsLogic {
	return &ListimagetagsLogic{
		Logger:      logx.WithContext(ctx),
		ctx:         ctx,
		svcCtx:      svcCtx,
		repoService: svc.NewRepoService(ctx, svcCtx),
	}
}

func (l *ListimagetagsLogic) Listimagetags(req *types.ListTagsReq) (resp *types.Response, err error) {
	// get application by app_name
	var application *commontypes.Application
	if err = l.svcCtx.Database.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.GetApplication")).
		Model(&commontypes.Application{}).
		Where("app_name = ?", req.AppName).
		Find(&application).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		err = errorx.NewDefaultError("Application \"%s\" not found", req.AppName)
		l.Logger.Error(err)
		return
	}
	var repo commontypes.ImageRepository
	l.svcCtx.Database.Model(&commontypes.ImageRepository{}).Where("id = ?", application.RepoId).First(&repo)
	var artifactList []commontypes.ArtifactListRes
	if repo.RepoType == constant.REPO_TYPE_ARTIFACTORY {
		var fileList []commontypes.JfrogFileItem
		if fileList, err = l.repoService.GetJrogArtifactList(repo, application.RepoName, application.ImageName); err != nil {
			l.Logger.Errorf("Jforg failed to list %s, %v", application.RepoName, err)
			return
		}
		for _, item := range fileList {
			if strings.Contains(item.Uri, "sha256") {
				continue
			}
			reg := regexp.MustCompile(`http[s]{0,1}://(.+)`)
			matches := reg.FindStringSubmatch(repo.RepoUrl)
			artifact_url := fmt.Sprintf("%s/%s/%s:%s", matches[1], application.RepoName, application.ImageName, item.Uri[1:])
			if application.DeployType != "容器" && application.DeployType != "Docker" {
				artifact_url = fmt.Sprintf("%s/artifactory/%s/%s%s", repo.RepoUrl, application.RepoName, application.ImageName, item.Uri)
			}
			artifactList = append(artifactList, commontypes.ArtifactListRes{
				ArtifactUrl:  artifact_url,
				LastModified: item.LastModified,
				Tag:          item.Uri[1:],
			})
		}
	}
	if repo.RepoType == constant.REPO_TYPE_HARBOR {
		var fileList []commontypes.HarborFileItem
		if fileList, err = l.repoService.GetHarborArtifactList(repo, application.RepoName, application.ImageName); err != nil {
			l.Logger.Errorf("Harbor failed to list %s, %v", application.RepoName, err)
			return
		}
		for _, item := range fileList {
			reg := regexp.MustCompile(`http[s]{0,1}://(.+)`)
			matches := reg.FindStringSubmatch(repo.RepoUrl)
			artifactList = append(artifactList, commontypes.ArtifactListRes{
				ArtifactUrl:  fmt.Sprintf("%s/%s/%s:%s", matches[1], application.RepoName, application.ImageName, item.Tags[0].Name),
				LastModified: item.Tags[0].PushTime,
				Tag:          item.Tags[0].Name,
			})
		}
	}
	if repo.RepoType == constant.REPO_TYPE_DOCKERHUB {
		var fileList []commontypes.DockerHubImageItem
		if fileList, err = l.repoService.GetDockerHubImages(repo, application.RepoName, application.ImageName, req.Tag); err != nil {
			l.Logger.Errorf("DockerHub failed to fetch %s/%s, %v", application.RepoName, application.ImageName, err)
			return
		}
		for _, item := range fileList {
			artifactUrl := fmt.Sprintf("%s/%s:%s", application.RepoName, application.ImageName, item.Name)
			if application.RepoName == "library" {
				artifactUrl = fmt.Sprintf("%s:%s", application.ImageName, item.Name)
			}
			artifactList = append(artifactList, commontypes.ArtifactListRes{
				ArtifactUrl:  artifactUrl,
				LastModified: item.LastUpdated,
				Tag:          item.Name,
			})
		}
	}
	if repo.RepoType == constant.REPO_TYPE_S3 {
		if artifactList, err = l.repoService.GetS3ArtifactList(repo, application.RepoName, application.ImageName); err != nil {
			return
		}
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: artifactList,
	}
	return
}
