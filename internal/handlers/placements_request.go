package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/ansedo/toptraffic/internal/logger"
	"github.com/ansedo/toptraffic/internal/models"
)

const (
	BidRequestPath    = "/bid_request"
	BidRequestTimeout = 200 * time.Millisecond
)

func PlacementsRequest(ctx context.Context, advDomains []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
			logger.FromCtx(ctx).Warn(err.Error(), zap.Int("status code", http.StatusBadRequest))
			return
		}

		var placementRequest models.PlacementRequest
		if err = json.Unmarshal(body, &placementRequest); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
			logger.FromCtx(ctx).Warn(err.Error(), zap.Int("status code", http.StatusBadRequest))
			return
		}
		if err = placementRequest.Validate(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
			logger.FromCtx(ctx).Warn(err.Error(), zap.Int("status code", http.StatusBadRequest))
			return
		}

		var bidRequest models.BidRequest
		bidRequest.LoadFromPlacementRequest(placementRequest)
		bidResponse, err := sendBidRequests(ctx, bidRequest, advDomains)
		if err != nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		bidResponse.LeaveOnlyMaxProfitImps()

		var placementResponse models.PlacementResponse
		placementResponse.LoadFromPlacementRequestAndBidResponse(placementRequest, bidResponse)
		if placementResponse.IsEmpty() {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(placementResponse)
	}
}

func sendBidRequests(ctx context.Context, bidRequest models.BidRequest, advDomains []string) (models.BidResponse, error) {
	var bidResponse models.BidResponse
	bidRequestJSON, err := json.Marshal(&bidRequest)
	if err != nil {
		logger.FromCtx(ctx).Warn(`json marshal bid request`, zap.Error(err))
		return bidResponse, err
	}

	wg := &sync.WaitGroup{}
	chBidResponses := make(chan models.BidResponse, len(bidRequest.Imp)*len(advDomains))
	for _, advDomain := range advDomains {
		wg.Add(1)
		go func(ctx context.Context, wg *sync.WaitGroup, advDomain string) {
			defer wg.Done()
			bidResponse, err = sendBidRequest(ctx, bidRequestJSON, advDomain)
			if err != nil {
				logger.FromCtx(ctx).Warn(`request to advDomain: `+advDomain, zap.Error(err))
				return
			}
			if bidResponse.IsEmpty() {
				return
			}

			chBidResponses <- bidResponse
		}(ctx, wg, advDomain)
	}
	wg.Wait()
	close(chBidResponses)

	bidResponse.ID = bidRequest.ID
	for b := range chBidResponses {
		bidResponse.Imp = append(bidResponse.Imp, b.Imp...)
	}
	return bidResponse, nil
}

func sendBidRequest(ctx context.Context, bidRequestJSON []byte, advDomain string) (models.BidResponse, error) {
	requestURL := advDomain + BidRequestPath
	var bidResponse models.BidResponse
	ctx, cancel := context.WithTimeout(ctx, BidRequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewBuffer(bidRequestJSON))
	if err != nil {
		return bidResponse, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.FromCtx(ctx).Warn(`request to url: `+requestURL, zap.Error(err))
		return bidResponse, err
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		logger.FromCtx(ctx).Warn(`request to url: `+requestURL, zap.Error(err))
		return bidResponse, err
	}
	if res.StatusCode > 299 {
		logger.FromCtx(ctx).Warn(`request to url: `+requestURL, zap.Error(err), zap.Int("status code", res.StatusCode))
		return bidResponse, err
	}
	if res.StatusCode == http.StatusNoContent {
		return bidResponse, nil
	}

	if err = json.Unmarshal(body, &bidResponse); err != nil {
		logger.FromCtx(ctx).Warn(`request to url: `+requestURL, zap.Error(err), zap.Int("status code", res.StatusCode))
		return bidResponse, err
	}

	return bidResponse, nil
}
