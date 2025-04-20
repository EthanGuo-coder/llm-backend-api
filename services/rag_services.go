package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/EthanGuo-coder/llm-backend-api/models"
	"github.com/EthanGuo-coder/llm-backend-api/rag/client"
)

// 默认配置
const (
	DefaultTopK = 5
)

// RAGService 提供 RAG 相关功能
type RAGService struct {
	client *client.RagClient
}

var ragService *RAGService

// GetRAGService 获取 RAG 服务单例
func GetRAGService() (*RAGService, error) {
	if ragService == nil {
		ragClient, err := client.NewRagClient()
		if err != nil {
			return nil, fmt.Errorf("无法创建 RAG 客户端: %w", err)
		}

		ragService = &RAGService{
			client: ragClient,
		}
	}

	return ragService, nil
}

// Close 关闭 RAG 服务
func (s *RAGService) Close() {
	if s.client != nil {
		s.client.Close()
	}
}

// CreateKnowledgeBase 创建知识库
func (s *RAGService) CreateKnowledgeBase(uid, kbName, embeddingModel string) (*models.CreateKnowledgeBaseResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := s.client.CreateKnowledgeBase(ctx, uid, kbName, embeddingModel)
	if err != nil {
		return nil, fmt.Errorf("创建知识库失败: %w", err)
	}

	return &models.CreateKnowledgeBaseResponse{
		Success: resp.Success,
		KBID:    resp.KbId,
		Message: resp.Message,
	}, nil
}

// ListKnowledgeBases 获取用户的知识库列表
func (s *RAGService) ListKnowledgeBases(uid string) (*models.ListKnowledgeBasesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := s.client.ListKnowledgeBases(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("获取知识库列表失败: %w", err)
	}

	kbs := make([]models.KnowledgeBase, 0, len(resp.Kbs))
	for _, kb := range resp.Kbs {
		kbs = append(kbs, models.KnowledgeBase{
			ID:             kb.KbId,
			Name:           kb.KbName,
			EmbeddingModel: kb.EmbeddingModel,
		})
	}

	return &models.ListKnowledgeBasesResponse{
		Success: resp.Success,
		KBs:     kbs,
		Message: resp.Message,
	}, nil
}

// UploadDocument 上传文档
func (s *RAGService) UploadDocument(uid, kbID, docName string, fileContent []byte, fileType string) (*models.UploadDocumentResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute) // 上传大文件可能需要更长时间
	defer cancel()

	resp, err := s.client.UploadDocument(ctx, uid, kbID, docName, fileContent, fileType)
	if err != nil {
		return nil, fmt.Errorf("上传文档失败: %w", err)
	}

	return &models.UploadDocumentResponse{
		Success: resp.Success,
		DocID:   resp.DocId,
		Message: resp.Message,
	}, nil
}

// RetrieveInfo 从知识库检索信息
func (s *RAGService) RetrieveInfo(kbID, query string, topK int) (*models.RetrieveResponse, error) {
	if topK <= 0 {
		topK = DefaultTopK
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := s.client.RetrieveInfo(ctx, kbID, query, topK)
	if err != nil {
		return nil, fmt.Errorf("检索信息失败: %w", err)
	}

	results := make([]models.RetrieveResult, 0, len(resp.Results))
	for _, result := range resp.Results {
		results = append(results, models.RetrieveResult{
			Content: result.Content,
			Score:   result.Score,
			DocID:   result.DocId,
			DocName: result.DocName,
		})
	}

	return &models.RetrieveResponse{
		Success: resp.Success,
		Results: results,
		Message: resp.Message,
	}, nil
}

// ListEmbeddingModels 获取支持的嵌入模型列表
func (s *RAGService) ListEmbeddingModels() (*models.ListEmbeddingModelsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := s.client.ListEmbeddingModels(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取嵌入模型列表失败: %w", err)
	}

	return &models.ListEmbeddingModelsResponse{
		Success: resp.Success,
		Models:  resp.Models,
		Message: resp.Message,
	}, nil
}

// ListDocuments 获取知识库中的文档列表
func (s *RAGService) ListDocuments(kbID string) (*models.ListDocumentsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := s.client.ListDocuments(ctx, kbID)
	if err != nil {
		return nil, fmt.Errorf("获取文档列表失败: %w", err)
	}

	docs := make([]models.Document, 0, len(resp.Docs))
	for _, doc := range resp.Docs {
		docs = append(docs, models.Document{
			ID:       doc.DocId,
			KBID:     doc.KbId,
			Name:     doc.DocName,
			FileType: doc.FileType,
		})
	}

	return &models.ListDocumentsResponse{
		Success: resp.Success,
		Docs:    docs,
		Message: resp.Message,
	}, nil
}

// DeleteKnowledgeBase 删除知识库
func (s *RAGService) DeleteKnowledgeBase(uid, kbID string) (*models.DeleteKnowledgeBaseResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := s.client.DeleteKnowledgeBase(ctx, uid, kbID)
	if err != nil {
		return nil, fmt.Errorf("删除知识库失败: %w", err)
	}

	return &models.DeleteKnowledgeBaseResponse{
		Success: resp.Success,
		Message: resp.Message,
	}, nil
}

// DeleteDocument 删除文档
func (s *RAGService) DeleteDocument(uid, kbID, docID string) (*models.DeleteDocumentResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := s.client.DeleteDocument(ctx, uid, kbID, docID)
	if err != nil {
		return nil, fmt.Errorf("删除文档失败: %w", err)
	}

	return &models.DeleteDocumentResponse{
		Success: resp.Success,
		Message: resp.Message,
	}, nil
}

// GenerateRagPrompt 生成基于知识库的对话提示
func (s *RAGService) GenerateRagPrompt(kbID, query string, topK int) (string, error) {
	// 从知识库中检索相关信息
	retrieveResp, err := s.RetrieveInfo(kbID, query, topK)
	if err != nil {
		return "", fmt.Errorf("从知识库检索信息失败: %w", err)
	}

	if !retrieveResp.Success || len(retrieveResp.Results) == 0 {
		return "", errors.New("知识库中未找到相关信息")
	}

	// 构建包含检索结果的提示
	var promptBuilder strings.Builder
	promptBuilder.WriteString("以下是与问题相关的背景信息：\n\n")

	for i, result := range retrieveResp.Results {
		promptBuilder.WriteString(fmt.Sprintf("[文档%d: %s]\n%s\n\n", i+1, result.DocName, result.Content))
	}

	promptBuilder.WriteString("请基于上述信息回答以下问题：\n\n")
	promptBuilder.WriteString(query)

	return promptBuilder.String(), nil
}
