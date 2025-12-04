package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.PLATFORM != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset only allowed in dev"))
		return
	} else {
		cfg.fileserverHits.Store(0)

		ctx := r.Context()
		err := cfg.db.DeleteAllUsers(ctx)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to delete users", err)
			return
		}
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hits reset to 0 and users reset"))
	}
}