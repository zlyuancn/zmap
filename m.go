/*
-------------------------------------------------
   Author :       zlyuan
   date：         2019/7/29
   Description :
-------------------------------------------------
*/

package zmap

import (
    "strings"
)

type M map[string]interface{}

// 添加一个值, 返回是否成功
func (m M) Add(k string, v interface{}, filters ...MapFilter) bool {
    if m == nil {
        return false
    }

    if len(filters) == 0 {
        return true
    }

    if filterKV(k, v, filters...) {
        return false
    }
    return true
}

// 过滤
func (m M) Filter(filters ...MapFilter) M {
    if len(m) == 0 || len(filters) == 0 {
        return m
    }

    for k, v := range m {
        if filterKV(k, v, filters...) {
            delete(m, k)
        }
    }
    return m
}

// 根据指定过滤函数过滤
func (m M) FilterOf(fn FilterFn) M {
    if len(m) == 0 || fn == nil {
        return m
    }

    for k, v := range m {
        if fn(k, v) {
            delete(m, k)
        }
    }
    return m
}

func filterKV(k string, v interface{}, filters ...MapFilter) bool {
    for _, filter := range filters {
        if filter.filter(k, v) {
            return true
        }
    }
    return false
}

// 遍历
func (m M) Foreach(fn func(k string, v interface{})) M {
    if len(m) == 0 || fn == nil {
        return m
    }
    for k, v := range m {
        fn(k, v)
    }
    return m
}

// 返回key数量
func (m M) Len() int {
    return len(m)
}

// 获取一个值, 如果key不存在返回默认值
func (m M) GetDefault(k string, default_value interface{}) interface{} {
    v, ok := m[k]
    if ok {
        return v
    }
    return default_value
}

// 弹出一个值
func (m M) Pop(k string) (interface{}, bool) {
    v, ok := m[k]
    if ok {
        delete(m, k)
        return v, true
    }
    return nil, false
}

// 弹出一个值, 如果key不存在返回默认值
func (m M) PopDefault(k string, default_value interface{}) interface{} {
    v, ok := m[k]
    if ok {
        delete(m, k)
        return v
    }
    return default_value
}

// 浅拷贝获取一个副本
func (m M) Copy() M {
    if len(m) == 0 {
        return make(M)
    }

    out := make(M, len(m))
    for k, v := range m {
        out[k] = v
    }
    return out
}

// 更新, 如果自己没初始化会panic
func (m M) Update(data M) {
    if m == nil {
        panic("未初始化的zmap.M")
    }
    for k, v := range data {
        m[k] = v
    }
}

// 更新忽略大小写, 如果自己没初始化会panic
func (m M) UpdateIgnoreCase(data M) {
    if m == nil || data == nil {
        panic("未初始化的zmap.M")
    }

    if len(data) == 0 {
        return
    }

    buf := make(M, len(data))
    for k, v := range data {
        buf[strings.ToLower(k)] = v
    }

    for k := range m {
        lk := strings.ToLower(k)
        if v, ok := buf[lk]; ok {
            m[k] = v
        }
    }
}
