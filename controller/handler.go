package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/renoelis/returnall-go/model"
	"github.com/renoelis/returnall-go/utils"
)

// 全局日志记录器
var logger *utils.Logger

// InitLogger 初始化日志记录器
func InitLogger(logPath string) error {
	var err error
	// 创建日志记录器，只输出到文件，不输出到控制台
	logger, err = utils.NewLogger("returnAll-api", logPath, false)
	return err
}

// parseJSONBody 解析请求体，尝试转换为JSON
func parseJSONBody(bodyBytes []byte) (interface{}, bool, string, string, map[string]interface{}) {
	if len(bodyBytes) == 0 {
		return nil, true, "", "", nil
	}

	// 保留原始字符串，包括换行符
	originalStr := string(bodyBytes)

	var jsonData interface{}
	err := json.Unmarshal(bodyBytes, &jsonData)
	if err == nil {
		return jsonData, true, "", originalStr, nil
	}

	// 如果请求体不是有效的JSON，则以紧凑字符串形式返回
	// 删除所有换行符和多余空格，保持紧凑格式
	re := regexp.MustCompile(`\s+`)
	compactStr := strings.TrimSpace(re.ReplaceAllString(originalStr, " "))

	// 提取错误详细信息
	errorDetails := make(map[string]interface{})
	errorDetails["error_type"] = "JSONDecodeError"

	// 尝试提取错误位置信息
	var line, column, position int
	var errorChar, lineContent, pointer, errorMessageBase string

	// 分析错误消息和位置
	syntaxError, ok := err.(*json.SyntaxError)
	if ok {
		position = int(syntaxError.Offset)

		// 计算行号和列号
		line, column = findLineAndColumn(originalStr, position)

		// 获取错误字符
		if position < len(originalStr) {
			errorChar = string(originalStr[position])
		}

		// 获取错误行内容和优化错误位置
		lines := strings.Split(originalStr, "\n")

		// 额外的错误位置优化：检查引号转义问题
		// 这种问题经常在位置定位上有偏差
		fixedPosition := false

		// 先检查是否有反斜杠后跟引号的问题 (如 "application/json\",)
		for i := 0; i < len(lines); i++ {
			// 针对常见的Content-Type等头部检查
			if strings.Contains(lines[i], "Content-Type") && strings.Contains(lines[i], "\\\"") {
				lineContent = lines[i]
				line = i + 1
				column = strings.LastIndex(lines[i], "\\\"") + 1
				fixedPosition = true
				break
			}

			// 检查任何键值对中存在的反斜杠引号问题
			backslashQuoteIndex := strings.Index(lines[i], "\\\"")
			if backslashQuoteIndex >= 0 {
				lineContent = lines[i]
				line = i + 1
				column = backslashQuoteIndex + 1
				fixedPosition = true
				break
			}
		}

		// 如果上面的检查没有发现问题，继续检查其他常见问题
		if !fixedPosition && line > 1 && (strings.Contains(err.Error(), "invalid character") || strings.Contains(err.Error(), "unexpected") || strings.Contains(err.Error(), "delimiter")) {
			// 检查前一行是否有未正确闭合的引号或转义字符问题
			prevLine := lines[line-2]

			// 检查是否有未正确转义的引号，如 "something\"
			badEscapeQuoteIndex := strings.LastIndex(prevLine, "\\\"")
			if badEscapeQuoteIndex >= 0 {
				// 更新错误位置到前一行
				line = line - 1
				column = badEscapeQuoteIndex + 1 // 指向问题字符的位置
				fixedPosition = true
			}

			// 检查是否有引号后面多了引号的情况，如 "something""
			badQuoteIndex := strings.LastIndex(prevLine, "\"\"")
			if badQuoteIndex >= 0 {
				// 更新错误位置到前一行
				line = line - 1
				column = badQuoteIndex + 1 // 指向第二个引号
				fixedPosition = true
			}

			// 检查是否有JSON属性值中未闭合的引号，如 "value": "something 没有闭合
			if strings.Count(prevLine, "\"")%2 != 0 {
				// 查找最后一个引号的位置
				lastQuoteIndex := strings.LastIndex(prevLine, "\"")
				if lastQuoteIndex >= 0 {
					// 更新错误位置到前一行
					line = line - 1
					column = lastQuoteIndex + 1 // 指向问题字符的位置
					fixedPosition = true
				}
			}
		}

		// 获取更新后的行内容和指针位置
		if line > 0 && line <= len(lines) {
			if !fixedPosition {
				lineContent = lines[line-1]
			}
			pointer = strings.Repeat(" ", column-1) + "^"
		}

		// 解析Go错误消息，提取有效部分
		goErrMsg := err.Error()

		// 尝试将Go的错误消息转换为Python风格
		if strings.Contains(goErrMsg, "invalid character") {
			if strings.Contains(goErrMsg, "after") {
				errorMessageBase = "Expecting ',' delimiter"
			} else if strings.Contains(goErrMsg, "looking for beginning of") {
				errorMessageBase = "Expecting value"
			} else {
				errorMessageBase = "Invalid character"
			}
		} else {
			// 其他错误，保留原始消息
			errorMessageBase = goErrMsg
		}
	} else if unmarshalTypeError, ok := err.(*json.UnmarshalTypeError); ok {
		// 处理类型错误
		position = int(unmarshalTypeError.Offset)
		line, column = findLineAndColumn(originalStr, position)

		lines := strings.Split(originalStr, "\n")
		if line > 0 && line <= len(lines) {
			lineContent = lines[line-1]
			pointer = strings.Repeat(" ", column-1) + "^"
		}

		errorMessageBase = fmt.Sprintf("Cannot unmarshal %v into type %v",
			unmarshalTypeError.Value, unmarshalTypeError.Type.String())
	} else {
		// 其他类型的错误，尝试提供一个适当的位置
		position = 0
		line = 1
		column = 1
		errorMessageBase = err.Error()
	}

	// 构建Python风格的错误消息
	pythonStyleError := fmt.Sprintf("%s: line %d column %d (char %d)",
		errorMessageBase, line, column, position)
	errorMsg := fmt.Sprintf("JSON解析错误: %s", pythonStyleError)

	// 将位置信息添加到错误详情中 - 按照指定顺序添加字段
	// 删除message字段，确保pointer在line_content后面
	errorDetails["error_type"] = "JSONDecodeError"
	errorDetails["line"] = line
	errorDetails["column"] = column
	errorDetails["position"] = position
	errorDetails["error_char"] = errorChar
	errorDetails["line_content"] = lineContent
	errorDetails["pointer"] = pointer

	return compactStr, false, errorMsg, originalStr, errorDetails
}

