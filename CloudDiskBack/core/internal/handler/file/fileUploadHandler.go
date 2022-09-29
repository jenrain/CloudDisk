package file

import (
	"core/internal/logic/file"
	"core/models"
	"core/tools"
	"crypto/md5"
	"fmt"
	"net/http"
	"path"
	"strconv"

	"core/internal/svc"
	"core/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FileUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileUploadRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		fmt.Println("开始解析file ")
		f, header, err := r.FormFile("file")

		// 获取parent_id
		r.ParseForm()
		parentId, _ := strconv.Atoi(r.FormValue("parent_id"))
		fmt.Println("parent_id: ", parentId)

		fmt.Println("file: ", header.Filename)
		if err != nil {
			return
		}
		// 判断文件是否存在
		b := make([]byte, header.Size)
		_, err = f.Read(b)
		if err != nil {
			return
		}
		hash := fmt.Sprintf("%x", md5.Sum(b))
		// 从数据库里面找
		var rp models.RepositoryPool
		svcCtx.DB.Where("hash = ?", hash).Find(&rp)
		// 文件存在
		if rp.Hash == hash {
			httpx.OkJson(w, &types.FileUploadResponse{Identity: rp.Identity, Ext: rp.Ext, Name: rp.Name})
			return
		}
		cosPath, err := tools.CosUpload(r)
		if err != nil {
			return
		}
		// 构造req传递给logic
		req.Name = header.Filename
		req.Ext = path.Ext(header.Filename)
		req.Size = header.Size
		req.Hash = hash
		req.Path = cosPath
		req.ParentId = parentId

		l := file.NewFileUploadLogic(r.Context(), svcCtx)
		resp, err := l.FileUpload(&req, r.Header.Get("UserIdentity"))
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
