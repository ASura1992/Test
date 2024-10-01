package wallet

import "sync"

// 用户
type User struct {
	Id     string
	Wallet *Wallet
}

// 用户集合
var userCollection sync.Map //userId:User{}
// 把user信息初始化加载进缓存
func SetUser(user User) {
	userCollection.Store(user.Id, user)
}
