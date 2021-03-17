package main

import (
	"bytes"
	"github.com/SEFI2/instagen/utils"
	"github.com/TheForgotten69/goinsta/v2"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/shurcooL/graphql"
	"golang.org/x/net/context"
	"image/jpeg"

	"fmt"
	"log"
	"time"
)

func instaHandler() (int, error) {
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
		"limit": graphql.Int(10),
		"skip": graphql.Int(0),
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

	fmt.Println("Upload images")
	// Create images
	for _, job := range q.Jobs {
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

		caption := fmt.Sprintf("%s, %s %s", string(job.Title), string(job.Phone), "#работа #москва #жумуш #jumush #работамосква #жумушмосква #москважумуш #жердеш")
		_, err = insta.UploadPhoto(buf, caption, 1, 1)
		if err != nil {
			fmt.Println("Cannot upload. Error:", err)
			continue
		}

		time.Sleep(30 * time.Second)
	}

	return 0, nil
}

func main() {
	lambda.Start(instaHandler)
}
