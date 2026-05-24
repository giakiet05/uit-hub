package service

import (
	"time"

	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/apperror"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/dto"
)

func (s *Service) GetRoomsAvailability(query dto.QueryRoomsAvailability) (dto.SuccessResponse, *apperror.AppError) {
	if isBlank(query.Date) || isBlank(query.Start) || isBlank(query.End) {
		return dto.SuccessResponse{}, newBadRequest("date, start, end are required")
	}
	if _, parseErr := time.Parse("2006-01-02", query.Date); parseErr != nil {
		return dto.SuccessResponse{}, newBadRequest("date is invalid")
	}

	startMin, err := parseTimeToMinutes(query.Start)
	if err != nil {
		return dto.SuccessResponse{}, newBadRequest("start is invalid")
	}
	endMin, err := parseTimeToMinutes(query.End)
	if err != nil {
		return dto.SuccessResponse{}, newBadRequest("end is invalid")
	}
	if startMin >= endMin {
		return dto.SuccessResponse{}, newBadRequest("start must be before end")
	}

	available := make([]dto.RoomItem, 0)
	rooms, repoErr := s.rooms.List()
	if repoErr != nil {
		return dto.SuccessResponse{}, newInternal("")
	}
	bookings, repoErr := s.roomBookings.List()
	if repoErr != nil {
		return dto.SuccessResponse{}, newInternal("")
	}
	for _, room := range rooms {
		if isRoomAvailable(room.ID, query.Date, startMin, endMin, bookings) {
			available = append(available, dto.RoomItem{ID: room.ID, Name: room.Name, Capacity: room.Capacity})
		}
	}

	result := dto.RoomAvailabilityData{Date: query.Date, Start: query.Start, End: query.End, Rooms: available}
	return success(result), nil
}
