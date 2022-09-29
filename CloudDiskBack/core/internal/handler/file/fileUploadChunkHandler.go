package file

import (
	"core/tools"
	"errors"
	"net/http"

	"core/internal/svc"
	"core/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FileUploadChunkHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileUploadChunkRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 参数必填判断
		if r.PostForm.Get("name") == "" {
			httpx.Error(w, errors.New("key is empty"))
			return
		}
		if r.PostForm.Get("upload_id") == "" {
			httpx.Error(w, errors.New("upload_id is empty"))
			return
		}
		if r.PostForm.Get("part_number") == "" {
			httpx.Error(w, errors.New("part_number is empty"))
			return
		}
		// 开始分片上传
		etag, err := tools.CosPartUpload(r)
		if err != nil {
			httpx.Error(w, err)
		}

		//l := logic.NewFileUploadChunkLogic(r.Context(), svcCtx)
		//resp, err := l.FileUploadChunk(&req)
		resp := &types.FileUploadChunkResponse{}
		resp.Etag = etag
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
