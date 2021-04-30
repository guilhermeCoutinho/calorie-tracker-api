package meal

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/guilhermeCoutinho/api-studies/controller/contextextensions"
	"github.com/guilhermeCoutinho/api-studies/dal"
	"github.com/guilhermeCoutinho/api-studies/messages"
	"github.com/guilhermeCoutinho/api-studies/mocks"
	"github.com/guilhermeCoutinho/api-studies/models"
	"github.com/guilhermeCoutinho/api-studies/server/http/wrapper"
	"github.com/stretchr/testify/assert"
)

type Mocks struct {
	MockMealDAL         *mocks.MockMealDAL
	MockUserDAL         *mocks.MockUserDAL
	MockCaloroeProvider *mocks.MockProvider
}

func Mock(ctrl *gomock.Controller) *Mocks {
	mealMock := mocks.NewMockMealDAL(ctrl)
	userMock := mocks.NewMockUserDAL(ctrl)
	providerMock := mocks.NewMockProvider(ctrl)

	return &Mocks{
		MockMealDAL:         mealMock,
		MockUserDAL:         userMock,
		MockCaloroeProvider: providerMock,
	}
}

func ptrToInt(val int) *int {
	return &val
}

func ptrToUUID(val uuid.UUID) *uuid.UUID {
	return &val
}

func ptrToString(val string) *string {
	return &val
}

func TestCreateMeal(t *testing.T) {
	t.Parallel()
	gomock.NewController(t)

	constantUserID := uuid.New()

	testTable := map[string]struct {
		payload *messages.CreateMealPayload
		vars    *messages.RouteVars
		ctx     func() context.Context
		mocks   func(context.Context, *messages.CreateMealPayload, *Mocks)

		response *messages.CreateMealResponse
		err      *wrapper.HandlerError
	}{
		"success": {
			payload: &messages.CreateMealPayload{
				Meal:     ptrToString("hamburguer"),
				Date:     ptrToString("2020-Jan-01"),
				Time:     ptrToString("12h"),
				Calories: ptrToInt(100),
				UserID:   ptrToUUID(constantUserID),
			},
			vars: &messages.RouteVars{},
			ctx: func() context.Context {
				ctx := context.Background()
				return contextextensions.ClaimsToCtx(ctx, &models.Claims{
					UserID:      constantUserID,
					AccessLevel: models.RegulerUser,
				})
			},
			mocks: func(ctx context.Context, args *messages.CreateMealPayload, m *Mocks) {
				claims, err := contextextensions.ClaimsFromCtx(ctx)
				assert.Nil(t, err)

				m.MockMealDAL.EXPECT().UpsertMeal(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, meal *models.Meal) {
					assert.Equal(t, 100, meal.Calories)
					assert.Equal(t, claims.UserID, meal.UserID)
					assert.Equal(t, args.Meal, meal.Meal)
				})
			},
			err: nil,
		},

		"wrong_access_level": {
			payload: &messages.CreateMealPayload{
				Meal:     ptrToString("hamburguer"),
				Date:     ptrToString("2020-Jan-01"),
				Time:     ptrToString("12h"),
				Calories: ptrToInt(100),
				UserID:   ptrToUUID(uuid.New()),
			},
			vars: &messages.RouteVars{},
			ctx: func() context.Context {
				ctx := context.Background()
				return contextextensions.ClaimsToCtx(ctx, &models.Claims{
					UserID:      uuid.New(),
					AccessLevel: models.RegulerUser,
				})
			},
			mocks: func(ctx context.Context, args *messages.CreateMealPayload, m *Mocks) {},
			err:   &wrapper.HandlerError{Err: fmt.Errorf("wrong access level"), StatusCode: http.StatusUnauthorized},
		},

		"success_provider": {
			payload: &messages.CreateMealPayload{
				Meal: ptrToString("hamburguer"),
				Date: ptrToString("2020-Jan-01"),
				Time: ptrToString("12h"),
			},
			vars: nil,
			ctx: func() context.Context {
				ctx := context.Background()
				return contextextensions.ClaimsToCtx(ctx, &models.Claims{
					UserID: uuid.New(),
				})
			},
			mocks: func(ctx context.Context, args *messages.CreateMealPayload, m *Mocks) {
				claims, err := contextextensions.ClaimsFromCtx(ctx)
				assert.Nil(t, err)
				m.MockCaloroeProvider.EXPECT().GetCalories(args.Meal).Return(99, nil)
				m.MockMealDAL.EXPECT().UpsertMeal(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, meal *models.Meal) {
					assert.Equal(t, 99, meal.Calories)
					assert.Equal(t, claims.UserID, meal.UserID)
					assert.Equal(t, args.Meal, meal.Meal)
				})
			},
			err: nil,
		},
	}

	for name, table := range testTable {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mocked := Mock(ctrl)
			dal := &dal.DAL{
				Meal: mocked.MockMealDAL,
				User: mocked.MockUserDAL,
			}

			ctx := table.ctx()
			table.mocks(ctx, table.payload, mocked)
			mealController := NewMeal(dal, nil, mocked.MockCaloroeProvider)

			_, err := mealController.Post(ctx, table.payload, table.vars)
			assert.Equal(t, table.err, err)
		})
	}

}
