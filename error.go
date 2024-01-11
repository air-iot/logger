package logger

import "github.com/pkg/errors"

type LogError struct {
	Suggest string // 建议
	Focus   Focus  // 关注
	ERR     error  // 响应错误
}

func (l *LogError) Error() string {
	if l.ERR != nil {
		return l.ERR.Error()
	}
	return ""
}

func NewError(focus Focus, suggest string, err error) error {
	res := &LogError{
		Focus:   focus,
		Suggest: suggest,
		ERR:     errors.WithStack(err),
	}
	return res
}

func NewErrorFocusNotice(suggest string, err error) error {
	res := &LogError{
		Focus:   FocusNotice,
		Suggest: suggest,
		ERR:     errors.WithStack(err),
	}
	return res
}

func NewErrorFocus(focus Focus, err error) error {
	res := &LogError{
		Focus: focus,
		ERR:   errors.WithStack(err),
	}
	return res
}

func NewErrorSuggest(suggest string, err error) error {
	res := &LogError{
		Suggest: suggest,
		ERR:     errors.WithStack(err),
	}
	return res
}
