package tui

import (
	"log"

	"github.com/redis/go-redis/v9"
)

func GenerateHashType(client *redis.Client, node *Node) *RedisHash {
	out, err := client.HGetAll(ctx, node.FullKey).Result()
	if err != nil {
		log.Fatal(err)
	}
	return &RedisHash{
		RedisType: HASH,
		Data:      out,
	}
}

func GenerateStringType(client *redis.Client, node *Node) *RedisString {
	out, err := client.Get(ctx, node.FullKey).Result()
	if err != nil {
		log.Fatal(err)
	}
	return &RedisString{
		RedisType: STRING,
		Data:      out,
	}
}

func GenerateListType(client *redis.Client, node *Node) *RedisList {
	out, err := client.LRange(ctx, node.FullKey, 0, -1).Result()
	if err != nil {
		log.Fatal(err)
	}
	return &RedisList{
		RedisType: LIST,
		Data:      out,
	}
}

func GenerateSetType(client *redis.Client, node *Node) *RedisSet {
	out, err := client.SMembers(ctx, node.FullKey).Result()
	if err != nil {
		log.Fatal(err)
	}
	return &RedisSet{
		RedisType: LIST,
		Data:      out,
	}
}

func GenerateStreamType(client *redis.Client, node *Node) *RedisStream {
	out, err := client.XRange(ctx, node.FullKey, "-", "+").Result()
	if err != nil {
		log.Fatal(err)
	}
	return &RedisStream{
		RedisType: LIST,
		Data:      out,
	}
}

func GenerateZSetType(client *redis.Client, node *Node) *RedisZSet {
	out, err := client.ZRange(ctx, node.FullKey, 0, -1).Result()
	if err != nil {
		log.Fatal(err)
	}
	return &RedisZSet{
		RedisType: LIST,
		Data:      out,
	}
}
