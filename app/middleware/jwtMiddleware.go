// 中间件层
package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/ljzh152535/go-framework/app/global"
	res "github.com/ljzh152535/go-framework/goadmin-gin/ReturnData"
	"github.com/ljzh152535/go-framework/goadmin-utils/jwtutil"
)

func JWTAuth(r *gin.Context) {
	// 1.除了login和logout之外的所有的接口,都要验证请求是否携带token,并且token是否合法
	requestUrl := r.FullPath()
	global.GVA_LOG.Debug(map[string]interface{}{"请求路径": requestUrl}, "")

	if requestUrl == "/api/auth/login" || requestUrl == "/api/auth/logout" {
		global.GVA_LOG.Debug(map[string]interface{}{"请求路径": requestUrl}, "登录和登出不需要验证token")
		r.Next()
		return
	}

	// 验证token的合法性，其他接口需要验证
	// 获取是否携带token
	tokenString := r.Request.Header.Get("Authorization")
	if tokenString == "" {
		// 说明请求没有携带token
		//returnData.Status = 401
		//returnData.Message = "请求未携带token,请登录后尝试"
		//r.JSON(200, returnData)
		res.FailWithCodeMessage(res.ErrorCode, "请求未携带token,请登录后尝试", r)
		//res.Result(401, nil, "请求未携带token,请登录后尝试", r)
		r.Abort()
		return
	}

	// token不为空,要去验证token是否合法
	claims, err := jwtutil.ParseToken(tokenString, global.GVA_CONFIG.System.JWT_SIGN_KEY)
	if err != nil {
		//returnData.Status = 401
		//returnData.Message = "token验证不通过"
		//r.JSON(200, returnData)
		res.FailWithCodeMessage(res.ErrorCode, "token验证不通过", r)
		//res.Result(401, nil, "token验证不通过", r)
		r.Abort() // 终止请求
		return
	}

	// 验证成功
	// 如果其他逻辑需要获取该值可以使用 r.Get()
	r.Set("claims", claims)
	r.Next()
}
