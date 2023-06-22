package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"taskManager/db"
	"taskManager/logger"
	"taskManager/user"
	"time"
)

const secret = "98d987f876865f6fa35bcf55f343962080064edfaa30763cf24b5c49cba7930754fb07df635a454b031fd927f607a5e5be15f7ad5306f721c1047b5c59952434"

var log = logger.New()

func ExtractToken(c *fiber.Ctx) error {

	authorization := c.Get("authorization")
	if len(authorization) == 0 {
		return c.Status(400).JSON(fiber.Map{"message": "You need a token to operate"})
	}

	token := strings.SplitN(authorization, " ", 2)[1]
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, &claims, func(tok *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		log.Error("%s", err)
		return c.Status(400).JSON(fiber.Map{"message": "An error occurred"})
	}

	userID, err := primitive.ObjectIDFromHex(claims["ID"].(string))
	if err != nil {
		log.Error("%s", err)
		return c.Status(500).JSON(fiber.Map{"message": "An error occurred"})
	}

	c.Locals("ownerID", userID)
	c.Locals("exp", claims["exp"])
	return c.Next()
}

func StaleCheck(c *fiber.Ctx) error {
	exp := int64(c.Locals("exp").(float64))

	if exp < time.Now().Unix() {
		log.Info("token expiration: %d, time: %d", exp, time.Now().Unix())
		return c.Status(403).JSON(fiber.Map{"message": "Token expired."})
	}
	isValid, err := isTokenValid(c.Locals("ownerID").(primitive.ObjectID))
	if err != nil {
		log.Error("%s", err)
		return c.Status(500).JSON(fiber.Map{"message": "An error occurred."})
	} else if !isValid {
		return c.Status(403).JSON(fiber.Map{"message": "Invalid token."})
	}

	return c.Next()
}

func isTokenValid(id primitive.ObjectID) (bool, error) {
	var usr user.User
	res := db.GetItemFromId(db.UsersCollection, id)

	if res.Err() != nil {
		return false, res.Err()
	} else if err := res.Decode(&usr); err != nil {
		return false, err
	}

	return usr.Valid, nil
}
