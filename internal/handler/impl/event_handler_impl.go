package impl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"ticket_goroutine/internal/domain"
	"ticket_goroutine/internal/domain/dto"
	"ticket_goroutine/internal/handler"
	"ticket_goroutine/internal/usecase"
	"ticket_goroutine/utils"
	"time"

	"github.com/go-playground/validator/v10"
)

type EventHandler struct {
	usecase usecase.EventUseCaseInterface
}

func NewEventHandler(usecase usecase.EventUseCaseInterface) (handler.EventHandlerInterface) {
	return EventHandler {
		usecase: usecase,
	}
}

func (h EventHandler) Save(responseWriter http.ResponseWriter, request *http.Request) {
	var event domain.Event

	err := json.NewDecoder(request.Body).Decode(&event)

	if err != nil {
		response := dto.GlobalResponse {
			StatusCode: http.StatusBadRequest,
			StatusDesc: http.StatusText(http.StatusBadRequest),
			Message: err.Error(),
			RequestCreated: time.Now().Format("2024-02-01 17:14:05"),
			ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
			Data: dto.EventResponse {},
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.Header().Set("X-Content-Type-Options", "nosniff")
		responseWriter.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(responseWriter).Encode(response)
		return
	}

	errValidate := utils.ValidateStruct(&event)

	if errValidate != nil {
		if _, ok := errValidate.(*validator.InvalidValidationError); ok {
			http.Error(responseWriter, errValidate.Error(), http.StatusInternalServerError)
            return
		}

		errors := make(map[string]string)
		for _, err := range errValidate.(validator.ValidationErrors) {
            errors[err.Field()] = fmt.Sprintf("Validation failed on '%s' tag", err.Tag())
        }

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.Header().Set("X-Content-Type-Options", "nosniff")
        responseWriter.WriteHeader(http.StatusBadRequest)

		response := dto.GlobalResponse {
			StatusCode: http.StatusBadRequest,
			StatusDesc: http.StatusText(http.StatusBadRequest),
			Message: "Failed to save Event because didn't pass the validation",
			RequestCreated: time.Now().Format("2024-02-01 17:14:05"),
			ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
			Data: errors,
		}

        json.NewEncoder(responseWriter).Encode(response)
		return
	}
	savedEvent, errSave := h.usecase.Save(event)

	if errSave != nil {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.Header().Set("X-Content-Type-Options", "nosniff")
        responseWriter.WriteHeader(http.StatusBadRequest)

		response := dto.GlobalResponse {
			StatusCode: http.StatusBadRequest,
			StatusDesc: http.StatusText(http.StatusBadRequest),
			Message: errSave.Error(),
			RequestCreated: time.Now().Format("2024-02-01 17:14:05"),
			ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
			Data: dto.EventResponse {},
		}

		json.NewEncoder(responseWriter).Encode(response)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusCreated)

	response := dto.GlobalResponse {
		StatusCode: http.StatusCreated,
		StatusDesc: http.StatusText(http.StatusCreated),
		Message: "Created",
		RequestCreated: time.Now().Format("2024-02-01 17:14:05"),
		ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
		Data: savedEvent,
	}

	json.NewEncoder(responseWriter).Encode(response)
}

func (h EventHandler) FindById(responseWriter http.ResponseWriter, request *http.Request) {
	idString := request.PathValue("id")

	id, errConv := strconv.Atoi(idString)

	if errConv != nil {
		response := dto.GlobalResponse {
			StatusCode: http.StatusBadRequest,
			StatusDesc: http.StatusText(http.StatusBadRequest),
			Message: errConv.Error(),
			RequestCreated: time.Now().Format("2024-02-01 17:14:05"),
			ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
			Data: dto.EventResponse {},
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.Header().Set("X-Content-Type-Options", "nosniff")
        responseWriter.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(responseWriter).Encode(response)
		return
	}

	foundEvent, errFound := h.usecase.FindById(id)

	if errFound != nil {
		response := dto.GlobalResponse {
			StatusCode: http.StatusBadRequest,
			StatusDesc: http.StatusText(http.StatusBadRequest),
			Message: errFound.Error(),
			RequestCreated: time.Now().Format("2024-02-01 17:14:05"),
			ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
			Data: dto.EventResponse {},
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.Header().Set("X-Content-Type-Options", "nosniff")
        responseWriter.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(responseWriter).Encode(response)
		return
	}

	response := dto.GlobalResponse {
		StatusCode: http.StatusOK,
		StatusDesc: http.StatusText(http.StatusOK),
		Message: "OK",
		RequestCreated: time.Now().Format("2024-02-01 17:14:05"),
		ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
		Data: foundEvent,
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(responseWriter).Encode(response)
}

func (h EventHandler) GetAll(responseWriter http.ResponseWriter, request *http.Request) {
	response := dto.GlobalResponse {
		StatusCode: http.StatusOK,
		StatusDesc: http.StatusText(http.StatusOK),
		Message: "OK",
		RequestCreated: time.Now().Format("2024-02-01 17:14:05"),
		ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
		Data: h.usecase.GetAll(),
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(responseWriter).Encode(response)
}