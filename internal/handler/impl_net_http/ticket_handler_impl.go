package implnethttp

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"net/http"
// 	"strconv"
// 	"ticket_goroutine/internal/domain"
// 	"ticket_goroutine/internal/domain/dto"
// 	"ticket_goroutine/internal/handler"
// 	"ticket_goroutine/internal/usecase"
// 	"ticket_goroutine/utils"
// 	"time"

// 	"github.com/go-playground/validator/v10"
// 	"github.com/rs/zerolog/log"
// )

// type TicketHandler struct {
// 	usecaseTicket usecase.TicketUseCaseInterface
// }

// func NewTicketHandler(usecaseTicket usecase.TicketUseCaseInterface) (handler.TicketHandlerInterface) {
// 	return TicketHandler {
// 		usecaseTicket: usecaseTicket,
// 	}
// }

// func (h TicketHandler) Save(responseWriter http.ResponseWriter, request *http.Request) {
// 	log.Trace().Msg("Entering ticket handler save")

// 	ctx, cancel := context.WithTimeout(request.Context(), 2 * time.Second)
// 	defer cancel()

// 	var ticket domain.Ticket

// 	log.Trace().Msg("Decoding json")
// 	err := json.NewDecoder(request.Body).Decode(&ticket)

// 	if err != nil {
// 		log.Trace().Msg("JSON decode error")
// 		log.Error().Str("Error message: ", err.Error())
// 		response := dto.GlobalResponse {
// 			StatusCode: http.StatusBadRequest,
// 			StatusDesc: http.StatusText(http.StatusBadRequest),
// 			Message: err.Error(),
// 			RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
// 			ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
// 			Data: "",
// 		}

// 		responseWriter.Header().Set("Content-Type", "application/json")
// 		responseWriter.Header().Set("X-Content-Type-Options", "nosniff")
// 		responseWriter.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(responseWriter).Encode(response)
// 		return
// 	}

// 	log.Trace().Msg("Validating user input")
// 	errValidate := utils.ValidateStruct(&ticket)

// 	if errValidate != nil {
// 		log.Trace().Msg("Validation error")
// 		if _, ok := errValidate.(*validator.InvalidValidationError); ok {
// 			log.Trace().Msg("Error with validator")
// 			log.Error().Str("Error message: ", errValidate.Error())
// 			http.Error(responseWriter, errValidate.Error(), http.StatusInternalServerError)
//             return
// 		}

// 		log.Trace().Msg("User input error")
// 		errors := make(map[string]string)
// 		for _, err := range errValidate.(validator.ValidationErrors) {
//             errors[err.Field()] = fmt.Sprintf("Validation failed on '%s' tag", err.Tag())
// 			log.Error().Msg(fmt.Sprintf("Validation failed on '%s' tag", err.Tag()))
//         }

// 		responseWriter.Header().Set("Content-Type", "application/json")
// 		responseWriter.Header().Set("X-Content-Type-Options", "nosniff")
//         responseWriter.WriteHeader(http.StatusBadRequest)

// 		response := dto.GlobalResponse {
// 			StatusCode: http.StatusBadRequest,
// 			StatusDesc: http.StatusText(http.StatusBadRequest),
// 			Message: "Failed to save ticket because didn't pass the validation",
// 			RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
// 			ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
// 			Data: errors,
// 		}

//         json.NewEncoder(responseWriter).Encode(response)
// 		return
// 	}

// 	log.Debug().
// 		Int("Ticket ID: ", ticket.TicketID).
// 		Int("Event ID: ", ticket.EventID).
// 		Str("Name: ", ticket.Name).
// 		Float64("Price: ", ticket.Price).
// 		Int("Stock: ", ticket.Stock).
// 		Str("Type: ", ticket.Type).
// 		Msg("Continuing ticket save process")

// 	done := make(chan struct{})
// 	log.Info().Msg("Channel created")

// 	go func() {
// 		defer close(done)
// 		log.Trace().Msg("Inside goroutine trying to save")
// 		savedTicket, errSave := h.usecaseTicket.Save(ctx, ticket)

// 		if errSave != nil {
// 			log.Trace().Msg("Checking error cause")
// 			responseWriter.Header().Set("Content-Type", "application/json")
// 			responseWriter.Header().Set("X-Content-Type-Options", "nosniff")

// 			var response dto.GlobalResponse

// 			if errors.Is(errSave, context.DeadlineExceeded) || errors.Is(errSave, context.Canceled) {
// 				log.Trace().Msg("Timeout error")
// 				log.Error().Str("Error message: ", errSave.Error())
// 				responseWriter.WriteHeader(http.StatusRequestTimeout)

