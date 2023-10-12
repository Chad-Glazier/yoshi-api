package test

import (
	"net/http"
	"yoshi-api/types"
	"yoshi-api/db"
	"strconv"
)

func ReadUser(w http.ResponseWriter, r *http.Request) {
	usrId, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil {
		(&types.ApiError{
			Status: http.StatusBadRequest,
			Message: "You need to provide an ID!",
		}).Send(w)
		return
	}

	connection, err := db.Connect()
	if err != nil {
		(&types.ApiError{
			Status: http.StatusInternalServerError,
			Message: err.Error(),
		}).Send(w)
		return
	}
	defer db.Close(connection)
	
	usr, err := db.ReadUser(connection, usrId)

	if err != nil {
		(&types.ApiError{
			Status: http.StatusInternalServerError,
			Message: err.Error(),
		}).Send(w)
		return
	}

	if usr == nil {
		(&types.ApiError{
			Status: http.StatusNotFound,
			Message: "That user wasn't found!",
		}).Send(w)
		return
	}

	usr.Send(w)
}