package script_env

import (
	"testing"
	"time"

	"github.com/yzabiyaka/expr_demo/pkg/script_env/mocks"
	"github.com/yzabiyaka/expr_demo/pkg/testsupport"

	"github.com/antonmedv/expr"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

// bench model
type auth struct {
	ID           uuid.UUID `db:"id"`
	Created_at   time.Time `db:"created_at"`
	Account      string    `db:"account"`
	Amount_cents uint64    `db:"amount_cents"`
	Pos          string    `db:"pos"`
	Country      string    `db:"country"`
}

func Benchmark_CodeRun(b *testing.B) {
	b.StopTimer()

	code := `e.Pos == '155' and e.Country in COUNTRY( "account==e.Account and pos == e.Pos and country IS NOT NULL","SUM(amount_cents)>3000000 and SUM(amount_cents)<60000000")`
	environment := NewEnvironment(auth{})
	compiledExpression, err := expr.Compile(code, expr.Env(environment))
	if err != nil {
		b.Fatalf("compilation failed: %v", err)
	}
	myModel := auth{
		Created_at:   time.Now(),
		Account:      testsupport.RandomAccount(),
		Amount_cents: 10_000_00,
		Pos:          "155",
		Country:      "RU",
	}

	ctrl := gomock.NewController(b)
	repo := mocks.NewMockDataRepo(ctrl)
	Setup(
		&environment,
		Model(myModel),
		Repo(repo),
	)
	repo.EXPECT().GetStringsFromData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]string{"RU", "CN"}, nil).AnyTimes()
	
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, err := expr.Run(compiledExpression, environment)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}
