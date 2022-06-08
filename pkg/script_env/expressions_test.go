package script_env

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// environment environment for test purposes
type environment map[string]interface{}

// NewEnvironment constructor for test env
func NewEnvironment(evt interface{}) environment {
	e := make(environment)
	Init(&e, evt)
	return e
}

// AddKey implements Modifier interface
func (e *environment) AddKey(key string, val interface{}) {
	(*e)[key] = val
}

// COUNTRY is a test method
func (e environment) COUNTRY(arg0, arg1 interface{}) ([]string, error) {
	return nil, nil
}

func Test_getSelectExpressions(t *testing.T) {

	type args struct {
		params []interface{}
	}

	type res struct {
		where, having string
		fields        []string
	}

	type test struct {
		name    string
		a       args
		wantErr bool
		wantRes res
	}

	tests := []test{
		{
			name: "1. Naive grouping",
			a: args{
				params: []interface{}{
					"*", "SUM(amount_cents)>1000000",
				},
			},
			wantErr: false,
			wantRes: res{
				where:  "",
				having: "sum(amount_cents)>1000000",
				fields: []string{},
			},
		},
		{
			name: "2. Where expression with one cond",
			a: args{
				params: []interface{}{
					`account = e.Account`, "SUM(amount_cents)>1000000",
				},
			},
			wantErr: false,
			wantRes: res{
				where:  "account = $1",
				having: "sum(amount_cents)>1000000",
				fields: []string{"account"},
			},
		},
		{
			name: "3. Where expression with two cond",
			a: args{
				params: []interface{}{
					`account = e.Account and pos=e.Pos`, "SUM(amount_cents)>1000000",
				},
			},
			wantErr: false,
			wantRes: res{
				where:  "account = $1 and pos=$2",
				having: "sum(amount_cents)>1000000",
				fields: []string{"account", "pos"},
			},
		},
		{
			name: "4. && with two cond",
			a: args{
				params: []interface{}{
					`account = e.Account && pos=e.Pos`, "SUM(amount_cents)>1000000",
				},
			},
			wantErr: false,
			wantRes: res{
				where:  "account = $1 and pos=$2",
				having: "sum(amount_cents)>1000000",
				fields: []string{"account", "pos"},
			},
		},
		{
			name: "5. Double = converts to single =",
			a: args{
				params: []interface{}{
					`account == e.Account`,
				},
			},
			wantErr: false,
			wantRes: res{
				where:  "account = $1",
				having: "",
				fields: []string{"account"},
			},
		},
		{
			name: "6. Complex expr",
			a: args{
				params: []interface{}{
					`(amount_cents > (e.Amount_cents + 123)) &&(pos==e.Pos)`,
				},
			},
			wantErr: false,
			wantRes: res{
				where:  "(amount_cents > ($1 + 123)) and (pos = $2)",
				having: "",
				fields: []string{"amount_cents", "pos"},
			},
		},
		{
			name: "7. || converts to or",
			a: args{
				params: []interface{}{
					`account = e.Account || pos=e.Pos || country==e.Country`, "",
				},
			},
			wantErr: false,
			wantRes: res{
				where:  "account = $1 or pos=$3 or country = $2",
				having: "",
				fields: []string{"account", "country", "pos"},
			},
		},
		{
			name: "8. Empty groups",
			a: args{
				params: []interface{}{},
			},
			wantErr: false,
			wantRes: res{
				where:  "",
				having: "",
				fields: []string{},
			},
		},
	}

	env := NewEnvironment(auth{})
	Setup(&env, Model(auth{}))

	for _, tt := range tests {
		ttt := tt
		t.Run(ttt.name, func(t *testing.T) {
			w, h, fields, e := getSelectComponents(env, ttt.a.params...)
			if (e != nil) != ttt.wantErr {
				t.Errorf("want err %v, got %v", ttt.wantErr, e)
			}
			assert.Equal(t, ttt.wantRes.where, strings.Join(strings.Fields(w), " "))
			assert.Equal(t, ttt.wantRes.having, strings.Join(strings.Fields(h), " "))
			assert.Equal(t, ttt.wantRes.fields, fields)
		})
	}
}
