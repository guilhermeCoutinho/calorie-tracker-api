package dal

import (
	"fmt"
	"strings"

	"github.com/go-pg/pg/v10/orm"
)

type SortingMode string

const (
	ASC  SortingMode = "ASC"
	DESC SortingMode = "DESC"
)

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
	Mode   SortingMode
}

type Filtering struct {
	Filter string
}

func DefaultSorting() *Sorting {
	return &Sorting{Mode: DESC}
}

func DefaultPagination() *Pagination {
	return &Pagination{Limit: 20}
}

func addQueryOptions(query *orm.Query, options *QueryOptions) (*orm.Query, error) {
	if options == nil {
		return query, nil
	}
	if options.Pagination != nil {
		query = query.Limit(options.Pagination.Limit).Offset(options.Pagination.Offset)
	}

	if options.Sorting != nil {
		query = query.Order(fmt.Sprintf("%s %s", options.Sorting.SortBy, options.Sorting.Mode))
	}

	if options.Filtering != nil {
		q, params, err := options.Filtering.getFormattedQuery()
		if err != nil {
			return nil, err
		}
		query = query.Where(q, params...)
	}
	return query, nil
}

func (f *Filtering) getFormattedQuery() (string, []interface{}, error) {
	f.Filter = strings.ReplaceAll(f.Filter, "(", " ( ")
	f.Filter = strings.ReplaceAll(f.Filter, ")", " ) ")
	args := strings.Fields(f.Filter)

	validator := NewValidator()

	for i, arg := range args {
		if arg == "(" || arg == ")" {
			continue
		} else if arg == "AND" || arg == "OR" {
			continue
		}

		replaceWith, success := validator.Validate(arg)
		if !success {
			return "", nil, fmt.Errorf("failed to parse string %v, %v, %v", args, i, arg)
		}
		args[i] = replaceWith
	}

	return strings.Join(args, " "), validator.params, nil
}
