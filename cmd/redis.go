package cmd

import (
	"context"
	"fmt"
	"log"
	"unicode/utf8"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/cobra"
)

// RedisClient 是一个结构体，用于封装 Redis 客户端
type RedisClient struct {
	client *redis.ClusterClient
}

// NewRedisClient 创建一个新的 Redis 客户端
func NewRedisClient(addr, password string, db int) *RedisClient {
	rdb := redis.NewFailoverClusterClient(&redis.FailoverOptions{
		MasterName:     "mymaster",
		SentinelAddrs:  []string{addr},
		RouteByLatency: true,
		Password:       password,
	})

	return &RedisClient{client: rdb}
}

// Ping 检查 Redis 是否连接成功
func (r *RedisClient) Ping(ctx context.Context) error {
	_, err := r.client.Ping(ctx).Result()
	return err
}

// Close 关闭 Redis 客户端连接
func (r *RedisClient) Close() error {
	return r.client.Close()
}

var RedisCmd = &cobra.Command{
	Use: "redis",
	Run: func(cmd *cobra.Command, args []string) {
		input, _ := cmd.Flags().GetString("input")
		bytesCount := utf8.RuneCountInString(input)

		fmt.Printf("输入的字符串 '%s' 的字节大小为: %d\n", input, bytesCount)

		// 从命令行标志获取 Redis 配置
		addr, _ := cmd.Flags().GetString("addr")
		password, _ := cmd.Flags().GetString("password")
		db, _ := cmd.Flags().GetInt("db")

		// 创建 Redis 客户端
		redisClient := NewRedisClient(addr, password, db)

		// 创建一个上下文
		ctx := context.Background()

		// 检查连接是否成功
		err := redisClient.Ping(ctx)
		if err != nil {
			log.Fatalf("无法连接到 Redis: %v", err)
		} else {
			fmt.Println("成功连接到 Redis")
		}

		// 示例：设置一个键值对
		err = redisClient.client.Set(ctx, "key", "value", 0).Err()
		if err != nil {
			log.Fatalf("设置键值对失败: %v", err)
		} else {
			fmt.Println("键值对设置成功")
		}

		// 关闭 Redis 客户端连接
		err = redisClient.Close()
		if err != nil {
			log.Fatalf("关闭 Redis 连接失败: %v", err)
		} else {
			fmt.Println("成功关闭 Redis 连接")
		}
	},
}

func init() {
	RedisCmd.Flags().StringP("input", "i", "", "输入字符串")
	RedisCmd.Flags().StringP("addr", "a", "localhost:26379", "Redis Sentinel 地址")
	RedisCmd.Flags().StringP("password", "p", "", "Redis 密码")
	RedisCmd.Flags().IntP("db", "d", 0, "Redis 数据库编号")
}
