package cache

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/go-redis/redis/v8"
    "golang.org/x/net/context"
	"filestore/internal/models"
     // Ensure this is used
)

// RedisClient is the Redis client instance
var RedisClient *redis.Client
var ctx = context.Background()

// InitRedis initializes the Redis client
func InitRedis() {
    RedisClient = redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_ADDR"),
        Password: os.Getenv("REDIS_PASSWORD"), // no password set
        DB:       0,                           // use default DB
    })

    _, err := RedisClient.Ping(ctx).Result()
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }

    log.Println("Connected to Redis")
}

// GetFilesFromCache retrieves file metadata from Redis cache
func GetFilesFromCache(userID uint) ([]models.FileMetadata, bool) {
    key := fmt.Sprintf("user:%d:files", userID)
    val, err := RedisClient.Get(ctx, key).Result()
    if err == redis.Nil {
        return nil, false
    } else if err != nil {
        log.Printf("Failed to get files from Redis: %v", err)
        return nil, false
    }

    var File []models.FileMetadata
    err = json.Unmarshal([]byte(val), &File)
    if err != nil {
        log.Printf("Failed to unmarshal files data: %v", err)
        return nil, false
    }

    return File, true
}

// SetFilesInCache stores file metadata in Redis cache
func SetFilesInCache(userID uint, File []models.FileMetadata) {
    key := fmt.Sprintf("user:%d:files", userID)
    data, err := json.Marshal(File)
    if err != nil {
        log.Printf("Failed to marshal files data: %v", err)
        return
    }

    err = RedisClient.Set(ctx, key, data, 5*time.Minute).Err()
    if err != nil {
        log.Printf("Failed to set files in Redis: %v", err)
    }
}