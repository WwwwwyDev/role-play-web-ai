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

# 启动服务
echo "📦 启动Docker服务..."
docker-compose up -d

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 15

# 检查服务状态
echo "🔍 检查服务状态..."
if [ $(docker-compose ps | wc -l) -eq 1 ]; then
    echo "No containers running, exiting..."
    exit 1
fi

# 等待Ollama服务完全启动
echo "⏳ 等待Ollama服务启动..."
while ! docker exec role-play-ai-ollama ollama list >/dev/null 2>&1; do
    echo "等待Ollama服务就绪..."
    sleep 5
done

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
echo ""
echo "💡 提示："
echo "   - 首次访问可能需要等待模型加载完成"
echo "   - 建议使用Chrome或Edge浏览器以获得最佳语音体验"
echo "   - 如需停止服务，请运行: docker-compose down"
echo ""
