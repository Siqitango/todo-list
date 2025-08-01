-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS todo_list;

-- 使用创建的数据库
USE todo_list;

-- 创建 todos 表
CREATE TABLE IF NOT EXISTS todos (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    priority INT NOT NULL DEFAULT 2,
    status INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- 插入测试数据（可选）
INSERT INTO todos (title, description, priority, status, created_at, updated_at)
VALUES 
('完成项目配置', '配置项目所需的所有依赖和环境', 3, 1, NOW(), NOW()),
('实现用户认证', '添加用户登录和注册功能', 2, 1, NOW(), NOW()),
('开发API文档', '为所有接口编写详细文档', 1, 1, NOW(), NOW());