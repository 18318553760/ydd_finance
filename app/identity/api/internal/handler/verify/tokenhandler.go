package verify

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ydd_finance/app/identity/api/internal/logic/verify"
	"ydd_finance/app/identity/api/internal/svc"
	"ydd_finance/app/identity/api/internal/types"
	"ydd_finance/common/response"
)

func TokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.VerifyTokenReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := verify.NewTokenLogic(r.Context(), svcCtx)
		resp, err := l.Token(req,r)
		response.HttpResult(r, w, resp, err)
		//if err != nil {
		//	httpx.Error(w, err)
		//} else {
		//	httpx.OkJson(w, resp)
		//}
	}
}
