package wallet

import (
	"fmt"
	"github.com/shopspring/decimal"
	"log"
	"testing"
)

func TestWallet(t *testing.T) {
	//假设有这些用户
	userIds := []string{"1", "2", "3", "4", "5"}
	userId1 := "1"                                //如当前登录的账户
	userId2 := "4"                                //汇款接受对象
	operationCoinType := CoinBtc                  //操作的币种
	saveAmount := decimal.NewFromFloat(20.5)      //user1存钱额度
	deductMoney := decimal.NewFromFloat(5.5)      //user1取钱额度
	remittanceAmount := decimal.NewFromFloat(5.8) //userId1向userId2汇款额度
	//读取出来写入缓存中
	for _, v := range userIds {
		SetUser(User{
			Id: v,
			Wallet: &Wallet{
				Balance: map[CoinType]decimal.Decimal{
					CoinBtc: decimal.Zero,
					CoinBth: decimal.Zero,
				},
			},
		})
	}
	fmt.Println("all users wallet is empty")
	val, ok := userCollection.Load(userId1)
	if !ok {
		log.Fatalln("user id :", userId1, "不存在")
	}
	user1, ok := val.(User)
	if !ok {
		log.Fatalln("assertion failure ")
	}
	val2, ok := userCollection.Load(userId2)
	if !ok {
		log.Fatalln("user id :", userId2, "不存在")
	}
	user2, ok := val2.(User)
	if !ok {
		log.Fatalln("assertion failure ")
	}
	user1.SaveMoney(saveAmount, operationCoinType)
	fmt.Println("userId:", userId1, "充值币种", operationCoinType, saveAmount, "后账户信息:", user1.Wallet.Balance)
	//user1取钱
	if err := user1.DeductMoney(deductMoney, operationCoinType); err != nil {
		log.Fatalln("userId:", userId1, "取出币种", operationCoinType, ":", deductMoney, "错误,", err)
	}
	fmt.Println("userId:", userId1, "取出币种", operationCoinType, ":", deductMoney, "后账户信息:", user1.Wallet.Balance)
	//userId1向userId2汇款
	if err := user1.RemittanceToOther(remittanceAmount, operationCoinType, userId2); err != nil {
		log.Fatalln("userId:", userId1, "向 userId:", userId2, "汇款", operationCoinType, ":", remittanceAmount, "错误,", err)
	}
	fmt.Println("userId:", userId1, "向 userId:", userId2, "汇款", operationCoinType, ":", remittanceAmount, "后账户信息:", user1.Wallet.Balance)
	fmt.Println("userId:", userId2, "接收到 userId:", userId1, "汇款", operationCoinType, ":", remittanceAmount, "后账户信息:", user2.Wallet.Balance)
}
