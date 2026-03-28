package templates

// RedisGo is the template for pkg/redis/redis.go
const RedisGo = `package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// ErrNotFound is returned when a key does not exist.
var ErrNotFound = errors.New("redis: key not found")

// Client wraps the go-redis client.
type Client struct {
	rdb *redis.Client
}

// New creates a new Redis client and pings with a 5s timeout.
func New(addr, password string, db int) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &Client{rdb: rdb}, nil
}

// Set stores a value with the given key and TTL.
func (c *Client) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if err := c.rdb.Set(ctx, key, value, ttl).Err(); err != nil {
		return fmt.Errorf("redis set %q: %w", key, err)
	}
	return nil
}

// Get retrieves a value by key. Returns ErrNotFound if the key does not exist.
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	val, err := c.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("redis get %q: %w", key, err)
	}
	return val, nil
}

// Delete removes a key.
func (c *Client) Delete(ctx context.Context, key string) error {
	if err := c.rdb.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis del %q: %w", key, err)
	}
	return nil
}

// Exists checks whether a key exists.
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("redis exists %q: %w", key, err)
	}
	return n > 0, nil
}

// Close closes the Redis connection.
func (c *Client) Close() error {
	return c.rdb.Close()
}
`

// KafkaGo is the template for pkg/kafka/kafka.go
const KafkaGo = `package kafka

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
)

// Producer wraps a franz-go Kafka client for producing messages.
type Producer struct {
	client *kgo.Client
}

// NewProducer creates a new Kafka producer.
func NewProducer(brokers []string) (*Producer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka producer: %w", err)
	}
	return &Producer{client: client}, nil
}

// Publish sends a message to the given topic.
func (p *Producer) Publish(ctx context.Context, topic string, key, value []byte) error {
	record := &kgo.Record{
		Topic: topic,
		Key:   key,
		Value: value,
	}

	result := p.client.ProduceSync(ctx, record)
	if err := result.FirstErr(); err != nil {
		return fmt.Errorf("kafka publish to %q: %w", topic, err)
	}
	return nil
}

// Close closes the producer.
func (p *Producer) Close() {
	p.client.Close()
}

// Consumer wraps a franz-go Kafka client for consuming messages.
type Consumer struct {
	client *kgo.Client
}

// NewConsumer creates a new Kafka consumer.
func NewConsumer(brokers []string, groupID string, topics ...string) (*Consumer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(groupID),
		kgo.ConsumeTopics(topics...),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka consumer: %w", err)
	}
	return &Consumer{client: client}, nil
}

// Poll continuously polls for messages and calls handler for each record.
// It respects ctx cancellation.
func (c *Consumer) Poll(ctx context.Context, handler func(*kgo.Record) error) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		fetches := c.client.PollFetches(ctx)
		if fetches.IsClientClosed() {
			return nil
		}

		if errs := fetches.Errors(); len(errs) > 0 {
			return fmt.Errorf("kafka poll error: %v", errs[0])
		}

		fetches.EachRecord(func(record *kgo.Record) {
			if err := handler(record); err != nil {
				// log or handle per-record errors here
				_ = err
			}
		})
	}
}

// Close closes the consumer.
func (c *Consumer) Close() {
	c.client.Close()
}
`

// NatsGo is the template for pkg/nats/nats.go
const NatsGo = `package nats

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

// Options configures the NATS client.
type Options struct {
	Name          string
	URL           string
	Timeout       time.Duration
	MaxReconnects int
}

// Client wraps the nats.go connection.
type Client struct {
	conn *nats.Conn
}

// New creates a new NATS client with the given options.
func New(opts Options) (*Client, error) {
	if opts.Timeout == 0 {
		opts.Timeout = 5 * time.Second
	}
	if opts.MaxReconnects == 0 {
		opts.MaxReconnects = 10
	}

	nc, err := nats.Connect(
		opts.URL,
		nats.Name(opts.Name),
		nats.Timeout(opts.Timeout),
		nats.MaxReconnects(opts.MaxReconnects),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			if err != nil {
				// log disconnect error
				_ = err
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			// log reconnect
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("nats connect: %w", err)
	}

	return &Client{conn: nc}, nil
}

// Publish publishes a message to the subject.
func (c *Client) Publish(subject string, data []byte) error {
	if err := c.conn.Publish(subject, data); err != nil {
		return fmt.Errorf("nats publish to %q: %w", subject, err)
	}
	return nil
}

// Subscribe subscribes to the subject and calls handler for each message.
func (c *Client) Subscribe(subject string, handler nats.MsgHandler) (*nats.Subscription, error) {
	sub, err := c.conn.Subscribe(subject, handler)
	if err != nil {
		return nil, fmt.Errorf("nats subscribe to %q: %w", subject, err)
	}
	return sub, nil
}

// QueueSubscribe subscribes to a queue group.
func (c *Client) QueueSubscribe(subject, queue string, handler nats.MsgHandler) (*nats.Subscription, error) {
	sub, err := c.conn.QueueSubscribe(subject, queue, handler)
	if err != nil {
		return nil, fmt.Errorf("nats queue subscribe to %q / %q: %w", subject, queue, err)
	}
	return sub, nil
}

// Close drains the connection and closes it.
func (c *Client) Close() error {
	if err := c.conn.Drain(); err != nil {
		return fmt.Errorf("nats drain: %w", err)
	}
	return nil
}
`

