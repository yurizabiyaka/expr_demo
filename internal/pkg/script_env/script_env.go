package script_env

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	keyEvent                 = "e"
	keyRepo                  = "_repo"
	keyFieldNames            = "_fields"
	keyFieldsIndexes         = "_fieldsIndexes"
	keyFieldNameMatchers     = "_matchers"
	keyScriptToSqlConverters = "_converters"

	mockPlaceholder     = "%" // regexp methods Replace* uses $ as a special symbol
	postgresPlaceholder = "$"
)

var (
	matchAnd      = regexp.MustCompile(`([^&])(&{2})([^&])`)
	matchOr       = regexp.MustCompile(`([^\|])(\|{2})([^_\|])`)
	matchDoubleEq = regexp.MustCompile(`([^=])(={2})([^=])`)
)

/*

хочется видеть что-то
'country in country("account=account", "SUM(amount_cents)>1000000")'
1) вне кавычек подставить "e." перед полем, поле с большой буквы
2) внутри кавычек справа от = < > <= >= подставить "{{e." перед полем и "}}" после названия, поле с большой буквы, например amount >= 2*(amount+16) and pos=pos => amount >= 2*(e.Amount+16) and pos=e.Pos
3) агрегационные функции country, account, pos, и тп формировать для каждой новой таблицы, чтобы возвращали массив строк или чисел - соответствующие столбцы

*/

// Environment expressions env
type Environment map[string]interface{}

type stringConverter func(s string) string

// Create makes default env
func Create(event interface{}) Environment {
	env := Environment{
		keyEvent: event,
	}
	fields, fieldsIndexes := getFieldNamesAsStringsSlice(event) // in lowercase
	env[keyFieldNames] = fields                                 // in lowercase
	env[keyFieldsIndexes] = fieldsIndexes                       // in lowercase
	// matchers contain lowercase fields as keys
	env[keyFieldNameMatchers] = getMatchers(env[keyFieldNames].([]string)) // in lowercase
	env[keyScriptToSqlConverters] = getScriptToSqlConverters()

	return env
}

// New fills environment with values using modification opts
func New(event interface{}, opts ...Opts) Environment {
	e := Create(event)
	for _, opt := range opts {
		opt(&e)
	}
	return e
}

// Opts modifications
type Opts func(env *Environment)

//go:generate mockgen -destination=./mocks/data_repo.go -package=mocks . DataRepo

// DataRepo data access
type DataRepo interface {
	GetStringsFromData(column, whereDef, havingDef string, eventValues []interface{}) ([]string, error)
}

// Repo fills dao
func Repo(repo DataRepo) Opts {
	return func(env *Environment) {
		(*env)[keyRepo] = repo
	}
}

func getFieldNamesAsStringsSlice(event interface{}) ([]string, map[string]int) {
	indexes := make(map[string]int)
	t := reflect.TypeOf(event)

	names := make([]string, t.NumField())
	for i := range names {
		// todo проверить, что очередное поле не является struct или interface{} или массивом struct или interface{}. В случае struct можно углубиться с точками ?
		names[i] = strings.ToLower(t.Field(i).Name)
		indexes[names[i]] = i
	}
	return names, indexes
}

func getMatchers(fieldNames []string) map[string]*regexp.Regexp {
	matchers := make(map[string]*regexp.Regexp)
	keyEvt := strings.ToLower(keyEvent)
	for _, field := range fieldNames {
		regexpTemplate := fmt.Sprintf(`(^|\W)(%s\.%s)($|\W)`, keyEvt, field)
		matchers[field] = regexp.MustCompile(regexpTemplate)
	}
	return matchers
}

