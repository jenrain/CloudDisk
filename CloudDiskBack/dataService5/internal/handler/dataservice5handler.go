package handler

import (
	"net/http"

	"CloudDIsk/dataService5/internal/logic"
	"CloudDIsk/dataService5/internal/svc"
	"CloudDIsk/dataService5/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DataService5Handler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewDataService5Logic(r.Context(), svcCtx)
		resp, err := l.DataService5(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
