package handler

import (
	"crtexBalance/internal/models"
	"log"
	"net/http"

	"github.com/mailru/easyjson"
)

// @Summary Get Balance
// @Tags balance
// @Description getting the user's balance
// @ID get-balance
// @Accept  json
// @Produce  json
// @Param input body models.User true "user id"
// @Success 200 {object} models.User
// @Failure 500 {object} models.Response
// @Router / [post]
func (h *Handler) getBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var err error
	var user models.User
	var newUser *models.User

	if err = easyjson.UnmarshalFromReader(r.Body, &user); err != nil {
		Error(err, w, http.StatusInternalServerError)
		return
	}

	if err = user.Validate(); err != nil {
		Error(err, w, http.StatusBadRequest)
		return
	}

	if newUser, err = h.services.GetBalance(user.Id); err != nil {
		Error(err, w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = easyjson.MarshalToWriter(newUser, w)
	if err != nil {
		Error(err, w, http.StatusInternalServerError)
		return
	}
}

// @Summary Replenishment Balance
// @Tags balance
// @Description replenishment of the user's balance
// @ID replenishment-balance
// @Accept  json
// @Produce  json
// @Param input body models.Replenishment true "replenishment information"
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /topup [post]
func (h *Handler) replenishmentBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var err error
	var replenishment models.Replenishment

	if err = easyjson.UnmarshalFromReader(r.Body, &replenishment); err != nil {
		Error(err, w, http.StatusInternalServerError)
		return
	}

	if err = replenishment.Validate(); err != nil {
		Error(err, w, http.StatusBadRequest)
		return
	}

	if err = h.services.ReplenishmentBalance(&replenishment); err != nil {
		Error(err, w, http.StatusInternalServerError)
		return
	}

	response := &models.Response{
		Message: "баланс пополнен",
	}
	w.WriteHeader(http.StatusOK)
	_, err = easyjson.MarshalToWriter(response, w)
	if err != nil {
		Error(err, w, http.StatusInternalServerError)
		return
	}
}

// @Summary Money transfer
// @Tags balance
// @Description money transfer between users
// @ID transfer
// @Accept  json
// @Produce  json
// @Param input body models.Money true "transfer information"
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /transfer [post]
func (h *Handler) transfer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var err error
	var money models.Money

	if err = easyjson.UnmarshalFromReader(r.Body, &money); err != nil {
		Error(err, w, http.StatusInternalServerError)
		return
	}

	if err = money.Validate(); err != nil {
		Error(err, w, http.StatusBadRequest)
		return
	}

	if err := h.services.Transfer(&money); err != nil {
		Error(err, w, http.StatusInternalServerError)
		return
	}

	response := &models.Response{
		Message: "перевод стредств выполнен",
	}
	w.WriteHeader(http.StatusOK)
	_, err = easyjson.MarshalToWriter(response, w)
	if err != nil {
		Error(err, w, http.StatusInternalServerError)
		return
	}
}

func Error(err error, w http.ResponseWriter, status int) {
	log.Println(err.Error())
	response := &models.Response{
		Message: err.Error(),
	}
	res, err := easyjson.Marshal(response)
	if err != nil {
		log.Println(err.Error())
		return
	}
	w.WriteHeader(status)
	w.Write(res)
}
