package esi

import "time"

type BaseEsiModel interface {
	Expires() time.Time
	Pages() int
	SetExpires(expiresAt time.Time)
	SetPages(pages int)
}

type baseEsiModel struct {
	expires time.Time
	pages   int
}

func (b *baseEsiModel) Expires() time.Time {
	return b.expires
}

func (b *baseEsiModel) Pages() int {
	return b.pages
}

func (b *baseEsiModel) SetExpires(expires time.Time) {
	b.expires = expires
}

func (b *baseEsiModel) SetPages(pages int) {
	b.pages = pages
}
