package environment

import "github.com/yzabiyaka/expr_demo/pkg/script_env"

// Environment expressions env
type Environment map[string]interface{}

// New makes a new Env
func New(model interface{}) Environment {
	e := make(Environment)
	script_env.Init(&e, model)
	return e
}

// AddKey adds a  key
func (e *Environment) AddKey(key string, val interface{}) {
	(*e)[key] = val
}

// COUNTRY is embedded func that has access to event and repo
func (e Environment) COUNTRY(lst ...interface{}) (interface{}, error) {
	return script_env.GetFromRepo(e, "country")
}
