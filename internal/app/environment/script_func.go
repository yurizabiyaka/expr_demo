package environment

import "github.com/yzabiyaka/expr_demo/pkg/script_env"

// Environment expressions env
type Environment map[string]interface{}

// COUNTRY is embedded func that has access to event and repo
func (e Environment) COUNTRY(lst ...interface{}) (interface{}, error) {
	return script_env.GetFromRepo(e, "country")
}
