package floodcontroller

import (
	"context"
	"task/app/database"
)

type FloodController struct {
	N        int
	K        int
	database *database.Database
}

func NewFloodController(N int, K int, db *database.Database) *FloodController {
	return &FloodController{
		N:        N,
		K:        K,
		database: db,
	}
}

func (floodController *FloodController) Check(ctx context.Context, userID int64) (bool, error) {
	reqCount, err := floodController.database.CheckAmountRequestInN(ctx, userID, floodController.N)
	if err != nil {
		return false, err
	}

	if reqCount > floodController.K {
		return false, nil
	}

	err = floodController.database.AddUserReq(ctx, userID)

	return true, err
}
