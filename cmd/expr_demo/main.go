package main

import (
	"fmt"
	"github.com/antonmedv/expr/vm"
	"log"
	"time"

	"expr_demo/internal/app/cmdline"
	"expr_demo/internal/app/model"
	"expr_demo/internal/app/repository"
	"expr_demo/pkg/dbconnect"

	"github.com/antonmedv/expr"
)

type Env struct {
	Event         model.Auth
	HistoryEvents []model.Auth
}

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
	if code, err = cmdline.Rule(); err != nil {
		panic(err)
	} else {
		if code != nil {
			compiledExpression, err = expr.Compile(*code, expr.Env(Env{}))
			if err != nil {
				panic(err)
			}
		}
	}

	myModel := model.Auth{
		CreatedAt:   created,
		Account:     cmdline.Account(),
		AmountCents: cmdline.AmountCents(),
		POS:         pos,
		CountryMnem: country,
	}
	repo := repository.NewDataRepo(db)
	if myModel.ID, err = repo.Save(myModel); err != nil {
		panic(err)
	}
	log.Printf("data inserted: %+v", myModel)

	if compiledExpression != nil {
		// load environment
		env := Env{
			Event:         myModel,
			HistoryEvents: []model.Auth{},
		}
		// run script
		fmt.Println("---------------------------")
		output, err := expr.Run(compiledExpression, env)
		fmt.Println("---------------------------")
		if err != nil {
			panic(err)
		}
		fmt.Println(output)
		fmt.Println("---------------------------")
	}
}
