package framework

import (
	"errors"
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
			fmt.Println("bind failure! key: {}, error: {}!", name, err)
			return errors.New(err.Error())
		}
		gobe.instances[name] = inst
	}
	return nil
}

func (gobe *GobeContainer) Is_bind(key string) bool {
	return (gobe.GetServiceProvider(key) != nil)
}

func (gobe *GobeContainer) GetServiceProvider(key string) ServiceProvider {
	gobe.lock.RLock()
	defer gobe.lock.RUnlock()
	if sp, ok := gobe.providers[key]; ok {
		return sp
	}
	return nil
}

func (gobe *GobeContainer) GetService(key string) (interface{}, error) {
	return gobe.GetServiceInstance(key, nil, false)
}

func (gobe *GobeContainer) MustGetService(key string) interface{} {
	inst, err := gobe.GetServiceInstance(key, nil, false)
	if err == nil {
		return inst
	} else {
		panic("Get Service failure!")
	}
}

func (gobe *GobeContainer) NewService(key string, params []any) (interface{}, error) {
	return gobe.GetServiceInstance(key, nil, true)
}

func (gobe *GobeContainer) GetServiceInstance(key string, params []any, force bool) (interface{}, error) {
	/*
		@params:
		key: 通过key得知服务的名字,用key去拿到sp,然后用于实例化服务.
		params: 用于在实例化服务时加入额外的参数.
		force: 是否需要强制重新实例化.
	*/
	gobe.lock.RLock()
	defer gobe.lock.RUnlock()
	sp := gobe.GetServiceProvider(key)
	if sp == nil {
		return nil, errors.New("contract" + key + "has not registerded yet")
	}

	if force == true {
		return gobe.newInstance(sp, params)
	}

	if inst, ok := gobe.instances[key]; ok {
		return inst, nil
	}

	inst, err := gobe.newInstance(sp, params)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return inst, nil
}

func (gobe *GobeContainer) newInstance(sp ServiceProvider, params []any) (interface{}, error) {
	if err := sp.Boot(gobe); err != nil {
		return nil, errors.New(err.Error())
	}
	if params == nil {
		params = sp.Params(gobe)
	}
	method := sp.Register(gobe)
	inst, err := method(params)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return inst, nil
}
