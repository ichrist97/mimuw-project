package model

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// device enum
type Device int

const (
	Mobile Device = iota + 1
	PC
	TV
)

func DeviceString(d Device) string {
	return []string{"Mobile", "PC", "TV"}[d]
}

// action enum
type Action int

const (
	VIEW Action = iota + 1
	BUY
)

func ActionString(a Action) string {
	return []string{"VIEW", "BUY"}[a]
}

type Product struct {
	ProductId  string `json:"product_id" validate:"required"`
	BrandId    string `json:"brand_id" validate:"required"`
	CategoryId string `json:"category_id" validate:"required"`
	Price      int    `json:"price" validate:"required,number"`
}

type UserTagEvent struct {
	Time        string  `json:"time" validate:"required"`
	Cookie      string  `json:"cookie" validate:"required"`
	Country     string  `json:"country" validate:"required"`
	Device      Device  `json:"device" binding:"required,enum"`
	Action      Action  `json:"action" binding:"required,enum"`
	Origin      string  `json:"origin" validate:"required"`
	ProductInfo Product `json:"product_info" validate:"dive"`
}

type ErrorResponse struct {
	Field string
	Tag   string
	Value string
}

var validate = validator.New()

func ValidateUserTagEvent(c *fiber.Ctx) error {
	var errors []*ErrorResponse
	body := new(UserTagEvent)
	c.BodyParser(&body)

	err := validate.Struct(body)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var el ErrorResponse
			el.Field = err.Field()
			el.Tag = err.Tag()
			el.Value = err.Param()
			errors = append(errors, &el)
		}
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	return c.Next()
}
