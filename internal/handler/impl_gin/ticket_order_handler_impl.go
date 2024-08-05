package implgin

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"ticket_goroutine/internal/domain"
	"ticket_goroutine/internal/domain/dto"
	"ticket_goroutine/internal/handler"
	"ticket_goroutine/internal/usecase"
	"ticket_goroutine/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

type TicketOrderHandler struct {
	usecaseTicketOrder usecase.TicketOrderUseCaseInterface
}

func NewTicketOrderHandler(usecaseTicketOrder usecase.TicketOrderUseCaseInterface) (handler.TicketOrderHandlerInterface) {
	return TicketOrderHandler {
		usecaseTicketOrder: usecaseTicketOrder,
	}
}

func (h TicketOrderHandler) Save(c *gin.Context) {
	log.Trace().Msg("Entering ticket order handler save")

	ctx, cancel := context.WithTimeout(context.Background(), 7 * time.Second)
	defer cancel()

	var ticketOrder domain.TicketOrder

	log.Trace().Msg("Decoding json")
	err := c.ShouldBindJSON(&ticketOrder)

	if err != nil {
		log.Trace().Msg("JSON decode error")
		log.Error().Str("Error message: ", err.Error())
		response := dto.GlobalResponse {
			StatusCode: http.StatusBadRequest,
			StatusDesc: http.StatusText(http.StatusBadRequest),
			Message: err.Error(),
			RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
			ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
			Data: "",
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	log.Trace().Msg("Validating user input")
	errValidate := utils.ValidateStruct(&ticketOrder)

	if errValidate != nil {
		log.Trace().Msg("Validation error")
		if _, ok := errValidate.(*validator.InvalidValidationError); ok {
			log.Trace().Msg("Error with validator")
			log.Error().Str("Error message: ", errValidate.Error())
			c.JSON(http.StatusInternalServerError, errValidate.Error())
            return
		}

		log.Trace().Msg("User input error")
		errors := make(map[string]string)
		for _, err := range errValidate.(validator.ValidationErrors) {
            errors[err.Field()] = fmt.Sprintf("Validation failed on '%s' tag", err.Tag())
			log.Error().Msg(fmt.Sprintf("Validation failed on '%s' tag", err.Tag()))
        }

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
        c.Writer.WriteHeader(http.StatusBadRequest)

		response := dto.GlobalResponse {
			StatusCode: http.StatusBadRequest,
			StatusDesc: http.StatusText(http.StatusBadRequest),
			Message: "Failed to save ticket order because didn't pass the validation",
			RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
			ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
			Data: errors,
		}

        c.JSON(http.StatusBadRequest, response)
		return
	}

	log.Debug().
		Int("Order ID: ", ticketOrder.OrderID).
		Int("Ticket ID: ", ticketOrder.TicketID).
		Int("User ID: ", ticketOrder.UserID).
		Int("Amount: ", ticketOrder.Amount).
		Float64("Total Price: ", ticketOrder.TotalPrice).
		Msg("Continuing ticket order save process")

	done := make(chan struct{})
	log.Info().Msg("Channel created")

	go func() {
		defer close(done)
		log.Trace().Msg("Inside goroutine trying to save")
		savedTicketOrder, errSave := h.usecaseTicketOrder.Save(ctx, ticketOrder)

		if errSave != nil {
			log.Trace().Msg("Checking error cause")
			c.Writer.Header().Set("Content-Type", "application/json")
			c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

			var response dto.GlobalResponse

			if errors.Is(errSave, context.DeadlineExceeded) || errors.Is(errSave, context.Canceled) {
				log.Trace().Msg("Timeout error")
				log.Error().Str("Error message: ", errSave.Error())
				c.Writer.WriteHeader(http.StatusRequestTimeout)

				response = dto.GlobalResponse {
					StatusCode: http.StatusRequestTimeout,
					StatusDesc: http.StatusText(http.StatusRequestTimeout),
					Message: "Request Timed Out",
					RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
					ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
					Data: "",
				}
				c.JSON(http.StatusRequestTimeout, response)
			} else if errSave.Error() == "not enough stock" {
				log.Trace().Msg("Save error")
				log.Error().Str("Error message: ", errSave.Error())
				c.Writer.WriteHeader(http.StatusUnprocessableEntity)

				response = dto.GlobalResponse {
					StatusCode: http.StatusUnprocessableEntity,
					StatusDesc: http.StatusText(http.StatusUnprocessableEntity),
					Message: errSave.Error(),
					RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
					ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
					Data: "",
				}
				c.JSON(http.StatusUnprocessableEntity, response)
			} else if errSave.Error() == "not enough balance" {
				log.Trace().Msg("Save error")
				log.Error().Str("Error message: ", errSave.Error())
				c.Writer.WriteHeader(http.StatusPaymentRequired)

				response = dto.GlobalResponse {
					StatusCode: http.StatusPaymentRequired,
					StatusDesc: http.StatusText(http.StatusPaymentRequired),
					Message: errSave.Error(),
					RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
					ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
					Data: "",
				}
				c.JSON(http.StatusPaymentRequired, response)
			} else {
				log.Trace().Msg("Save error")
				log.Error().Str("Error message: ", errSave.Error())
				c.Writer.WriteHeader(http.StatusBadRequest)

				response = dto.GlobalResponse {
					StatusCode: http.StatusBadRequest,
					StatusDesc: http.StatusText(http.StatusBadRequest),
					Message: errSave.Error(),
					RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
					ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
					Data: "",
				}
				c.JSON(http.StatusBadRequest, response)
			}
			return
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusCreated)

		response := dto.GlobalResponse {
			StatusCode: http.StatusCreated,
			StatusDesc: http.StatusText(http.StatusCreated),
			Message: "Created",
			RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
			ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
			Data: savedTicketOrder,
		}

		log.Info().Msg("Ticket order created successfully and returning json")
		c.JSON(http.StatusCreated, response)
	}()

	select {
	case <- ctx.Done():
		log.Trace().Msg("Request timeout channel")
		c.Writer.WriteHeader(http.StatusRequestTimeout)

		response := dto.GlobalResponse {
			StatusCode: http.StatusRequestTimeout,
			StatusDesc: http.StatusText(http.StatusRequestTimeout),
			Message: "Request Timed Out",
			RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
			ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
			Data: "",
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.JSON(http.StatusRequestTimeout, response)
	case <- done:
		log.Trace().Msg("Request completed")
	}
}

func (h TicketOrderHandler) FindById(c *gin.Context) {
	log.Trace().Msg("Entering ticket order handler find by id")

	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
	defer cancel()

	idString := c.Param("id")
	log.Debug().Str("Received Id is: ", idString)

	log.Trace().Msg("Trying to convert id in string to int")
	id, errConv := strconv.Atoi(idString)

	if errConv != nil {
		log.Trace().Msg("Error happens when converting id string to int")
		log.Error().Str("Error message: ", errConv.Error())
		response := dto.GlobalResponse {
			StatusCode: http.StatusBadRequest,
			StatusDesc: http.StatusText(http.StatusBadRequest),
			Message: errConv.Error(),
			RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
			ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
			Data: "",
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
        c.Writer.WriteHeader(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	done := make(chan struct{})
	log.Info().Msg("Channel created")
	go func() {
		log.Trace().Msg("Inside goroutine trying to fetch data by id")
		defer close(done)
		foundTicketOrder, errFound := h.usecaseTicketOrder.FindById(ctx, id)

		if errFound != nil {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

			var response dto.GlobalResponse

			log.Trace().Msg("Checking error cause")
			if errors.Is(errFound, context.DeadlineExceeded) || errors.Is(errFound, context.Canceled) {
				log.Trace().Msg("Timeout error")
				log.Error().Str("Error message: ", errFound.Error())

				c.Writer.WriteHeader(http.StatusRequestTimeout)

				response = dto.GlobalResponse {
					StatusCode: http.StatusRequestTimeout,
					StatusDesc: http.StatusText(http.StatusRequestTimeout),
					Message: "Request Timed Out",
					RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
					ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
					Data: "",
				}
				c.JSON(http.StatusRequestTimeout, response)
			} else {
				log.Trace().Msg("Fetch error")
				log.Error().Str("Error message: ", errFound.Error())
				c.Writer.WriteHeader(http.StatusBadRequest)

				response = dto.GlobalResponse {
					StatusCode: http.StatusBadRequest,
					StatusDesc: http.StatusText(http.StatusBadRequest),
					Message: errFound.Error(),
					RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
					ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
					Data: "",
				}
				c.JSON(http.StatusBadRequest, response)
			}
			return
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusOK)

		response := dto.GlobalResponse {
			StatusCode: http.StatusOK,
			StatusDesc: http.StatusText(http.StatusOK),
			Message: "OK",
			RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
			ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
			Data: foundTicketOrder,
		}
		log.Info().Msg("Ticket order fetched successfully and returning json")
		c.JSON(http.StatusOK, response)
	}()

	select {
	case <- ctx.Done():
		log.Trace().Msg("Request timeout channel")
		c.Writer.WriteHeader(http.StatusRequestTimeout)

		response := dto.GlobalResponse {
			StatusCode: http.StatusRequestTimeout,
			StatusDesc: http.StatusText(http.StatusRequestTimeout),
			Message: "Request Timed Out",
			RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
			ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
			Data: "",
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.JSON(http.StatusRequestTimeout, response)
	case <- done:
		log.Info().Msg("Request completed")
	}
}

func (h TicketOrderHandler) GetAll(c *gin.Context) {
	log.Trace().Msg("Entering ticket order get all handler")
	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
	defer cancel()

	done := make(chan struct{})
	log.Info().Msg("Channel created")
	go func() {
		log.Trace().Msg("Inside goroutine trying to get all data")
		defer close(done)
		allTicketOrders, err := h.usecaseTicketOrder.GetAll(ctx)

		if err != nil {
			log.Trace().Msg("Error happens when trying to fetch data")
			log.Error().Str("Error message: ", err.Error())
			c.Writer.Header().Set("Content-Type", "application/json")
			c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

			c.Writer.WriteHeader(http.StatusRequestTimeout)

			response := dto.GlobalResponse {
				StatusCode: http.StatusRequestTimeout,
				StatusDesc: http.StatusText(http.StatusRequestTimeout),
				Message: "Request Timed Out",
				RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
				ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
				Data: "",
			}
			c.JSON(http.StatusRequestTimeout, response)
			return
		}

		response := dto.GlobalResponse {
			StatusCode: http.StatusOK,
			StatusDesc: http.StatusText(http.StatusOK),
			Message: "OK",
			RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
			ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
			Data: allTicketOrders,
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusOK)
		log.Info().Msg("Event fetched and returning json")
		c.JSON(http.StatusOK, response)
	}()

	select {
	case <- ctx.Done():
		log.Trace().Msg("Request timeout channel")
		c.Writer.WriteHeader(http.StatusRequestTimeout)

		response := dto.GlobalResponse {
			StatusCode: http.StatusRequestTimeout,
			StatusDesc: http.StatusText(http.StatusRequestTimeout),
			Message: "Request Timed Out",
			RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
			ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
			Data: "",
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.JSON(http.StatusRequestTimeout, response)
	case <- done:
		log.Info().Msg("Request completed")
	}
}