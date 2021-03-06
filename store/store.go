package store

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/alicebob/miniredis/v2"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/go-redis/redis/v8"
)

const (
	pending               = "pending"
	transit               = "transit"
	complete              = "complete"
	failed                = "failed"
	shadow                = "shadow"
	ibcReceiveFailed      = "IBC_receive_failed"
	ibcReceiveSuccess     = "IBC_receive_success"
	tokensUnlockedTimeout = "Tokens_unlocked_timeout"
	tokensUnlockedAck     = "Tokens_unlocked_ack"

	// pool swap fees is stored only for one hour(12 * defaultExpiry)
	poolExpiryMul = 12
)

var defaultExpiry = 300 * time.Second

type Store struct {
	Client        *redis.Client
	ConnectionURL string
	Config        struct{ ExpiryTime time.Duration }
}

type TxHashEntry struct {
	Chain  string
	Status string
	TxHash string
}

type Ticket struct {
	Owner    string        `json:"owner,omitempty"`
	Info     string        `json:"info,omitempty"`
	Height   int64         `json:"height,omitempty"`
	Status   string        `json:"status,omitempty"`
	TxHashes []TxHashEntry `json:"tx_hashes,omitempty"`
	Error    string        `json:"error,omitempty"`
}

func (t *Ticket) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, t)
}
func (t Ticket) MarshalBinary() (data []byte, err error) {
	return json.Marshal(t)
}

// NewClient creates a new redis client
func NewClient(connUrl string) (*Store, error) {

	var store Store

	store.Client = redis.NewClient(&redis.Options{
		Addr: connUrl,
		DB:   0,
	})

	store.ConnectionURL = connUrl

	store.Config.ExpiryTime = defaultExpiry

	return &store, nil

}

func (s *Store) CreateTicket(chain, txHash, owner string) error {
	owner = hex.EncodeToString([]byte(owner))
	data := Ticket{
		Owner:  owner,
		Status: pending,
	}

	key := GetKey(chain, txHash)
	if err := s.CreateShadowKey(key); err != nil {
		return err
	}

	if err := s.SetWithExpiry(key, data, 0); err != nil {
		return err
	}

	return s.sAdd(owner, key)
}

func (s *Store) SetComplete(key string, height int64) error {
	ticket, err := s.Get(key)
	if err != nil {
		return err
	}

	if err := s.SetWithExpiry(key, Ticket{Status: complete,
		Height: height}, 2); err != nil {
		return err
	}

	if err := s.DeleteShadowKey(key); err != nil {
		return err
	}

	return s.sRemove(ticket.Owner, key)
}

func (s *Store) SetIBCReceiveFailed(key string, txHashes []TxHashEntry, height int64) error {
	if err := s.CreateShadowKey(key); err != nil {
		return err
	}

	return s.SetWithExpiry(key, Ticket{Status: ibcReceiveFailed,
		TxHashes: txHashes, Height: height}, 0)
}

func (s *Store) SetIBCReceiveSuccess(key, owner string, txHashes []TxHashEntry, height int64) error {
	if err := s.SetWithExpiry(key, Ticket{
		Status:   ibcReceiveSuccess,
		TxHashes: txHashes,
		Height:   height}, 2); err != nil {
		return err
	}

	if err := s.DeleteShadowKey(key); err != nil {
		return err
	}

	return s.sRemove(owner, key)
}

func (s *Store) SetUnlockTimeout(key, owner string, txHashes []TxHashEntry, height int64) error {
	if err := s.SetWithExpiry(key, Ticket{Status: tokensUnlockedTimeout,
		Height:   height,
		TxHashes: txHashes}, 2); err != nil {
		return err
	}

	if err := s.DeleteShadowKey(key); err != nil {
		return err
	}

	return s.sRemove(owner, key)
}

func (s *Store) SetUnlockAck(key, owner string, txHashes []TxHashEntry, height int64) error {
	if err := s.SetWithExpiry(key, Ticket{Status: tokensUnlockedAck,
		Height:   height,
		TxHashes: txHashes}, 2); err != nil {
		return err
	}

	if err := s.DeleteShadowKey(key); err != nil {
		return err
	}

	return s.sRemove(owner, key)
}

func (s *Store) SetFailedWithErr(key, error string, height int64) error {
	if !s.Exists(key) {
		return fmt.Errorf("key doesn't exists")
	}

	prev, err := s.Get(key)
	if err != nil {
		return err
	}

	data := Ticket{
		Height: height,
		Status: failed,
		Error:  error,
	}

	if err := s.SetWithExpiry(key, data, 2); err != nil {
		return err
	}

	return s.sRemove(prev.Owner, key)
}

