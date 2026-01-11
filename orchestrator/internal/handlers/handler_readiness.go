package handlers

import (
	"net/http"
	"github.com/souls-syntax/Templates/internal/utils"
)

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJson(w,200,struct{}{})
}
