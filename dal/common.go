package dal

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/spf13/viper"
)

type DAL struct {
	Meal MealDAL
	User UserDAL
}

func NewDAL(
	config *viper.Viper,
	db *pg.DB,
) *DAL {
	return &DAL{
		Meal: NewMeal(config, db),
		User: NewUser(config, db),
	}
}

type QueryOptions struct {
	Pagination *Pagination
	Sorting    *Sorting
	Filtering  *Filtering
}
type Pagination struct {
	Limit  int
	Offset int
}

type Sorting struct {
	SortBy string
	DESC   bool
}

type Filtering struct {
	Filter string
}

func addQueryOptions(query *orm.Query, options *QueryOptions) *orm.Query {
	if options == nil {
		return query
	}
	if options.Pagination != nil {
		query = query.Limit(options.Pagination.Limit).Offset(options.Pagination.Offset)
	}

	debug, _ := json.Marshal(options)
	fmt.Println(string(debug))

	if options.Sorting != nil {
		order := "ASC"
		if options.Sorting.DESC {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", options.Sorting.SortBy, order))
	}

	if options.Filtering != nil {
		q := options.Filtering.getFormattedQuery()
		fmt.Println("Query is ", q)
		query = query.Where(q)
	}
	return query
}

func (f *Filtering) getFormattedQuery() string {
	//(date eq '2016-05-01') AND ((number_of_calories gt 20) OR (number_of_calories lt 10)).
	replaceMap := map[string]string{
		"eq": "=",
		"ne": "!=",
		"gt": ">",
		"lt": "<",
	}
	queryParts := strings.Split(f.Filter, " ")

	finalQuery := []string{}
	for _, part := range queryParts {
		if replaceWith, ok := replaceMap[part]; ok {
			finalQuery = append(finalQuery, replaceWith)
			continue
		}
		finalQuery = append(finalQuery, part)
	}

	return strings.Join(finalQuery, " ")
}

func upsertAllFields(q *orm.Query, v interface{}) error {
	val := reflect.ValueOf(v).Elem()

	for i := 0; i < val.NumField(); i++ {
		fieldNameRaw := val.Type().Field(i).Tag.Get("sql")

		if fieldNameRaw == "" {
			continue
		}

		fieldName := strings.Split(fieldNameRaw, ",")[0]
		q = q.Set(fmt.Sprintf("%s = EXCLUDED.%s", fieldName, fieldName))
	}
	_, err := q.Insert()
	return err
}
