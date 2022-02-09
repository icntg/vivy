package initialize

import "github.com/gin-gonic/gin"

func Initialize(c *gin.Context) {
	// todo:
	// 如果数据库连接测试失败。则自动生成授权码。
	// 验证授权码之后，才能进行初始化工作。
	// 初始化保存的账号密码，保存在另外的配置文件中。
}
