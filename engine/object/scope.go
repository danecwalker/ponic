// Copyright (c) 2022 DevDane <dane@danecwalker.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package object

type BindType int

const (
	CONST BindType = iota
	LET
	FUNC
)

type ValueBinding struct {
	Object Object
	Type   BindType
}

type Scope struct {
	Parent *Scope
	Values map[string]ValueBinding
}

func NewScope() *Scope {
	return &Scope{
		Values: make(map[string]ValueBinding),
	}
}

func (s *Scope) Get(name string) (Object, bool) {
	bind, ok := s.Values[name]
	if !ok && s.Parent != nil {
		return s.Parent.Get(name)
	}
	return bind.Object, ok
}

func (s *Scope) Set(name string, val Object, bindType BindType) {
	bind, ok := s.Values[name]
	if ok {
		if bind.Type == CONST {
			panic("Cannot reassign constant")
		} else {
			s.Values[name] = ValueBinding{val, bind.Type}
		}
	} else {
		s.Values[name] = ValueBinding{val, bindType}
	}
}
