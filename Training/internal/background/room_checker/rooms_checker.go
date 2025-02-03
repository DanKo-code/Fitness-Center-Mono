package room_checker

import (
	"Training/internal/model"
	"Training/pkg/logger"
	"context"
	"log"
	"time"
)

type roomUseCase interface {
	UpdateRoomsList(ctx context.Context, roomMap *model.RoomMap) error
}

type RoomChecker struct {
	roomMap     *model.RoomMap
	roomUseCase roomUseCase
	broadcast   chan model.BroadcastMsg
}

func NewRoomChecker(roomMap *model.RoomMap, roomUseCase roomUseCase, broadcast chan model.BroadcastMsg) *RoomChecker {
	return &RoomChecker{
		roomMap,
		roomUseCase,
		broadcast,
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

func (rc RoomChecker) RunBroadcaster(ctx context.Context) {
	for {

		select {
		case <-ctx.Done():
			logger.Logger.Error("Ending broadcaster")
			return
		default:
			logger.Logger.Info("start get from chanal broadcast: ")
			msg := <-rc.broadcast
			logger.Logger.Info("end get from chanal broadcast: ", msg)

			for _, client := range rc.roomMap.Map[msg.RoomKey] {
				if client.Conn != msg.Client {
					client.Mutex.Lock()
					err := client.Conn.WriteJSON(msg.Message)
					client.Mutex.Unlock()

					if err != nil {
						err := client.Conn.Close()
						if err != nil {
							return
						}
						log.Fatal(err)
					}
				}
			}
		}
	}
}
