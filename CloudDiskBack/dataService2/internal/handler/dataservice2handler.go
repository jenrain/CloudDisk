package handler

import (
	"net/http"

	"CloudDIsk/dataService2/internal/logic"
	"CloudDIsk/dataService2/internal/svc"
	"CloudDIsk/dataService2/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DataService2Handler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewDataService2Logic(r.Context(), svcCtx)
		resp, err := l.DataService2(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
