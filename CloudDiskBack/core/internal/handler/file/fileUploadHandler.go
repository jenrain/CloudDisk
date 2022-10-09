package file

import (
	"core/internal/logic/file"
	"core/internal/svc"
	"core/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func FileUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileUploadRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		l := file.NewFileUploadLogic(r.Context(), svcCtx)
		resp, err := l.FileUpload(&req, r)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