func (s *Store) SetInTransit(key, destChain, sourceChannel, sendPacketSequence, txHash, chainName string, height int64) error {

	if !s.Exists(key) {
		return fmt.Errorf("key doesn't exists")
	}

	if err := s.CreateShadowKey(key); err != nil {
		return err
	}

	ticket, err := s.Get(key)
	if err != nil {
		return err
	}

	ticket.Status = transit
	if err := s.SetWithExpiry(key, ticket, 2); err != nil {
		return err
	}

	newKey := GetIBCKey(destChain, sourceChannel, sendPacketSequence)

	if err := s.SetWithExpiry(newKey, Ticket{Info: key,
		Owner: ticket.Owner,
		TxHashes: []TxHashEntry{{
			Chain:  chainName,
			Status: transit,
			TxHash: txHash,
		}}}, 2); err != nil {
		return err
	}

	return nil
}

func (s *Store) SetIbcTimeoutUnlock(key, txHash, chainName string, height int64) error {

	prev, err := s.Get(key)

	if err != nil {
		return err
	}

	txHashes := append(prev.TxHashes, TxHashEntry{
		Chain:  chainName,
		Status: tokensUnlockedTimeout,
		TxHash: txHash,
	})

	return s.SetUnlockTimeout(prev.Info, prev.Owner, txHashes, height)
}

func (s *Store) SetIbcAckUnlock(key, txHash, chainName string, height int64) error {

	prev, err := s.Get(key)

	if err != nil {
		return err
	}

	txHashes := append(prev.TxHashes, TxHashEntry{
		Chain:  chainName,
		Status: tokensUnlockedAck,
		TxHash: txHash,
	})

	return s.SetUnlockAck(prev.Info, prev.Owner, txHashes, height)
}

func (s *Store) SetIbcReceived(key, txHash, chainName string, height int64) error {

	prev, err := s.Get(key)

	if err != nil {
		return err
	}

	txHashes := append(prev.TxHashes, TxHashEntry{
		Chain:  chainName,
		Status: ibcReceiveSuccess,
		TxHash: txHash,
	})

	return s.SetIBCReceiveSuccess(prev.Info, prev.Owner, txHashes, height)
}

func (s *Store) SetIbcFailed(key, txHash, chainName string, height int64) error {

	prev, err := s.Get(key)

	if err != nil {
		return err
	}

	txHashes := append(prev.TxHashes, TxHashEntry{
		Chain:  chainName,
		Status: ibcReceiveFailed,
		TxHash: txHash,
	})
	return s.SetIBCReceiveFailed(prev.Info, txHashes, height)
}

func (s *Store) SetPoolSwapFees(poolId, offerCoinAmount, offerCoinDenom string) error {
	poolTicket := fmt.Sprintf("pool/%s/%d", poolId, time.Now().Unix())

	offerCoinAmountInt, ok := sdk.NewIntFromString(offerCoinAmount)
	if !ok {
		return fmt.Errorf("unable to convert offerCoinAmout to sdk Int")
	}

	coin := sdk.NewCoin(offerCoinDenom, offerCoinAmountInt)
	return s.SetWithExpiry(poolTicket, coin.String(), poolExpiryMul) //  mul is 12 as time out is set to 5minutes by default
}

func (s *Store) CreateShadowKey(key string) error {
	shadowKey := shadow + key
	return s.SetWithExpiry(shadowKey, "", 1)
}

func (s *Store) Exists(key string) bool {
	exists, _ := s.Client.Exists(context.Background(), key).Result()

	return exists == 1
}

func (s *Store) SetWithExpiry(key string, value interface{}, mul int64) error {
	return s.Client.Set(context.Background(), key, value, time.Duration(mul)*(s.Config.ExpiryTime)).Err()
}

func (s *Store) SetWithExpiryTime(key string, value interface{}, duration time.Duration) error {
	return s.Client.Set(context.Background(), key, value, duration).Err()
}

func (s *Store) Get(key string) (Ticket, error) {
	var res Ticket
	if err := s.Client.Get(context.Background(), key).Scan(&res); err != nil {
		return Ticket{}, err
	}

	return res, nil
}

