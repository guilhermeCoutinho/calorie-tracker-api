package contextextensions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/guilhermeCoutinho/api-studies/dal"
	"github.com/guilhermeCoutinho/api-studies/models"
	"github.com/guilhermeCoutinho/api-studies/server/http/wrapper"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const ctxKey = "ctxKey"

func ClaimsToCtx(ctx context.Context, claims *models.Claims) context.Context {
	rawBytes, _ := json.Marshal(claims)
	return context.WithValue(ctx, ctxKey, rawBytes)
}

func ClaimsFromCtx(ctx context.Context) (*models.Claims, error) {
	val, ok := ctx.Value(ctxKey).([]byte)
	if !ok {
		return nil, fmt.Errorf("failed to assert context")
	}
	var claim models.Claims
	err := json.Unmarshal(val, &claim)
	if err != nil {
		return nil, err
	}

	return &claim, nil
}

func LoggerFromCtx(ctx context.Context) logrus.FieldLogger {
	return ctx.Value(wrapper.LoggerCtxKey).(logrus.FieldLogger)
}

func GetQueryOptions(ctx context.Context) *dal.QueryOptions {
	params := ctx.Value(wrapper.URLParamsCtxKey).(url.Values)
	options := &dal.QueryOptions{}

	if val, ok := params["pagination"]; ok {
		options.Pagination = dal.DefaultPagination()
		json.Unmarshal([]byte(val[0]), options.Pagination)
	}

	if val, ok := params["sorting"]; ok {
		options.Sorting = dal.DefaultSorting()
		json.Unmarshal([]byte(val[0]), options.Sorting)
	}

	if val, ok := params["filtering"]; ok {
		options.Filtering = &dal.Filtering{}
		json.Unmarshal([]byte(val[0]), options.Filtering)
	}

	return options
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
