package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/Adnanoff029/url-shortener/api/database"
	"github.com/Adnanoff029/url-shortener/api/helpers"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type request struct {
	URL      string        `json:"url"`
	ShortURL string        `json:"short_url"`
	Expiry   time.Duration `json:"expiry"`
}

type response struct {
	URL             string        `json:"url"`
	ShortURL        string        `json:"short_url"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

func ShortenURL(c *fiber.Ctx) error {
	body := new(request)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Cannot parse JSON"})
	}
	// Rate limiting
	r2 := database.CreateClient(1)
	defer r2.Close()
	val, err := r2.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		_ = r2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else {
		val, _ = r2.Get(database.Ctx, c.IP()).Result()
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"message":         "Rate limit exceeded",
				"rate_limit_rest": limit / time.Nanosecond / time.Minute,
			})
		}
	}
	// Check if the input is of url format abc.xyz
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid URL"})
	}
	// Handle domain error
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"message": "The service is unavailable."})
	}
	// enforce https for security
	body.URL = helpers.EnforceHTTP(body.URL)
	var id string
	if body.ShortURL == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.ShortURL
	}

	r := database.CreateClient(0)
	defer r.Close()
	val, _ = r.Get(database.Ctx, id).Result()

	if val != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Short URL already exists.",
		})
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}
	err = r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Unable to connect to the server",
		})
	}

	resp := response{
		URL:             body.URL,
		ShortURL:        "",
		Expiry:          body.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: 30,
	}

	r2.Decr(database.Ctx, c.IP())

	val, _ = r2.Get(database.Ctx, c.IP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)
	ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()
	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute
	resp.ShortURL = os.Getenv("DOMAIN") + "/" + id

	return c.Status(fiber.StatusOK).JSON(resp)
}
