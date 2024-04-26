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
	if err = l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, "SpanName", "tidb.GetApplication")).
		Model(&commontypes.Application{}).
		Where("app_name = ?", req.AppName).
		First(&application).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		err = errorx.NewDefaultError(fmt.Sprintf("Application \"%s\" not found", req.AppName))
		l.Logger.Error(err)
		return
	}
	var artifactList []commontypes.ArtifactListRes
	if application.Repo.RepoType == constant.REPO_TYPE_ARTIFACTORY {
		var fileList []commontypes.JfrogFileItem
		if fileList, err = l.repoService.GetJrogArtifactList(application.Repo, application.RepoName, application.ImageName); err != nil {
			l.Logger.Errorf("Jforg failed to list %s, %v", application.RepoName, err)
			return
		}
		for _, item := range fileList {
			if strings.Contains(item.Uri, "sha256") {
				continue
			}
			reg := regexp.MustCompile(`http[s]{0,1}://(.+)`)
			matches := reg.FindStringSubmatch(application.Repo.RepoUrl)
			artifactList = append(artifactList, commontypes.ArtifactListRes{
				ArtifactUrl:  fmt.Sprintf("%s/%s/%s:%s", matches[1], application.RepoName, application.ImageName, item.Uri[1:]),
				LastModified: item.LastModified,
				Tag:          item.Uri[1:],
			})
		}
	}
	if application.Repo.RepoType == constant.REPO_TYPE_HARBOR {
		var fileList []commontypes.HarborFileItem
		if fileList, err = l.repoService.GetHarborArtifactList(application.Repo, application.RepoName, application.ImageName); err != nil {
			l.Logger.Errorf("Harbor failed to list %s, %v", application.RepoName, err)
			return
		}
		for _, item := range fileList {
			reg := regexp.MustCompile(`http[s]{0,1}://(.+)`)
			matches := reg.FindStringSubmatch(application.Repo.RepoUrl)
			artifactList = append(artifactList, commontypes.ArtifactListRes{
				ArtifactUrl:  fmt.Sprintf("%s/%s/%s:%s", matches[1], application.RepoName, application.ImageName, item.Tags[0].Name),
				LastModified: item.Tags[0].PushTime,
				Tag:          item.Tags[0].Name,
			})
		}
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: artifactList,
	}
	return
}
