package controller

import (
	"bytes"
	"fmt"
	"go-sqlx-gin/db_client"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
)

type Profile struct {
	UUID        string `json:"uuid"`
	UserID      string `json:"user_id"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

func CreateProfile(c *gin.Context) {
	u, err := uuid.NewRandom()
	if err != nil {
		fmt.Println(err)
		return
	}
	uu := u.String()
	var profile Profile
	c.BindJSON(&profile)
	ins, err := db_client.DBClient.Prepare("INSERT INTO profile(uuid, user_id ,description, image) VALUES(?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	ins.Exec(uu, profile.UserID, profile.Description, profile.Image)

	c.JSON(http.StatusOK, gin.H{"uuid": uu})
}

func EditProfile(c *gin.Context) {
	creds := credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, token)
	cfg := aws.NewConfig().WithRegion("ap-northeast-1").WithCredentials(creds)
	svc := s3.New(session.New(), cfg)

	form, _ := c.MultipartForm()

	files := form.File["images[]"]

	var imageNames []ImageName
	imageName := ImageName{}

	for _, file := range files {

		f, err := file.Open()

		if err != nil {
			log.Println(err)
		}

		defer f.Close()

		size := file.Size
		buffer := make([]byte, size)

		f.Read(buffer)
		fileBytes := bytes.NewReader(buffer)
		fileType := http.DetectContentType(buffer)
		path := "/media/" + file.Filename
		params := &s3.PutObjectInput{
			Bucket:        aws.String("article-s3-jpskgc"),
			Key:           aws.String(path),
			Body:          fileBytes,
			ContentLength: aws.Int64(size),
			ContentType:   aws.String(fileType),
		}
		resp, err := svc.PutObject(params)

		fmt.Printf("response %s", awsutil.StringValue(resp))

		imageName.NAME = file.Filename

		imageNames = append(imageNames, imageName)
	}

	c.JSON(http.StatusOK, imageNames)
}