func (s *Store) GetUserTickets(user string) (map[string][]string, error) {
	var keys []string
	keys, err := s.sMembers(hex.EncodeToString([]byte(user)))
	if err != nil {
		return map[string][]string{}, err
	}

	res := make(map[string][]string)
	for _, key := range keys {
		s := strings.Split(key, "/")
		if len(s) != 2 {
			return map[string][]string{}, fmt.Errorf("unable to resolve chain name and tx hash")
		}

		res[s[0]] = append(res[s[0]], s[1])
	}

	return res, nil
}

func (s *Store) GetPools() ([]byte, error) {
	bz, err := s.Client.Get(context.Background(), "pools").Bytes()
	if err != nil {
		return nil, fmt.Errorf("cannot fetch pools from cache, %w", err)
	}

	return bz, nil
}

func (s *Store) GetParams() ([]byte, error) {
	bz, err := s.Client.Get(context.Background(), "params").Bytes()
	if err != nil {
		return bz, fmt.Errorf("cannot fetch params from cache, %w", err)
	}

	return bz, nil
}

func (s *Store) GetSupply() ([]byte, error) {
	bz, err := s.Client.Get(context.Background(), "supply").Bytes()
	if err != nil {
		return nil, fmt.Errorf("cannot fetch total supply from cache, %w", err)
	}

	return bz, nil
}

func (s *Store) GetNodeInfo() ([]byte, error) {
	bz, err := s.Client.Get(context.Background(), "node_info").Bytes()
	if err != nil {
		return nil, fmt.Errorf("cannot fetch node info from cache, %w", err)
	}

	return bz, nil
}

func (s *Store) Delete(key string) error {
	return s.Client.Del(context.Background(), key).Err()
}

func (s *Store) DeleteShadowKey(key string) error {
	shadowKey := shadow + key
	return s.Delete(shadowKey)
}
func (s *Store) sAdd(user, key string) error {
	return s.Client.SAdd(context.Background(), user, key).Err()
}

func (s *Store) sMembers(user string) ([]string, error) {
	var keys []string
	err := s.Client.SMembers(context.Background(), user).ScanSlice(&keys)
	if err != nil {
		return []string{}, err
	}

	return keys, err
}

func (s *Store) sRemove(user, key string) error {
	return s.Client.SRem(context.Background(), user, key).Err()
}

func (s *Store) GetSwapFees(poolId string) (sdk.Coins, error) {
	values, err := s.scan(fmt.Sprintf("pool/%s/*", poolId))
	if err != nil {
		return sdk.Coins{}, err
	}

	var coins sdk.Coins
	for _, value := range values {
		coin, err := sdk.ParseCoinNormalized(value)
		if err != nil {
			return sdk.Coins{}, err
		}

		coins = coins.Add(coin)
	}

	return coins, nil
}

func (s *Store) scan(prefix string) ([]string, error) {
	keys, nextCur, err := s.Client.Scan(context.Background(), 0, prefix, 10).Result()
	if err != nil {
		return nil, err
	}

	values, err := s.getValues(keys)
	if err != nil {
		return nil, err
	}

	if nextCur == 0 {
		return values, nil
	}

	for nextCur != 0 {
		var nextKeys []string
		nextKeys, nextCur, err = s.Client.Scan(context.Background(), nextCur, prefix, 100).Result()
		if err != nil {
			return nil, err
		}

		newValues, err := s.getValues(nextKeys)
		if err != nil {
			return nil, err
		}

		values = append(values, newValues...)
	}
	return values, nil
}

func (s *Store) getValues(keys []string) ([]string, error) {
	values := make([]string, 0, len(keys))

	for _, k := range keys {
		value, err := s.Client.Get(context.Background(), k).Result()
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return values, nil
}

func GetKey(chain, txHash string) string {
	return fmt.Sprintf("%s/%s", chain, txHash)
}

func GetIBCKey(chain, packetSrcChannel, packetSequence string) string {
	return fmt.Sprintf("%s-%s-%s", chain, packetSrcChannel, packetSequence)
}

func SetupTestStore() (*miniredis.Miniredis, *Store) {
	m, err := miniredis.Run()
	if err != nil {
		log.Fatalf("got error: %s when running miniredis", err)
	}

	s, err := NewClient(m.Addr())
	if err != nil {
		log.Fatalf("got error: %s when creating new store client", err)
	}

	return m, s
}

func ResetTestStore(m *miniredis.Miniredis, s *Store) {
	m.DB(s.Client.Options().DB).FlushDB()
}
