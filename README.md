# 🎭 AI角色扮演聊天网站

<div align="center">

![Vue](https://img.shields.io/badge/Vue-3.x-4FC08D?style=flat-square&logo=vue.js)
![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=go)
![Docker](https://img.shields.io/badge/Docker-支持-2496ED?style=flat-square&logo=docker)
![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)

一个基于AI的角色扮演聊天网站，与历史人物、文学角色、科学家等AI角色进行深度对话



</div>

---

## ✨ 功能特性

<table>
<tr>
<td width="50%">

### 🎭 角色扮演
- **丰富角色库**：哈利·波特、苏格拉底、爱因斯坦等
- **智能对话**：基于Ollama大语言模型
- **个性化体验**：每个角色都有独特的性格和说话方式

### 🎤 语音交互
- **语音输入**：支持实时语音转文字
- **语音输出**：AI回复可转换为语音播放
- **多语言支持**：支持多种语言的语音识别

</td>
<td width="50%">

### 🔐 安全认证
- **JWT认证**：安全的用户认证系统
- **会话管理**：基于Redis的会话存储
- **API限流**：智能限流保护系统安全

### ⚡ 高性能
- **Redis缓存**：多层缓存优化响应速度
- **智能限流**：防止API滥用
- **响应式设计**：完美适配移动端和桌面端

</td>
</tr>
</table>

## 📺 演示视频

[视频下载](./video.mp4) | [问题回答](./question.txt)

<video src="./video.mp4" autoplay="true" controls="controls">
</video>

## 🛠️ 技术栈

<div align="center">

| 层级 | 技术 | 说明 |
|------|------|------|
| **前端** | Vue 3 + Vite | 现代化前端框架 |
| | Tailwind CSS | 原子化CSS框架 |
| | Pinia | 状态管理 |
| | Yarn | 包管理器 |
| **后端** | Go + Gin | 高性能Web框架 |
| | MySQL 8.0 | 关系型数据库 |
| | Redis 7.0 | 缓存和会话存储 |
| | JWT | 用户认证 |
| **AI** | Ollama | 大语言模型服务 |
| | qwen2.5:1.5b | 轻量级AI模型 |
| **部署** | Docker | 容器化部署 |
| | Docker Compose | 多容器编排 |

</div>

## 🚀 快速开始

### 📋 环境要求

- **Docker** & **Docker Compose** (推荐)
- 或手动安装：Go 1.23+, Node.js 18+, MySQL 8.0+, Redis 7.0+

### 🎯 一键部署

```bash
# 克隆项目
git clone https://github.com/WwwwwyDev/role-play-web-ai.git
cd role-play-web-ai

```

#### 🔧 启动服务
```bash

./start.sh

# 访问地址
# 前端: http://localhost:3000
# 后端API: http://localhost:8080
# Swagger文档: http://localhost:8080/swagger/index.html
```

#### 🛑 停止服务
```bash
sh ./stop.sh
```

## 📖 API文档

项目集成了Swagger自动生成API文档，启动服务后访问：

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **JSON格式**: http://localhost:8080/swagger/doc.json
- **YAML格式**: http://localhost:8080/swagger/doc.yaml

### 🔑 认证说明

在Swagger UI中测试需要认证的API：
1. 点击右上角"Authorize"按钮
2. 输入JWT token（格式：`Bearer your-token`）
3. 点击"Authorize"完成认证

## 🏗️ 项目结构

```
role-play-web-ai/
├── 📁 backend/                 # 后端Go服务
│   ├── 📁 internal/           # 内部包
│   │   ├── 📁 config/        # 配置管理
│   │   ├── 📁 database/      # 数据库连接
│   │   ├── 📁 handlers/      # HTTP处理器
│   │   ├── 📁 middleware/    # 中间件
│   │   ├── 📁 models/        # 数据模型
│   │   └── 📁 services/      # 业务逻辑
│   ├── 📁 docs/              # Swagger文档
│   └── 📄 Dockerfile         # 后端Docker配置
├── 📁 frontend/              # 前端Vue应用
│   ├── 📁 src/               # 源代码
│   │   ├── 📁 components/    # Vue组件
│   │   ├── 📁 views/         # 页面组件
│   │   ├── 📁 stores/        # Pinia状态管理
│   │   └── 📁 services/      # API服务
│   └── 📄 Dockerfile         # 前端Docker配置
├── 📁 database/              # 数据库脚本
├── 📄 docker-compose.yml     # 开发环境配置
├── 📄 start.sh              # 启动脚本
└── 📄 stop.sh               # 停止脚本
```

## 🔧 开发指南

### 本地开发

```bash
# 后端开发
cd backend
go mod tidy
go run main.go

# 前端开发
cd frontend
yarn install
yarn dev
```

### 重新生成API文档

```bash
cd backend
go run github.com/swaggo/swag/cmd/swag@latest init
```

### 数据库迁移

```bash
# 查看数据库状态
docker-compose exec mysql mysql -u app_user -p role_play_ai

# 执行SQL脚本
docker-compose exec mysql mysql -u app_user -p role_play_ai < database/schema.sql
```

## 🐳 Docker部署

```bash
docker-compose up -d
```
### 清理环境
```bash
# 停止并删除所有容器
docker-compose down

# 清理Docker缓存
docker system prune -f
```

## 📊 性能优化

- **Redis缓存**：角色数据、对话列表、AI响应缓存
- **API限流**：防止恶意请求，保护系统稳定
- **智能缓存**：AI响应智能缓存，减少重复计算
- **会话管理**：基于Redis的分布式会话存储

## 🤝 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 🙏 致谢

- [Vue.js](https://vuejs.org/) - 渐进式JavaScript框架
- [Gin](https://gin-gonic.com/) - Go Web框架
- [Ollama](https://ollama.ai/) - 本地大语言模型
- [Tailwind CSS](https://tailwindcss.com/) - 原子化CSS框架

---

<div align="center">

**⭐ 如果这个项目对你有帮助，请给它一个星标！**

Made with ❤️ by WwwyDev

</div>