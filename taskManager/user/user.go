package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"taskManager/db"
	"taskManager/logger"
	"time"
)

const secret = "98d987f876865f6fa35bcf55f343962080064edfaa30763cf24b5c49cba7930754fb07df635a454b031fd927f607a5e5be15f7ad5306f721c1047b5c59952434"

var log = logger.New()

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `bson:"username" json:"username"`
	Email    string             `bson:"email" json:"email"`
	Password []byte             `bson:"password" json:"password"`
	Valid    bool               `bson:"token_valid" json:"token_valid"`
}

type userRegisterDetails struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type userLoginDetails struct {
	Field1   string `json:"field1"` // username or email
	Password string `json:"password"`
}

func createUser(username string, email string, password []byte) *User {
	return &User{primitive.NewObjectID(), username, email, password, false}
}

func Register(c *fiber.Ctx) error {
	var userCredentials userRegisterDetails
	err := c.BodyParser(&userCredentials)

	if err != nil {
		log.Error("Body Parser: %s", err)
		return c.Status(400).JSON(fiber.Map{"message": "An error occurred."})
	}

	if db.SearchUserFromUsername(userCredentials.Username).Err() == nil {
		return c.Status(409).JSON(fiber.Map{"message": "Username already taken"})
	} else if db.SearchUserFromEmail(userCredentials.Email).Err() == nil {
		return c.Status(409).JSON(fiber.Map{"message": "Email already taken"})
	}

	pw, err := bcrypt.GenerateFromPassword([]byte(userCredentials.Password), 14)
	if err != nil {
		log.Error("GenerateFromPassword: %s", err)
		return c.Status(500).JSON(fiber.Map{"message": "An error occurred"})
	}

	user := createUser(userCredentials.Username, userCredentials.Email, pw)
	_, err = db.InsertItem(db.UsersCollection, user)

	if err != nil {
		log.Error("Insert Item: %s", err)
		return c.Status(500).JSON(fiber.Map{"message": "Could not insert item to the database"})
	}

	log.Info("User '%s' registered.", user.ID)
	return c.Status(200).JSON(fiber.Map{"message": "Registered successfully"})
}

func Login(c *fiber.Ctx) error {
	var userCredentials userLoginDetails
	err := c.BodyParser(&userCredentials)
	if err != nil {
		log.Error("Body Parser: %s", err.Error())
		return c.Status(400).JSON(fiber.Map{"message": "An error occurred."})
	}

	if res := db.SearchUserFromEmail(userCredentials.Field1); res.Err() == nil {
		var user User
		err := res.Decode(&user)
		if err != nil {
			log.Error("idk man %s", err)
			return c.Status(400).JSON(fiber.Map{"message": "An error occurred"})
		}
		if bcrypt.CompareHashAndPassword(user.Password, []byte(userCredentials.Password)) == nil {
			tok, err := createToken(user)
			if err != nil {
				log.Error("Create Token: %s", err)
				return c.Status(500).JSON(fiber.Map{"message": "An error occurred."})
			}
			log.Info("%s", user.ID)
			_, err = db.ValidateToken(user.ID)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"message": "An error occurred."})
			}

			return c.Status(200).JSON(fiber.Map{"token": tok})
		}
		return c.Status(401).JSON(fiber.Map{"message": "Incorrect password"})
	} else if res := db.SearchUserFromUsername(userCredentials.Field1); res.Err() == nil {
		var user User
		err := res.Decode(&user)
		if err != nil {
			log.Error("idk man %s", err)
			return c.Status(400).JSON(fiber.Map{"message": "An error occurred"})
		}
		if bcrypt.CompareHashAndPassword(user.Password, []byte(userCredentials.Password)) == nil {
			tok, err := createToken(user)
			if err != nil {
				log.Error("Create Token: %s", err)
				return c.Status(500).JSON(fiber.Map{"message": "An error occurred"})
			}
			_, err = db.ValidateToken(user.ID)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"message": "An error occurred."})
			}

			return c.Status(200).JSON(fiber.Map{"token": tok})
		}
		return c.Status(401).JSON(fiber.Map{"message": "Incorrect password"})
	}
	return c.Status(400).JSON(fiber.Map{"message": "No users found"})
}

func Logout(c *fiber.Ctx) error {
	userID := c.Locals("ownerID").(primitive.ObjectID)
	_, err := db.InvalidateToken(userID)
	if err != nil {
		log.Error("%s", err)
		return c.Status(500).JSON(fiber.Map{"message": "An error occurred"})
	}

	return c.Status(200).RedirectToRoute("http://localhost:3000/", fiber.Map{})
}

func createToken(user User) (string, error) {
	claims := jwt.MapClaims{
		"ID":       user.ID,
		"email":    user.Email,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 6).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tok, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}
	log.Info(tok)
	return tok, nil
}
