package dao

import (
	"reflect"

	"github.com/gocraft/dbr"
)

func contains(arr []string, a string) bool {
	for _, b := range arr {
		if a == b {
			return true
		}
	}
	return false
}

func commonInsert(sr dbr.SessionRunner, tableName string, obj interface{}, omitField []string) (interface{}, error) {
	builder := sr.InsertInto(tableName)
	v := reflect.Indirect(reflect.ValueOf(obj))
	var idValue reflect.Value
	for i := 0; i < v.NumField(); i++ {
		f := v.Type().Field(i)
		dbTag := f.Tag.Get("db")
		if dbTag == "id" {
			idValue = v.Field(i)
		}
		if contains(omitField, dbTag) {
			continue
		}
		val := v.Field(i).Interface()
		builder.Pair(dbTag, val)
	}
	result, err := builder.Exec()
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	idValue.Set(reflect.ValueOf(id))
	return obj, nil
}
