package businessDay

import (
	"net/http"
)

type service struct{
	RecordHandler
}

func NewService(recordHandler RecordHandler) *service {
	return &service{
		recordHandler,
	}
}

func (s *service) GetCurrent(r *http.Request, request *GetCurrentRequest, response *GetCurrentResponse) error {
	return s.RecordHandler.GetCurrent(request, response)
}

func (s *service) GetAfter(r *http.Request, request *GetAfterRequest, response *GetAfterResponse) error {
	return s.RecordHandler.GetAfter(request, response)
}

func (s *service) GetBefore(r *http.Request, request *GetBeforeRequest, response *GetBeforeResponse) error {
	return s.RecordHandler.GetBefore(request, response)
}

func (s *service) GetSelected(r *http.Request, request *GetSelectedRequest, response *GetSelectedResponse) error {
	return s.RecordHandler.GetSelected(request, response)
}

func (s* service) UpdateShifts(r *http.Request, request *UpdateShiftsRequest, response *UpdateShiftsResponse) error {
	return s.RecordHandler.UpdateShifts(request, response)
}

func (s* service) AssignEmployeesToShift(r *http.Request, request *AssignEmployeesToShiftRequest, response *AssignEmployeesToShiftResponse) error {
	return s.RecordHandler.AssignEmployeesToShift(request, response)
}