# llm-backend-api

## Overview

`llm-backend-api` is a robust and scalable backend solution designed to facilitate seamless interactions with large language models (LLMs). Leveraging the power of Golang and Redis, this project offers a clear and user-friendly API for managing conversations, handling user authentication, and streaming chat responses from AI models. Whether you're building a chatbot, an interactive assistant, or any application requiring intelligent dialogue capabilities, `llm-backend-api` provides the necessary tools to manage and streamline these interactions efficiently.

## Highlights

- **Clear and Usable API**: Intuitive endpoints for managing conversations, users, and streaming chat messages.
- **Scalable Architecture**: Built with Golang and Redis to ensure high performance and scalability.
- **Secure Authentication**: Robust JWT-based authentication to protect user data and interactions.
- **Streaming Responses**: Efficiently stream AI responses in real-time, enhancing user experience.
- **Flexible Configuration**: Easily configurable through YAML files to suit various deployment environments.
- **Persistent Storage**: Utilizes SQLite for reliable data persistence and Redis for fast access to session data.

## Tech Stack

- **Language**: Golang
- **Framework**: Gin
- **Database**: SQLite
- **Cache**: Redis
- **Authentication**: JWT (JSON Web Tokens)
- **Password Security**: bcrypt

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

## Example `curl` Commands

### 1. **Register a User**

```bash
curl -X POST http://localhost:8080/api/users/register \
-H "Content-Type: application/json" \
-d '{
    "username": "john_doe",
    "password": "SecureP@ssw0rd!"
}'
```

### 2. **Login a User**

```bash
curl -X POST http://localhost:8080/api/users/login \
-H "Content-Type: application/json" \
-d '{
    "username": "john_doe",
    "password": "SecureP@ssw0rd!"
}'
```

### 3. **Create a Conversation**

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

### 4. **Stream Chat Messages**

```bash
curl -X POST http://localhost:8080/api/chat/329629 \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_JWT_TOKEN" \
-d '{
    "message": "介绍一下RUST"
}'
```

### 5. **Get Conversation History**

```bash
curl -X GET http://localhost:8080/api/conversations/history/329629 \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 6. **List User Conversations**

```bash
curl -X GET http://localhost:8080/api/conversations/list \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 7. **Delete a Conversation**

```bash
curl -X POST http://localhost:8080/api/conversations/del/329629 \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_JWT_TOKEN"
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

---

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request for any enhancements or bug fixes.

## License

This project is licensed under the [MIT License](LICENSE).

---

## Contact

For any inquiries or support, please contact [Ethan Guo](mailto:ethan@example.com).