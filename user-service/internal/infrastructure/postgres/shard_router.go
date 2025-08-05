package postgres

import (
	"fmt"
	"hash/fnv"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Wrestler094/scalable-ecommerce-platform/user-service/internal/config"
)

type ShardRouter struct {
	shards []*sqlx.DB
	count  int
}

func NewShardRouter(pgShards []config.PGShard) (*ShardRouter, error) {
	const op = "postgres.NewShardRouter"

	var shards []*sqlx.DB

	for i, shard := range pgShards {
		db, err := sqlx.Open("postgres", shard.DSN())
		if err != nil {
			return nil, fmt.Errorf("%s: open shard %d (%s): %w", op, i, shard.Name, err)
		}
		if err := db.Ping(); err != nil {
			return nil, fmt.Errorf("%s: ping shard %d (%s): %w", op, i, shard.Name, err)
		}
		shards = append(shards, db)
	}

	return &ShardRouter{
		shards: shards,
		count:  len(shards),
	}, nil
}

// GetShardByUserID returns sqlx.DB based on userID.
func (r *ShardRouter) GetShardByUserID(userID int64) *sqlx.DB {
	h := fnv.New32()
	h.Write([]byte(strconv.FormatInt(userID, 10)))
	index := int(h.Sum32() % uint32(r.count))
	return r.shards[index]
}

// Close all shard connections.
func (r *ShardRouter) Close() error {
	const op = "postgres.ShardRouter.Close"

	for i, db := range r.shards {
		if err := db.Close(); err != nil {
			return fmt.Errorf("%s: failed to close shard %d: %w", op, i, err)
		}
	}

	return nil
}
