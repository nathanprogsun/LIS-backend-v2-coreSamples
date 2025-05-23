package common

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	// "strconv" // Not strictly needed for RedisConfig part, but for GetCurrent if parsing port separately
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	goredis "github.com/go-redsync/redsync/v4/redis/goredis/v8"
	capi "github.com/hashicorp/consul/api"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

// RedisConfig struct definition (ensure it's present)
type RedisConfig struct {
	MasterName    string   `json:"mastername,omitempty"`
	SentinelAddrs []string `json:"sentinel_addrs,omitempty"` // it may just be the redis server address, if we are in dev or dev_docker_compose
	Password      string   `json:"password,omitempty"`
}

// RedisClient struct definition (ensure it's present)
type RedisClient struct {
	RedisReadClient  *redis.Client
	RedisWriteClient *redis.Client
}

// GetRedisConfigFromConsul retrieves Redis configuration from Consul.
// Updated to return error for better error handling.
func GetRedisConfigFromConsul(client *capi.Client, prefix string, key string) (*RedisConfig, error) {
	val, _, err := client.KV().Get(prefix+"/"+key, nil)
	if err != nil {
		return nil, fmt.Errorf("Consul KV().Get failed for %s/%s: %w", prefix, key, err)
	}
	if val == nil || val.Value == nil {
		return nil, fmt.Errorf("Redis config key not found in Consul: %s/%s. Value is nil", prefix, key)
	}
	config := &RedisConfig{}
	err = json.Unmarshal(val.Value, config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal Redis config from Consul for key %s/%s: %w", prefix, key, err)
	}
	return config, nil
}

// GetCurrentRedisConfig determines if Redis config should be loaded from ENV or Consul.
func GetCurrentRedisConfig(client *capi.Client, consulPrefix string) (*RedisConfig, error) {
	if Env.RunEnv == DevDockerComposeEnv {
		Info("Running in dev_docker_compose mode, attempting to load Redis config from ENV variables.")
		redisAddr := os.Getenv("REDIS_ADDR") // Expected format "host:port"
		redisPassword := os.Getenv("REDIS_PASSWORD")

		if redisAddr != "" {
			InfoFields("Loaded Redis config from ENV", zap.String("address", redisAddr))
			return &RedisConfig{
				SentinelAddrs: []string{redisAddr}, // Store single address here for standalone mode
				Password:      redisPassword,
				MasterName:    "", // Not used in standalone mode
			}, nil
		}
		Warn("dev_docker_compose mode: REDIS_ADDR ENV VAR not set. Attempting fallback to Consul.")
	}

	var consulKey string = "redisSentinel" // Default key
	// Add logic if different envs use different redis keys in consul, e.g.:
	// if Env.RunEnv == SpecificEnvForDifferentRedisKey { consulKey = "otherRedisKey" }

	InfoFields("Loading Redis config from Consul", zap.String("key", consulKey), zap.String("consulPrefix", consulPrefix))
	cfg, err := GetRedisConfigFromConsul(client, consulPrefix, consulKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get Redis config from Consul (key: %s): %w", consulKey, err)
	}
	return cfg, nil
}

func NewRedisClient(readClient *redis.Client, writeClient *redis.Client) *RedisClient {
	return &RedisClient{
		RedisReadClient:  readClient,
		RedisWriteClient: writeClient,
	}
}

