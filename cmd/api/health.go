package main

import (
	"net/http"
)

// healthcheckHandler godoc
//
//	@Summary		Healthcheck
//	@Description	Healthcheck endpoint
//	@Tags			ops
//	@Produce		json
//	@Success		200	{object}	string	"ok"
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"status": "ok"}

	// log.Println("Sleeping for test before sending the health check response")
	// time.Sleep(time.Second * 4)

	err := app.writeResponse(w, http.StatusOK, data)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
	}

}
