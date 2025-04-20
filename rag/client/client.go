package client

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/EthanGuo-coder/llm-backend-api/config"
	pb "github.com/EthanGuo-coder/llm-backend-api/rag/protos"
)

// RagClient 封装了对 RAG 服务的调用
type RagClient struct {
	client pb.KnowledgeBaseServiceClient
	conn   *grpc.ClientConn
}

// NewRagClient 创建新的 RAG 客户端
func NewRagClient() (*RagClient, error) {
	// 从配置中获取 RAG 服务地址
	ragServiceAddr := config.AppConfig.RAG.ServiceAddr

	// 创建 gRPC 连接
	conn, err := grpc.Dial(
		ragServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(100*1024*1024)), // 100MB
	)
	if err != nil {
		return nil, err
	}

	// 创建 gRPC 客户端
	client := pb.NewKnowledgeBaseServiceClient(conn)

	return &RagClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close 关闭 gRPC 连接
func (c *RagClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// CreateKnowledgeBase 创建知识库
func (c *RagClient) CreateKnowledgeBase(ctx context.Context, uid, kbName, embeddingModel string) (*pb.CreateKBResponse, error) {
	req := &pb.CreateKBRequest{
		Uid:            uid,
		KbName:         kbName,
		EmbeddingModel: embeddingModel,
	}

	return c.client.CreateKnowledgeBase(ctx, req)
}

// ListKnowledgeBases 获取用户的知识库列表
func (c *RagClient) ListKnowledgeBases(ctx context.Context, uid string) (*pb.ListKBsResponse, error) {
	req := &pb.ListKBsRequest{
		Uid: uid,
	}

	return c.client.ListKnowledgeBases(ctx, req)
}

// UploadDocument 上传文档
func (c *RagClient) UploadDocument(ctx context.Context, uid, kbID, docName string, fileContent []byte, fileType string) (*pb.UploadDocResponse, error) {
	req := &pb.UploadDocRequest{
		Uid:         uid,
		KbId:        kbID,
		DocName:     docName,
		FileContent: fileContent,
		FileType:    fileType,
	}

	return c.client.UploadDocument(ctx, req)
}

// RetrieveInfo 从知识库检索信息
func (c *RagClient) RetrieveInfo(ctx context.Context, kbID, query string, topK int) (*pb.RetrieveResponse, error) {
	req := &pb.RetrieveRequest{
		KbId:  kbID,
		Query: query,
		TopK:  int32(topK),
	}

	return c.client.RetrieveInfo(ctx, req)
}

// ListEmbeddingModels 获取支持的嵌入模型列表
func (c *RagClient) ListEmbeddingModels(ctx context.Context) (*pb.ListEmbeddingModelsResponse, error) {
	req := &pb.ListEmbeddingModelsRequest{}

	return c.client.ListEmbeddingModels(ctx, req)
}

// ListDocuments 获取知识库中的文档列表
func (c *RagClient) ListDocuments(ctx context.Context, kbID string) (*pb.ListDocsResponse, error) {
	req := &pb.ListDocsRequest{
		KbId: kbID,
	}

	return c.client.ListDocuments(ctx, req)
}

// DeleteKnowledgeBase 删除知识库
func (c *RagClient) DeleteKnowledgeBase(ctx context.Context, uid, kbID string) (*pb.DeleteKBResponse, error) {
	req := &pb.DeleteKBRequest{
		Uid:  uid,
		KbId: kbID,
	}

	return c.client.DeleteKnowledgeBase(ctx, req)
}

// DeleteDocument 删除文档
func (c *RagClient) DeleteDocument(ctx context.Context, uid, kbID, docID string) (*pb.DeleteDocResponse, error) {
	req := &pb.DeleteDocRequest{
		Uid:   uid,
		KbId:  kbID,
		DocId: docID,
	}

	return c.client.DeleteDocument(ctx, req)
}
