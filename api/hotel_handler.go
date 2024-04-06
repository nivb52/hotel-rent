package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	e "github.com/nivb52/hotel-rent/api/errors"
	"github.com/nivb52/hotel-rent/db"
	"github.com/nivb52/hotel-rent/types"
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
	types.HotelFilter
	db.Pagination
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var qParams HotelQueryParams
	err := c.QueryParser(&qParams)
	if err != nil {
		return e.ErrConflict(c, "Those missing filter data")
	}

	hotels, err := h.store.Hotel.GetHotels(c.Context(), &qParams.HotelFilter, &qParams.Pagination)
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

// func find the free dates of a room
func (h *HotelHandler) HandleGetHotelRooms(c *fiber.Ctx) error {
	var (
		id          = c.Params("id")
		whereClause types.HotelFilter
	)

	err := c.BodyParser(&whereClause)
	if err != nil {
		fmt.Println("Get Hotel Params Filter - Failed to parse, due: ", err)
		return c.Status(fiber.ErrBadRequest.Code).SendString(fiber.ErrBadRequest.Message)
	}

	hotel, err := h.store.Room.GetHotelRooms(c.Context(), id, &whereClause)
	if err != nil {
		fmt.Println("GetHotelRooms with filter failed, due: ", err)
		return c.Status(fiber.StatusInternalServerError).SendString(fiber.ErrInternalServerError.Message)
	}

	return c.JSON(hotel)
}
