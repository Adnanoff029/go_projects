package main

import (
	// "context"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	// "github.com/slack-io/slacker"
	"github.com/shomali11/slacker"
)

func getEnvValue(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv(key)
}

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	example := []string{"my birth year is 2002"}
	bot.Command("my birth year is <year>", &slacker.CommandDefinition{
		Description: "We get year of birth",
		Examples:    example,
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			year := request.Param("year")
			yob, err := strconv.Atoi(year)
			if err != nil {
				println(err)
			}
			age := 2025 - yob
			r := fmt.Sprintf("age is %d", age)
			response.Reply(r)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