// 				response = dto.GlobalResponse {
// 					StatusCode: http.StatusRequestTimeout,
// 					StatusDesc: http.StatusText(http.StatusRequestTimeout),
// 					Message: "Request Timed Out",
// 					RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
// 					ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
// 					Data: "",
// 				}
// 			} else {
// 				log.Trace().Msg("Save error")
// 				log.Error().Str("Error message: ", errSave.Error())
// 				responseWriter.WriteHeader(http.StatusBadRequest)

// 				response = dto.GlobalResponse {
// 					StatusCode: http.StatusBadRequest,
// 					StatusDesc: http.StatusText(http.StatusBadRequest),
// 					Message: errSave.Error(),
// 					RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
// 					ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
// 					Data: "",
// 				}
// 			}

// 			json.NewEncoder(responseWriter).Encode(response)
// 			return
// 		}

// 		responseWriter.Header().Set("Content-Type", "application/json")
// 		responseWriter.WriteHeader(http.StatusCreated)

// 		response := dto.GlobalResponse {
// 			StatusCode: http.StatusCreated,
// 			StatusDesc: http.StatusText(http.StatusCreated),
// 			Message: "Created",
// 			RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
// 			ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
// 			Data: savedTicket,
// 		}

// 		log.Info().Msg("Ticket created successfully and returning json")
// 		json.NewEncoder(responseWriter).Encode(response)
// 	}()

// 	select {
// 	case <- ctx.Done():
// 		log.Trace().Msg("Request timeout channel")
// 		responseWriter.WriteHeader(http.StatusRequestTimeout)

// 		response := dto.GlobalResponse {
// 			StatusCode: http.StatusRequestTimeout,
// 			StatusDesc: http.StatusText(http.StatusRequestTimeout),
// 			Message: "Request Timed Out",
// 			RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
// 			ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
// 			Data: "",
// 		}

// 		responseWriter.Header().Set("Content-Type", "application/json")
// 		responseWriter.Header().Set("X-Content-Type-Options", "nosniff")
// 		json.NewEncoder(responseWriter).Encode(response)
// 	case <- done:
// 		log.Trace().Msg("Request completed")
// 	}
// }

// func (h TicketHandler) FindById(responseWriter http.ResponseWriter, request *http.Request) {
// 	log.Trace().Msg("Entering ticket handler find by id")

// 	ctx, cancel := context.WithTimeout(request.Context(), 2 * time.Second)
// 	defer cancel()

// 	idString := request.PathValue("id")
// 	log.Debug().Str("Received Id is: ", idString)

// 	log.Trace().Msg("Trying to convert id in string to int")
// 	id, errConv := strconv.Atoi(idString)

// 	if errConv != nil {
// 		log.Trace().Msg("Error happens when converting id string to int")
// 		log.Error().Str("Error message: ", errConv.Error())
// 		response := dto.GlobalResponse {
// 			StatusCode: http.StatusBadRequest,
// 			StatusDesc: http.StatusText(http.StatusBadRequest),
// 			Message: errConv.Error(),
// 			RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
// 			ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
// 			Data: "",
// 		}

// 		responseWriter.Header().Set("Content-Type", "application/json")
// 		responseWriter.Header().Set("X-Content-Type-Options", "nosniff")
//         responseWriter.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(responseWriter).Encode(response)
// 		return
// 	}

// 	done := make(chan struct{})
// 	log.Info().Msg("Channel created")
// 	go func() {
// 		log.Trace().Msg("Inside goroutine trying to fetch data by id")
// 		defer close(done)
// 		foundTicket, errFound := h.usecaseTicket.FindById(ctx, id)

// 		if errFound != nil {
// 			responseWriter.Header().Set("Content-Type", "application/json")
// 			responseWriter.Header().Set("X-Content-Type-Options", "nosniff")

// 			var response dto.GlobalResponse

// 			log.Trace().Msg("Checking error cause")
// 			if errors.Is(errFound, context.DeadlineExceeded) || errors.Is(errFound, context.Canceled) {
// 				log.Trace().Msg("Timeout error")
// 				log.Error().Str("Error message: ", errFound.Error())

// 				responseWriter.WriteHeader(http.StatusRequestTimeout)

