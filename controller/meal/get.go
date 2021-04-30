package meal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/guilhermeCoutinho/api-studies/controller/contextextensions"
	"github.com/guilhermeCoutinho/api-studies/messages"
	"github.com/guilhermeCoutinho/api-studies/models"
	"github.com/guilhermeCoutinho/api-studies/server/http/wrapper"
)

func (m *Meal) Get(ctx context.Context, args *struct{}, vars *messages.RouteVars) (*messages.GetMealsResponse, *wrapper.HandlerError) {
	raw, _ := json.Marshal(vars)
	fmt.Println(string(raw))

	claims, err := contextextensions.ClaimsFromCtx(ctx)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusInternalServerError}
	}

	accessFilter := func(al models.AccessLevel) bool { return al == models.Admin }
	userID, wrapperErr := contextextensions.ValidatetGetAccessLevel(claims, vars, accessFilter)
	if wrapperErr != nil {
		return nil, wrapperErr
	}

	mealIDForQuery := m.getMealIDFromURL(nil, vars)
	meals, err := m.dal.Meal.GetMeals(ctx, mealIDForQuery, userID, contextextensions.GetQueryOptions(ctx))
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusInternalServerError}
	}

	return &messages.GetMealsResponse{
		Meals: meals,
	}, nil
}
