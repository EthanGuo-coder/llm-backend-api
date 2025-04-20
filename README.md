# 🤖 llm-backend-api

#### 🚀 A Robust and Scalable Backend Solution for Large Language Models

#### *"llm-backend-api > [Gin Framework](https://github.com/gin-gonic/gin) + [Redis](https://redis.io)"*

🌐 [English](./README.md) · [简体中文](./README_zh-CN.md)

## Overview

`llm-backend-api` is a **robust** and **scalable** backend solution designed to facilitate seamless interactions with large language models (LLMs). Leveraging the power of **Golang** and **Redis**, this project offers a clear and user-friendly API for managing conversations, handling user authentication, and streaming chat responses from AI models. Whether you're building a chatbot, an interactive assistant, or any application requiring intelligent dialogue capabilities, `llm-backend-api` provides the necessary tools to manage and streamline these interactions efficiently.

## 🌟 Highlights

- **✨ Clear and Usable API**: Intuitive endpoints for managing conversations, users, and streaming chat messages.
- **⚡ Scalable Architecture**: Built with Golang and Redis to ensure high performance and scalability.
- **🔒 Secure Authentication**: Robust JWT-based authentication to protect user data and interactions.
- **📡 Streaming Responses**: Efficiently stream AI responses in real-time, enhancing user experience.
- **🛠️ Flexible Configuration**: Easily configurable through YAML files to suit various deployment environments.
- **💾 Persistent Storage**: Utilizes SQLite for reliable data persistence and Redis for fast access to session data.
- **🧠 RAG Service Integration**: Retrieval-Augmented Generation capabilities for knowledge-based conversations.

## 🛠️ Tech Stack

- **📝 Language**: Golang
- **🏗️ Framework**: Gin
- **🗄️ Database**: SQLite
- **⚙️ Cache**: Redis
- **🔑 Authentication**: JWT (JSON Web Tokens)
- **🔐 Password Security**: bcrypt

---

## Installation

### Prerequisites

