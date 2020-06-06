package errs

type Errs interface {
	Add(ie Error)
	Error() []Error
}

// swagger:model errors
type Errors struct {
	Errors []Error `json:"errors"`
}

// swagger:model error
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *Errors) Add(ie Error) {
	e.Errors = append(e.Errors, ie)
}

func (e *Errors) Error() []Error {
	return e.Errors
}

func New(ie Error) Errs {
	return &Errors{[]Error{ie}}
}

func NewGeneralError(err error) Errs {
	return &Errors{[]Error{{"400", err.Error()}}}
}

var (

	//Date Error code and messages
	InvalidStartDate = Error{"101", "Invalid start date.Start date cannot be in the future date."}
	InvalidEndDate = Error{"102", "Invalid end date.End date cannot be in the future date."}
	InvalidStartEndDate   = Error{"103", `Invalid start date. End date cannot be in the past`}
	//Location Error code and message
	NasIdIsEmpty   = Error{"201", "Invalid nas Id"}
)
