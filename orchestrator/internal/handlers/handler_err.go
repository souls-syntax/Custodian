package handlers

import (
	"net/http"
	"github.com/souls-syntax/Templates/internal/utils"
)

func HandlerErr(w http.ResponseWriter, r *http.Request) {	
	utils.RespondWithError(w,400,"Something Went Wrong")
}
