package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spinup-host/config"
)

// config provides a way for users to change configuration values (config.yaml) without interacting
// directly with the yaml file


// ToggleMonitoring toggles the current "system-wide" monitoring status for Spinup
func ToggleMonitoring(w http.ResponseWriter, r *http.Request) {
	if (*r).Method != "POST" {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}

	var newConfig config.MutableConfig
	byteArray, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error %v", err)
		ErrorResponse(w, "error reading from request body", 500)
		return
	}
	if err = json.Unmarshal(byteArray, &newConfig); err != nil {
		fmt.Printf("error %v", err)
		ErrorResponse(w, "error parsing request body", 500)
		return
	}

	if err = config.Cfg.UpdateConfig(config.MutableConfig{
		Monitoring: newConfig.Monitoring,
	}); err != nil {
		fmt.Printf("error %v", err)
		ErrorResponse(w, "could not update Spinup configuration", 500)
		return
	}
	respond(http.StatusOK, w, map[string]string{
		"message": "Configuration has been updated",
	})
}