package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alexchebotarsky/heatpump-api/model/heatpump"
)

type HeatpumpStateFetcher interface {
	FetchHeatpumpState() (*heatpump.State, error)
}

func GetHeatpumpState(fetcher HeatpumpStateFetcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state, err := fetcher.FetchHeatpumpState()
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

type HeatpumpStateUpdater interface {
	UpdateHeatpumpState(*heatpump.State) (*heatpump.State, error)
}

type IRTransmitter interface {
	TransmitIRSignal(ctx context.Context, binaryString string) error
}

func UpdateHeatpumpState(updater HeatpumpStateUpdater, irTransmitter IRTransmitter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var state heatpump.State
		err := json.NewDecoder(r.Body).Decode(&state)
		if err != nil {
			HandleError(w, fmt.Errorf("error decoding heatpump state: %v", err), http.StatusBadRequest, false)
			return
		}

		err = state.Validate()
		if err != nil {
			HandleError(w, fmt.Errorf("error validating heatpump state: %v", err), http.StatusBadRequest, false)
			return
		}

		updatedState, err := updater.UpdateHeatpumpState(&state)
		if err != nil {
			HandleError(w, fmt.Errorf("error updating heatpump state: %v", err), http.StatusInternalServerError, true)
			return
		}

		binaryString, err := updatedState.ToBinary()
		if err != nil {
			HandleError(w, fmt.Errorf("error converting heatpump state to binary: %v", err), http.StatusInternalServerError, true)
			return
		}

		err = irTransmitter.TransmitIRSignal(r.Context(), binaryString)
		if err != nil {
			HandleError(w, fmt.Errorf("error publishing binary heatpump state: %v", err), http.StatusInternalServerError, true)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(updatedState)
		handleWritingErr(err)
	}
}
