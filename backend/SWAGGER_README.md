# Swagger API 文档

本项目已集成Swagger自动生成API文档功能。

## 访问文档

启动服务后，可以通过以下地址访问Swagger文档：

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **JSON格式**: http://localhost:8080/swagger/doc.json
- **YAML格式**: http://localhost:8080/swagger/doc.yaml

## 重新生成文档

当修改了API注释后，需要重新生成Swagger文档：

```bash
# 进入backend目录
cd backend

# 重新生成文档
go run github.com/swaggo/swag/cmd/swag@latest init

# 或者如果已安装swag工具
swag init
```

## 添加新的API注释

为新的API端点添加Swagger注释，例如：

```go
// GetUsers 获取用户列表
// @Summary 获取用户列表
// @Description 获取所有用户的列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "用户列表"
// @Failure 401 {object} map[string]string "未授权"
// @Router /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
    // 处理逻辑
}
```

## 注释说明

- `@Summary`: API简要描述
- `@Description`: API详细描述
- `@Tags`: API分组标签
- `@Accept`: 接受的请求格式
- `@Produce`: 返回的响应格式
- `@Security`: 安全认证方式
- `@Param`: 请求参数说明
- `@Success`: 成功响应说明
- `@Failure`: 失败响应说明
- `@Router`: 路由路径和HTTP方法

## 认证说明

本项目使用JWT Bearer Token认证，在Swagger UI中：

1. 点击右上角的"Authorize"按钮
2. 在Bearer Token字段中输入JWT token
3. 点击"Authorize"完成认证
4. 现在可以测试需要认证的API端点

## 注意事项

- 修改API注释后需要重新生成文档
- 确保所有导入的包都正确
- 生成的docs包会自动导入到main.go中
