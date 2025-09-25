# AI角色扮演聊天网站

一个基于AI的角色扮演聊天网站，用户可以与各种历史人物、文学角色、科学家等AI角色进行深度对话。支持文本和语音聊天功能。

## 功能特性

- 🎭 **丰富角色库**：哈利·波特、苏格拉底、爱因斯坦、达芬奇等
- 💬 **智能对话**：基于Ollama大语言模型的智能对话
- 🎤 **语音聊天**：支持语音输入和输出
- 🔐 **用户认证**：JWT认证，安全的用户系统
- 📱 **响应式设计**：现代化的UI设计，支持移动端
- 🚀 **Docker部署**：一键部署，开箱即用

## 技术栈

### 后端
- **Go** - 主要后端语言
- **Gin** - Web框架
- **MySQL** - 数据库
- **JWT** - 用户认证
- **Ollama** - AI大语言模型

### 前端
- **Vue 3** - 前端框架
- **Vite** - 构建工具
- **Tailwind CSS** - 样式框架
- **Pinia** - 状态管理
- **Vue Router** - 路由管理

## 快速开始

### 环境要求

- Docker & Docker Compose
- 或者手动安装：Go 1.21+, Node.js 18+, MySQL 8.0+

### 使用Docker部署（推荐）

1. **克隆项目**
```bash
git clone <repository-url>
cd role-play-web-ai
```

2. **启动服务**
```bash
docker-compose up -d
```

3. **初始化Ollama模型**
```bash
# 拉取模型（首次运行需要）
docker exec role-play-ai-ollama ollama pull llama2

# 或者使用其他模型
docker exec role-play-ai-ollama ollama pull qwen2.5:7b
```

4. **访问应用**
- 前端：http://localhost:3000
- 后端API：http://localhost:8080
- Ollama：http://localhost:11434

### 手动部署

#### 后端部署

1. **安装依赖**
```bash
cd backend
go mod download
```

2. **配置环境变量**
```bash
cp config.env.example .env
# 编辑.env文件，配置数据库和JWT密钥
```

3. **初始化数据库**
```bash
mysql -u root -p < ../database/schema.sql
```

4. **启动服务**
```bash
go run main.go
```

#### 前端部署

1. **安装依赖**
```bash
cd frontend
npm install
```

2. **启动开发服务器**
```bash
npm run dev
```

3. **构建生产版本**
```bash
npm run build
```

## 配置说明

### 环境变量

#### 后端配置
```env
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=role_play_ai

# JWT配置
JWT_SECRET=your-secret-key-here

# 服务器配置
PORT=8080

# Ollama配置
OLLAMA_BASE_URL=http://localhost:11434
OLLAMA_MODEL=llama2
```

### 数据库结构

- `users` - 用户表
- `characters` - 角色表
- `conversations` - 对话会话表
- `messages` - 消息表

## API文档

### 认证接口

- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `GET /api/v1/auth/me` - 获取用户信息

### 角色接口

- `GET /api/v1/characters` - 获取角色列表
- `GET /api/v1/characters/search?q=关键词` - 搜索角色
- `GET /api/v1/characters/:id` - 获取角色详情

### 对话接口

- `GET /api/v1/conversations` - 获取对话列表
- `POST /api/v1/conversations` - 创建新对话
- `GET /api/v1/conversations/:id` - 获取对话详情
- `POST /api/v1/conversations/:id/messages` - 发送消息
- `DELETE /api/v1/conversations/:id` - 删除对话

## 开发指南

### 项目结构

```
role-play-web-ai/
├── backend/                 # Go后端
│   ├── internal/
│   │   ├── config/         # 配置管理
│   │   ├── database/       # 数据库连接
│   │   ├── handlers/       # HTTP处理器
│   │   ├── middleware/     # 中间件
│   │   ├── models/         # 数据模型
│   │   └── services/       # 业务逻辑
│   ├── main.go            # 主入口
│   └── Dockerfile         # Docker配置
├── frontend/               # Vue前端
│   ├── src/
│   │   ├── components/    # 组件
│   │   ├── views/         # 页面
│   │   ├── stores/        # 状态管理
│   │   ├── services/      # API服务
│   │   └── router/        # 路由配置
│   └── Dockerfile         # Docker配置
├── database/              # 数据库脚本
└── docker-compose.yml     # Docker编排
```

### 添加新角色

1. 在数据库中添加角色记录：
```sql
INSERT INTO characters (name, description, avatar_url, system_prompt, category) 
VALUES ('角色名', '角色描述', '/avatars/avatar.svg', '系统提示词', '分类');
```

2. 系统会自动加载新角色到前端界面。

### 自定义AI模型

修改环境变量中的`OLLAMA_MODEL`：
```env
OLLAMA_MODEL=qwen2.5:7b  # 使用通义千问模型
OLLAMA_MODEL=llama2      # 使用Llama2模型
```

## 故障排除

### 常见问题

1. **Ollama连接失败**
   - 确保Ollama服务正在运行
   - 检查模型是否已下载：`ollama list`

2. **数据库连接失败**
   - 检查MySQL服务状态
   - 验证数据库配置信息

3. **前端无法连接后端**
   - 检查后端服务是否启动
   - 确认端口配置正确

### 日志查看

```bash
# 查看所有服务日志
docker-compose logs

# 查看特定服务日志
docker-compose logs backend
docker-compose logs frontend
docker-compose logs mysql
docker-compose logs ollama
```

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 联系方式

如有问题或建议，请通过以下方式联系：

- 提交 Issue

---