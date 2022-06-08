package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Auth struct {
	ID           uuid.UUID `db:"id"`
	Created_at   time.Time `db:"created_at"`
	Account      string    `db:"account"`
	Amount_cents uint64    `db:"amount_cents"`
	Pos          string    `db:"pos"`
	Country      string    `db:"country"`
	// ----------------
	// Poses - можно сделать специальное поле в модели, через которое извлекать данные. В функцию можно передать необходимые параметры заранее, в том числе экземпляр подготовленного окружения, где хранятся названия полей и регэкспы
	PosesBlackList func(args ...interface{}) ([]string, error)
}

// COUNTRY - можно методы извлечения исторических данных привязать к модели, но тогда требуется обратная трансляция параметров в библиотеку для получения выражений where, group by и having. Также необходим экземпляр подготовленного окружения, где хранятся названия полей и регэкспы
func (a Auth) COUNTRY(agrs ...interface{}) ([]string, error) {
	for i := range agrs {
		fmt.Println(agrs[i])
	}
	return nil, nil
}
