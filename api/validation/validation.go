package validation

type Validator interface {
	Validate() (*ValidationResult, error)
}

type ValidationResult struct {
	Errors map[string]string `json:"errors"`
	OK     bool              `json:"ok"`
}

func (result *ValidationResult) Error(name, msg string) {
	result.Errors[name] = msg
	result.OK = false
}

func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		make(map[string]string),
		true,
	}
}
