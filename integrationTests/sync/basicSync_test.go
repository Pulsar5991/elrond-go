package sync

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-go/integrationTests"
	"github.com/ElrondNetwork/elrond-go/sharding"
	"github.com/stretchr/testify/assert"
)

var stepDelay = time.Second
var delayP2pBootstrap = time.Second * 2

func TestSyncWorksInShard_EmptyBlocksNoForks(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	maxShards := uint32(1)
	shardId := uint32(0)
	numNodesPerShard := 6

	advertiser := integrationTests.CreateMessengerWithKadDht(context.Background(), "")
	_ = advertiser.Bootstrap()
	advertiserAddr := integrationTests.GetConnectableAddress(advertiser)

	nodes := make([]*integrationTests.TestProcessorNode, numNodesPerShard+1)
	for i := 0; i < numNodesPerShard; i++ {
		nodes[i] = integrationTests.NewTestSyncNode(
			maxShards,
			shardId,
			shardId,
			advertiserAddr,
		)
	}

	metachainNode := integrationTests.NewTestSyncNode(
		maxShards,
		sharding.MetachainShardId,
		shardId,
		advertiserAddr,
	)
	nodes[numNodesPerShard] = metachainNode

	idxProposerShard0 := 0
	idxProposerMeta := numNodesPerShard
	idxProposers := []int{idxProposerShard0, idxProposerMeta}

	defer func() {
		_ = advertiser.Close()
		for _, n := range nodes {
			_ = n.Messenger.Close()
		}
	}()

	for _, n := range nodes {
		_ = n.Messenger.Bootstrap()
		_ = n.StartSync()
	}

	fmt.Println("Delaying for nodes p2p bootstrap...")
	time.Sleep(delayP2pBootstrap)

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	updateRound(nodes, round)
	nonce++

	numRoundsToTest := 5
	for i := 0; i < numRoundsToTest; i++ {
		integrationTests.ProposeBlock(nodes, idxProposers, round, nonce)

		time.Sleep(stepDelay)

		round = integrationTests.IncrementAndPrintRound(round)
		updateRound(nodes, round)
		nonce++
	}

	time.Sleep(stepDelay)

	testAllNodesHaveTheSameBlockHeightInBlockchain(t, nodes)
}

func testAllNodesHaveTheSameBlockHeightInBlockchain(t *testing.T, nodes []*integrationTests.TestProcessorNode) {
	expectedNonce := nodes[0].BlockChain.GetCurrentBlockHeader().GetNonce()
	for i := 1; i < len(nodes); i++ {
		if nodes[i].BlockChain.GetCurrentBlockHeader() == nil {
			assert.Fail(t, fmt.Sprintf("Node with idx %d does not have a current block", i))
		} else {
			assert.Equal(t, expectedNonce, nodes[i].BlockChain.GetCurrentBlockHeader().GetNonce())
		}
	}
}

func updateRound(nodes []*integrationTests.TestProcessorNode, round uint64) {
	for _, n := range nodes {
		n.Rounder.IndexField = int64(round)
	}
}