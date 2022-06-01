package script_env

import (
	"expr_demo/internal/pkg/model"
	"fmt"
	"regexp"
)

type Environment map[string]interface{}

func Create() Environment {
	env := Environment{
		"e": model.Auth{},
		"h": []model.Auth{},
		"RECORDS": func(lst ...interface{}) []string {
			selDef := "*"
			whereDef := ""
			if len(lst) > 0 {
				str := lst[0].(string)
				if str != "*" {
					whereDef = "WHERE" + str
				}
			}
			if len(lst) > 1 {
				//fmt.Println(lst[1])
				regx := regexp.MustCompile(`(\w+)\s+(SUM|COUNT)\s*\((\w+|\*)\)`)
				str := lst[1].(string)
				vals := regx.FindStringSubmatch(str)
				//fmt.Println(vals)
				if len(vals) > 3 {
					selDef = fmt.Sprintf("%s, %s(%s)", vals[1], vals[2], vals[3])
					if vals[3] != "*" {
						whereDef += " GROUP BY " + vals[3]
					}
				}
			}
			if len(lst) > 2 {
				whereDef += " HAVING " + lst[2].(string)
			}
			fmt.Printf("SELECT %s FROM data %s", selDef, whereDef)
			return []string{
				"CN", "CZ", "PL",
			}
		},
	}
	return env
}

func New(opts ...Opts) Environment {
	e := Create()
	for _, opt := range opts {
		opt(&e)
	}
	return e
}

type Opts func(env *Environment)

func Event(evt model.Auth) Opts {
	return func(env *Environment) {
		(*env)["e"] = evt
	}
}
