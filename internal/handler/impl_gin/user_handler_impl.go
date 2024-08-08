package implgin

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"ticket_goroutine/internal/domain"
	"ticket_goroutine/internal/domain/dto"
	"ticket_goroutine/internal/handler"
	"ticket_goroutine/internal/usecase"
	"ticket_goroutine/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

type UserHandler struct {
	usecase usecase.UserUseCaseInterface
}

func NewUserHandler(usecase usecase.UserUseCaseInterface) (handler.UserHandlerInterface) {
	return UserHandler {
		usecase: usecase,
	}
}

func (h UserHandler) Save(c *gin.Context) {
	log.Trace().Msg("Entering user handler save")

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	var user domain.User

	log.Trace().Msg("Decoding json")
	err := c.ShouldBindJSON(&user)

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
	errValidate := utils.ValidateStruct(&user)
	sevenOrMore, number, upper := utils.VerifyPassword(user.Password)
	var errPassword []string

	if !sevenOrMore {
		errPassword = append(errPassword, "seven or more")
	}

	if !number {
		errPassword = append(errPassword, "number")
	}

	if !upper {
		errPassword = append(errPassword, "upper")
	}

	if errValidate != nil || len(errPassword) > 0 {
		log.Trace().Msg("Validation error")
		errors := make(map[string]string)
		log.Trace().Msg("User input error")
		if errValidate != nil {
			if _, ok := errValidate.(*validator.InvalidValidationError); ok {
				log.Trace().Msg("Error with validator")
				log.Error().Str("Error message: ", errValidate.Error())
				c.JSON(http.StatusInternalServerError, errValidate.Error())
				return
			}

			for _, err := range errValidate.(validator.ValidationErrors) {
				errors[err.Field()] = fmt.Sprintf("Validation failed on '%s' tag", err.Tag())
				log.Error().Msg(fmt.Sprintf("Validation failed on '%s' tag", err.Tag()))
			}
		}

		if len(errPassword) > 0 {
			errors["Password"] = fmt.Sprintf("Password is not valid on [%v] validation", errPassword)
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
        c.Writer.WriteHeader(http.StatusBadRequest)

		response := dto.GlobalResponse {
			StatusCode: http.StatusBadRequest,
			StatusDesc: http.StatusText(http.StatusBadRequest),
			Message: "Failed to save user because didn't pass the validation",
			RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
			ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
			Data: errors,
		}

        c.JSON(http.StatusBadRequest, response)
		return
	}

	log.Debug().
		Int("User ID: ", user.UserID).
		Str("Email: ", user.Email).
		Str("Name: ", user.Name).
		Str("Phone Number: ", user.PhoneNumber).
		Float64("Balance: ", user.Balance).
		Msg("Continuing user save process")

	done := make(chan struct{})
	log.Info().Msg("Channel created")

	go func() {
		defer close(done)
		log.Trace().Msg("Inside goroutine trying to save")

		savedUser, errSave := h.usecase.Save(ctx, user)
		
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
			} else if errors.Is(errSave, utils.ErrHash) {
				log.Error().Msg(fmt.Sprintf("Error hashing password with message: %s", err.Error()))

				c.Writer.WriteHeader(http.StatusInternalServerError)

				response := dto.GlobalResponse {
					StatusCode: http.StatusInternalServerError,
					StatusDesc: http.StatusText(http.StatusInternalServerError),
					Message: err.Error(),
					RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
					ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
					Data: "",
				}
				c.JSON(http.StatusInternalServerError, response)
				return
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
			Data: savedUser,
		}

		log.Info().Msg("User created successfully and returning json")
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

func (h UserHandler) FindById(c *gin.Context) {
	log.Trace().Msg("Entering user handler find by id")

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
		foundUser, errFound := h.usecase.FindById(ctx, id)

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

		response := dto.GlobalResponse {
			StatusCode: http.StatusOK,
			StatusDesc: http.StatusText(http.StatusOK),
			Message: "OK",
			RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
			ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
			Data: foundUser,
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusOK)
		log.Info().Msg("User fetched and returning json")
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

func (h UserHandler) GetAll(c *gin.Context) {
	log.Trace().Msg("Entering user get all handler")
	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
	defer cancel()

	done := make(chan struct{})
	log.Info().Msg("Channel created")
	go func() {
		log.Trace().Msg("Inside goroutine trying to get all data")
		defer close(done)
		allUsers, err := h.usecase.GetAll(ctx)

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
			Data: allUsers,
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusOK)
		log.Info().Msg("User fetched and returning json")
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

func (h UserHandler) Login(c *gin.Context) {
	var loginRequest dto.Login

	err := c.ShouldBindJSON(&loginRequest)
	
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
	
	foundUser, err := h.usecase.FindByEmail(c.Request.Context(), loginRequest.Email)

	if err != nil {
		log.Trace().Msg("Found user error")
		log.Error().Str("Error message: ", err.Error())
		response := dto.GlobalResponse {
			StatusCode: http.StatusNotFound,
			StatusDesc: http.StatusText(http.StatusNotFound),
			Message: err.Error(),
			RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
			ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
			Data: "",
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.WriteHeader(http.StatusNotFound)
		c.JSON(http.StatusNotFound, response)
		return
	}

	isSame := utils.CheckPasswordHash(loginRequest.Password, foundUser.Password)

	if !isSame {
		log.Trace().Msg("Password mismatch error")
		log.Error().Str("Error message: ", errors.New("wrong password").Error())
		response := dto.GlobalResponse {
			StatusCode: http.StatusUnauthorized,
			StatusDesc: http.StatusText(http.StatusUnauthorized),
			Message: errors.New("wrong password").Error(),
			RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
			ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
			Data: "",
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.JSON(http.StatusNotFound, response)
		return
	}
	
	claims := domain.Claims{
		Email: foundUser.Email,
		Name: foundUser.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := []byte(os.Getenv("SECRET"))
	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		log.Trace().Msg("Error creating signature")
		log.Error().Str("Error message: ", err.Error())
		response := dto.GlobalResponse {
			StatusCode: http.StatusInternalServerError,
			StatusDesc: http.StatusText(http.StatusInternalServerError),
			Message: err.Error(),
			RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
			ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
			Data: "",
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.JSON(http.StatusNotFound, response)
		return
	}

	// c.SetCookie("token", tokenString, time.Now().Add(time.Minute * 1).Second(), "", "", false, true)
	response := dto.GlobalResponse {
		StatusCode: http.StatusInternalServerError,
		StatusDesc: http.StatusText(http.StatusInternalServerError),
		Message: "OK",
		RequestCreated: time.Now().Format("2006-01-02 15:04:05"),
		ProcessTime: time.Duration(time.Since(time.Now()).Milliseconds()),
		Data: tokenString,
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	c.JSON(http.StatusOK, response)
}