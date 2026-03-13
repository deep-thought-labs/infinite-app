package mempool

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"go.opentelemetry.io/otel"

<<<<<<< HEAD
=======
	cmttypes "github.com/cometbft/cometbft/types"

<<<<<<< HEAD
	"github.com/cosmos/evm/mempool/internal/heightsync"
	"github.com/cosmos/evm/mempool/internal/queue"
>>>>>>> 84e141c6 (feat!: Krakatoa (#1030))
=======
>>>>>>> 5a7f5358 (feat(krakatoa): `operate_exclusively` toggles krakatoa features (#1060))
	"github.com/cosmos/evm/mempool/miner"
	"github.com/cosmos/evm/mempool/reserver"
	"github.com/cosmos/evm/mempool/txpool"
	"github.com/cosmos/evm/mempool/txpool/legacypool"
<<<<<<< HEAD
	"github.com/cosmos/evm/x/precisebank/types"
=======
	"github.com/cosmos/evm/rpc/stream"
>>>>>>> 18163be6 (refactor(config): remove evmAppOptions entirely (#661))
	evmtypes "github.com/cosmos/evm/x/vm/types"

<<<<<<< HEAD
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/log"
=======
	"cosmossdk.io/log/v2"
>>>>>>> 1651ab92 (chore: update to main and updated ibc (#1003))
	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdkmempool "github.com/cosmos/cosmos-sdk/types/mempool"
)

var (
	meter = otel.Meter("github.com/cosmos/evm/mempool")

	_ sdkmempool.ExtMempool = (*ExperimentalEVMMempool)(nil)
)

<<<<<<< HEAD
<<<<<<< HEAD
=======
=======
// AllowUnsafeSyncInsert indicates whether to perform synchronous inserts into the mempool
// for testing purposes. When true, Insert will block until the transaction is fully processed.
// This should be used only in tests to ensure deterministic behavior
var AllowUnsafeSyncInsert = false

>>>>>>> ad11ddc9 (fix: experimental EVM mempool not gossiping nonce + n transactions (#890))
const (
	// SubscriberName is the name of the event bus subscriber for the EVM mempool
	SubscriberName = "evm"
	// fallbackBlockGasLimit is the default block gas limit is 0 or missing in genesis file
	fallbackBlockGasLimit = 100_000_000
)

>>>>>>> 67950ad9 (feat: add block max_gas and min_tip ingestion to app.go mempool config (#582))
type (
	// ExperimentalEVMMempool is a unified mempool that manages both EVM and Cosmos SDK transactions.
	// It provides a single interface for transaction insertion, selection, and removal while
	// maintaining separate pools for EVM and Cosmos transactions. The mempool handles
	// fee-based transaction prioritization and manages nonce sequencing for EVM transactions.
	ExperimentalEVMMempool struct {
		/** Keepers **/
		vmKeeper VMKeeperI

		/** Mempools **/
		txPool       *txpool.TxPool
		legacyTxPool *legacypool.LegacyPool
		cosmosPool   sdkmempool.ExtMempool

		/** Utils **/
		logger        log.Logger
		txConfig      client.TxConfig
		clientCtx     client.Context
		blockchain    *Blockchain
		blockGasLimit uint64 // Block gas limit from consensus parameters
		minTip        *uint256.Int

<<<<<<< HEAD
		/** Verification **/
		anteHandler sdk.AnteHandler

<<<<<<< HEAD
		/** Concurrency **/
		mtx sync.Mutex
=======
=======
>>>>>>> a5aa523b (feat(mempool): shared ctx beween cosmos tx rechecks and inserts (#1031) (#1055))
		eventBus *cmttypes.EventBus
<<<<<<< HEAD

		/** Transaction Reaping **/
		reapList *ReapList

		/** Transaction Tracking **/
		txTracker *txTracker

		/** Transaction Inserting **/
		cosmosQueue *queue.Queue[sdk.Tx]
		evmQueue    *queue.Queue[ethtypes.Transaction]
>>>>>>> 84e141c6 (feat!: Krakatoa (#1030))
=======
>>>>>>> 5a7f5358 (feat(krakatoa): `operate_exclusively` toggles krakatoa features (#1060))
	}
)

// EVMMempoolConfig contains configuration options for creating an EVMsdkmempool.
// It allows customization of the underlying mempools, verification functions,
// and broadcasting functions used by the sdkmempool.
type EVMMempoolConfig struct {
	LegacyPoolConfig *legacypool.Config
	CosmosPoolConfig *sdkmempool.PriorityNonceMempoolConfig[math.Int]
	AnteHandler      sdk.AnteHandler
	BroadCastTxFn    func(txs []*ethtypes.Transaction) error
	// Block gas limit from consensus parameters
	BlockGasLimit uint64
	MinTip        *uint256.Int
}

// NewExperimentalEVMMempool creates a new unified mempool for EVM and Cosmos transactions.
// It initializes both EVM and Cosmos transaction pools, sets up blockchain interfaces,
// and configures fee-based prioritization. The config parameter allows customization
// of pools and verification functions, with sensible defaults created if not provided.
func NewExperimentalEVMMempool(
	getCtxCallback func(height int64, prove bool) (sdk.Context, error),
	logger log.Logger,
	vmKeeper VMKeeperI,
	feeMarketKeeper FeeMarketKeeperI,
	txConfig client.TxConfig,
	config *EVMMempoolConfig,
	cosmosPoolMaxTx int,
) *ExperimentalEVMMempool {
	var (
		cosmosPool sdkmempool.ExtMempool
		blockchain *Blockchain
	)

	// add the mempool name to the logger
	logger = logger.With(log.ModuleKey, "ExperimentalEVMMempool")

	logger.Debug("creating new EVM mempool")

	if config == nil {
		panic("config must not be nil")
	}

	if config.BlockGasLimit == 0 {
		logger.Warn("block gas limit is 0, setting to fallback", "fallback_limit", fallbackBlockGasLimit)
		config.BlockGasLimit = fallbackBlockGasLimit
	}

	blockchain = NewBlockchain(getCtxCallback, logger, vmKeeper, feeMarketKeeper, config.BlockGasLimit)

	// Create txPool from configuration
	legacyConfig := legacypool.DefaultConfig
	if config.LegacyPoolConfig != nil {
		legacyConfig = *config.LegacyPoolConfig
	}

	legacyPool := legacypool.New(legacyConfig, logger, blockchain)

	tracker := reserver.NewReservationTracker()
	txPool, err := txpool.New(uint64(0), blockchain, tracker, []txpool.SubPool{legacyPool})
	if err != nil {
		panic(err)
	}

	if len(txPool.Subpools) != 1 {
		panic("tx pool should contain one subpool")
	}
	if _, ok := txPool.Subpools[0].(*legacypool.LegacyPool); !ok {
		panic("tx pool should contain only legacypool")
	}

	// TODO: move this logic to evmd.createMempoolConfig and set the max tx there
	// Create Cosmos Mempool from configuration
	cosmosPoolConfig := config.CosmosPoolConfig
	if cosmosPoolConfig == nil {
		// Default configuration
		defaultConfig := sdkmempool.PriorityNonceMempoolConfig[math.Int]{}
		defaultConfig.TxPriority = sdkmempool.TxPriority[math.Int]{
			GetTxPriority: func(goCtx context.Context, tx sdk.Tx) math.Int {
				ctx := sdk.UnwrapSDKContext(goCtx)
				cosmosTxFee, ok := tx.(sdk.FeeTx)
				if !ok {
					return math.ZeroInt()
				}
				found, coin := cosmosTxFee.GetFee().Find(vmKeeper.GetEvmCoinInfo(ctx).Denom)
				if !found {
					return math.ZeroInt()
				}

				gasPrice := coin.Amount.Quo(math.NewIntFromUint64(cosmosTxFee.GetGas()))

				return gasPrice
			},
			Compare: func(a, b math.Int) int {
				return a.BigInt().Cmp(b.BigInt())
			},
			MinValue: math.ZeroInt(),
		}
		cosmosPoolConfig = &defaultConfig
	}

	cosmosPoolConfig.MaxTx = cosmosPoolMaxTx
	cosmosPool = sdkmempool.NewPriorityMempool(*cosmosPoolConfig)

	evmMempool := &ExperimentalEVMMempool{
		vmKeeper:      vmKeeper,
		txPool:        txPool,
		legacyTxPool:  txPool.Subpools[0].(*legacypool.LegacyPool),
		cosmosPool:    cosmosPool,
		logger:        logger,
		txConfig:      txConfig,
		blockchain:    blockchain,
		blockGasLimit: config.BlockGasLimit,
		minTip:        config.MinTip,
	}

	legacyPool.OnTxPromoted = evmMempool.onEVMTxPromoted(config.BroadCastTxFn)

	vmKeeper.SetEvmMempool(evmMempool)

	return evmMempool
}

// IsExclusive returns true if this mempool is the ONLY mempool in the chain.
func (m *ExperimentalEVMMempool) IsExclusive() bool {
	return false
}

// GetBlockchain returns the blockchain interface used for chain head event notifications.
// This is primarily used to notify the mempool when new blocks are finalized.
func (m *ExperimentalEVMMempool) GetBlockchain() *Blockchain {
	return m.blockchain
}

// GetTxPool returns the underlying EVM txpool.
// This provides direct access to the EVM-specific transaction management functionality.
func (m *ExperimentalEVMMempool) GetTxPool() *txpool.TxPool {
	return m.txPool
}

// SetClientCtx sets the client context provider for broadcasting transactions
func (m *ExperimentalEVMMempool) SetClientCtx(clientCtx client.Context) {
	m.clientCtx = clientCtx
}

// Insert adds a transaction to the appropriate mempool (EVM or Cosmos).
// EVM transactions are routed to the EVM transaction pool, while all other
// transactions are inserted into the Cosmos sdkmempool.
func (m *ExperimentalEVMMempool) Insert(goCtx context.Context, tx sdk.Tx) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	blockHeight := ctx.BlockHeight()

<<<<<<< HEAD
<<<<<<< HEAD
	ctx := sdk.UnwrapSDKContext(goCtx)
	blockHeight := ctx.BlockHeight()

	m.logger.Debug("inserting transaction into mempool", "block_height", blockHeight)

	if blockHeight < 2 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidHeight, "Mempool is not ready. Please wait for block 1 to finalize.")
	}

	ethMsg, err := m.getEVMMessage(tx)
	if err == nil {
		// Insert into EVM pool
		m.logger.Debug("inserting EVM transaction", "tx_hash", ethMsg.Hash)
		ethTxs := []*ethtypes.Transaction{ethMsg.AsTransaction()}
		errs := m.txPool.Add(ethTxs, AllowUnsafeSyncInsert)
		if len(errs) > 0 && errs[0] != nil {
			m.logger.Error("failed to insert EVM transaction", "error", errs[0], "tx_hash", ethMsg.Hash)
			return errs[0]
		}
		m.logger.Debug("EVM transaction inserted successfully", "tx_hash", ethMsg.Hash)
=======
	if errC != nil {
		// if we got back a non nil async error channel, wait for that to
		// resolve
		select {
		case err := <-errC:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

// InsertAsync adds a transaction to the appropriate mempool (EVM or Cosmos). EVM
// transactions are routed to the EVM transaction pool, while all other
// transactions are inserted into the Cosmos sdkmempool. EVM transactions are
// inserted async, i.e. they are scheduled for promotion only, we do not wait
// for it to complete.
func (m *ExperimentalEVMMempool) InsertAsync(ctx context.Context, tx sdk.Tx) error {
	errC, err := m.insert(ctx, tx)
	if err != nil {
		return fmt.Errorf("inserting tx: %w", err)
	}

	select {
	case err := <-errC:
		// if we have a result immediately, ready on the channel returned from
		// insert, return that (cosmos tx or unable to try and insert the tx
		// due to parsing error).
		return err
	case <-ctx.Done():
		return ctx.Err()
	default:
		// result was not ready immediately, return nil while async things happen
>>>>>>> 84e141c6 (feat!: Krakatoa (#1030))
		return nil
	}
}

// insert inserts a tx into its respective mempool, returning a channel for any
// async errors that may happen later upon actual mempool insertion, and an
// error for any errors that occurred synchronously.
func (m *ExperimentalEVMMempool) insert(_ context.Context, tx sdk.Tx) (<-chan error, error) {
=======
	m.logger.Debug("inserting transaction into mempool", "block_height", blockHeight)
>>>>>>> 5a7f5358 (feat(krakatoa): `operate_exclusively` toggles krakatoa features (#1060))
	ethMsg, err := evmTxFromCosmosTx(tx)
	switch {
	case err == nil:
		// Insert into EVM pool
		hash := ethMsg.Hash()
		m.logger.Debug("inserting EVM transaction", "tx_hash", hash)
		ethTxs := []*ethtypes.Transaction{ethMsg.AsTransaction()}
		errs := m.txPool.Add(ethTxs, AllowUnsafeSyncInsert)
		if len(errs) > 0 && errs[0] != nil {
			m.logger.Error("failed to insert EVM transaction", "error", errs[0], "tx_hash", hash)
			return errs[0]
		}
		m.logger.Debug("EVM transaction inserted successfully", "tx_hash", hash)
		return nil
	case errors.Is(err, ErrMultiMsgEthereumTransaction):
		// there are multiple messages in this tx and one or more of them is an
		// evm tx, this is invalid
		return err
	default:
		// Insert into cosmos pool for non-EVM transactions
		m.logger.Debug("inserting Cosmos transaction")
		if err = m.cosmosPool.Insert(goCtx, tx); err != nil {
			m.logger.Error("failed to insert Cosmos transaction", "error", err)
			return err
		}

		m.logger.Debug("Cosmos transaction inserted successfully")
		return nil
	}
}

// InsertInvalidNonce handles transactions that failed with nonce gap errors.
// It attempts to insert EVM transactions into the pool as non-local transactions,
// allowing them to be queued for future execution when the nonce gap is filled.
// Non-EVM transactions are discarded as regular Cosmos flows do not support nonce gaps.
func (m *ExperimentalEVMMempool) InsertInvalidNonce(txBytes []byte) error {
	tx, err := m.txConfig.TxDecoder()(txBytes)
	if err != nil {
		return err
	}

	var ethTxs []*ethtypes.Transaction
	msgs := tx.GetMsgs()
	if len(msgs) != 1 {
		return fmt.Errorf("%w, got %d", ErrExpectedOneMessage, len(msgs))
	}
	for _, msg := range tx.GetMsgs() {
		ethMsg, ok := msg.(*evmtypes.MsgEthereumTx)
		if ok {
			ethTxs = append(ethTxs, ethMsg.AsTransaction())
			continue
		}
	}
	errs := m.txPool.Add(ethTxs, false)
	if errs != nil {
		if len(errs) != 1 {
			return fmt.Errorf("%w, got %d", ErrExpectedOneError, len(errs))
		}
		return errs[0]
	}
	return nil
}

// Select returns a unified iterator over both EVM and Cosmos transactions.
// The iterator prioritizes transactions based on their fees and manages proper
// sequencing. The i parameter contains transaction hashes to exclude from selection.
func (m *ExperimentalEVMMempool) Select(goCtx context.Context, i [][]byte) sdkmempool.Iterator {
	return m.buildIterator(goCtx, i)
}

// SelectBy iterates through transactions until the provided filter function returns false.
// It uses the same unified iterator as Select but allows early termination based on
// custom criteria defined by the filter function.
func (m *ExperimentalEVMMempool) SelectBy(goCtx context.Context, txs [][]byte, filter func(sdk.Tx) bool) {
	defer func(t0 time.Time) { telemetry.MeasureSince(t0, "expmempool_selectby_duration") }(time.Now()) //nolint:staticcheck

	iter := m.buildIterator(goCtx, txs)

	for iter != nil && filter(iter.Tx()) {
		iter = iter.Next()
	}
}

// buildIterator ensures that EVM mempool has checked txs for reorgs up to COMMITTED
// block height and then returns a combined iterator over EVM & Cosmos txs.
func (m *ExperimentalEVMMempool) buildIterator(ctx context.Context, txs [][]byte) sdkmempool.Iterator {
	defer func(t0 time.Time) { telemetry.MeasureSince(t0, "expmempool_builditerator_duration") }(time.Now()) //nolint:staticcheck

	evmIterator, cosmosIterator := m.getIterators(ctx, txs)

	return NewEVMMempoolIterator(
		evmIterator,
		cosmosIterator,
		m.logger,
		m.txConfig,
		m.vmKeeper.GetEvmCoinInfo(sdk.UnwrapSDKContext(ctx)).Denom,
		m.blockchain,
	)
}

// CountTx returns the total number of transactions in both EVM and Cosmos pools.
// This provides a combined count across all mempool types.
func (m *ExperimentalEVMMempool) CountTx() int {
	pending, _ := m.txPool.Stats()
	return m.cosmosPool.CountTx() + pending
}

// Remove fallbacks for RemoveWithReason
func (m *ExperimentalEVMMempool) Remove(tx sdk.Tx) error {
	return m.RemoveWithReason(context.Background(), tx, sdkmempool.RemoveReason{
		Caller: "remove",
		Error:  nil,
	})
}

// RemoveWithReason removes a transaction from the appropriate sdkmempool.
// For EVM transactions, removal is typically handled automatically by the pool
// based on nonce progression. Cosmos transactions are removed from the Cosmos pool.
func (m *ExperimentalEVMMempool) RemoveWithReason(ctx context.Context, tx sdk.Tx, reason sdkmempool.RemoveReason) error {
	chainCtx, err := m.blockchain.GetLatestContext()
	if err != nil || chainCtx.BlockHeight() == 0 {
		m.logger.Warn("Failed to get latest context, skipping removal")
		return nil
	}

<<<<<<< HEAD
	m.logger.Debug("removing transaction from mempool")

	msg, err := m.getEVMMessage(tx)
	if err == nil {
		// Comet will attempt to remove transactions from the mempool after completing successfully.
		// We should not do this with EVM transactions because removing them causes the subsequent ones to
		// be dequeued as temporarily invalid, only to be requeued a block later.
		// The EVM mempool handles removal based on account nonce automatically.
		if m.shouldRemoveFromEVMPool(tx) {
			m.logger.Debug("manually removing EVM transaction", "tx_hash", msg.Hash())
			m.legacyTxPool.RemoveTx(msg.Hash(), false, true)
		} else {
			m.logger.Debug("skipping manual removal of EVM transaction, leaving to mempool to handle", "tx_hash", msg.Hash)
		}
		return nil
=======
	msgEthereumTx, err := evmTxFromCosmosTx(tx)
	switch {
	case errors.Is(err, ErrNoMessages):
		return err
	case err != nil:
<<<<<<< HEAD
		// unable to parse evm tx -> process as cosmos tx
		return m.removeCosmosTx(ctx, tx, reason)
>>>>>>> 84e141c6 (feat!: Krakatoa (#1030))
=======
		m.logger.Debug("Removing Cosmos transaction")

		if err := sdkmempool.RemoveWithReason(ctx, m.cosmosPool, tx, reason); err != nil {
			m.logger.Error("Failed to remove Cosmos transaction", "error", err)
			return err
		}

		m.logger.Debug("Cosmos transaction removed successfully")
		return nil
>>>>>>> 5a7f5358 (feat(krakatoa): `operate_exclusively` toggles krakatoa features (#1060))
	}

	hash := msgEthereumTx.Hash()

	if m.shouldRemoveFromEVMPool(hash, reason) {
		m.logger.Debug("Manually removing EVM transaction", "tx_hash", hash)
		m.legacyTxPool.RemoveTx(hash, false, true, convertRemovalReason(reason.Caller))
	}

	return nil
}

// convertRemovalReason converts a removal caller to a removal reason
func convertRemovalReason(caller sdkmempool.RemovalCaller) txpool.RemovalReason {
	switch caller {
	case sdkmempool.CallerRunTxRecheck:
		return legacypool.RemovalReasonRunTxRecheck
	case sdkmempool.CallerRunTxFinalize:
		return legacypool.RemovalReasonRunTxFinalize
	case sdkmempool.CallerPrepareProposalRemoveInvalid:
		return legacypool.RemovalReasonPreparePropsoalInvalid
	default:
		return txpool.RemovalReason("")
	}
}

// shouldRemoveFromEVMPool determines whether an EVM transaction should be manually removed.
func (m *ExperimentalEVMMempool) shouldRemoveFromEVMPool(hash common.Hash, reason sdkmempool.RemoveReason) bool {
	if reason.Error == nil {
		return false
	}
	// Comet will attempt to remove transactions from the mempool after completing successfully.
	// We should not do this with EVM transactions because removing them causes the subsequent ones to
	// be dequeued as temporarily invalid, only to be requeued a block later.
	// The EVM mempool handles removal based on account nonce automatically.
	isKnown := errors.Is(reason.Error, ErrNonceGap) ||
		errors.Is(reason.Error, sdkerrors.ErrInvalidSequence) ||
		errors.Is(reason.Error, sdkerrors.ErrOutOfGas)

	if isKnown {
		m.logger.Debug("Transaction validation succeeded, should be kept", "tx_hash", hash, "caller", reason.Caller)
		return false
	}

	m.logger.Debug("Transaction validation failed, should be removed", "tx_hash", hash, "caller", reason.Caller)
	return true
}

<<<<<<< HEAD
=======
// SetEventBus sets CometBFT event bus to listen for new block header event.
func (m *ExperimentalEVMMempool) SetEventBus(eventBus *cmttypes.EventBus) {
	if m.HasEventBus() {
		m.eventBus.Unsubscribe(context.Background(), SubscriberName, stream.NewBlockHeaderEvents) //nolint: errcheck
	}
	m.eventBus = eventBus
	sub, err := eventBus.Subscribe(context.Background(), SubscriberName, stream.NewBlockHeaderEvents)
	if err != nil {
		panic(err)
	}
	go func() {
		bc := m.GetBlockchain()
		for range sub.Out() {
			bc.NotifyNewBlock()
		}
	}()
}

// HasEventBus returns true if the blockchain is configured to use an event bus for block notifications.
func (m *ExperimentalEVMMempool) HasEventBus() bool {
	return m.eventBus != nil
}

func (m *ExperimentalEVMMempool) Close() error {
	var errs []error
	if m.eventBus != nil {
		if err := m.eventBus.Unsubscribe(context.Background(), SubscriberName, stream.NewBlockHeaderEvents); err != nil {
			errs = append(errs, fmt.Errorf("failed to unsubscribe from event bus: %w", err))
		}
	}

	if err := m.txPool.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close txpool: %w", err))
	}

	return errors.Join(errs...)
}

>>>>>>> 00328124 (fix: allow mempool interrupt while it is still ongoing (#495))
// getEVMMessage validates that the transaction contains exactly one message and returns it if it's an EVM message.
// Returns an error if the transaction has no messages, multiple messages, or the single message is not an EVM transaction.
func evmTxFromCosmosTx(tx sdk.Tx) (*evmtypes.MsgEthereumTx, error) {
	msgs := tx.GetMsgs()
	if len(msgs) == 0 {
		return nil, ErrNoMessages
	}

	// ethereum txs should only contain a single msg that is a MsgEthereumTx
	// type
	if len(msgs) > 1 {
		// transaction has > 1 msg, will be treated as a cosmos tx by the
		// mempool. validate that none of the msgs are a MsgEthereumTx since
		// those should only be used in the single msg case
		for _, msg := range msgs {
			if _, ok := msg.(*evmtypes.MsgEthereumTx); ok {
				return nil, ErrMultiMsgEthereumTransaction
			}
		}

		// transaction has > 1 msg, but none were ethereum txs, this is
		// still not a valid eth tx
		return nil, fmt.Errorf("%w, got %d", ErrExpectedOneMessage, len(msgs))
	}

	ethMsg, ok := msgs[0].(*evmtypes.MsgEthereumTx)
	if !ok {
		return nil, ErrNotEVMTransaction
	}
	return ethMsg, nil
}

// getIterators prepares iterators over pending EVM and Cosmos transactions.
// It configures EVM transactions with proper base fee filtering and priority ordering,
// while setting up the Cosmos iterator with the provided exclusion list.
func (m *ExperimentalEVMMempool) getIterators(ctx context.Context, txs [][]byte) (evm *miner.TransactionsByPriceAndNonce, cosmos sdkmempool.Iterator) {
	var (
		evmIterator    *miner.TransactionsByPriceAndNonce
		cosmosIterator sdkmempool.Iterator
		wg             sync.WaitGroup
	)

	wg.Go(func() {
		evmIterator = m.evmIterator(ctx)
	})

	wg.Go(func() {
		cosmosIterator = m.cosmosPool.Select(ctx, txs)
	})

	wg.Wait()

	return evmIterator, cosmosIterator
}

// evmIterator returns an iterator over the current valid txs in the evm
// mempool at height.
func (m *ExperimentalEVMMempool) evmIterator(ctx context.Context) *miner.TransactionsByPriceAndNonce {
	sdkctx := sdk.UnwrapSDKContext(ctx)
	baseFee := m.vmKeeper.GetBaseFee(sdkctx)
	var baseFeeUint *uint256.Int
	if baseFee != nil {
		baseFeeUint = uint256.MustFromBig(baseFee)
	}

	filter := txpool.PendingFilter{
		MinTip:       m.minTip,
		BaseFee:      baseFeeUint,
		BlobFee:      nil,
		OnlyPlainTxs: true,
		OnlyBlobTxs:  false,
	}

	evmPendingTxs := m.txPool.Pending(ctx, filter)
	return miner.NewTransactionsByPriceAndNonce(nil, evmPendingTxs, baseFee)
}

func (m *ExperimentalEVMMempool) onEVMTxPromoted(broadcastTxFn func(txs []*ethtypes.Transaction) error) func(tx *ethtypes.Transaction) {
	if broadcastTxFn != nil {
		return func(tx *ethtypes.Transaction) {
			if err := broadcastTxFn(ethtypes.Transactions{tx}); err != nil {
				m.logger.Error("Failed to broadcast transaction", "err", err, "tx_hash", tx.Hash())
			}
		}
	}

	return func(tx *ethtypes.Transaction) {
		if err := m.broadcastEVMTransaction(m.clientCtx, tx); err != nil {
			m.logger.Error("Failed to broadcast transaction", "err", err, "tx_hash", tx.Hash())
		}
	}
}

// broadcastEVMTransaction converts an Ethereum transaction to Cosmos SDK format and broadcasts them.
// This function wraps EVM transactions in MsgEthereumTx messages and submits them to the network
// using the provided client context. It handles encoding and error reporting for each transaction.
func (m *ExperimentalEVMMempool) broadcastEVMTransaction(clientCtx client.Context, ethTx *ethtypes.Transaction) error {
	msg := &evmtypes.MsgEthereumTx{}
	ethSigner := ethtypes.LatestSigner(evmtypes.GetEthChainConfig())
	if err := msg.FromSignedEthereumTx(ethTx, ethSigner); err != nil {
		return fmt.Errorf("failed to convert ethereum transaction: %w", err)
	}

	cosmosTx, err := msg.BuildTx(clientCtx.TxConfig.NewTxBuilder(), evmtypes.GetEVMCoinDenom())
	if err != nil {
		return fmt.Errorf("failed to build cosmos tx: %w", err)
	}

	txBytes, err := clientCtx.TxConfig.TxEncoder()(cosmosTx)
	if err != nil {
		return fmt.Errorf("failed to encode transaction: %w", err)
	}

	res, err := clientCtx.BroadcastTxSync(txBytes)
	if err != nil {
		return fmt.Errorf("failed to broadcast transaction %s: %w", ethTx.Hash().Hex(), err)
	}
	if res.Code != 0 && res.Code != 19 && res.RawLog != "already known" {
		return fmt.Errorf("transaction %s rejected by mempool: code=%d, log=%s", ethTx.Hash().Hex(), res.Code, res.RawLog)
	}
	return nil
}
