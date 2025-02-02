package room_checker

import (
	"Training/internal/model"
	"Training/pkg/logger"
	"context"
	"time"
)

type roomUseCase interface {
	UpdateRoomsList(ctx context.Context, roomMap *model.RoomMap) error
}

type RoomChecker struct {
	roomMap     *model.RoomMap
	roomUseCase roomUseCase
}

func NewRoomChecker(roomMap *model.RoomMap, roomUseCase roomUseCase) *RoomChecker {
	return &RoomChecker{
		roomMap,
		roomUseCase,
	}
}

func (rc RoomChecker) Run(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ticker.C:
			err := rc.roomUseCase.UpdateRoomsList(ctx, rc.roomMap)
			if err != nil {
				return err
			}
			logger.Logger.Info("Successfully UpdateRoomsList")

		case <-ctx.Done():
			logger.Logger.Info("Stopping RoomChecker...")
			return nil
		}
	}

}
