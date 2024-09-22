package piano

import (
	"regexp"

	"github.com/google/uuid"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2"
	"github.com/jindongh/home/common"
)
var TITLE_PAT = regexp.MustCompile(`\nT: (?P<title>.*)\n`)

// SetupRoutes func
func Route(app *fiber.App, store *session.Store, config *common.Config) {
	app.Get("/piano", func(c *fiber.Ctx) error {
		sess, _ := store.Get(c)
		return c.Render("templates/piano", fiber.Map{
			"Email": sess.Get("email"),
			"Config": config,
			"Action": "piano",
		})
	})
	app.Get("/piano/:id", func(c *fiber.Ctx) error {
		sess, _ := store.Get(c)
		return c.Render("templates/piano", fiber.Map{
			"Email": sess.Get("email"),
			"Config": config,
			"Action": "piano",
			"ID": c.Params("id"),
		})
	})
	app.Get("/api/piano", func(c *fiber.Ctx) error {
		return listSongs(c, getUser(c, store))
	})
	app.Get("/api/piano/:id", func (c *fiber.Ctx) error {
		return getSong(c, getUser(c, store), c.Params("id"))
	})
	app.Post("/api/piano", func(c *fiber.Ctx) error {
		return addSong(c, getUser(c, store))
	})
	app.Delete("/api/piano/:id", func(c *fiber.Ctx) error {
		return removeSong(c, getUser(c, store), c.Params("id"))
	})
}

func getUser(c *fiber.Ctx, store *session.Store) string {
	sess, _ := store.Get(c)
	email := sess.Get("email")
	if email == nil {
		return "guest"
	}
	return email.(string)
}
func addSong(c *fiber.Ctx, user string) error {
	db := DB.Db
	song := new(Song)
	// Store the body in the user and return error if encountered
	err := c.BodyParser(song)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "Something's wrong with your input",
			"data": err})
	}
	match := TITLE_PAT.FindStringSubmatch(song.Content)
	if len(match) > 1 {
		song.Title = match[1]
	}
	song.Username = user
	err = db.Create(&song).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"message":  "Could not create user",
			"data": err})
	}
	// Return the created user
	return c.Status(201).JSON(fiber.Map{
		"status": "success",
		"message":  "User has created",
		"data": user})
}
func getSong(c *fiber.Ctx, user, id string) error {
	db := DB.Db
	var song Song
	// find all users in the database
	db.Find(&song, "username = ? and id = ?", user, id)
	if song.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"status": "error",
			"message": "Users not found",
			"data": nil})
	}
	return c.Status(200).JSON(fiber.Map{
		"status": "sucess",
		"message": "Song Found",
		"data": song})
}
func removeSong(c *fiber.Ctx, user, id string) error {
	db := DB.Db
	var song Song
	db.Find(&song, "username = ? and id = ?", user, id)
	if song.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"status": "error",
			"message": "Song not found",
			"data": nil})
	}
	err := db.Delete(&song, "username = ? and id = ?", user, id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status": "error",
			"message": "Failed to delete user",
			"data": nil})
	}
	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"message": "Song deleted"})
}

func listSongs(c *fiber.Ctx, user string) error {
	db := DB.Db
	var songs []Song
	db.Find(&songs, "username = ?", user)
	return c.Status(200).JSON(fiber.Map{
		"status": "sucess",
		"data": songs})
}
