package idgenerator

import (
	"fmt"

	"github.com/bwmarrin/snowflake"
)

type Generator interface {
	Generate() int64
}

type SnowflakeGenerator struct {
	node *snowflake.Node
}

func NewSnowflakeGenerator(nodeID int64, epoch int64) (*SnowflakeGenerator, error) {
	snowflake.Epoch = epoch
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		return nil, fmt.Errorf("idgen: failed to create node: %w", err)
	}
	return &SnowflakeGenerator{node: node}, nil
}

func (g *SnowflakeGenerator) Generate() int64 {
	return g.node.Generate().Int64()
}
