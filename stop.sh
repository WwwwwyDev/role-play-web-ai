#!/bin/bash

echo "🛑 停止AI角色扮演聊天网站..."

# 检查Docker是否安装
if ! command -v docker &> /dev/null; then
    echo "❌ Docker未安装"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose未安装"
    exit 1
fi

# 停止服务
echo "📦 停止Docker服务..."
docker-compose down

echo "✅ 服务已停止"
echo ""
echo "💡 提示："
echo "   - 如需重新启动，请运行: ./start.sh"
echo ""
