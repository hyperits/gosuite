// Package errors 提供统一的错误处理
package errors

import (
	"errors"
	"fmt"
)

// 标准错误变量
var (
	// ErrNotConfigured 未配置错误
	ErrNotConfigured = errors.New("not configured")

	// ErrNotConnected 未连接错误
	ErrNotConnected = errors.New("not connected")

	// ErrNilConfig 配置为空错误
	ErrNilConfig = errors.New("config is nil")

	// ErrNilClient 客户端为空错误
	ErrNilClient = errors.New("client is nil")

	// ErrInvalidParameter 参数无效错误
	ErrInvalidParameter = errors.New("invalid parameter")

	// ErrNotFound 资源未找到错误
	ErrNotFound = errors.New("not found")

	// ErrTimeout 超时错误
	ErrTimeout = errors.New("timeout")

	// ErrAlreadyClosed 已关闭错误
	ErrAlreadyClosed = errors.New("already closed")
)

// New 创建一个新错误
func New(text string) error {
	return errors.New(text)
}

// Wrap 包装错误，添加上下文信息
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// Wrapf 包装错误，支持格式化
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
}

// Is 判断错误是否匹配
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As 将错误转换为目标类型
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Unwrap 解包错误
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// Join 合并多个错误（Go 1.20+）
func Join(errs ...error) error {
	return errors.Join(errs...)
}

// OpError 操作错误，包含操作名称和原因
type OpError struct {
	Op   string // 操作名称，如 "mysql.connect", "redis.get"
	Kind string // 错误类型，如 "connection", "timeout", "validation"
	Err  error  // 原始错误
}

func (e *OpError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("%s: %s error", e.Op, e.Kind)
	}
	return fmt.Sprintf("%s: %s error: %v", e.Op, e.Kind, e.Err)
}

func (e *OpError) Unwrap() error {
	return e.Err
}

// NewOpError 创建操作错误
func NewOpError(op, kind string, err error) *OpError {
	return &OpError{
		Op:   op,
		Kind: kind,
		Err:  err,
	}
}
