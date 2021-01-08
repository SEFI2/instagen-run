package main

import (
	"fmt"
	"github.com/SEFI2/instagen/utils"
	"github.com/TheForgotten69/goinsta/v2"
	"github.com/shurcooL/graphql"
	"golang.org/x/net/context"
	"log"
	"os"
	"time"
)

func main() {
	client := graphql.NewClient("http://localhost:8080/graphql", nil)
	if client == nil {
		fmt.Println("Client is null")
		return
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
		"limit": graphql.Int(1),
		"skip": graphql.Int(0),
	}

	// Query
	err := client.Query(context.Background(), &q, variables)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := os.Mkdir("results-png", os.ModePerm); err != nil {
		fmt.Println(err)
		return
	}

	if err := os.Mkdir("results-jpg", os.ModePerm); err != nil {
		fmt.Println(err)
		return
	}

	// Create images
	for i, job := range q.Jobs {
		title := fmt.Sprintf("%s. %s. %s", string(job.Title), string(job.Location), string(job.Phone))
		content := string(job.Description)
		author := "@tabyshmak.ru"
		resultName := fmt.Sprintf("results-png/result%d.png", i)
		if err := utils.CreateInstagramPost(title, content, author, resultName); err != nil {
			fmt.Println("Cannot create an image. Error:", err)
			continue
		}

		jpgFilename := fmt.Sprintf("results-jpg/result%d.jpg", i)
		if err := utils.PNGtoJPG(resultName, jpgFilename); err != nil {
			fmt.Println("Cannot convert to jpg. Error:", err)
			continue
		}
		time.Sleep(4 * time.Second)
	}

	// Login to instagram
	insta := goinsta.New("tabyshmak.ru", "firdavs13")
	if err := insta.Login(); err != nil {
		log.Fatal(err)
		return
	}
	defer insta.Logout()

	// Upload photos to instagram
	fmt.Println("Uploading...")
	for i, job := range q.Jobs {
		file, err := os.Open(fmt.Sprintf("results-jpg/result%d.jpg", i))
		if err != nil {
			fmt.Println("Cannot read file. Error:", err)
			continue
		}
		caption := fmt.Sprintf("%s, %s %s", string(job.Title), string(job.Phone), "#работа #москва #жумуш #jumush #работамосква #жумушмосква #москважумуш #жердеш")
		_, err = insta.UploadPhoto(file, caption, 1, 1)
		if err != nil {
			fmt.Println("Cannot upload. Error:", err)
			continue
		}
		time.Sleep(4 * time.Second)
	}
}
