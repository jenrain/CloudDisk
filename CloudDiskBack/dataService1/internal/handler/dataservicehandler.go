package handler

import (
	"net/http"

	"CloudDIsk/dataService1/internal/logic"
	"CloudDIsk/dataService1/internal/svc"
	"CloudDIsk/dataService1/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DataServiceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewDataServiceLogic(r.Context(), svcCtx)
		resp, err := l.DataService(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
