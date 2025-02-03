package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alexchebotarsky/heatpump-api/heatpump"
)

type StateFetcher interface {
	FetchState() (*heatpump.State, error)
}

func GetHeatpumpState(stateFetcher StateFetcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state, err := stateFetcher.FetchState()
		if err != nil {
			HandleError(w, fmt.Errorf("error fetching state: %v", err), http.StatusInternalServerError, true)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(state)
		handleWritingErr(err)
	}
}

type StateUpdater interface {
	UpdateState(state heatpump.State) (*heatpump.State, error)
}

func UpdateHeatpumpState(StateUpdater StateUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var state heatpump.State
		err := json.NewDecoder(r.Body).Decode(&state)
		if err != nil {
			HandleError(w, fmt.Errorf("error decoding POST body: %v", err), http.StatusBadRequest, false)
			return
		}

		err = state.Validate()
		if err != nil {
			HandleError(w, fmt.Errorf("error validating POST body: %v", err), http.StatusBadRequest, false)
			return
		}

		updatedState, err := StateUpdater.UpdateState(state)
		if err != nil {
			HandleError(w, fmt.Errorf("error updating heatpump state: %v", err), http.StatusInternalServerError, true)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(updatedState)
		handleWritingErr(err)
	}
}
