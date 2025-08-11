package main

import "time"

type AnyRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Slug   string    `json:"slug"`
	Expiry time.Time `json:"expiry"`
}
