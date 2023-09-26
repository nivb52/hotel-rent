package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nivb52/hotel-rent/db"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

type HotelQueryParams struct {
	Rooms  bool
	Rating int
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	// var qParams HotelQueryParams
	// err := c.QueryParser(qParams)
	// if err != nil {
	// 	return err
	// }

	hotels, err := h.store.Hotel.GetHotels(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(hotels)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	var id = c.Params("id")
	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), id)

	if err != nil {
		return err
	}

	return c.JSON(hotel)
}

func (h *HotelHandler) HandleGetHotelRooms(c *fiber.Ctx) error {
	var id = c.Params("id")
	hotel, err := h.store.Room.GetHotelRooms(c.Context(), id)

	if err != nil {
		return err
	}

	return c.JSON(hotel)
}
