package framework

import (
	"fmt"
	"sync"
)

type Container interface {
	Bind(comer ServiceProvider) error
	Is_bind(key string) bool
	GetService(key string) (interface{}, error)
	MustGetService(key string) interface{}
	GetServiceByParams(key string, params []any) (interface{}, error)
}

// 服务容器实例
type GobeContainer struct {
	Container
	providers map[string]ServiceProvider
	instances map[string]any
	lock      sync.RWMutex
}

func NewGobeContainer() *GobeContainer {
	return &GobeContainer{
		providers: map[string]ServiceProvider{},
		instances: map[string]any{},
		lock:      sync.RWMutex{},
	}
}

func (gobe *GobeContainer) GetProviderList() []string {
	ret := []string{}
	for _, provider := range gobe.providers {
		name := provider.Name()
		line := fmt.Sprintf(name)
		ret = append(ret, line)
	}
	return ret
}

/*----为GobeContainer实现Container接口----*/
