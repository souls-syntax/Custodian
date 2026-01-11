package handlers

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/souls-syntax/Templates/internal/utils"
	"github.com/souls-syntax/Templates/internal/models"
)

func (c *ApiConfig) HandlerVerify(w http.ResponseWriter, r *http.Request) {
	var req models.Query
	
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		utils.RespondWithError(w,400,"Incorrect format of JSON query")
		return
	}
	
	ctx := r.Context()

	resp, err := c.Verifier.Verify(ctx,req.Query)
	if err != nil {
		fmt.Printf("❌ CRITICAL FAILURE: %v\n", err)
		utils.RespondWithError(w,500,"Verification failed")
		return
	}

	utils.RespondWithJson(w,200,resp)
}
