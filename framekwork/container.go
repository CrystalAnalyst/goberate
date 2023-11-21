package framework

type Container interface {
	Bind(comer ServiceProvider) error
	Is_bind(key string) bool
	GetService(key string) (interface{}, error)
	MustGetService(key string) interface{}
	GetServiceByParams(key string, params []any) (interface{}, error)
}
