package session

import (
	"context"
	"encoding/json"
	"net/http"

	frame "nxy/testsocket/agent/service"

	httptransport "github.com/go-kit/kit/transport/http"
)

func MakeHandler(bs Service) http.Handler {
	return httptransport.NewServer(makeTestEndpoint(bs), decodeRequest, encodeResponse)
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var f frame.Frame
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		return nil, err
	}
	return f, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
