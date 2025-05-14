package model

import (
	"time"
)

// RequestInfo 请求信息结构体
type RequestInfo struct {
	RequestID    string                 `json:"request_id"`
	Timestamp    string                 `json:"timestamp"`
	Method       string                 `json:"method"`
	URL          string                 `json:"url"`
	Path         string                 `json:"path,omitempty"`
	PathParts    []string               `json:"path_parts,omitempty"`
	PathParams   map[string]string      `json:"path_params"`
	QueryParams  map[string]string      `json:"query_params"`
	Headers      map[string]string      `json:"headers"`
	Client       ClientInfo             `json:"client"`
	Body         interface{}            `json:"body"`
	IsValidJSON  bool                   `json:"is_valid_json"`
	JSONError    string                 `json:"json_error,omitempty"`
	OriginalBody string                 `json:"original_body,omitempty"`
	ErrorDetails map[string]interface{} `json:"error_details,omitempty"`
}

// ClientInfo 客户端信息结构体
type ClientInfo struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

// Response API响应结构体
type Response struct {
	Message     string      `json:"message"`
	RequestInfo RequestInfo `json:"request_info"`
	LogID       string      `json:"log_id"`
}

// NewRequestInfo 创建请求信息对象
func NewRequestInfo(requestID string) RequestInfo {
	return RequestInfo{
		RequestID:  requestID,
		Timestamp:  time.Now().Format(time.RFC3339),
		PathParams: make(map[string]string),
		Headers:    make(map[string]string),
		Client:     ClientInfo{},
	}
}

// NewResponse 创建API响应对象
func NewResponse(requestInfo RequestInfo) Response {
	return Response{
		Message:     "success",
		RequestInfo: requestInfo,
		LogID:       requestInfo.RequestID,
	}
}

// JSONErrorDetails JSON解析错误的详细信息
type JSONErrorDetails struct {
	ErrorType   string `json:"error_type"`
	Message     string `json:"message"`
	Line        int    `json:"line,omitempty"`
	Column      int    `json:"column,omitempty"`
	Position    int    `json:"position,omitempty"`
	ErrorChar   string `json:"error_char,omitempty"`
	LineContent string `json:"line_content,omitempty"`
	Pointer     string `json:"pointer,omitempty"`
} 