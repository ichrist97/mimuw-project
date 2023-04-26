package model

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ProductId  string `json:"product_id" validate:"required"`
	BrandId    string `json:"brand_id" validate:"required"`
	CategoryId string `json:"category_id" validate:"required"`
	Price      int    `json:"price" validate:"required,number"`
}

type UserTagEvent struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Time        time.Time          `json:"time" validate:"required"`
	Cookie      string             `json:"cookie" validate:"required"`
	Country     string             `json:"country" validate:"required"`
	Device      string             `json:"device" validate:"required"`
	Action      string             `json:"action" validate:"required"`
	Origin      string             `json:"origin" validate:"required"`
	ProductInfo Product            `json:"product_info" validate:"dive"`
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
