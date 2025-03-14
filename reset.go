package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment"))
		return
	}
	cfg.fileserverHits.Store(0)
	cfg.db.ResetDB(req.Context())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 & database reset"))
}
