package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/guilhermeCoutinho/api-studies/dal"
	"github.com/guilhermeCoutinho/api-studies/server/http/wrapper"
	"github.com/sirupsen/logrus"
)

const ctxKey = "ctxKey"

func ClaimsToCtx(ctx context.Context, claims *Claims) context.Context {
	rawBytes, _ := json.Marshal(claims)
	return context.WithValue(ctx, ctxKey, rawBytes)
}

func ClaimsFromCtx(ctx context.Context) (*Claims, error) {
	val, ok := ctx.Value(ctxKey).([]byte)
	if !ok {
		return nil, fmt.Errorf("failed to assert context")
	}
	var claim Claims
	err := json.Unmarshal(val, &claim)
	if err != nil {
		return nil, err
	}

	return &claim, nil
}

func LoggerFromCtx(ctx context.Context) logrus.FieldLogger {
	return ctx.Value(wrapper.LoggerCtxKey).(logrus.FieldLogger)
}

func getQueryOptions(ctx context.Context) *dal.QueryOptions {
	params := ctx.Value(wrapper.URLParamsCtxKey).(url.Values)
	options := &dal.QueryOptions{}

	if val, ok := params["pagination"]; ok {
		options.Pagination = &dal.Pagination{}
		json.Unmarshal([]byte(val[0]), options.Pagination)
	}

	if val, ok := params["sorting"]; ok {
		options.Sorting = &dal.Sorting{}
		json.Unmarshal([]byte(val[0]), options.Sorting)
	}

	if val, ok := params["filtering"]; ok {
		options.Filtering = &dal.Filtering{
			Filter: val[0],
		}
	}

	return options
}
