package handler

import (
	"net/http"

	"CloudDIsk/dataService3/internal/logic"
	"CloudDIsk/dataService3/internal/svc"
	"CloudDIsk/dataService3/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DataService3Handler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewDataService3Logic(r.Context(), svcCtx)
		resp, err := l.DataService3(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
