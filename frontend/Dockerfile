# 构建阶段
FROM node:22.19-alpine AS builder

WORKDIR /app

# 复制package文件
COPY package*.json yarn.lock ./

# 清理可能的平台特定依赖
RUN rm -rf node_modules package-lock.json

# 配置yarn镜像源
RUN yarn config set registry https://registry.npmmirror.com

# 安装所有依赖（包括开发依赖，因为构建需要）
RUN yarn install --frozen-lockfile --network-timeout 1000000

# 复制源代码
COPY . .

# 构建应用
RUN yarn build

# 运行阶段
FROM nginx:1.29.1-alpine-slim

# 复制构建的文件到nginx
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制nginx配置
COPY nginx.conf /etc/nginx/nginx.conf

# 暴露端口
EXPOSE 80

# 启动nginx
CMD ["nginx", "-g", "daemon off;"]
