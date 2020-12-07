package service

import (
	"context"

	"github.com/pkg/errors"
)

type bracerType int

const (
	btUnknown     bracerType = 0
	btCommonOpen  bracerType = 1
	btCommonClose bracerType = -1
	btSquareOpen  bracerType = 2
	btSquareClose bracerType = -2
	btFigureOpen  bracerType = 3
	btFigureClose bracerType = -3
)

var btypes = map[rune]bracerType{
	'(': btCommonOpen,
	')': btCommonClose,
	'[': btSquareOpen,
	']': btSquareClose,
	'{': btFigureOpen,
	'}': btFigureClose,
}

type Service interface {
	Validate(ctx context.Context, input string) (valid bool, err error)
}

type service struct {
}

func (s *service) Validate(ctx context.Context, input string) (valid bool, err error) {
	if len(input)%2 != 0 {
		return
	}

	stack := make([]bracerType, 0, len(input)/2)

	for _, elem := range input {
		bt, ok := btypes[elem]
		if !ok {
			err = errors.New("invalid input")
			return
		}
		switch bt {
		case btCommonOpen, btSquareOpen, btFigureOpen:
			stack = append(stack, bt)
		case btCommonClose, btSquareClose, btFigureClose:
			l := len(stack)
			last := stack[l-1]
			if -last != bt {
				return
			}
			stack = stack[:l-1]
		default:
			err = errors.New("internal error")
			return
		}
	}
	valid = true
	return
}

func NewService() (svc Service) {
	svc = &service{}
	return
}
