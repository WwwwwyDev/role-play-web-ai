#!/bin/bash

echo "🚀 启动AI角色扮演聊天网站..."

# 检查Docker是否安装
if ! command -v docker &> /dev/null; then
    echo "❌ Docker未安装，请先安装Docker"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose未安装，请先安装Docker Compose"
    exit 1
fi

# 更新镜像并启动服务
echo "📦 更新并启动Docker服务..."

# 停止现有服务
echo "🛑 停止现有服务..."
docker-compose down

# 重新构建前端和后端镜像
echo "🔨 重新构建前端和后端镜像..."
docker-compose build --no-cache frontend backend

# 启动服务
echo "🚀 启动服务..."
docker-compose up -d

# 等待服务启动
echo "⏳ 等待服务启动..."
echo "🔍 检查服务状态..."

# 等待所有服务启动完成
max_attempts=30
attempt=0

while [ $attempt -lt $max_attempts ]; do
    # 检查所有容器是否都在运行
    running_containers=$(docker-compose ps --services --filter "status=running" | wc -l)
    total_services=$(docker-compose config --services | wc -l)
    
    if [ "$running_containers" -eq "$total_services" ] && [ "$running_containers" -gt 0 ]; then
        echo "✅ 所有服务已启动完成"
        break
    fi
    
    attempt=$((attempt + 1))
    echo "⏳ 等待服务启动... (尝试 $attempt/$max_attempts)"
    sleep 2
done

if [ $attempt -eq $max_attempts ]; then
    echo "❌ 服务启动超时，请检查Docker日志"
    echo "🔍 当前服务状态："
    docker-compose ps
    exit 1
fi

# 检查各个服务的健康状态
echo "🔍 检查服务健康状态..."

# 检查MySQL
echo "📊 检查MySQL数据库..."
mysql_ready=false
for i in {1..10}; do
    if docker exec role-play-ai-mysql mysqladmin ping -h localhost --silent; then
        echo "✅ MySQL数据库已就绪"
        mysql_ready=true
        break
    fi
    echo "⏳ 等待MySQL数据库... (尝试 $i/10)"
    sleep 2
done

if [ "$mysql_ready" = false ]; then
    echo "⚠️ MySQL数据库可能未完全就绪，但继续执行"
fi

# 检查Redis
echo "📊 检查Redis缓存..."
redis_ready=false
for i in {1..10}; do
    if docker exec role-play-ai-redis redis-cli ping | grep -q "PONG"; then
        echo "✅ Redis缓存已就绪"
        redis_ready=true
        break
    fi
    echo "⏳ 等待Redis缓存... (尝试 $i/10)"
    sleep 2
done

if [ "$redis_ready" = false ]; then
    echo "⚠️ Redis缓存可能未完全就绪，但继续执行"
fi

# 检查后端API
echo "📊 检查后端API..."
backend_ready=false
for i in {1..15}; do
    if curl -s http://localhost:8080/api/v1/characters >/dev/null 2>&1; then
        echo "✅ 后端API已就绪"
        backend_ready=true
        break
    fi
    echo "⏳ 等待后端API... (尝试 $i/15)"
    sleep 2
done

if [ "$backend_ready" = false ]; then
    echo "⚠️ 后端API可能未完全就绪，但继续执行"
fi

# 等待Ollama服务完全启动
echo "⏳ 等待Ollama服务启动..."
ollama_attempt=0
ollama_max_attempts=20

while [ $ollama_attempt -lt $ollama_max_attempts ]; do
    if docker exec role-play-ai-ollama ollama list >/dev/null 2>&1; then
        echo "✅ Ollama服务已就绪"
        break
    fi
    
    ollama_attempt=$((ollama_attempt + 1))
    echo "⏳ 等待Ollama服务就绪... (尝试 $ollama_attempt/$ollama_max_attempts)"
    sleep 3
done

if [ $ollama_attempt -eq $ollama_max_attempts ]; then
    echo "❌ Ollama服务启动超时，但继续执行其他步骤"
    echo "💡 提示：您可以稍后手动检查Ollama服务状态"
fi

# 初始化Ollama模型
echo "🤖 初始化AI模型..."
echo "正在检查qwen2.5:1.5b模型..."

# 检查模型是否已存在
if docker exec role-play-ai-ollama ollama list | grep -q "qwen2.5:1.5b"; then
    echo "✅ qwen2.5:1.5b模型已存在"
else
    echo "📥 下载qwen2.5:1.5b模型，这可能需要几分钟时间..."
    docker exec role-play-ai-ollama ollama pull qwen2.5:1.5b
    if [ $? -eq 0 ]; then
        echo "✅ qwen2.5:1.5b模型下载完成"
    else
        echo "❌ qwen2.5:1.5b模型下载失败，请手动执行: docker exec role-play-ai-ollama ollama pull qwen2.5:1.5b"
    fi
fi

echo ""
echo "🎉 启动完成！"
echo ""
echo "📱 前端地址: http://localhost:3000"
echo "🔧 后端API: http://localhost:8080"
echo "🤖 Ollama: http://localhost:11434"
echo "🔍 Swagger UI: http://localhost:8080/swagger/index.html"
echo ""
echo "💡 提示："
echo "   - 首次访问可能需要等待模型加载完成"
echo "   - 建议使用Chrome或Edge浏览器以获得最佳语音体验"
echo "   - 如需停止服务，请运行: docker-compose down"
echo ""
