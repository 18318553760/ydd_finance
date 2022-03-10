package verify

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ydd_finance/app/identity/api/internal/logic/verify"
	"ydd_finance/app/identity/api/internal/svc"
)

func PingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := verify.NewPingLogic(r.Context(), svcCtx)
		resp, err := l.Ping()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
