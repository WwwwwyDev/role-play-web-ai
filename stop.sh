#!/bin/bash

echo "🛑 停止AI角色扮演聊天网站..."

# 停止所有服务
docker-compose down

echo "✅ 所有服务已停止"
echo ""
echo "💡 提示："
echo "   - 如需重新启动，请运行: ./start.sh"
echo "   - 如需完全清理数据，请运行: docker-compose down -v"
echo ""
