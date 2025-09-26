-- 数据库索引优化脚本
-- 基于查询模式分析，为常用查询字段添加索引以提升性能

-- 用户表索引
-- email字段用于登录验证，需要快速查找
CREATE INDEX idx_users_email ON users(email);

-- 角色表索引
-- category字段用于角色分类筛选
CREATE INDEX idx_characters_category ON characters(category);

-- 对话表索引
-- user_id字段用于获取用户的对话列表，这是最频繁的查询
CREATE INDEX idx_conversations_user_id ON conversations(user_id);

-- user_id + updated_at 复合索引，用于按更新时间排序获取用户对话
CREATE INDEX idx_conversations_user_updated ON conversations(user_id, updated_at DESC);

-- character_id字段用于角色相关的对话查询
CREATE INDEX idx_conversations_character_id ON conversations(character_id);

-- 消息表索引
-- conversation_id字段用于获取对话的所有消息，这是最频繁的查询
CREATE INDEX idx_messages_conversation_id ON messages(conversation_id);

-- conversation_id + created_at 复合索引，用于按时间顺序获取消息
CREATE INDEX idx_messages_conversation_created ON messages(conversation_id, created_at ASC);

-- role字段用于按角色类型筛选消息（用户消息 vs AI消息）
CREATE INDEX idx_messages_role ON messages(role);

-- created_at字段用于时间范围查询
CREATE INDEX idx_messages_created_at ON messages(created_at);

-- 外键索引（如果MySQL没有自动创建）
-- 这些索引有助于JOIN操作的性能
CREATE INDEX idx_conversations_user_id_fk ON conversations(user_id);
CREATE INDEX idx_conversations_character_id_fk ON conversations(character_id);
CREATE INDEX idx_messages_conversation_id_fk ON messages(conversation_id);

-- 全文搜索索引（可选，用于角色搜索功能）
-- 如果角色搜索功能使用频繁，可以考虑添加全文索引
-- ALTER TABLE characters ADD FULLTEXT(name, description, category);

-- 显示索引创建结果
SHOW INDEX FROM users;
SHOW INDEX FROM characters;
SHOW INDEX FROM conversations;
SHOW INDEX FROM messages;