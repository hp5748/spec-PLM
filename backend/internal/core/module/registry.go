package module

import "sort"

// Registry 模块注册中心
type Registry struct {
	modules map[string]Module
	order   []string // 模块初始化顺序
}

func NewRegistry() *Registry {
	return &Registry{
		modules: make(map[string]Module),
		order:   make([]string, 0),
	}
}

// Register 注册模块
func (r *Registry) Register(m Module) {
	r.modules[m.Name()] = m
	r.order = append(r.order, m.Name())
}

// Get 获取模块
func (r *Registry) Get(name string) Module {
	return r.modules[name]
}

// All 获取所有模块（按注册顺序）
func (r *Registry) All() []Module {
	modules := make([]Module, 0, len(r.modules))
	for _, name := range r.order {
		if m, ok := r.modules[name]; ok {
			modules = append(modules, m)
		}
	}
	return modules
}

// Names 获取所有模块名称
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.modules))
	for name := range r.modules {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
