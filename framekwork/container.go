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
func (gobe *GobeContainer) Bind(comer ServiceProvider) error {
	gobe.lock.Lock()
	name := comer.Name()
	gobe.providers[name] = comer
	gobe.lock.Unlock()
	if comer.Is_defer() == false {
		// 注册即实例化(必须返回一个可用的服务实例).
		// 而服务实例需要先调用Boot完成参数、配置的初始化.
		if err := comer.Boot(gobe); err != nil {
			return err
		}
		param := comer.Params(gobe)
		newInstanceFunc := comer.Register(gobe)
		inst, err := newInstanceFunc(param...)
		if err != nil {
			// 实例化失败.
		}
		gobe.instances[name] = inst
	}
	return nil
}