// DBTemplate returns the pkg/db/db.go template for the given driver keyword.
// Returns an empty string for "none" or unknown drivers.
func DBTemplate(db string) string {
	switch db {
	case "postgres":
		return PostgresGo
	case "mysql":
		return MySQLGo
	case "sqlite":
		return SQLiteGo
	case "mongo":
		return MongoGo
	case "gorm":
		return GormGo
	default:
		return ""
	}
}

// PostgresGo is the template for pkg/db/db.go (pgx driver).
const PostgresGo = `package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB wraps a pgx connection pool.
type DB struct {
	pool *pgxpool.Pool
}

// New opens a connection pool and verifies connectivity with a 5s ping.
func New(ctx context.Context, dsn string) (*DB, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("db: parse config: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(pingCtx, cfg)
	if err != nil {
		return nil, fmt.Errorf("db: connect: %w", err)
	}

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("db: ping: %w", err)
	}

	return &DB{pool: pool}, nil
}

// Pool returns the underlying pgxpool.Pool for raw query access.
func (d *DB) Pool() *pgxpool.Pool { return d.pool }

// Close closes the connection pool.
func (d *DB) Close() { d.pool.Close() }
`

// MySQLGo is the template for pkg/db/db.go (go-sql-driver/mysql).
const MySQLGo = `package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DB wraps the standard *sql.DB for MySQL.
type DB struct {
	sql *sql.DB
}

// New opens a MySQL connection and pings with a 5s timeout.
func New(dsn string) (*DB, error) {
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("db: open: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("db: ping: %w", err)
	}

	return &DB{sql: sqlDB}, nil
}

// SQL returns the underlying *sql.DB for raw query access.
func (d *DB) SQL() *sql.DB { return d.sql }

// Close closes the database connection.
func (d *DB) Close() error { return d.sql.Close() }
`

// SQLiteGo is the template for pkg/db/db.go (modernc.org/sqlite).
const SQLiteGo = `package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// DB wraps the standard *sql.DB for SQLite.
type DB struct {
	sql *sql.DB
}

// New opens a SQLite database file and pings it.
func New(path string) (*DB, error) {
	sqlDB, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("db: open: %w", err)
	}

	// SQLite performs best with a single writer connection.
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetConnMaxLifetime(0)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("db: ping: %w", err)
	}

	return &DB{sql: sqlDB}, nil
}

// SQL returns the underlying *sql.DB for raw query access.
func (d *DB) SQL() *sql.DB { return d.sql }

// Close closes the database file.
func (d *DB) Close() error { return d.sql.Close() }
`

// MongoGo is the template for pkg/db/db.go (mongo-driver).
const MongoGo = `package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB wraps the mongo.Client.
type DB struct {
	client *mongo.Client
}

// New connects to MongoDB and verifies the connection with a 5s ping.
func New(ctx context.Context, uri string) (*DB, error) {
	opts := options.Client().ApplyURI(uri).
		SetConnectTimeout(5 * time.Second).
		SetServerSelectionTimeout(5 * time.Second)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("db: connect: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		_ = client.Disconnect(ctx)
		return nil, fmt.Errorf("db: ping: %w", err)
	}

	return &DB{client: client}, nil
}

// Client returns the underlying *mongo.Client.
func (d *DB) Client() *mongo.Client { return d.client }

// Collection is a convenience helper to get a collection.
func (d *DB) Collection(database, collection string) *mongo.Collection {
	return d.client.Database(database).Collection(collection)
}

// Close disconnects from MongoDB.
func (d *DB) Close(ctx context.Context) error {
	if err := d.client.Disconnect(ctx); err != nil {
		return fmt.Errorf("db: disconnect: %w", err)
	}
	return nil
}
`

// GormGo is the template for pkg/db/db.go (GORM + postgres driver).
const GormGo = `package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB wraps a gorm.DB instance.
type DB struct {
	gorm *gorm.DB
}

// New opens a GORM connection (postgres driver) from the given DSN.
func New(dsn string) (*DB, error) {
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("db: open: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("db: get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("db: ping: %w", err)
	}

	return &DB{gorm: gormDB}, nil
}

// GORM returns the underlying *gorm.DB for query building.
func (d *DB) GORM() *gorm.DB { return d.gorm }

// Close closes the underlying database connection.
func (d *DB) Close() error {
	sqlDB, err := d.gorm.DB()
	if err != nil {
		return fmt.Errorf("db: get underlying sql.DB: %w", err)
	}
	return sqlDB.Close()
}
`

// SwaggerYAML is the template for docs/swagger.yaml
const SwaggerYAML = `openapi: "3.0.3"
info:
  title: {{.ProjectName}} API
  version: "1.0.0"
  description: API documentation for {{.ProjectName}}
servers:
  - url: http://localhost:8080/api/v1
    description: Local development
paths:
  /health:
    get:
      summary: Health check
      tags:
        - Health
      responses:
        "200":
          description: Service is healthy
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
                    example: ok
                  app:
                    type: string
                    example: {{.ProjectName}}
                  env:
                    type: string
                    example: development
`
