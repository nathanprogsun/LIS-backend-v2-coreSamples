package common

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	goredis "github.com/go-redsync/redsync/v4/redis/goredis/v8"
	capi "github.com/hashicorp/consul/api"
	"github.com/hibiken/asynq"
	"time"
)

type RedisConfig struct {
	MasterName    string   `json:"mastername,omitempty"`
	SentinelAddrs []string `json:"sentinel_addrs,omitempty"` // it may just be the redis server address, if we are in dev
	Password      string   `json:"password,omitempty"`
}

type RedisClient struct {
	RedisReadClient  *redis.Client
	RedisWriteClient *redis.Client
}

func GetRedisConfigFromConsul(client *capi.Client, prefix string, key string) *RedisConfig {
	redisConfig := &RedisConfig{}
	val, _, err := client.KV().Get(prefix+"/"+key, nil)
	if err != nil {
		Fatal(err)
	}
	err = json.Unmarshal(val.Value, redisConfig)
	if err != nil {
		Fatal(err)
	}
	return redisConfig
}

func NewRedisClient(readClient *redis.Client, writeClient *redis.Client) *RedisClient {
	return &RedisClient{
		RedisReadClient:  readClient,
		RedisWriteClient: writeClient,
	}
}

func GetRedisReadWriteClient(config *RedisConfig) *RedisClient {
	if Env.RunEnv == DevEnv {
		client := redis.NewClient(&redis.Options{
			Addr:     config.SentinelAddrs[0],
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		return &RedisClient{
			RedisReadClient:  client,
			RedisWriteClient: client,
		}
	}
	wopts := &redis.FailoverOptions{
		MasterName:    config.MasterName,
		SentinelAddrs: config.SentinelAddrs,
	}
	ropts := &redis.FailoverOptions{
		MasterName:    config.MasterName,
		SentinelAddrs: config.SentinelAddrs,
		SlaveOnly:     true,
	}
	if Env.RunEnv == AksProductionEnv || Env.RunEnv == AksStagingEnv {
		wopts.Password = config.Password
		wopts.SentinelPassword = config.Password
		ropts.Password = config.Password
		ropts.SentinelPassword = config.Password
	}
	writeClient := redis.NewFailoverClient(wopts)
	readClient := redis.NewFailoverClient(ropts)
	return &RedisClient{
		RedisReadClient:  readClient,
		RedisWriteClient: writeClient,
	}
}

func (c *RedisClient) Close() error {
	if c.RedisReadClient != nil {
		if err := c.RedisReadClient.Close(); err != nil {
			return err
		}
	}

	if c.RedisWriteClient != nil {
		if err := c.RedisWriteClient.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (c *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return c.RedisWriteClient.Set(ctx, key, value, expiration)
}

func (c *RedisClient) SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return c.RedisWriteClient.SetEX(ctx, key, value, expiration)
}

func (c *RedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return c.RedisReadClient.Get(ctx, key)
}

func (c *RedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return c.RedisWriteClient.Del(ctx, keys...)
}

func GetAsynqClient(config *RedisConfig) *asynq.Client {
	if Env.RunEnv == DevEnv {
		opts := asynq.RedisClientOpt{
			Addr: config.SentinelAddrs[0],
		}
		return asynq.NewClient(opts)
	}
	opts := asynq.RedisFailoverClientOpt{
		MasterName:    config.MasterName,
		SentinelAddrs: config.SentinelAddrs,
	}
	if Env.RunEnv == AksProductionEnv || Env.RunEnv == AksStagingEnv {
		opts.Password = config.Password
		opts.SentinelPassword = config.Password
	}
	return asynq.NewClient(opts)
}

func GetAsyncServer(config *RedisConfig) *asynq.Server {
	if Env.RunEnv == DevEnv {
		opts := asynq.RedisClientOpt{
			Addr: config.SentinelAddrs[0],
		}
		return asynq.NewServer(opts, asynq.Config{
			Concurrency: 2,
			Queues: map[string]int{
				Env.AsynqQueueName: 1,
			},
		})
	}
	opts := asynq.RedisFailoverClientOpt{
		MasterName:    config.MasterName,
		SentinelAddrs: config.SentinelAddrs,
	}
	if Env.RunEnv == AksProductionEnv || Env.RunEnv == AksStagingEnv {
		opts.Password = config.Password
		opts.SentinelPassword = config.Password
	}
	return asynq.NewServer(opts, asynq.Config{
		Concurrency: 2,
		Queues: map[string]int{
			Env.AsynqQueueName: 1,
		},
	})
}

func GetRedisSync(client *redis.Client) *redsync.Redsync {
	pool := goredis.NewPool(client)
	rs := redsync.New(pool)
	return rs
}
