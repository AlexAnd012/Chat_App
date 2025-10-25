package storage

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	cli *redis.Client
}

const keyMessages = "chat:messages"

func NewRedisRepo(addr string) (*RedisRepo, error) {
	cli := redis.NewClient(&redis.Options{Addr: addr})
	ctx := context.Background()
	if err := cli.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}
	return &RedisRepo{cli: cli}, nil
}

func (r *RedisRepo) Close() error { return r.cli.Close() }

func (r *RedisRepo) Ping(ctx context.Context) error { return r.cli.Ping(ctx).Err() }

func (r *RedisRepo) AppendMessage(ctx context.Context, jsonBytes []byte) error {
	return r.cli.LPush(ctx, keyMessages, jsonBytes).Err()
}

// Ограничиваем длину списка, оставляем только первые keep элементов
func (r *RedisRepo) TrimHistory(ctx context.Context, keep int) error {
	return r.cli.LTrim(ctx, keyMessages, 0, int64(keep-1)).Err()
}

// Чтение из Redis последних N сообщений
func (r *RedisRepo) RecentMessages(ctx context.Context, limit int) ([]string, error) {
	vals, err := r.cli.LRange(ctx, keyMessages, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}
	return vals, nil
}
