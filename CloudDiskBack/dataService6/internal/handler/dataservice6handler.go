package handler

import (
	"net/http"

	"CloudDIsk/dataService6/internal/logic"
	"CloudDIsk/dataService6/internal/svc"
	"CloudDIsk/dataService6/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DataService6Handler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewDataService6Logic(r.Context(), svcCtx)
		resp, err := l.DataService6(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