// findLineAndColumn 根据字节偏移计算行号和列号
func findLineAndColumn(text string, offset int) (int, int) {
	// 确保偏移在有效范围内
	if offset < 0 {
		offset = 0
	}
	if offset > len(text) {
		offset = len(text)
	}

	// 计算到偏移位置之前的所有字符
	beforeOffset := text[:offset]

	// 计算换行符的数量，即行号
	lines := strings.Split(beforeOffset, "\n")
	line := len(lines)

	// 最后一行的内容长度，即列号
	var column int
	if line > 0 {
		column = len(lines[line-1]) + 1 // +1 因为列号通常从1开始
	} else {
		column = 1
	}

	return line, column
}

// ReturnAllRequest 处理 POST /returnAll 请求
func ReturnAllRequest(c *gin.Context) {
	// 生成请求ID
	requestID := uuid.New().String()

	// 读取请求体
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取请求体失败"})
		return
	}
	// 由于body已被读取，需要重新设置，以便其他中间件或处理函数可以再次读取
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// 获取请求信息
	requestInfo := model.NewRequestInfo(requestID)
	requestInfo.Method = c.Request.Method

	// 修正URL信息，这里重新构建URL而不是直接使用gin提供的URL
	// 确保URL中不包含多余的斜杠并使用正确的协议
	scheme := "http"
	// 从请求头获取转发的协议信息
	forwardedProto := c.Request.Header.Get("X-Forwarded-Proto")
	if forwardedProto != "" {
		scheme = forwardedProto // 优先使用X-Forwarded-Proto
	} else if c.Request.TLS != nil {
		scheme = "https"
	}

	host := c.Request.Host

	// 规范化URL路径，确保不包含多余斜杠
	urlPath := c.Request.URL.Path
	reg := regexp.MustCompile(`/+`)
	urlPath = reg.ReplaceAllString(urlPath, "/")

	query := c.Request.URL.RawQuery
	if query != "" {
		query = "?" + query
	}

	// 构建标准URL
	requestInfo.URL = scheme + "://" + host + urlPath + query

	// 获取路径参数
	for _, param := range c.Params {
		requestInfo.PathParams[param.Key] = param.Value
	}

	// 获取查询参数
	queryParams := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		queryParams[key] = values[0] // 只取第一个值
	}
	requestInfo.QueryParams = queryParams

	// 获取请求头信息，确保统一格式并避免重复
	headers := make(map[string]string)

	// 首先获取Host头
	headers["Host"] = c.Request.Host

	// 然后获取所有其他标准头信息
	for key, values := range c.Request.Header {
		if len(values) > 0 {
			// 统一使用原始头名称格式，避免重复
			headers[key] = values[0] // 只取第一个值
		}
	}

	// 设置请求头
	requestInfo.Headers = headers

	// 获取客户端信息
	clientIP := c.ClientIP()
	if clientIP != "" {
		// 直接使用Gin框架提供的ClientIP方法获取IP
		requestInfo.Client.Host = clientIP

		// 尝试从X-Forwarded-For或X-Real-IP头获取更多信息
		// 优先使用X-Real-IP
		realIP := c.Request.Header.Get("X-Real-IP")
		if realIP != "" {
			requestInfo.Client.Host = realIP
		} else {
			// 如果没有X-Real-IP，尝试使用X-Forwarded-For的第一个IP
			forwardedFor := c.Request.Header.Get("X-Forwarded-For")
			if forwardedFor != "" {
				ips := strings.Split(forwardedFor, ",")
				if len(ips) > 0 {
					requestInfo.Client.Host = strings.TrimSpace(ips[0])
				}
			}
		}

		// 如果需要端口信息，从RemoteAddr中提取
		_, port, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err == nil && port != "" {
			requestInfo.Client.Port = port
		} else {
			// 默认端口信息
			requestInfo.Client.Port = "unknown"
		}
	} else {
		// 直接解析RemoteAddr
		host, port, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err == nil {
			requestInfo.Client.Host = host
			requestInfo.Client.Port = port
		} else {
			// 如果解析失败，直接使用RemoteAddr
			requestInfo.Client.Host = c.Request.RemoteAddr
			requestInfo.Client.Port = "unknown"
		}
	}

	// 解析请求体
	bodyContent, isValidJSON, errorMsg, originalStr, errorDetails := parseJSONBody(bodyBytes)
	requestInfo.Body = bodyContent
	requestInfo.IsValidJSON = isValidJSON
	if !isValidJSON && errorMsg != "" {
		requestInfo.JSONError = errorMsg
		requestInfo.OriginalBody = originalStr
		requestInfo.ErrorDetails = errorDetails
	}

	// 记录日志
	if logger != nil {
		logger.Info(fmt.Sprintf("Received request: %s", requestID), requestInfo)
	}

	// 返回请求信息
	c.JSON(http.StatusOK, model.NewResponse(requestInfo))
}

