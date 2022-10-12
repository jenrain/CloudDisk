package handler

import (
	"net/http"

	"CloudDIsk/dataService4/internal/logic"
	"CloudDIsk/dataService4/internal/svc"
	"CloudDIsk/dataService4/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DataService4Handler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewDataService4Logic(r.Context(), svcCtx)
		resp, err := l.DataService4(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
