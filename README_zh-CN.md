# 🤖 llm-backend-api

#### 🚀 面向大型语言模型的强大且可扩展的后端解决方案

#### *"llm-backend-api > [Gin 框架](https://github.com/gin-gonic/gin) + [Redis](https://redis.io)"*

🌐 [English](./README.md) · [简体中文](./README_zh-CN.md)

## 概述

`llm-backend-api` 是一个**强大**且**可扩展**的后端解决方案，旨在促进与大型语言模型（LLMs）的无缝交互。该项目利用 **Golang** 和 **Redis** 的强大功能，提供了一个清晰且用户友好的 API，用于管理对话、处理用户身份验证以及流式传输来自 AI 模型的聊天响应。无论您是在构建聊天机器人、交互式助手，还是任何需要智能对话功能的应用，`llm-backend-api` 都提供了必要的工具，以高效地管理和简化这些交互。

## 🌟 亮点

- **✨ 清晰且易用的 API**：直观的端点，用于管理对话、用户和流式聊天消息。
- **⚡ 可扩展的架构**：采用 Golang 和 Redis 构建，确保高性能和可扩展性。
- **🔒 安全的身份验证**：强大的基于 JWT 的身份验证，保护用户数据和交互。
- **📡 流式响应**：高效地实时流式传输 AI 响应，提升用户体验。
- **🛠️ 灵活的配置**：通过 YAML 文件轻松配置，适应各种部署环境。
- **💾 持久存储**：利用 SQLite 进行可靠的数据持久化，Redis 用于快速访问会话数据。

## 🛠️ 技术栈

- **📝 语言**：Golang
- **🏗️ 框架**：Gin
- **🗄️ 数据库**：SQLite
- **⚙️ 缓存**：Redis
- **🔑 身份验证**：JWT（JSON Web Tokens）
- **🔐 密码安全**：bcrypt

---

## 安装

### 前置条件

