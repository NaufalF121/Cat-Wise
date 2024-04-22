package server

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/gofiber/fiber/v2"
	"github.com/shareed2k/goth_fiber"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Get("/", s.HelloWorldHandler)
	s.App.Get("/login/:provider", goth_fiber.BeginAuthHandler)
	s.App.Get("/auth/callback/:provider", s.getAuthCallback)
	s.App.Get("/logout", s.logout)
	s.App.Get("/health", s.healthHandler)
	s.App.Post("/wise", s.CreateWise)
	s.App.Get("/wise", s.GetWise)

}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Id    string `json:"id"`
	Wise  string `json:"wise"`
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) getAuthCallback(c *fiber.Ctx) error {

	user, err := goth_fiber.CompleteUserAuth(c)
	if err != nil {
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(user)
}

func (s *FiberServer) logout(ctx *fiber.Ctx) error {
	if err := goth_fiber.Logout(ctx); err != nil {
		log.Fatal(err)
	}

	return ctx.Redirect("/")
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}

func (s *FiberServer) CreateWise(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	query := "INSERT INTO accounts (username, email, user_id, wise) VALUES ($1, $2, $3, $4) RETURNING user_id"
	result, err := s.db.Exec(query, user.Name, user.Email, user.Id, user.Wise)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	fmt.Println("New user added ", result)
	return c.Status(fiber.StatusCreated).JSON(user)
}

func (s *FiberServer) GetWise(c *fiber.Ctx) error {
	query := "SELECT wise FROM accounts"
	rows, err := s.db.Query(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer rows.Close()
	var wiseList []string
	for rows.Next() {
		var wise string
		if err := rows.Scan(&wise); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		wiseList = append(wiseList, wise)
	}
	randomNumber := rand.Intn(len(wiseList))
	randomWise := wiseList[randomNumber]
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"wise": randomWise,
	})

}
