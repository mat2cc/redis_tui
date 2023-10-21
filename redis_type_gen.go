package main

import (
	"log"

	"github.com/mat2cc/redis_tui/redis_type"
	"github.com/redis/go-redis/v9"
)

func GenerateHashType(client *redis.Client, node *Node) *redis_type.RedisHash {
	out, err := client.HGetAll(ctx, node.FullKey).Result()
	if err != nil {
		log.Fatal(err)
	}
	return &redis_type.RedisHash{
		RedisType: redis_type.HASH,
		Data:      out,
	}
}

func GenerateStringType(client *redis.Client, node *Node) *redis_type.RedisString {
	out, err := client.Get(ctx, node.FullKey).Result()
	if err != nil {
		log.Fatal(err)
	}
	return &redis_type.RedisString{
		RedisType: redis_type.STRING,
		Data:      out,
	}
}

func GenerateListType(client *redis.Client, node *Node) *redis_type.RedisList {
	out, err := client.LRange(ctx, node.FullKey, 0, -1).Result()
	if err != nil {
		log.Fatal(err)
	}
	return &redis_type.RedisList{
		RedisType: redis_type.LIST,
		Data:      out,
	}
}

func GenerateSetType(client *redis.Client, node *Node) *redis_type.RedisSet {
	out, err := client.SMembers(ctx, node.FullKey).Result()
	if err != nil {
		log.Fatal(err)
	}
	return &redis_type.RedisSet{
		RedisType: redis_type.LIST,
		Data:      out,
	}
}

func GenerateStreamType(client *redis.Client, node *Node) *redis_type.RedisStream {
	out, err := client.XRange(ctx, node.FullKey, "-", "+").Result()
	if err != nil {
		log.Fatal(err)
	}
	return &redis_type.RedisStream{
		RedisType: redis_type.LIST,
		Data:      out,
	}
}

func GenerateZSetType(client *redis.Client, node *Node) *redis_type.RedisZSet {
	out, err := client.ZRange(ctx, node.FullKey, 0, -1).Result()
	if err != nil {
		log.Fatal(err)
	}
	return &redis_type.RedisZSet{
		RedisType: redis_type.LIST,
		Data:      out,
	}
}
