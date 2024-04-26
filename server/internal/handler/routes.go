// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	auth "github.com/hongyuxuan/lizardcd/server/internal/handler/auth"
	db "github.com/hongyuxuan/lizardcd/server/internal/handler/db"
	kubernetes "github.com/hongyuxuan/lizardcd/server/internal/handler/kubernetes"
	lizardcd "github.com/hongyuxuan/lizardcd/server/internal/handler/lizardcd"
	static "github.com/hongyuxuan/lizardcd/server/internal/handler/static"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/services",
				Handler: lizardcd.ListservicesHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/services/:service_name",
				Handler: lizardcd.GetserviceHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/clusters",
				Handler: lizardcd.ListclustersHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/repo/image/tags",
				Handler: lizardcd.ListimagetagsHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/lizardcd"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPatch,
				Path:    "/cluster/:cluster/namespace/:namespace/deployments/:workload_name",
				Handler: kubernetes.PatchDeploymentHandler(serverCtx),
			},
			{
				Method:  http.MethodPatch,
				Path:    "/cluster/:cluster/namespace/:namespace/statefulsets/:workload_name",
				Handler: kubernetes.PatchStatefulsetHandler(serverCtx),
			},
			{
				Method:  http.MethodPatch,
				Path:    "/cluster/:cluster/namespace/:namespace/deployments/:workload_name/rollout",
				Handler: kubernetes.RolloutDeploymentHandler(serverCtx),
			},
			{
				Method:  http.MethodPatch,
				Path:    "/cluster/:cluster/namespace/:namespace/statefulsets/:workload_name/rollout",
				Handler: kubernetes.RolloutStatefulsetHandler(serverCtx),
			},
			{
				Method:  http.MethodPatch,
				Path:    "/cluster/:cluster/namespace/:namespace/deployments/scale",
				Handler: kubernetes.ScaleDeploymentHandler(serverCtx),
			},
			{
				Method:  http.MethodPatch,
				Path:    "/cluster/:cluster/namespace/:namespace/statefulsets/scale",
				Handler: kubernetes.ScaleStatefulsetHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/cluster/:cluster/namespace/:namespace/:resource_type/:resource_name/yaml",
				Handler: kubernetes.GetYamlHandler(serverCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/cluster/:cluster/namespace/:namespace/:resource_type/:resource_name",
				Handler: kubernetes.DeleteResourceHandler(serverCtx),
			},
			{
				Method:  http.MethodPatch,
				Path:    "/cluster/:cluster/namespace/:namespace/apply/yaml",
				Handler: kubernetes.PatchYamlHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/cluster/:cluster/namespace/:namespace/deployments",
				Handler: kubernetes.ListDeploymentHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/cluster/:cluster/namespace/:namespace/statefulsets",
				Handler: kubernetes.ListStatefulsetHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/cluster/:cluster/namespace/:namespace/deployments/:workload_name/pods",
				Handler: kubernetes.DeploymentPodsHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/cluster/:cluster/namespace/:namespace/statefulsets/:workload_name/pods",
				Handler: kubernetes.StatefulsetPodsHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/cluster/:cluster/namespace/:namespace/pods/:pod_name/events",
				Handler: kubernetes.PodEventsHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/kubernetes"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/:tablename",
				Handler: db.ListdataHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/:tablename/:id",
				Handler: db.GetdataHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/:tablename",
				Handler: db.CreatedataHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/:tablename/:id",
				Handler: db.UpdatedataHandler(serverCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/:tablename/:id",
				Handler: db.DeletedataHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/db"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/login",
				Handler: auth.LoginHandler(serverCtx),
			},
		},
		rest.WithPrefix("/auth"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/user/info",
				Handler: auth.UserinfoHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/chpasswd",
				Handler: auth.ChpasswdHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/auth"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/docs/:filename",
				Handler: static.DocfileHandler(serverCtx),
			},
		},
		rest.WithPrefix("/server-static"),
	)
}
