package main

import (
	"fmt"
	"log"
	"time"

	"github.com/yzabiyaka/expr_demo/internal/app/cmdline"
	my_script "github.com/yzabiyaka/expr_demo/internal/app/environment"
	"github.com/yzabiyaka/expr_demo/internal/app/repository"
	"github.com/yzabiyaka/expr_demo/internal/pkg/model"
	"github.com/yzabiyaka/expr_demo/pkg/dbconnect"
	"github.com/yzabiyaka/expr_demo/pkg/script_env"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
)

func main() {
	db, closer, err := dbconnect.Connect()
	defer closer()

	if err != nil {
		panic(err)
	}

	var created time.Time
	var pos, country string

	if created, err = cmdline.Time(); err != nil {
		panic(err)
	}
	if pos, err = cmdline.POS(); err != nil {
		panic(err)
	}
	if country, err = cmdline.Country(); err != nil {
		panic(err)
	}

	var code *string
	var compiledExpression *vm.Program
	environment := my_script.New(model.Auth{})

	if code, err = cmdline.Rule(); err != nil {
		panic(err)
	} else {
		if code != nil {
			compiledExpression, err = expr.Compile(*code, expr.Env(environment))
			if err != nil {
				panic(err)
			}
		}
	}

	myModel := model.Auth{
		Created_at:   created,
		Account:      cmdline.Account(),
		Amount_cents: cmdline.AmountCents(),
		Pos:          pos,
		Country:      country,
		// также можно сделать специальное поле в модели, через которое извлекать данные или использовать предподготовленные списки
		PosesBlackList: func(a ...interface{}) ([]string, error) {
			fmt.Println("In PosesBlackList")
			for i := range a {
				fmt.Println(a[i])
			}
			return []string{"155", "151"}, nil
		},
	}
	repo := repository.NewDataRepo(db)
	if myModel.ID, err = repo.Save(myModel); err != nil {
		panic(err)
	}
	log.Printf("data inserted: %+v", myModel)

	if compiledExpression != nil {
		// load environment

		script_env.Setup(
			&environment,
			script_env.Model(myModel),
			script_env.Repo(repo),
		)
		// run script
		fmt.Println("---------------------------")
		output, err := expr.Run(compiledExpression, environment)
		fmt.Println("\n---------------------------")
		if err != nil {
			panic(err)
		}
		fmt.Println(output)
		fmt.Println("---------------------------")
	}
}
