package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/ljzh152535/go-framework/app/global"
	res "github.com/ljzh152535/go-framework/goadmin-gin/ReturnData"
	"github.com/ljzh152535/go-framework/goadmin-utils/jwtutil"
)

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 登录逻辑
func Login(r *gin.Context) {
	// 1.获取前端传递的用户名和密码
	userInfo := UserInfo{}
	if err := r.ShouldBindJSON(&userInfo); err != nil {
		res.FailWithCodeMessage(res.ErrorCode, err.Error(), r)
		//res.Result(401, nil, err.Error(), r)
		return
	}
	global.GVA_LOG.Debug(map[string]interface{}{"用户名": userInfo.Username, "密码": userInfo.Password}, "开始验证登录信息")

	if userInfo.Username == global.GVA_CONFIG.System.Username && userInfo.Password == global.GVA_CONFIG.System.Password {
		// 认证成功，生成jwt的token
		token, err := jwtutil.GenToken(userInfo.Username, global.GVA_CONFIG.System.JWT_SIGN_KEY, global.GVA_CONFIG.System.JWT_EXPIRE_TIME)
		if err != nil {
			global.GVA_LOG.Error(map[string]interface{}{"用户名": userInfo.Username, "错误信息": err.Error()}, "用户名密码正确,但生成token失败")
			res.FailWithCodeMessage(res.ErrorCode, "生成token失败", r)
			//res.Result(401, nil, "生成token失败", r)
		}
		// token正常生成,返回给前端
		global.GVA_LOG.Info(map[string]interface{}{"用户名": userInfo.Username}, "登录成功")

		resMap := make(map[string]interface{})
		resMap["token"] = token
		res.OK("登录成功", resMap, r)
		//res.Result(200, resMap, "登录成功", r)
		return
	} else {
		// 用户名或密码错误
		//returnData.Status = 401
		//returnData.Message = "用户名或密码错误"
		//r.JSON(200, returnData)
		res.FailWithCodeMessage(res.ErrorCode, "用户名或密码错误", r)
		//res.Result(401, nil, "用户名或密码错误", r)
	}
}

// 登出逻辑
func Logout(r *gin.Context) {
	// 退出
	// 实现退出逻辑
	res.OKWithMessage("退出成功", r)
	//res.Result(200, nil, "退出成功", r)
	global.GVA_LOG.Debug(nil, "用户已退出")
}
