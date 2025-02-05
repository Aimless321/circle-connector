package account_link

import (
	"connector-backend/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
	"strconv"
)

func GetAll(c *fiber.Ctx) error {
	var accounts []Account
	result := models.Db.Find(&accounts)

	if result.RowsAffected == 0 {
		body := models.Body{ErrorField: models.Error{
			Status:  404,
			Message: "Cannot find any accounts",
		}}
		return c.Status(body.ErrorField.Status).JSON(body)
	}

	body := models.Body{Data: accounts}
	return c.Status(body.ErrorField.Status).JSON(body)
}

func LinkAccount(c *fiber.Ctx) error {
	payload := struct {
		SteamId   string `json:"steamId"`
		DiscordId string `json:"discordId"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	if len(payload.SteamId) == 0 || len(payload.DiscordId) == 0 {
		body := models.Body{ErrorField: models.Error{
			Status:  400,
			Message: "Empty steamId or discordId provided",
		}}
		return c.Status(body.ErrorField.Status).JSON(body)
	}

	steamId, err := parseUint64(payload.SteamId, "steamId")
	if err != nil {
		return c.Status(err.(models.Body).ErrorField.Status).JSON(err)
	}
	discordId, err := parseUint64(payload.DiscordId, "discordId")
	if err != nil {
		return c.Status(err.(models.Body).ErrorField.Status).JSON(err)
	}

	account := Account{
		DiscordId: discordId,
		SteamID64: steamId,
	}
	result := models.Db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "discordid"}, {Name: "steamid"}},
		DoUpdates: clause.AssignmentColumns([]string{"discordid", "steamid"}),
	}).Create(&account)
	if result.Error != nil {
		body := models.Body{ErrorField: models.Error{
			Status:  500,
			Message: result.Error.Error(),
		}}
		return c.Status(body.ErrorField.Status).JSON(body)
	}

	body := models.Body{Data: models.MessageResponse{
		Message: "Succefully linked account",
	}}
	return c.Status(body.ErrorField.Status).JSON(body)
}

func GetAllAccountsByDiscordId(c *fiber.Ctx) error {
	payload := struct {
		DiscordIds []uint64 `json:"discordIds,string"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	if len(payload.DiscordIds) > 256 {
		body := models.Body{ErrorField: models.Error{
			Status:  400,
			Message: "Cannot query more than 256 accounts at the same time",
		}}
		return c.Status(body.ErrorField.Status).JSON(body)
	}

	var accounts []Account
	result := models.Db.Where("discordid IN ?", payload.DiscordIds).Find(&accounts)

	if result.RowsAffected == 0 {
		body := models.Body{ErrorField: models.Error{
			Status:  404,
			Message: "Cannot find any users with the provided discordids",
		}}
		return c.Status(body.ErrorField.Status).JSON(body)
	}

	body := models.Body{Data: accounts}
	return c.Status(body.ErrorField.Status).JSON(body)
}

func GetAccountBySteamId(c *fiber.Ctx) error {
	steamId, err := parseUint64(c.Params("steamid"), "steamid")
	if err != nil {
		return c.Status(err.(models.Body).ErrorField.Status).JSON(err)
	}

	var account Account
	result := models.Db.Take(&account, Account{SteamID64: steamId})

	if result.RowsAffected == 0 {
		body := models.Body{ErrorField: models.Error{
			Status:  404,
			Message: "Cannot find user with steamid " + c.Params("steamid"),
		}}
		return c.Status(body.ErrorField.Status).JSON(body)
	}

	body := models.Body{Data: account}
	return c.Status(body.ErrorField.Status).JSON(body)
}

func GetAllAccountsBySteamId(c *fiber.Ctx) error {
	payload := struct {
		SteamIds []uint64 `json:"steamIds,string"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	if len(payload.SteamIds) > 256 {
		body := models.Body{ErrorField: models.Error{
			Status:  400,
			Message: "Cannot query more than 256 accounts at the same time",
		}}
		return c.Status(body.ErrorField.Status).JSON(body)
	}

	var accounts []Account
	result := models.Db.Where("steamid IN ?", payload.SteamIds).Find(&accounts)

	if result.RowsAffected == 0 {
		body := models.Body{ErrorField: models.Error{
			Status:  404,
			Message: "Cannot find any users with the provided steamids",
		}}
		return c.Status(body.ErrorField.Status).JSON(body)
	}

	body := models.Body{Data: accounts}
	return c.Status(body.ErrorField.Status).JSON(body)
}

func GetAccountByDiscordId(c *fiber.Ctx) error {
	discordId, err := parseUint64(c.Params("discordid"), "discordid")
	if err != nil {
		return c.Status(err.(models.Body).ErrorField.Status).JSON(err)
	}

	var account Account
	result := models.Db.Take(&account, Account{DiscordId: discordId})

	if result.RowsAffected == 0 {
		body := models.Body{ErrorField: models.Error{
			Status:  404,
			Message: "Cannot find user with discordid " + c.Params("discordid"),
		}}
		return c.Status(body.ErrorField.Status).JSON(body)
	}

	body := models.Body{Data: account}
	return c.Status(body.ErrorField.Status).JSON(body)
}

func parseUint64(string string, field string) (uint64, error) {
	intResult, err := strconv.ParseUint(string, 10, 64)
	if err != nil {
		body := models.Body{ErrorField: models.Error{
			Status:  400,
			Message: fmt.Sprintf("Invalid %v provided", field),
		}}

		return 0, body
	}

	return intResult, nil
}