- **Go**：确保已安装 Go。您可以从 [这里](https://golang.org/dl/) 下载。
- **Redis**：安装并运行 Redis。安装说明请参见 [这里](https://redis.io/download)。
- **SQLite**：SQLite 用于数据持久化。安装说明可在 [这里](https://www.sqlite.org/download.html) 找到。

### 克隆仓库

```bash
git clone https://github.com/EthanGuo-coder/llm-backend-api.git
cd llm-backend-api
```

### 安装依赖

```bash
go mod download
```

---

## 配置

应用程序使用位于根目录的 `config.yaml` 文件进行配置。以下是一个示例配置：

```yaml
server:
  port: "8080"

redis:
  address: "localhost:6379"
  password: ""
  db: 0

sqlite:
  path: "./llm_backend.db"
  max_open_conns: 10
  max_idle_conns: 5
  conn_max_lifetime: 300 # 以秒为单位

jwt:
  secret: "S3cureK3y#2024!AIsafety"
```

### 配置参数

- **服务器**
    - `port`：服务器运行的端口。

- **Redis**
    - `address`：Redis 服务器地址。
    - `password`：Redis 服务器密码（如果有）。
    - `db`：Redis 数据库编号。

- **SQLite**
    - `path`：SQLite 数据库文件的路径。
    - `max_open_conns`：数据库的最大打开连接数。
    - `max_idle_conns`：最大空闲连接数。
    - `conn_max_lifetime`：连接的最大生命周期（以秒为单位）。

- **JWT**
    - `secret`：用于签署 JWT 令牌的密钥。

---

## 运行项目

1. **加载配置**

   确保正确配置了 `config.yaml` 文件。

2. **初始化 Redis 和 SQLite**

   应用程序将根据提供的配置自动初始化 Redis 和 SQLite。

3. **启动服务器**

   ```bash
   go run main.go
   ```

   服务器将在 `config.yaml` 中指定的端口上启动（默认端口为 `8080`）。

   ```
   Connected to Redis successfully!
   SQLite initialized successfully!
   Server is running on port 8080
   ```

---

## API 文档

### 身份验证端点

#### 1. **注册用户**

- **端点**：`POST /api/users/register`
- **描述**：使用用户名和密码注册新用户。

##### **请求**

- **头部**
    - `Content-Type`：`application/json`

- **主体**

  ```json
  {
      "username": "john_doe",
      "password": "SecureP@ssw0rd!"
  }
  ```

##### **响应**

- **状态码**
    - `201 Created`：用户注册成功。
    - `400 Bad Request`：输入无效或用户名已存在。

- **主体**

  ```json
  {
      "message": "用户注册成功"
  }
  ```

---

#### 2. **用户登录**

- **端点**：`POST /api/users/login`
- **描述**：验证用户身份并返回 JWT 令牌。

##### **请求**

- **头部**
    - `Content-Type`：`application/json`

- **主体**

  ```json
  {
      "username": "john_doe",
      "password": "SecureP@ssw0rd!"
  }
  ```

##### **响应**

- **状态码**
    - `200 OK`：身份验证成功。
    - `401 Unauthorized`：用户名或密码无效。

- **主体**

  ```json
  {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6..."
  }
  ```

---

### 对话端点

#### 1. **创建对话**

- **端点**：`POST /api/conversations/create`
- **描述**：使用给定的标题和模型创建新的对话。

##### **请求**

- **头部**
    - `Content-Type`：`application/json`
    - `Authorization`：`Bearer <JWT 令牌>`

- **主体**

  ```json
  {
      "title": "我的新对话",
      "model": "gpt-4o",
      "api_key": "your-api-key-here" // 如果不同模型需要特定 API 密钥，则需要
  }
  ```

##### **响应**

- **状态码**
    - `200 OK`：对话创建成功。
    - `400 Bad Request`：请求主体无效。
    - `401 Unauthorized`：缺少或无效的 JWT 令牌。

- **主体**

  ```json
  {
      "conversation_id": 329629,
      "title": "我的新对话",
      "model": "gpt-4o",
      "api_key": "your-api-key-here",
      "created_time": 1731851729
  }
  ```

---

#### 2. **获取对话历史**

- **端点**：`GET /api/conversations/history/:conversation_id`
- **描述**：检索指定对话中的消息历史记录。

##### **请求**

- **头部**
    - `Content-Type`：`application/json`
    - `Authorization`：`Bearer <JWT 令牌>`

- **路径参数**
    - `conversation_id`（整数）：对话的 ID。

##### **响应**

- **状态码**
    - `200 OK`：成功检索历史记录。
    - `404 Not Found`：对话 ID 不存在。
    - `401 Unauthorized`：缺少或无效的 JWT 令牌。

- **主体**

  ```json
  {
      "conversation_id": 329629,
      "title": "我的新对话",
      "model": "gpt-4o",
      "messages": [
          {
              "role": "user",
              "content": "介绍一下RUST",
              "message_id": 1
          },
          {
              "role": "assistant",
              "content": "Rust 是一种系统编程语言，由 Graydon Hoare 设计...",
              "message_id": 2
          }
      ],
      "created_time": 1731851729
  }
  ```

---

#### 3. **列出用户对话**

- **端点**：`GET /api/conversations/list`
- **描述**：检索认证用户的所有对话列表。

##### **请求**

- **头部**
    - `Content-Type`：`application/json`
    - `Authorization`：`Bearer <JWT 令牌>`

##### **响应**

- **状态码**
    - `200 OK`：成功检索对话。
    - `401 Unauthorized`：缺少或无效的 JWT 令牌。

- **主体**

  ```json
  [
      {
          "conversation_id": 329629,
          "title": "我的新对话",
          "created_time": 1731851729
      },
      {
          "conversation_id": 329630,
          "title": "另一个对话",
          "created_time": 1731851730
      }
  ]
  ```

---

#### 4. **删除对话**

- **端点**：`POST /api/conversations/del/:conversation_id`
- **描述**：删除指定的对话。

##### **请求**

- **头部**
    - `Content-Type`：`application/json`
    - `Authorization`：`Bearer <JWT 令牌>`

- **路径参数**
    - `conversation_id`（整数）：要删除的对话 ID。

##### **响应**

- **状态码**
    - `200 OK`：对话删除成功。
    - `404 Not Found`：对话 ID 不存在。
    - `401 Unauthorized`：缺少或无效的 JWT 令牌。

- **主体**

  ```json
  {
      "message": "对话删除成功"
  }
  ```

---

### 聊天端点

#### 1. **流式聊天消息**

- **端点**：`POST /api/chat/:conversation_id`
- **描述**：向指定的对话发送消息，并流式传输来自 AI 模型的响应。

##### **请求**

- **头部**
    - `Content-Type`：`application/json`
    - `Authorization`：`Bearer <JWT 令牌>`

- **路径参数**
    - `conversation_id`（整数）：对话的 ID。

- **主体**

  ```json
  {
      "message": "介绍一下RUST"
  }
  ```

##### **响应**

- **状态码**
    - `200 OK`：消息已处理并开始流式传输响应。
    - `400 Bad Request`：对话 ID 或请求主体无效。
    - `401 Unauthorized`：缺少或无效的 JWT 令牌。
    - `404 Not Found`：对话 ID 不存在。
    - `500 Internal Server Error`：服务器遇到错误。

- **流式响应格式**

  ```json
  {"event":"message", "data":"R"}
  
  {"event":"message", "data":"ust"}
  
  {"event":"message", "data":" 是一种系统编程语言，由 Graydon Hoare 设计..."}
  
  {"event":"done", "data":"Stream finished"}
  
  {"event":"full_response", "data":"Complete AI response in a single message."}
  ```

  **事件说明：**

    - `message`：来自 AI 模型的增量响应块。
    - `done`：表示流式响应结束。
    - `full_response`：包含完整的拼接响应。

---

## 示例 `curl` 命令

### 1. **注册用户**

```bash
curl -X POST http://localhost:8080/api/users/register \
-H "Content-Type: application/json" \
-d '{
    "username": "john_doe",
    "password": "SecureP@ssw0rd!"
}'
```

### 2. **用户登录**

```bash
curl -X POST http://localhost:8080/api/users/login \
-H "Content-Type: application/json" \
-d '{
    "username": "john_doe",
    "password": "SecureP@ssw0rd!"
}'
```

### 3. **创建对话**

```bash
curl -X POST http://localhost:8080/api/conversations/create \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_JWT_TOKEN" \
-d '{
    "title": "我的新对话",
    "model": "gpt-4o",
    "api_key": "your-api-key-here"
}'
```

### 4. **流式聊天消息**

```bash
curl -X POST http://localhost:8080/api/chat/329629 \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_JWT_TOKEN" \
-d '{
    "message": "介绍一下RUST"
}'
```

### 5. **获取对话历史**

```bash
curl -X GET http://localhost:8080/api/conversations/history/329629 \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 6. **列出用户对话**

```bash
curl -X GET http://localhost:8080/api/conversations/list \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 7. **删除对话**

```bash
curl -X POST http://localhost:8080/api/conversations/del/329629 \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## 错误代码

| 状态码    | 描述                                               |
| --------- | -------------------------------------------------- |
| 200       | 请求成功。                                         |
| 201       | 资源创建成功。                                     |
| 400       | 请求无效（例如，缺少/无效的参数）。               |
| 401       | 未授权（无效或缺少 JWT 令牌）。                   |
| 404       | 资源未找到（例如，Invalid conversation ID）。     |
| 500       | 服务器内部错误。                                   |

---

## 注意事项

1. **身份验证**：所有端点，除了用户注册和登录，均需要在 `Authorization` 头部提供有效的 JWT 令牌。
2. **API 密钥**：创建对话时，可以指定 `api_key`，如果不同模型需要特定的认证。
3. **流式响应**：`流式聊天消息` 端点会增量流式传输响应。确保您的客户端能够适当处理 SSE（服务器发送事件）。
4. **数据持久化**：对话同时存储在 SQLite（用于持久化）和 Redis（用于快速访问）中。删除对话会同时从这两种存储系统中移除。
5. **安全性**：密码使用 bcrypt 安全哈希。确保配置中的 `jwt.secret` 保密。
6. **自定义**：修改 `config.yaml` 以适应您的部署环境，包括更改端口、数据库路径和 Redis 配置。
7. **可扩展性**：该项目是模块化的，允许轻松扩展功能，如添加新模型、集成额外服务或增强现有功能。

---

## 贡献

欢迎贡献！请 Fork 仓库并提交 Pull Request 以进行任何增强或修复。

## 许可证

本项目基于 [MIT 许可证](LICENSE) 许可。

---

## 联系方式

如有任何询问或需要支持，请联系 [Ethan Guo](mailto:ethanguo2003@163.com)。