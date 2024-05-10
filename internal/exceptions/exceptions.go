package exceptions

import "fmt"

type BaseError struct {
	StatusCode int
	Message    string
}

type ClientSideError struct {
	BaseError
}

type ValidationError struct {
	ClientSideError
}

type AuthenticationError struct {
	ClientSideError
}

type ResourceNotFoundError struct {
	ClientSideError
}

type BadRequestError struct {
	ClientSideError
}

type StatusMethodNotAllowed struct {
	ClientSideError
}

type ConnectionError struct {
	ClientSideError
}

type StatusConflictError struct {
	ClientSideError
}

type ForbiddenError struct {
	ClientSideError
}

type ServerSideError struct {
	BaseError
}

type InternalServerError struct {
	ServerSideError
}

func (e BaseError) Error() string {
	return fmt.Sprintf("%v %v", e.StatusCode, e.Message)
}

func NewValidationError(msg string) ValidationError {
	return ValidationError{
		ClientSideError: ClientSideError{
			BaseError: BaseError{
				StatusCode: 422,
				Message:    fmt.Sprintf("%v | %v", HTTPStatusCodes[422], msg),
			},
		},
	}
}

func NewAuthenticationError(msg string) AuthenticationError {
	return AuthenticationError{
		ClientSideError: ClientSideError{
			BaseError: BaseError{
				StatusCode: 401,
				Message:    fmt.Sprintf("%v | %v", HTTPStatusCodes[401], msg),
			},
		},
	}
}

func NewResourceNotFoundError(msg string) ResourceNotFoundError {
	return ResourceNotFoundError{
		ClientSideError: ClientSideError{
			BaseError: BaseError{
				StatusCode: 404,
				Message:    fmt.Sprintf("%v | %v", HTTPStatusCodes[404], msg),
			},
		},
	}
}

func NewBadRequestError(msg string) BadRequestError {
	return BadRequestError{
		ClientSideError: ClientSideError{
			BaseError: BaseError{
				StatusCode: 400,
				Message:    fmt.Sprintf("%v | %v", HTTPStatusCodes[400], msg),
			},
		},
	}
}

func NewStatusMethodNotAllowed() StatusMethodNotAllowed {
	return StatusMethodNotAllowed{
		ClientSideError: ClientSideError{
			BaseError: BaseError{
				StatusCode: 405,
				Message:    fmt.Sprintf("%v", HTTPStatusCodes[405]),
			},
		},
	}
}

func NewStatusConflicError(msg string) StatusConflictError {
	return StatusConflictError{
		ClientSideError: ClientSideError{
			BaseError: BaseError{
				StatusCode: 409,
				Message:    fmt.Sprintf("%v | %v", HTTPStatusCodes[409], msg),
			},
		},
	}
}

func NewForbiddenError() ForbiddenError {
	return ForbiddenError{
		ClientSideError: ClientSideError{
			BaseError: BaseError{
				StatusCode: 403,
				Message:    fmt.Sprintf("%v", HTTPStatusCodes[403]),
			},
		},
	}
}

func NewInternalServerError() InternalServerError {
	return InternalServerError{
		ServerSideError: ServerSideError{
			BaseError: BaseError{
				StatusCode: 500,
				Message:    fmt.Sprintf("%v", HTTPStatusCodes[500]),
			},
		},
	}
}
