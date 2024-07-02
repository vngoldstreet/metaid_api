package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var rdb *redis.Client
var mongoCLient *mongo.Client

func SetupLogger() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{}) // Định dạng log JSON (tùy chọn)

	// Mở tệp log để lưu
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	// Cấu hình logger cho Gin
	gin.SetMode(gin.ReleaseMode) // Đặt chế độ Gin (ReleaseMode hoặc DebugMode)
	gin.DefaultWriter = log.Writer()
}

func ConnectToRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "192.168.50.131:6379", // Thay thế bằng địa chỉ Redis thực tế
		Password: "",                    // Mật khẩu (nếu có)
		DB:       0,                     // Chọn cơ sở dữ liệu
	})
	ping, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("Ping err: %v\n", err)
	}
	fmt.Printf("ping: %v\n", ping)
}

func GetRedisClient() *redis.Client {
	return rdb
}

func GetMongoClient() *mongo.Client {
	return mongoCLient
}

func ConnectToMongoDB() {
	clientOptions := options.Client().ApplyURI("mongodb://192.168.50.131:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	mongoCLient = client
}

// Ta muốn sinh ra 1 mảng bao gồm tất cả các ký tự có thể kết hợp
func Logger() gin.HandlerFunc {
	logger := logrus.New()
	return func(c *gin.Context) {
		// Bắt đầu thời gian
		startTime := time.Now()

		// Xử lý request
		c.Next()

		// Kết thúc thời gian
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)

		// Trạng thái
		statusCode := c.Writer.Status()

		// Phương thức
		reqMethod := c.Request.Method

		// Đường dẫn
		reqUri := c.Request.RequestURI

		// Địa chỉ IP
		clientIP := c.ClientIP()

		// Log chi tiết
		logger.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"method":       reqMethod,
			"uri":          reqUri,
		}).Info()
	}
}
