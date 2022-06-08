package cmdline

import (
	"flag"
	"fmt"
	"io/ioutil"
	mrand "math/rand"
	"sync"
	"time"

	"github.com/yzabiyaka/expr_demo/pkg/testsupport"

	"github.com/pkg/errors"
)

const configDefault = "config.cfg"

var (
	once                               sync.Once
	accPtr, opTime, posPtr, countryPtr *string
	amountPtr                          *int
	ruleFile, rulePtr                  *string
)

func init() {
	accPtr = flag.String("acc", "", "account number 20 digs, default random")
	opTime = flag.String("time", "", "operation time, default random in last 30 days")
	amountPtr = flag.Int("cents", 0, "amoint in cents, default is up to 10000_00")
	posPtr = flag.String("pos", "", "pos, required")
	countryPtr = flag.String("cn", "", "country, required")
	ruleFile = flag.String("rulefile", "", "rule file name, default if config.cfg")
	rulePtr = flag.String("rule", "", "rule file name")
}

func parseFlags() {
	flag.Parse()
}

func Account() string {
	once.Do(parseFlags)
	if *accPtr == "" {
		return testsupport.RandomAccount()
	}
	return *accPtr
}

func Time() (time.Time, error) {
	once.Do(parseFlags)
	if *opTime == "" {
		return time.Now(), nil
	}
	operationTime, err := time.Parse("2006-01-02 15:04:05", *opTime)
	if err != nil {
		return time.Now(), errors.Wrap(err, "incorrect time value")
	}
	return operationTime, nil
}

func AmountCents() uint64 {
	once.Do(parseFlags)
	if *amountPtr == 0 {
		x := mrand.Int63n(100)
		return uint64(x * 10000)
	}
	return uint64(*amountPtr)
}

func POS() (string, error) {
	once.Do(parseFlags)
	if *posPtr == "" {
		return "", errors.New("POS is required")
	}
	return *posPtr, nil
}

func Country() (string, error) {
	once.Do(parseFlags)
	if *countryPtr == "" {
		return "", errors.New("Country is required")
	}
	return *countryPtr, nil
}

func Rule() (*string, error) {
	once.Do(parseFlags)
	if *ruleFile == "" {
		if *rulePtr != "" {
			return rulePtr, nil
		}
		return tryOpen(configDefault, true)
	}
	return tryOpen(*ruleFile, false)
}

func tryOpen(fileName string, ignoreError bool) (*string, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		if ignoreError {
			return nil, nil
		}
		return nil, errors.Wrap(err, fmt.Sprintf("cannot open rule file %s", fileName))
	}
	rule := string(bytes)
	return &rule, nil
}
