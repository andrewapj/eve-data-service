package esi

import "time"

// BaseEsiModel represents functionality available to all ESI types.
type BaseEsiModel interface {
	Expires() time.Time
	Pages() int
	SetExpires(expiresAt time.Time)
	SetPages(pages int)
}

// baseEsiModel is an implementation of a BaseEsiModel.
type baseEsiModel struct {
	expires time.Time
	pages   int
}

// Expires gets the expiry time of the data.
func (b *baseEsiModel) Expires() time.Time {
	return b.expires
}

// Pages gets the total number of pages available.
func (b *baseEsiModel) Pages() int {
	return b.pages
}

// SetExpires sets the expiry time.
func (b *baseEsiModel) SetExpires(expires time.Time) {
	b.expires = expires
}

// SetPages sets the total number of pages available.
func (b *baseEsiModel) SetPages(pages int) {
	b.pages = pages
}
