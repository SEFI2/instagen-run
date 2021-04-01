package main

import (
	"bytes"
	"github.com/SEFI2/instagen/utils"
	"github.com/TheForgotten69/goinsta/v2"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/aws/aws-lambda-go/lambda"
	"github.com/shurcooL/graphql"
	"golang.org/x/net/context"
	"image/jpeg"
	"io"
	"math/rand"

	"fmt"
	"log"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)



var followerCount = 0
var seed = rand.NewSource(time.Now().UnixNano())
var random = rand.New(seed)

func instaUploader() (int, error) {
	fmt.Println("Init client")
	client := graphql.NewClient("http://www.tabyshmak.ru:8080/graphql", nil)
	if client == nil {
		fmt.Println("Client is null")
		return 0, fmt.Errorf("client is null")
	}

	// GraphQL query
	var q struct {
		Jobs []struct {
			Title graphql.String
			Description graphql.String
			JobDate graphql.String
			Phone graphql.String
			Location graphql.String
			Created graphql.String
		} `graphql:"jobs(searchValue: $searchValue, skip: $skip, limit: $limit)"`
	}

	// Variables
	variables := map[string]interface{}{
		"searchValue":   graphql.String(""),
		"limit": graphql.Int(20),
		"skip": graphql.Int(10),
	}

	fmt.Println("Make query...")
	// Query
	err := client.Query(context.Background(), &q, variables)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	// Login to instagram
	insta := goinsta.New("tabyshmak.ru", "firdavs13")
	if err := insta.Login(); err != nil {
		log.Fatal(err)
		return 0, err
	}
	defer insta.Logout()

	// Telegram group upload
	var groupChanID int64 = -1001494836905
	bot, err := tgbotapi.NewBotAPI("1721591920:AAH5pHPAENB4tWuFN60IjBGrmIpzgMt_YNM")
	if err != nil {
		log.Panic(err)
		return 0, err
	}
	bot.Debug = true

	fmt.Println("Upload images")
	fmt.Println(len(q.Jobs))
	// Create images
	for i := 0; i < len(q.Jobs); i += 10 {
		var toBePublished []io.Reader

		for j := i; j < len(q.Jobs) && j < i + 10; j ++ {
			job := q.Jobs[j]
			title := fmt.Sprintf("%s. %s. %s", string(job.Title), string(job.Location), string(job.Phone))
			content := string(job.Description)
			author := "@tabyshmak.ru"
			img, err := utils.CreateInstagramPost(title, content, author)
			if err != nil {
				fmt.Println("Cannot create an image. Error:", err)
				continue
			}

			buf := new(bytes.Buffer)
			err = jpeg.Encode(buf, img, nil)
			if err != nil {
				fmt.Println("Cannot encode. Error:", err)
				continue
			}

			b := tgbotapi.FileBytes{Name: "image.jpg", Bytes: buf.Bytes()}
			msg := tgbotapi.NewPhotoUpload(groupChanID, b)
			msg.Caption = fmt.Sprintf("%s, %s",string(job.Title), string(job.Phone))
			_, err = bot.Send(msg)
			toBePublished = append(toBePublished, buf)
			time.Sleep(10 * time.Second)
		}


		caption := fmt.Sprintf("%s", "#работа #москва #жумуш #jumush #работамосква #жумушмосква #москважумуш #жердеш")
		_, err = insta.UploadAlbum(toBePublished, caption, 1, 1)
		// _, err = insta.UploadPhoto(buf, caption, 1, 1)
		if err != nil {
			fmt.Println("Cannot upload. Error:", err)
			continue
		}

		time.Sleep(30 * time.Second)
	}

	return 0, nil
}

func main() {
	// instaUploader()
	lambda.Start(instaUploader)

	/*
	var groupChanID int64 = -1001494836905
	bot, _ := tgbotapi.NewBotAPI("1721591920:AAH5pHPAENB4tWuFN60IjBGrmIpzgMt_YNM")

	msg := tgbotapi.NewMessage(groupChanID, "" +
		"Удобная платформа для поиска вакансий и подработок \n" +
		"Instagram: https://www.instagram.com/tabyshmak.ru/\n" +
		"Android: play.google.com/store/apps/details?id=tabyshmak.app.mobile\n")
	bot.Send(msg)
	 */
}



