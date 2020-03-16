/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/3/15
   Description :
-------------------------------------------------
*/

package zmap

import (
    "strings"
)

// 过滤器函数
type FilterFn func(k string, v interface{}) (filter bool)

// 过滤器
type MapFilter interface {
    // 将map的每个值传入这个方法, 返回bool值决定是否过滤掉
    filter(k string, v interface{}) (filter bool)
}

var _ MapFilter = (*mapFilterWarpper)(nil)

type mapFilterWarpper struct {
    fn FilterFn
}

// 将一个过滤器函数包装为过滤器
func WrapFnToFilter(fn FilterFn) MapFilter {
    return &mapFilterWarpper{fn: fn}
}

func (m mapFilterWarpper) filter(k string, v interface{}) (filter bool) {
    return m.fn(k, v)
}

var _ MapFilter = (*PrefixFilter)(nil)

// key前缀过滤器
type PrefixFilter struct {
    texts      []string
    check_case bool
}

// 创建key前缀过滤器
func NewPrefixFilter(text ...string) *PrefixFilter {
    return &PrefixFilter{
        texts:      append(([]string)(nil), text...),
        check_case: true,
    }
}

// key是否大小写敏感, 默认true
func (m *PrefixFilter) SetCheckCase(b bool) *PrefixFilter {
    m.check_case = b
    return m
}

// 设置需要过滤的前缀
func (m *PrefixFilter) SetText(text ...string) *PrefixFilter {
    m.texts = append(([]string)(nil), text...)
    return m
}

// 添加需要过滤的前缀
func (m *PrefixFilter) AddText(text ...string) *PrefixFilter {
    m.texts = append(m.texts, text...)
    return m
}

func (m *PrefixFilter) filter(k string, v interface{}) (filter bool) {
    if m.check_case {
        for _, t := range m.texts {
            if strings.HasPrefix(k, t) {
                return true
            }
        }
        return false
    }

    for _, t := range m.texts {
        if strings.HasPrefix(strings.ToLower(k), strings.ToLower(t)) {
            return true
        }
    }
    return false
}

// key后缀过滤器
type SufixFilter struct {
    texts      []string
    check_case bool
}

// 创建key后缀过滤器
func NewSufixFilter(text ...string) *SufixFilter {
    return &SufixFilter{
        texts:      append(([]string)(nil), text...),
        check_case: true,
    }
}

// key是否大小写敏感, 默认true
func (m *SufixFilter) SetCheckCase(b bool) *SufixFilter {
    m.check_case = b
    return m
}

// 设置需要过滤的后缀
func (m *SufixFilter) SetText(text ...string) *SufixFilter {
    m.texts = append(([]string)(nil), text...)
    return m
}

// 添加需要过滤的后缀
func (m *SufixFilter) AddText(text ...string) *SufixFilter {
    m.texts = append(m.texts, text...)
    return m
}

func (m *SufixFilter) filter(k string, v interface{}) (filter bool) {
    if m.check_case {
        for _, t := range m.texts {
            if strings.HasPrefix(k, t) {
                return true
            }
        }
        return false
    }

    for _, t := range m.texts {
        if strings.HasSuffix(strings.ToLower(k), strings.ToLower(t)) {
            return true
        }
    }
    return false
}

// key匹配过滤器
type MatchFilter struct {
    texts      []string
    check_case bool
    match      bool
}

// 创建key匹配过滤器
func NewMatchFilter(text ...string) *MatchFilter {
    return &MatchFilter{
        texts:      append(([]string)(nil), text...),
        check_case: true,
        match:      false,
    }
}

// key是否大小写敏感, 默认true
func (m *MatchFilter) SetCheckCase(b bool) *MatchFilter {
    m.check_case = b
    return m
}

// key是否完全匹配, 表示k和指定的文本相等, 默认false
func (m *MatchFilter) SetMatch(b bool) *MatchFilter {
    m.match = b
    return m
}

// 设置需要过滤的文本
func (m *MatchFilter) SetText(text ...string) *MatchFilter {
    m.texts = append(([]string)(nil), text...)
    return m
}

// 添加需要过滤的文本
func (m *MatchFilter) AddText(text ...string) *MatchFilter {
    m.texts = append(m.texts, text...)
    return m
}

func (m *MatchFilter) filter(k string, v interface{}) (filter bool) {
    var fn func(k, t string) bool
    if m.check_case {
        if m.match {
            fn = func(k, t string) bool {
                return k == t
            }
        } else {
            fn = func(k, t string) bool {
                return strings.Contains(k, t)
            }
        }
    } else if m.match {
        fn = func(k, t string) bool {
            return strings.ToLower(k) == strings.ToLower(t)
        }
    } else {
        fn = func(k, t string) bool {
            return strings.Contains(strings.ToLower(k), strings.ToLower(t))
        }
    }

    for _, t := range m.texts {
        if fn(k, t) {
            return true
        }
    }
    return false
}
