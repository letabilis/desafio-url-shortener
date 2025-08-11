package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type API struct {
	svc *Service
}

func NewAPI(svc *Service) *API {
	return &API{
		svc: svc,
	}
}

func (api *API) ShortenURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var long AnyRequest
		err := json.NewDecoder(r.Body).Decode(&long)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		shortenResponse, err := api.svc.GetSlug(r.Context(), long.URL)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		WriteJSON(w, http.StatusOK, shortenResponse)
	}
}

func (api *API) ResolveURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		longURL, err := api.svc.GetLongURL(r.Context(), slug)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Redirect(w, r, longURL, http.StatusFound)
	}
}