// 				response = dto.GlobalResponse {
// 					StatusCode: http.StatusRequestTimeout,
// 					StatusDesc: http.StatusText(http.StatusRequestTimeout),
// 					Message: "Request Timed Out",
// 					RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
// 					ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
// 					Data: "",
// 				}
// 			} else {
// 				log.Trace().Msg("Fetch error")
// 				log.Error().Str("Error message: ", errFound.Error())
// 				responseWriter.WriteHeader(http.StatusBadRequest)

// 				response = dto.GlobalResponse {
// 					StatusCode: http.StatusBadRequest,
// 					StatusDesc: http.StatusText(http.StatusBadRequest),
// 					Message: errFound.Error(),
// 					RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
// 					ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
// 					Data: "",
// 				}
// 			}
// 			json.NewEncoder(responseWriter).Encode(response)
// 			return
// 		}

// 		responseWriter.Header().Set("Content-Type", "application/json")
// 		responseWriter.WriteHeader(http.StatusOK)

// 		response := dto.GlobalResponse {
// 			StatusCode: http.StatusOK,
// 			StatusDesc: http.StatusText(http.StatusOK),
// 			Message: "OK",
// 			RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
// 			ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
// 			Data: foundTicket,
// 		}
// 		log.Info().Msg("Ticket fetched successfully and returning json")
// 		json.NewEncoder(responseWriter).Encode(response)
// 	}()

// 	select {
// 	case <- ctx.Done():
// 		log.Trace().Msg("Request timeout channel")
// 		responseWriter.WriteHeader(http.StatusRequestTimeout)

// 		response := dto.GlobalResponse {
// 			StatusCode: http.StatusRequestTimeout,
// 			StatusDesc: http.StatusText(http.StatusRequestTimeout),
// 			Message: "Request Timed Out",
// 			RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
// 			ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
// 			Data: "",
// 		}

// 		responseWriter.Header().Set("Content-Type", "application/json")
// 		responseWriter.Header().Set("X-Content-Type-Options", "nosniff")
// 		json.NewEncoder(responseWriter).Encode(response)
// 	case <- done:
// 		log.Info().Msg("Request completed")
// 	}
// }

// func (h TicketHandler) GetAll(responseWriter http.ResponseWriter, request *http.Request) {
// 	log.Trace().Msg("Entering ticket get all handler")
// 	ctx, cancel := context.WithTimeout(request.Context(), 30 * time.Second)
// 	defer cancel()

// 	done := make(chan struct{})
// 	log.Info().Msg("Channel created")
// 	go func() {
// 		log.Trace().Msg("Inside goroutine trying to get all data")
// 		defer close(done)
// 		allTickets, err := h.usecaseTicket.GetAll(ctx)

// 		if err != nil {
// 			log.Trace().Msg("Error happens when trying to fetch data")
// 			log.Error().Str("Error message: ", err.Error())
// 			responseWriter.Header().Set("Content-Type", "application/json")
// 			responseWriter.Header().Set("X-Content-Type-Options", "nosniff")

// 			responseWriter.WriteHeader(http.StatusRequestTimeout)

// 			response := dto.GlobalResponse {
// 				StatusCode: http.StatusRequestTimeout,
// 				StatusDesc: http.StatusText(http.StatusRequestTimeout),
// 				Message: "Request Timed Out",
// 				RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
// 				ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
// 				Data: "",
// 			}
// 			json.NewEncoder(responseWriter).Encode(response)
// 			return
// 		}

// 		response := dto.GlobalResponse {
// 			StatusCode: http.StatusOK,
// 			StatusDesc: http.StatusText(http.StatusOK),
// 			Message: "OK",
// 			RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
// 			ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
// 			Data: allTickets,
// 		}

// 		responseWriter.Header().Set("Content-Type", "application/json")
// 		responseWriter.WriteHeader(http.StatusOK)
// 		log.Info().Msg("Event fetched and returning json")
// 		json.NewEncoder(responseWriter).Encode(response)
// 	}()

// 	select {
// 	case <- ctx.Done():
// 		log.Trace().Msg("Request timeout channel")
// 		responseWriter.WriteHeader(http.StatusRequestTimeout)

// 		response := dto.GlobalResponse {
// 			StatusCode: http.StatusRequestTimeout,
// 			StatusDesc: http.StatusText(http.StatusRequestTimeout),
// 			Message: "Request Timed Out",
// 			RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
// 			ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
// 			Data: "",
// 		}

// 		responseWriter.Header().Set("Content-Type", "application/json")
// 		responseWriter.Header().Set("X-Content-Type-Options", "nosniff")
// 		json.NewEncoder(responseWriter).Encode(response)
// 	case <- done:
// 		log.Info().Msg("Request completed")
// 	}
// }