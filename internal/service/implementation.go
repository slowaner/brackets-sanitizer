package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type bracketType int

const (
	btUnknown     bracketType = 0
	btCommonOpen  bracketType = 1
	btCommonClose bracketType = -1
	btSquareOpen  bracketType = 2
	btSquareClose bracketType = -2
	btFigureOpen  bracketType = 3
	btFigureClose bracketType = -3
)

var btypes = map[rune]bracketType{
	'(': btCommonOpen,
	')': btCommonClose,
	'[': btSquareOpen,
	']': btSquareClose,
	'{': btFigureOpen,
	'}': btFigureClose,
}

var btypeRuneStrings = map[bracketType]string{
	btUnknown:     "",
	btCommonOpen:  "(",
	btCommonClose: ")",
	btSquareOpen:  "[",
	btSquareClose: "]",
	btFigureOpen:  "{",
	btFigureClose: "}",
}

type Service interface {
	Validate(ctx context.Context, input string) (valid bool, err error)
	Sanitize(ctx context.Context, input string) (result string, err error)
}

type service struct {
}

type sanitizerChildren []*sanitizerElem

func (e sanitizerChildren) String() string {
	if len(e) == 0 {
		return ""
	}

	buf := strings.Builder{}
	for _, elem := range e {
		buf.WriteString(elem.String())
	}

	return buf.String()
}

type sanitizerElem struct {
	bt       bracketType
	parent   *sanitizerElem
	children sanitizerChildren
	closed   bool
}

func (e *sanitizerElem) String() string {
	var (
		oStr string
		cStr string
	)
	if e.closed {
		oStr = btypeRuneStrings[e.bt]
		cStr = btypeRuneStrings[-e.bt]
	}
	return fmt.Sprintf("%s%s%s", oStr, e.children, cStr)
}

func (s *service) Sanitize(ctx context.Context, input string) (result string, err error) {
	root := &sanitizerElem{
		bt:       btUnknown,
		children: sanitizerChildren{},
		closed:   false, // Must be false while processing
	}

	cur := root

	for _, elem := range input {
		bt, ok := btypes[elem]
		if !ok {
			continue
		}
		switch bt {
		case btCommonOpen, btSquareOpen, btFigureOpen:
			child := &sanitizerElem{
				bt:       bt,
				parent:   cur,
				children: sanitizerChildren{},
				closed:   false,
			}
			cur.children = append(cur.children, child)
			cur = child
		case btCommonClose, btSquareClose, btFigureClose:
			openner := -bt
			// Processing close stack elem
			{
				// FIXME: wrong order
				//  "input":  "{[{]}"
				//  "result": "[]{}"
				//  possible: "{[]}"
				c := cur
				for c.parent != nil {
					if c.bt == openner {
						// Found bracket to close
						c.closed = true
						restChildren := make([]*sanitizerElem, 0, len(c.children))
						for _, child := range c.children {
							if child.closed {
								restChildren = append(restChildren, child)
								continue
							}
							child.parent = c.parent
							child.parent.children = append(child.parent.children, child)
						}
						c.children = restChildren
						break
					}

					c = c.parent
				}
			}

			for cur.closed {
				cur = cur.parent
			}
		default:
			err = errors.New("internal error")
			return
		}
	}

	root.closed = true
	result = root.String()

	return
}

func (s *service) Validate(ctx context.Context, input string) (valid bool, err error) {
	if len(input)%2 != 0 {
		return
	}

	stack := make([]bracketType, 0, len(input)/2)

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
			if l == 0 {
				return
			}
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
