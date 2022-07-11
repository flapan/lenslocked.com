package controllers

import (
	"github.com/flapan/lenslocked.com/models"
	"github.com/flapan/lenslocked.com/views"
)

func NewGalleries(gs models.GalleryService) *Galleries {
	return &Galleries{
		New: views.NewView("bootstrap", "galleries/new"),
	}
}

type Galleries struct {
	New *views.View
	gs  models.GalleryService
}
