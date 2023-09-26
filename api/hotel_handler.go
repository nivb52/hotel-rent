package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nivb52/hotel-rent/db"
)

type HotelHandler struct {
	hotelStore db.HotelStore
	roomStore  db.RoomStore
}

func NewHotelHandler(hs db.HotelStore, rs db.RoomStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hs,
		roomStore:  rs,
	}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {

	hotels, err := h.hotelStore.GetHotels(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(hotels)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	var id = c.Params("id")
	hotel, err := h.hotelStore.GetHotel(c.Context(), id)

	if err != nil {
		return err
	}

	return c.JSON(hotel)
}
