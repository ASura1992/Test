package wallet

import (
	"errors"
	"github.com/shopspring/decimal"
	"sync"
)

// 钱包
type Wallet struct {
	Mutex   sync.RWMutex
	Balance map[CoinType]decimal.Decimal
}

type CoinType string //货币类型
const (
	CoinBtc CoinType = "btc"
	CoinBth CoinType = "bth"
)

// 存钱
func (u *User) SaveMoney(amount decimal.Decimal, coinType CoinType) {
	u.Wallet.Mutex.Lock()
	defer u.Wallet.Mutex.Unlock()
	u.Wallet.Balance[coinType] = u.Wallet.Balance[coinType].Add(amount)
	return
}

// 扣除钱
func (u *User) DeductMoney(amount decimal.Decimal, coinType CoinType) error {
	u.Wallet.Mutex.Lock()
	defer u.Wallet.Mutex.Unlock()
	_, ok := u.Wallet.Balance[coinType]
	if !ok {
		return errors.New("no such coin type")
	}
	if amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("the withdrawal limit must be greater than 0")
	}
	if amount.GreaterThan(u.Wallet.Balance[coinType]) {
		return errors.New("your balance is insufficient")
	}
	u.Wallet.Balance[coinType] = u.Wallet.Balance[coinType].Sub(amount)
	return nil
}

// 向其他账号汇款
func (u *User) RemittanceToOther(quota decimal.Decimal, coinType CoinType, toUserId string) error {
	if u.Id == toUserId {
		return nil
	}
	val, ok := userCollection.Load(toUserId)
	if !ok {
		return errors.New("the remittance recipient does not exist")
	}
	toUser, ok := val.(User)
	if !ok {
		return errors.New("assertion user err")
	}
	//tips:这里下面的操作实际上可能需要数据库事务解决才可以，不然可能会导致数据前后不一致性
	//先扣钱
	if err := u.DeductMoney(quota, coinType); err != nil {
		return err
	}
	//再存入他人账号
	toUser.SaveMoney(quota, coinType)
	return nil
}

// 查询个人钱包
func (u *User) QueryWallet(coinTypes []CoinType) map[CoinType]decimal.Decimal {
	if len(coinTypes) == 0 {
		return nil
	}
	u.Wallet.Mutex.RLock()
	defer u.Wallet.Mutex.RUnlock()
	w := make(map[CoinType]decimal.Decimal, len(coinTypes))
	for _, coinType := range coinTypes {
		if val, ok := u.Wallet.Balance[coinType]; ok {
			w[coinType] = val
		} else {
			w[coinType] = decimal.NewFromFloat(0)
		}
	}
	return w
}
