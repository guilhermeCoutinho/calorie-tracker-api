package dal

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-pg/pg/v10/orm"
)

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
