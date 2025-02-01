package simple

import "flag"

type SimplePlugin struct {
	name  string
	value string
}

func (s *SimplePlugin) GetPrefix() string {
	return s.name
}

func NewSimplePlugin(name string) *SimplePlugin {
	return &SimplePlugin{name: name}
}

func (s *SimplePlugin) Get() interface{} {
	return s
}

func (s *SimplePlugin) Name() string {
	return s.name
}

func (s *SimplePlugin) InitFlags() {
	flag.StringVar(&s.value, "simple-value", "default", "Simple plugin value")
}

func (s *SimplePlugin) Configure() error {
	return nil
}

func (s *SimplePlugin) Run() error {
	return nil
}

func (s *SimplePlugin) Stop() <-chan bool {
	c := make(chan bool)
	go func() {
		c <- true
	}()
	return c
}

func (s *SimplePlugin) GetValue() string {
	return s.value
}
