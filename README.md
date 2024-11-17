# llm-backend-api

## Overview
This documentation provides detailed information about the three APIs for managing conversations and streaming chat responses. The APIs are implemented using the Gin framework and enable users to create conversations, retrieve conversation histories, and send/receive messages in a streaming fashion.

- **Tech Stack**: Golang + Redis

---

## 1. **Create a Conversation**

### **Endpoint**
`POST /api/conversations/create/`

### **Description**
Creates a new conversation with a given title and model.

### **Request**
#### **Headers**
- `Content-Type`: `application/json`

#### **Body**
```json
{
    "title": "My New Conversation", // Title of the conversation
    "model": "glm-4-air"           // Model used for the conversation
}
```

### **Response**
#### **Status Code**
- `200 OK`: The conversation was successfully created.
- `400 Bad Request`: The request body is invalid.

#### **Body**
```json
{
    "conversation_id": "329629",   // Unique ID of the conversation
    "title": "My New Conversation", // Title of the conversation
    "model": "gpt-4o",           // Model used
    "messages": [
        {
            "role": "system",       // Initial message role
            "content": "你是一个乐于回答各种问题的小助手" // System-provided context
        }
    ],
    "created_time": 1731851729      // Timestamp of creation
}
```

---

## 2. **Stream Chat Messages**

### **Endpoint**
`POST /api/conversations/:conversation_id/chat/`

### **Description**
Sends a message to the specified conversation and streams the response from the AI model.

### **Request**
#### **Headers**
- `Content-Type`: `application/json`

#### **Path Parameters**
- `conversation_id` (string): The ID of the conversation.

#### **Body**
```json
{
    "api_key": "8f61083c72a10210e865c3b5ce35ecef.9ZKwsInQfHCjVlwQ", // API key for authentication
    "message": "介绍一下RUST"                                        // Message to send
}
```

### **Response**
#### **Status Code**
- `200 OK`: The message was successfully processed, and the response is streamed.
- `404 Not Found`: The specified conversation ID does not exist.
- `401 Unauthorized`: Invalid or missing API key.

#### **Streamed Response Format**
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

## 3. **Get Conversation History**

### **Endpoint**
`GET /api/conversations/history/:conversation_id/`

### **Description**
Retrieves the history of messages in the specified conversation.

### **Request**
#### **Headers**
- `Content-Type`: `application/json`

#### **Path Parameters**
- `conversation_id` (string): The ID of the conversation.

### **Response**
#### **Status Code**
- `200 OK`: The history was successfully retrieved.
- `404 Not Found`: The specified conversation ID does not exist.

#### **Body**
```json
{
    "conversation_id": "329629",   // Unique ID of the conversation
    "title": "My New Conversation", // Title of the conversation
    "model": "glm-4-air",           // Model used
    "messages": [
        {
            "role": "system",       // Role of the message sender
            "content": "你是一个乐于回答各种问题的小助手" // System-provided context
        },
        {
            "role": "user",         // User's input
            "content": "介绍一下RUST"
        },
        {
            "role": "assistant",    // AI's response
            "content": "Rust 是一种系统编程语言，由 Graydon Hoare 设计..."
        }
    ],
    "created_time": 1731851729      // Timestamp of conversation creation
}
```

---

## Example `curl` Commands

### 1. **Create a Conversation**
```bash
curl -X POST http://localhost:8080/api/conversations/create/ \
-H "Content-Type: application/json" \
-d '{
    "title": "My New Conversation",
    "model": "glm-4-air"
}'
```

### 2. **Stream Chat Messages**
```bash
curl -X POST http://localhost:8080/api/conversations/329629/chat/ \
-H "Content-Type: application/json" \
-d '{
    "api_key": "xxxxxx",  // your openai apikey
    "message": "介绍一下RUST"
}'
```

### 3. **Get Conversation History**
```bash
curl -X GET http://localhost:8080/api/conversations/history/329629/ \
-H "Content-Type: application/json"
```

---

## Error Codes

| Status Code | Description                                         |
| ----------- | --------------------------------------------------- |
| 200         | Request succeeded.                                  |
| 400         | Invalid request (e.g., missing/invalid parameters). |
| 401         | Unauthorized (invalid or missing API key).          |
| 404         | Resource not found (e.g., invalid conversation ID). |

---

## Notes
1. Ensure the API key provided is valid and has sufficient permissions.
2. The server must be running on `http://localhost:8080` or replace with the actual hostname and port if deployed elsewhere.
3. The `streamSendMessage` endpoint sends responses incrementally; clients must handle streaming appropriately.