// ReturnAllWithAnyPath 处理 POST /returnAll/* 任意路径请求
func ReturnAllWithAnyPath(c *gin.Context) {
	// 获取完整路径
	path := c.Param("path")

	// 规范化路径：使用正则替换连续的斜杠为单个斜杠
	reg := regexp.MustCompile(`/+`)
	path = reg.ReplaceAllString(path, "/")

	// 生成请求ID
	requestID := uuid.New().String()

	// 读取请求体
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取请求体失败"})
		return
	}
	// 由于body已被读取，需要重新设置，以便其他中间件或处理函数可以再次读取
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// 获取请求信息
	requestInfo := model.NewRequestInfo(requestID)
	requestInfo.Method = c.Request.Method

	// 修正URL信息，这里重新构建URL而不是直接使用gin提供的URL
	// 确保URL中不包含多余的斜杠并使用正确的协议
	scheme := "http"
	// 从请求头获取转发的协议信息
	forwardedProto := c.Request.Header.Get("X-Forwarded-Proto")
	if forwardedProto != "" {
		scheme = forwardedProto // 优先使用X-Forwarded-Proto
	} else if c.Request.TLS != nil {
		scheme = "https"
	}

	host := c.Request.Host

	// 规范化URL路径，确保不包含多余斜杠
	urlPath := c.Request.URL.Path
	urlPath = reg.ReplaceAllString(urlPath, "/")

	query := c.Request.URL.RawQuery
	if query != "" {
		query = "?" + query
	}

	// 构建标准URL
	requestInfo.URL = scheme + "://" + host + urlPath + query

	// 设置规范化后的路径 - 该字段保留在顶层
	requestInfo.Path = path
	// 不再在path_params中设置path字段
	// requestInfo.PathParams["path"] = path

	// 将路径拆分成多个部分
	// 保留原始路径，但过滤掉空路径部分
	splitParts := strings.Split(path, "/")

	// 过滤掉空字符串部分
	var pathParts []string
	for _, part := range splitParts {
		if part != "" {
			pathParts = append(pathParts, part)
		}
	}

	requestInfo.PathParts = pathParts

	// 设置路径参数 - 移除"path"参数，避免冗余
	for _, param := range c.Params {
		if param.Key != "path" { // 不添加path参数
			requestInfo.PathParams[param.Key] = param.Value
		}
	}

	// 添加路径各部分作为参数，只添加非空部分
	for i, part := range pathParts {
		requestInfo.PathParams[fmt.Sprintf("path_part_%d", i)] = part
	}

	// 获取查询参数
	queryParams := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		queryParams[key] = values[0] // 只取第一个值
	}
	requestInfo.QueryParams = queryParams

	// 获取请求头信息，确保统一格式并避免重复
	headers := make(map[string]string)

	// 首先获取Host头
	headers["Host"] = c.Request.Host

	// 然后获取所有其他标准头信息
	for key, values := range c.Request.Header {
		if len(values) > 0 {
			// 统一使用原始头名称格式，避免重复
			headers[key] = values[0] // 只取第一个值
		}
	}

	// 设置请求头
	requestInfo.Headers = headers

	// 获取客户端信息
	clientIP := c.ClientIP()
	if clientIP != "" {
		// 直接使用Gin框架提供的ClientIP方法获取IP
		requestInfo.Client.Host = clientIP

		// 优先使用X-Real-IP
		realIP := c.Request.Header.Get("X-Real-IP")
		if realIP != "" {
			requestInfo.Client.Host = realIP
		} else {
			// 如果没有X-Real-IP，尝试使用X-Forwarded-For的第一个IP
			forwardedFor := c.Request.Header.Get("X-Forwarded-For")
			if forwardedFor != "" {
				ips := strings.Split(forwardedFor, ",")
				if len(ips) > 0 {
					requestInfo.Client.Host = strings.TrimSpace(ips[0])
				}
			}
		}

		// 如果需要端口信息，从RemoteAddr中提取
		_, port, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err == nil && port != "" {
			requestInfo.Client.Port = port
		} else {
			// 默认端口信息
			requestInfo.Client.Port = "unknown"
		}
	} else {
		// 直接解析RemoteAddr
		host, port, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err == nil {
			requestInfo.Client.Host = host
			requestInfo.Client.Port = port
		} else {
			// 如果解析失败，直接使用RemoteAddr
			requestInfo.Client.Host = c.Request.RemoteAddr
			requestInfo.Client.Port = "unknown"
		}
	}

	// 解析请求体
	bodyContent, isValidJSON, errorMsg, originalStr, errorDetails := parseJSONBody(bodyBytes)
	requestInfo.Body = bodyContent
	requestInfo.IsValidJSON = isValidJSON
	if !isValidJSON && errorMsg != "" {
		requestInfo.JSONError = errorMsg
		requestInfo.OriginalBody = originalStr
		requestInfo.ErrorDetails = errorDetails
	}

	// 记录日志
	if logger != nil {
		logger.Info(fmt.Sprintf("Received any-path request: %s", requestID), requestInfo)
	}

	// 返回请求信息
	c.JSON(http.StatusOK, model.NewResponse(requestInfo))
}

// RootHandler 处理 GET / 根路径请求
func RootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "returnAll API服务正在运行",
		"endpoints": gin.H{
			"returnAll": "/returnAll (POST) - 返回请求的所有内容",
		},
		"version": "1.0.0",
	})
}