func GetRedisReadWriteClient(config *RedisConfig) *RedisClient {
	if config == nil {
		Fatal(fmt.Errorf("received nil RedisConfig in GetRedisReadWriteClient. Ensure config is loaded correctly"))
		return nil
	}

	if Env.RunEnv == DevEnv || Env.RunEnv == DevDockerComposeEnv {
		if len(config.SentinelAddrs) == 0 || config.SentinelAddrs[0] == "" {
			Fatal(fmt.Errorf("Redis address (config.SentinelAddrs[0]) is empty for %s mode", Env.RunEnv))
			return nil // Or handle error more gracefully
		}
		addr := config.SentinelAddrs[0]
		InfoFields("Initializing standalone Redis client", zap.String("address", addr), zap.String("runEnv", Env.RunEnv))

		client := redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: config.Password,
			DB:       0,
		})
		return &RedisClient{ // Assuming NewRedisClient helper is not strictly needed here
			RedisReadClient:  client,
			RedisWriteClient: client,
		}
	}

	// Existing Sentinel logic for other environments
	InfoFields("Initializing Redis client with Sentinel",
		zap.Strings("sentinelAddrs", config.SentinelAddrs),
		zap.String("masterName", config.MasterName),
		zap.String("runEnv", Env.RunEnv))

	// Ensure MasterName and SentinelAddrs are valid for FailoverClient
	if config.MasterName == "" || len(config.SentinelAddrs) == 0 {
		Fatal(fmt.Errorf("Redis Sentinel config incomplete for %s mode: MasterName or SentinelAddrs is empty", Env.RunEnv))
		return nil
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

	if Env.RunEnv == AksProductionEnv || Env.RunEnv == AksStagingEnv { // Assuming these constants exist
		wopts.Password = config.Password
		wopts.SentinelPassword = config.Password
		ropts.Password = config.Password
		ropts.SentinelPassword = config.Password
	}
	writeClient := redis.NewFailoverClient(wopts)
	readClient := redis.NewFailoverClient(ropts)
	return &RedisClient{ // Assuming NewRedisClient helper is not strictly needed here
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
	if config == nil {
		Fatal(fmt.Errorf("received nil RedisConfig in GetAsynqClient"))
		return nil
	}
	if Env.RunEnv == DevEnv || Env.RunEnv == DevDockerComposeEnv {
		if len(config.SentinelAddrs) == 0 || config.SentinelAddrs[0] == "" {
			Fatal(fmt.Errorf("Redis address (config.SentinelAddrs[0]) is empty for Asynq in %s mode", Env.RunEnv))
			return nil
		}
		addr := config.SentinelAddrs[0]
		InfoFields("Initializing Asynq client (standalone Redis)", zap.String("address", addr), zap.String("runEnv", Env.RunEnv))
		opts := asynq.RedisClientOpt{
			Addr:     addr,
			Password: config.Password,
		}
		return asynq.NewClient(opts)
	}

	InfoFields("Initializing Asynq client (Sentinel Redis)", zap.Strings("sentinelAddrs", config.SentinelAddrs), zap.String("masterName", config.MasterName))
	// Ensure MasterName and SentinelAddrs are valid for FailoverClientOpt
	if config.MasterName == "" || len(config.SentinelAddrs) == 0 {
		Fatal(fmt.Errorf("Redis Sentinel config incomplete for Asynq in %s mode: MasterName or SentinelAddrs is empty", Env.RunEnv))
		return nil
	}
	opts := asynq.RedisFailoverClientOpt{
		MasterName:    config.MasterName,
		SentinelAddrs: config.SentinelAddrs,
		Password:      config.Password, // Assuming password applies to sentinel connections too if set
	}
	if Env.RunEnv == AksProductionEnv || Env.RunEnv == AksStagingEnv {
		opts.SentinelPassword = config.Password // If sentinel has a password too
	}
	return asynq.NewClient(opts)
}

func GetAsyncServer(config *RedisConfig) *asynq.Server {
	if config == nil {
		Fatal(fmt.Errorf("received nil RedisConfig in GetAsyncServer"))
		return nil
	}
	asynqConf := asynq.Config{
		Concurrency: 2, // Or from Env Var
		Queues: map[string]int{
			Env.AsynqQueueName: 1, // Or from Env Var
		},
	}
	if Env.RunEnv == DevEnv || Env.RunEnv == DevDockerComposeEnv {
		if len(config.SentinelAddrs) == 0 || config.SentinelAddrs[0] == "" {
			Fatal(fmt.Errorf("Redis address (config.SentinelAddrs[0]) is empty for Asynq Server in %s mode", Env.RunEnv))
			return nil
		}
		addr := config.SentinelAddrs[0]
		InfoFields("Initializing Asynq server (standalone Redis)", zap.String("address", addr), zap.String("runEnv", Env.RunEnv))
		opts := asynq.RedisClientOpt{
			Addr:     addr,
			Password: config.Password,
		}
		return asynq.NewServer(opts, asynqConf)
	}

	InfoFields("Initializing Asynq server (Sentinel Redis)", zap.Strings("sentinelAddrs", config.SentinelAddrs), zap.String("masterName", config.MasterName))
	if config.MasterName == "" || len(config.SentinelAddrs) == 0 {
		Fatal(fmt.Errorf("Redis Sentinel config incomplete for Asynq Server in %s mode: MasterName or SentinelAddrs is empty", Env.RunEnv))
		return nil
	}
	opts := asynq.RedisFailoverClientOpt{
		MasterName:    config.MasterName,
		SentinelAddrs: config.SentinelAddrs,
		Password:      config.Password,
	}
	if Env.RunEnv == AksProductionEnv || Env.RunEnv == AksStagingEnv {
		opts.SentinelPassword = config.Password
	}
	return asynq.NewServer(opts, asynqConf)
}

func GetRedisSync(client *redis.Client) *redsync.Redsync {
	pool := goredis.NewPool(client)
	rs := redsync.New(pool)
	return rs
}
