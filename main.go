package main

import (
    "fmt"
    "log"
    "time"
    "github.com/qwerty-dvorak/goaswses/helper"
    "github.com/qwerty-dvorak/goaswses/myses"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ses"
    "github.com/joho/godotenv"
)

func LoadEnv() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }
}

func main() {
    LoadEnv()

    records, err := helper.ReadCsvFile("list.csv")
    if err != nil {
        fmt.Println(err)
        return
    }

    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("ap-south-1"),
    })

    if err != nil {
        log.Fatalf("Failed to create AWS session: %v", err)
    }

    sesClient := ses.New(sess)

    for _, email := range records {
        data := map[string]interface{}{
            "Name": "Kairav",
        }
        myses.SendSESEmail(email, sesClient, data)
        time.Sleep(75 * time.Millisecond)
    }
    fmt.Println(records)
}