package response

import (
	"encoding/json"
	"log"
	"net/http"
)

// ErrorResponse representa uma resposta de erro padronizada
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}

// SuccessResponse representa uma resposta de sucesso padronizada
type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
	Total   int         `json:"total"`
}

type SucessResponseWithoutTotal struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

// SendError envia uma resposta de erro padronizada
func SendError(w http.ResponseWriter, statusCode int, err error, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	var errorMsg string
	if err != nil {
		errorMsg = err.Error()
	} else {
		errorMsg = message
	}

	errorResp := ErrorResponse{
		Error:   errorMsg,
		Message: message,
		Code:    statusCode,
	}

	if err := json.NewEncoder(w).Encode(errorResp); err != nil {
		log.Printf("Erro ao codificar resposta de erro: %v", err)
	}

	// Log do erro
	if err != nil {
		log.Printf("Erro %d: %s - %s", statusCode, message, err.Error())
	} else {
		log.Printf("Erro %d: %s", statusCode, message)
	}
}

func SendSuccess(w http.ResponseWriter, statusCode int, data any, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Erro ao codificar resposta de sucesso: %v", err)
		SendError(w, http.StatusInternalServerError, err, "Erro interno do servidor")
		return
	}
}

func SendSuccessWithTotal(w http.ResponseWriter, statusCode int, data interface{}, total int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	successResp := SuccessResponse{
		Data:    data,
		Message: message,
		Total:   total,
	}

	if err := json.NewEncoder(w).Encode(successResp); err != nil {
		log.Printf("Erro ao codificar resposta de sucesso: %v", err)
		SendError(w, http.StatusInternalServerError, err, "Erro interno do servidor")
		return
	}
}

func SendJSONResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Erro ao codificar resposta de sucesso: %v", err)

	}
}
