package framework

type Instance func(...any) (interface{}, error)

// Service: Register -> Instantialize -> Boot -> Work.
type ServiceProvider interface {
	Register(Container) Instance
	Is_defer() bool         //用于控制实例化的时机,为false的时候“注册即实例化”.
	Boot(Container) error   //在实例化后立马调用Boot进行参数初始化、配置等工作.
	Name() string           //服务名
	Params(Container) []any //传递给服务实例的参数
}
