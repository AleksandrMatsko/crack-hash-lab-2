package contracts

import "github.com/google/uuid"

type TaskRequest struct {
	StartIndex uint64    `json:"startIndex"`
	PartCount  uint64    `json:"partCount"`
	Alphabet   string    `json:"alphabet"`
	MaxLength  int       `json:"maxLength"`
	ToCrack    string    `json:"toCrack"`
	RequestID  uuid.UUID `json:"requestId"`
}

func (req TaskRequest) Validate() error {
	if req.Alphabet == "" {
		return ErrEmptyAlphabet
	}
	if req.MaxLength < 0 {
		return ErrNegativeMaxLength
	}
	if req.ToCrack == "" {
		return ErrEmptyHashToCrack
	}
	return nil
}

type TaskResultRequest struct {
	StartIndex uint64    `json:"startIndex"`
	RequestID  uuid.UUID `json:"requestId"`
	Cracks     []string  `json:"cracks"`
}
