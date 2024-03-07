package main

import (
	"ClipCon-Server/routes"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// 加载.env文件
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 启动Gin服务器
	router := routes.SetupRouter()
	router.Run(":" + os.Getenv("PORT"))
}