func (e Environment) getSelectComponents(lst ...interface{}) (whereDef, havingDef string, fieldsInPhOrder []string, err error) {
	fieldsInPhOrder = []string{}
	placeholders := make(map[string]string) // key is field name, value is placeholder
	matchers := e[keyFieldNameMatchers].(map[string]*regexp.Regexp)
	replacers := e[keyScriptToSqlConverters].([]stringConverter)

	whereArg := strings.ToLower(safeGetAsStr(lst, 0, "*"))
	havingArg := strings.ToLower(safeGetAsStr(lst, 1, ""))

	if whereArg != "*" {
		preparePlaceholders(&placeholders, &fieldsInPhOrder, whereArg, matchers)
		whereArg = replaceFieldsWithPlaceholders(whereArg, placeholders, matchers)
		for _, replacer := range replacers {
			whereArg = replacer(whereArg)
		}
	} else {
		whereArg = ""
	}

	if havingArg != "" {
		preparePlaceholders(&placeholders, &fieldsInPhOrder, havingArg, matchers)
		havingArg = replaceFieldsWithPlaceholders(havingArg, placeholders, matchers)
		for _, replacer := range replacers {
			whereArg = replacer(whereArg)
		}
	}

	return whereArg, havingArg, fieldsInPhOrder, nil
}

func (e Environment) getFromRepo(fieldName string, lst ...interface{}) (interface{}, error) {
	whereDef, havingDef, fieldNames, err := e.getSelectComponents(lst...)
	if err != nil {
		return nil, errors.Wrap(err, "cannot build select")
	}

	values, err := getValuesByFieldNames(fieldNames, e[keyFieldsIndexes].(map[string]int), e[keyEvent])
	if err != nil {
		return nil, errors.Wrap(err, "cannot extract values for select")
	}

	if repoIf, ok := e[keyRepo]; ok {
		if repo, ok := repoIf.(DataRepo); ok {
			return repo.GetStringsFromData(fieldName, whereDef, havingDef, values)
		}
	}

	return []string{}, nil
}

func safeGetAsStr(lst []interface{}, idx int, defVal string) string {
	if len(lst) > idx {
		if str, ok := lst[idx].(string); ok {
			return str
		}
	}
	return defVal
}

func preparePlaceholders(placeholders *map[string]string, fieldsInPhOrder *[]string, s string, matchers map[string]*regexp.Regexp) {
	// preserve the fields order for unit tests
	correspFields := []string{}
	for fieldName := range matchers {
		correspFields = append(correspFields, fieldName)
	}
	sort.Strings(correspFields)

	for _, fieldName := range correspFields {
		matcher := matchers[fieldName]
		if !matcher.MatchString(s) { // skip if the field is not contained in s
			continue
		}
		placeholder, ok := "", false
		// create or pick up a placeholder for the field
		if placeholder, ok = (*placeholders)[fieldName]; !ok {
			placeholder = mockPlaceholder + strconv.Itoa(len(*placeholders)+1)
			(*placeholders)[fieldName] = placeholder
			*fieldsInPhOrder = append(*fieldsInPhOrder, fieldName)
		}
	}
}

func replaceFieldsWithPlaceholders(s string, placeholders map[string]string, matchers map[string]*regexp.Regexp) string {
	for fieldName, matcher := range matchers {
		if !matcher.MatchString(s) {
			continue
		}
		if ph, ok := placeholders[fieldName]; ok {
			// replace the field with the placeholder
			s = matcher.ReplaceAllString(s, "$1"+ph+"$3")
		}
	}
	return s
}

func getValuesByFieldNames(fields []string, fieldsIndex map[string]int, event interface{}) (values []interface{}, err error) {
	t := reflect.TypeOf(event)
	val := reflect.ValueOf(event)

	for _, field := range fields {
		if index, ok := fieldsIndex[field]; ok && index < t.NumField() {
			fval := val.Field(index)
			values = append(values, fval.Interface())
		} else {
			return nil, errors.New(fmt.Sprintf("should insert field %s, but has no such field in %s", field, t.Name()))
		}
	}

	return
}

func getScriptToSqlConverters() []stringConverter {
	return []stringConverter{
		func(s string) string {
			return strings.ReplaceAll(s, mockPlaceholder, postgresPlaceholder)
		},
		func(s string) string {
			return matchAnd.ReplaceAllString(s, `$1 and $3`)
		},
		func(s string) string {
			return matchOr.ReplaceAllString(s, `$1 or $3`)
		},
		func(s string) string {
			return matchDoubleEq.ReplaceAllString(s, `$1 = $3`)
		},
	}
}