- **Go**: Ensure you have Go installed. You can download it from [here](https://golang.org/dl/).
- **Redis**: Install and run Redis. Instructions can be found [here](https://redis.io/download).
- **SQLite**: SQLite is used for data persistence. Installation instructions are available [here](https://www.sqlite.org/download.html).

### Clone the Repository

```bash
git clone https://github.com/EthanGuo-coder/llm-backend-api.git
cd llm-backend-api
```

### Install Dependencies

```bash
go mod download
```

---

## Configuration

The application is configured using the `config.yaml` file located in the root directory. Below is an example configuration:

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
  conn_max_lifetime: 300 # in seconds

jwt:
  secret: "S3cureK3y#2024!AIsafety"
```

### Configuration Parameters

- **Server**
  - `port`: The port on which the server will run.

- **Redis**
  - `address`: Redis server address.
  - `password`: Redis server password (if any).
  - `db`: Redis database number.

- **SQLite**
  - `path`: Path to the SQLite database file.
  - `max_open_conns`: Maximum number of open connections to the database.
  - `max_idle_conns`: Maximum number of idle connections.
  - `conn_max_lifetime`: Maximum lifetime of a connection in seconds.

- **JWT**
  - `secret`: Secret key for signing JWT tokens.

---

## Running the Project

1. **Load Configuration**

   Ensure the `config.yaml` file is properly configured.

2. **Initialize Redis and SQLite**

   The application will automatically initialize Redis and SQLite based on the provided configuration.

3. **Start the Server**

   ```bash
   go run main.go
   ```

   The server will start on the port specified in `config.yaml` (default is `8080`).

   ```
   Connected to Redis successfully!
   SQLite initialized successfully!
   Server is running on port 8080
   ```

---

## API Documentation

### Authentication Endpoints

#### 1. **Register User**

- **Endpoint**: `POST /api/users/register`
- **Description**: Registers a new user with a username and password.

##### **Request**

- **Headers**
  - `Content-Type`: `application/json`

- **Body**

  ```json
  {
      "username": "john_doe",
      "password": "SecureP@ssw0rd!"
  }
  ```

##### **Response**

- **Status Codes**
  - `201 Created`: User registered successfully.
  - `400 Bad Request`: Invalid input or username already exists.

- **Body**

  ```json
  {
      "message": "User registered successfully"
  }
  ```

---

#### 2. **Login User**

- **Endpoint**: `POST /api/users/login`
- **Description**: Authenticates a user and returns a JWT token.

##### **Request**

- **Headers**
  - `Content-Type`: `application/json`

- **Body**

  ```json
  {
      "username": "john_doe",
      "password": "SecureP@ssw0rd!"
  }
  ```

##### **Response**

- **Status Codes**
  - `200 OK`: Authentication successful.
  - `401 Unauthorized`: Invalid username or password.

- **Body**

  ```json
  {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6..."
  }
  ```

---

### Conversation Endpoints

#### 1. **Create a Conversation**

- **Endpoint**: `POST /api/conversations/create`
- **Description**: Creates a new conversation with a given title and model.

##### **Request**

- **Headers**
  - `Content-Type`: `application/json`
  - `Authorization`: `Bearer <JWT Token>`

- **Body**

  ```json
  {
      "title": "My New Conversation",
      "model": "gpt-4o",
      "api_key": "your-api-key-here" // Required if different models need specific API keys
  }
  ```

##### **Response**

- **Status Codes**
  - `200 OK`: Conversation created successfully.
  - `400 Bad Request`: Invalid request body.
  - `401 Unauthorized`: Missing or invalid JWT token.

- **Body**

  ```json
  {
      "conversation_id": 329629,
      "title": "My New Conversation",
      "model": "gpt-4o",
      "api_key": "your-api-key-here",
      "created_time": 1731851729
  }
  ```

---

#### 2. **Get Conversation History**

- **Endpoint**: `GET /api/conversations/history/:conversation_id`
- **Description**: Retrieves the history of messages in the specified conversation.

##### **Request**

- **Headers**
  - `Content-Type`: `application/json`
  - `Authorization`: `Bearer <JWT Token>`

- **Path Parameters**
  - `conversation_id` (integer): The ID of the conversation.

##### **Response**

- **Status Codes**
  - `200 OK`: History retrieved successfully.
  - `404 Not Found`: Conversation ID does not exist.
  - `401 Unauthorized`: Missing or invalid JWT token.

- **Body**

  ```json
  {
      "conversation_id": 329629,
      "title": "My New Conversation",
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

#### 3. **List User Conversations**

- **Endpoint**: `GET /api/conversations/list`
- **Description**: Retrieves a list of all conversations for the authenticated user.

##### **Request**

- **Headers**
  - `Content-Type`: `application/json`
  - `Authorization`: `Bearer <JWT Token>`

##### **Response**

- **Status Codes**
  - `200 OK`: Conversations retrieved successfully.
  - `401 Unauthorized`: Missing or invalid JWT token.

- **Body**

  ```json
  [
      {
          "conversation_id": 329629,
          "title": "My New Conversation",
          "created_time": 1731851729
      },
      {
          "conversation_id": 329630,
          "title": "Another Conversation",
          "created_time": 1731851730
      }
  ]
  ```

---

#### 4. **Delete a Conversation**

- **Endpoint**: `POST /api/conversations/del/:conversation_id`
- **Description**: Deletes a specified conversation.

##### **Request**

- **Headers**
  - `Content-Type`: `application/json`
  - `Authorization`: `Bearer <JWT Token>`

- **Path Parameters**
  - `conversation_id` (integer): The ID of the conversation to delete.

##### **Response**

- **Status Codes**
  - `200 OK`: Conversation deleted successfully.
  - `404 Not Found`: Conversation ID does not exist.
  - `401 Unauthorized`: Missing or invalid JWT token.

- **Body**

  ```json
  {
      "message": "Conversation deleted successfully"
  }
  ```

---

### Chat Endpoints

#### 1. **Stream Chat Messages**

- **Endpoint**: `POST /api/chat/:conversation_id`
- **Description**: Sends a message to the specified conversation and streams the response from the AI model.

##### **Request**

- **Headers**
  - `Content-Type`: `application/json`
  - `Authorization`: `Bearer <JWT Token>`

- **Path Parameters**
  - `conversation_id` (integer): The ID of the conversation.

- **Body**

  ```json
  {
      "message": "介绍一下RUST"
  }
  ```

##### **Response**

- **Status Codes**
  - `200 OK`: Message processed and response streamed.
  - `400 Bad Request`: Invalid conversation ID or request body.
  - `401 Unauthorized`: Missing or invalid JWT token.
  - `404 Not Found`: Conversation ID does not exist.
  - `500 Internal Server Error`: Server encountered an error.

- **Streamed Response Format**

  ```json
  {"event":"message", "data":"R"}
  
  {"event":"message", "data":"ust"}
  
  {"event":"message", "data":" 是一种系统编程语言，由 Graydon Hoare 设计..."}
  
  {"event":"done", "data":"Stream finished"}
  
  {"event":"full_response", "data":"Complete AI response in a single message."}
  ```

  **Explanation of Events:**

  - `message`: Incremental response chunks from the AI model.
  - `done`: Indicates the end of the streamed response.
  - `full_response`: Contains the full concatenated response.

---

### RAG Service Endpoints

#### RAG Knowledge Base Management

##### 1. **Create Knowledge Base**

- **Endpoint**: `POST /api/rag/kb/create`
- **Description**: Creates a new knowledge base with specified embedding model.

##### **Request**

- **Headers**
  - `Content-Type`: `application/json`
  - `Authorization`: `Bearer <JWT Token>`

- **Body**

  ```json
  {
      "kb_name": "法律知识库",
      "embedding_model": "zhipu-embedding-3"  // Optional, uses default model if not specified
  }
  ```

##### **Response**

- **Status Codes**
  - `200 OK`: Knowledge base created successfully.
  - `400 Bad Request`: Invalid request parameters.
  - `401 Unauthorized`: Missing or invalid JWT token.
  - `500 Internal Server Error`: Server error.

- **Body**

  ```json
  {
      "success": true,
      "kb_id": "a1b2c3d4-5678-90ab-cdef-123456789abc",
      "message": "知识库创建成功，使用模型: zhipu-embedding-3"
  }
  ```

---

##### 2. **Get Knowledge Base List**

- **Endpoint**: `GET /api/rag/kb/list`
- **Description**: Retrieves all knowledge bases for the current user.

##### **Request**

- **Headers**
  - `Authorization`: `Bearer <JWT Token>`

##### **Response**

- **Status Codes**
  - `200 OK`: Knowledge bases retrieved successfully.
  - `401 Unauthorized`: Missing or invalid JWT token.
  - `500 Internal Server Error`: Server error.

- **Body**

  ```json
  {
      "success": true,
      "kbs": [
          {
              "kb_id": "a1b2c3d4-5678-90ab-cdef-123456789abc",
              "kb_name": "法律知识库",
              "embedding_model": "zhipu-embedding-3"
          },
          {
              "kb_id": "b2c3d4e5-6789-01bc-defg-2345678901de",
              "kb_name": "技术文档知识库",
              "embedding_model": "zhipu-embedding-2"
          }
      ],
      "message": "知识库列表获取成功"
  }
  ```

---

##### 3. **Delete Knowledge Base**

- **Endpoint**: `POST /api/rag/kb/delete`
- **Description**: Deletes a specified knowledge base.

##### **Request**

- **Headers**
  - `Content-Type`: `application/json`
  - `Authorization`: `Bearer <JWT Token>`

- **Body**

  ```json
  {
      "kb_id": "a1b2c3d4-5678-90ab-cdef-123456789abc"
  }
  ```

##### **Response**

- **Status Codes**
  - `200 OK`: Knowledge base deleted successfully.
  - `400 Bad Request`: Invalid request parameters.
  - `401 Unauthorized`: Missing or invalid JWT token.
  - `500 Internal Server Error`: Server error.

- **Body**

  ```json
  {
      "success": true,
      "message": "知识库删除成功"
  }
  ```

---

#### RAG Document Management

##### 1. **Upload Document**

- **Endpoint**: `POST /api/rag/doc/upload`
- **Description**: Uploads a document to the specified knowledge base.

##### **Request**

- **Headers**
  - `Content-Type`: `multipart/form-data`
  - `Authorization`: `Bearer <JWT Token>`

- **Form Parameters**
  - `kb_id`: Knowledge base ID
  - `file`: Document file

##### **Response**

- **Status Codes**
  - `200 OK`: Document uploaded successfully.
  - `400 Bad Request`: Invalid request parameters.
  - `401 Unauthorized`: Missing or invalid JWT token.
  - `500 Internal Server Error`: Server error.

- **Body**

  ```json
  {
      "success": true,
      "doc_id": "d1e2f3g4-5678-90ab-cdef-123456789abc",
      "message": "文档上传处理成功，共分割并添加 25 个文本块"
  }
  ```

---

##### 2. **Get Document List**

- **Endpoint**: `GET /api/rag/doc/list`
- **Description**: Retrieves all documents in the specified knowledge base.

##### **Request**

- **Headers**
  - `Authorization`: `Bearer <JWT Token>`

- **Query Parameters**
  - `kb_id`: Knowledge base ID

##### **Response**

- **Status Codes**
  - `200 OK`: Document list retrieved successfully.
  - `400 Bad Request`: Invalid request parameters.
  - `401 Unauthorized`: Missing or invalid JWT token.
  - `500 Internal Server Error`: Server error.

- **Body**

  ```json
  {
      "success": true,
      "docs": [
          {
              "doc_id": "d1e2f3g4-5678-90ab-cdef-123456789abc",
              "kb_id": "a1b2c3d4-5678-90ab-cdef-123456789abc",
              "doc_name": "合同协议.docx",
              "file_type": "docx"
          },
          {
              "doc_id": "e2f3g4h5-6789-01bc-defg-2345678901de",
              "kb_id": "a1b2c3d4-5678-90ab-cdef-123456789abc",
              "doc_name": "法律条款.txt",
              "file_type": "txt"
          }
      ],
      "message": "文档列表获取成功"
  }
  ```

---

##### 3. **Delete Document**

- **Endpoint**: `POST /api/rag/doc/delete`
- **Description**: Deletes a specified document from the knowledge base.

##### **Request**

- **Headers**
  - `Content-Type`: `application/json`
  - `Authorization`: `Bearer <JWT Token>`

- **Body**

  ```json
  {
      "kb_id": "a1b2c3d4-5678-90ab-cdef-123456789abc",
      "doc_id": "d1e2f3g4-5678-90ab-cdef-123456789abc"
  }
  ```

##### **Response**

- **Status Codes**
  - `200 OK`: Document deleted successfully.
  - `400 Bad Request`: Invalid request parameters.
  - `401 Unauthorized`: Missing or invalid JWT token.
  - `500 Internal Server Error`: Server error.

- **Body**

  ```json
  {
      "success": true,
      "message": "文档删除成功"
  }
  ```

---

#### RAG Retrieval and Chat

##### 1. **Knowledge Base Retrieval**

- **Endpoint**: `POST /api/rag/retrieve`
- **Description**: Retrieves information related to a query from the knowledge base.

##### **Request**

- **Headers**
  - `Content-Type`: `application/json`
  - `Authorization`: `Bearer <JWT Token>`

- **Body**

  ```json
  {
      "kb_id": "a1b2c3d4-5678-90ab-cdef-123456789abc",
      "query": "什么是不可抗力条款？",
      "top_k": 5  // Optional, default is 5
  }
  ```

##### **Response**

- **Status Codes**
  - `200 OK`: Retrieval successful.
  - `400 Bad Request`: Invalid request parameters.
  - `401 Unauthorized`: Missing or invalid JWT token.
  - `500 Internal Server Error`: Server error.

- **Body**

  ```json
  {
      "success": true,
      "results": [
          {
              "content": "不可抗力条款是指合同中约定的因不可预见、不可避免、不可克服的客观情况，导致合同无法履行或无法完全履行时，免除当事人部分或全部责任的条款...",
              "score": 0.85,
              "doc_id": "d1e2f3g4-5678-90ab-cdef-123456789abc",
              "doc_name": "合同协议.docx"
          },
          {
              "content": "在法律实践中，不可抗力通常包括自然灾害（如地震、洪水、台风等）和社会异常事件（如战争、罢工、政府行为等）...",
              "score": 0.72,
              "doc_id": "e2f3g4h5-6789-01bc-defg-2345678901de",
              "doc_name": "法律条款.txt"
          }
      ],
      "message": "检索成功"
  }
  ```

---

##### 2. **Knowledge Base Chat**

- **Endpoint**: `POST /api/rag/chat`
- **Description**: Conducts a chat based on knowledge base retrieval.

##### **Request**

- **Headers**
  - `Content-Type`: `application/json`
  - `Authorization`: `Bearer <JWT Token>`

- **Body**

  ```json
  {
      "conversation_id": 329629,
      "kb_id": "a1b2c3d4-5678-90ab-cdef-123456789abc",
      "message": "什么是不可抗力条款？",
      "top_k": 3  // Optional, default is 5
  }
  ```

##### **Response**

- Same as regular chat endpoint, using SSE streaming format:
  ```
  {"event":"message", "data":"不可抗力条款是指"}
  {"event":"message", "data":"合同中约定的因"}
  {"event":"message", "data":"不可预见、不可避免..."}
  {"event":"done", "data":"Stream finished"}
  {"event":"full_response", "data":"不可抗力条款是指合同中约定的因不可预见、不可避免、不可克服的客观情况，导致合同无法履行或无法完全履行时，免除当事人部分或全部责任的条款。在法律实践中，不可抗力通常包括自然灾害（如地震、洪水、台风等）和社会异常事件（如战争、罢工、政府行为等）。"}
  ```

---

#### RAG Metadata

##### 1. **Get Supported Embedding Models**

- **Endpoint**: `GET /api/rag/models`
- **Description**: Retrieves all embedding models supported by the system.

##### **Request**

- **Headers**
  - `Authorization`: `Bearer <JWT Token>`

##### **Response**

- **Status Codes**
  - `200 OK`: Models retrieved successfully.
  - `401 Unauthorized`: Missing or invalid JWT token.
  - `500 Internal Server Error`: Server error.

- **Body**

  ```json
  {
      "success": true,
      "models": [
          "zhipu-embedding-3",
          "zhipu-embedding-2"
      ],
      "message": "支持的Embedding模型列表获取成功"
  }
  ```

---

## Example `curl` Commands

### Basic API Commands

#### 1. **Register a User**

```bash
curl -X POST http://localhost:8080/api/users/register \
-H "Content-Type: application/json" \
-d '{
    "username": "john_doe",
    "password": "SecureP@ssw0rd!"
}'
```

#### 2. **Login a User**

```bash
curl -X POST http://localhost:8080/api/users/login \
-H "Content-Type: application/json" \
-d '{
    "username": "john_doe",
    "password": "SecureP@ssw0rd!"
}'
```

#### 3. **Create a Conversation**

```bash
curl -X POST http://localhost:8080/api/conversations/create \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_JWT_TOKEN" \
-d '{
    "title": "My New Conversation",
    "model": "gpt-4o",
    "api_key": "your-api-key-here"
}'
```

#### 4. **Stream Chat Messages**

```bash
curl -X POST http://localhost:8080/api/chat/329629 \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_JWT_TOKEN" \
-d '{
    "message": "介绍一下RUST"
}'
```

#### 5. **Get Conversation History**

```bash
curl -X GET http://localhost:8080/api/conversations/history/329629 \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### 6. **List User Conversations**

```bash
curl -X GET http://localhost:8080/api/conversations/list \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### 7. **Delete a Conversation**

```bash
curl -X POST http://localhost:8080/api/conversations/del/329629 \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### RAG Service Commands

#### 1. **Create Knowledge Base**

```bash
curl -X POST http://localhost:8080/api/rag/kb/create \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_JWT_TOKEN" \
-d '{
    "kb_name": "法律知识库",
    "embedding_model": "zhipu-embedding-3"
}'
```

#### 2. **Get Knowledge Base List**

```bash
curl -X GET http://localhost:8080/api/rag/kb/list \
-H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### 3. **Upload Document to Knowledge Base**

```bash
curl -X POST http://localhost:8080/api/rag/doc/upload \
-H "Authorization: Bearer YOUR_JWT_TOKEN" \
-F "kb_id=a1b2c3d4-5678-90ab-cdef-123456789abc" \
-F "file=@/path/to/your/document.docx"
```

#### 4. **Retrieve Information from Knowledge Base**

```bash
curl -X POST http://localhost:8080/api/rag/retrieve \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_JWT_TOKEN" \
-d '{
    "kb_id": "a1b2c3d4-5678-90ab-cdef-123456789abc",
    "query": "什么是不可抗力条款？",
    "top_k": 5
}'
```

#### 5. **Chat with Knowledge Base**

```bash
curl -X POST http://localhost:8080/api/rag/chat \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_JWT_TOKEN" \
-d '{
    "conversation_id": 329629,
    "kb_id": "a1b2c3d4-5678-90ab-cdef-123456789abc",
    "message": "什么是不可抗力条款？",
    "top_k": 3
}'
```

---

## Error Codes

| Status Code | Description                                         |
| ----------- | --------------------------------------------------- |
| 200         | Request succeeded.                                  |
| 201         | Resource created successfully.                     |
| 400         | Invalid request (e.g., missing/invalid parameters). |
| 401         | Unauthorized (invalid or missing JWT token).        |
| 404         | Resource not found (e.g., invalid conversation ID). |
| 500         | Internal server error.                              |

---

## Notes

1. **Authentication**: All endpoints, except for user registration and login, require a valid JWT token in the `Authorization` header.
2. **API Keys**: When creating a conversation, you can specify an `api_key` if different models require specific authentication.
3. **Streaming Responses**: The `Stream Chat Messages` endpoint streams responses incrementally. Ensure your client can handle SSE (Server-Sent Events) appropriately.
4. **Data Persistence**: Conversations are stored in both SQLite (for persistence) and Redis (for quick access). Deleting a conversation removes it from both storage systems.
5. **Security**: Passwords are securely hashed using bcrypt. Ensure your `jwt.secret` in the configuration is kept confidential.
6. **Customization**: Modify the `config.yaml` to suit your deployment environment, including changing ports, database paths, and Redis configurations.
7. **Extensibility**: The project is modular, allowing for easy extension of features such as adding new models, integrating additional services, or enhancing existing functionalities.
8. **RAG Service**: The Retrieval-Augmented Generation service enables knowledge-based conversations by retrieving relevant information from uploaded documents.

---

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request for any enhancements or bug fixes.

## License

This project is licensed under the [MIT License](LICENSE).

---

## Contact

For any inquiries or support, please contact [Ethan Guo](mailto:ethanguo2003@163.com).