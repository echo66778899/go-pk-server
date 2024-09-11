package engine

import (
	msgpb "go-pk-server/gen"
	mylog "go-pk-server/log"
)

const BUY_IN_SIZE = 2000

type BalanceManager struct {
	balance map[string]int
	inUse   map[string]int
}

func NewBalanceManager() *BalanceManager {
	return &BalanceManager{
		balance: make(map[string]int),
	}
}

func (bm *BalanceManager) TakeOneBuyIn(name string) int {
	bm.balance[name] -= BUY_IN_SIZE
	mylog.Infof("Player %s has taken one buy-in, balance is now %d.", name, bm.balance[name])
	mylog.Debugf("All players' balances: %+v", bm.balance)
	return BUY_IN_SIZE
}

func (bm *BalanceManager) PaybackOneBuyIn(name string) {
	bm.balance[name] += BUY_IN_SIZE
	mylog.Infof("Player %s has paid back one buy-in, balance is now %d.", name, bm.balance[name])
	mylog.Debugf("All players' balances: %+v", bm.balance)
}

// Amount can be negative
func (bm *BalanceManager) ReturnStack(name string, amount int) {
	mylog.Infof("Player %s has returned %d chips, balance is now %d.", name, amount, bm.balance[name])
	bm.balance[name] += amount
	if bm.balance[name] == 0 {
		mylog.Infof("Player %s is busted.", name)
		delete(bm.balance, name)
	}
}

func (bm *BalanceManager) UpdateCurrentPlayerChip(name string, amount int) {
	bm.inUse[name] = amount
}

func (bm *BalanceManager) GetBalance(name string) int {
	mylog.Infof("Player %s has balance %d.", name, bm.balance[name])
	return bm.balance[name] + bm.inUse[name]
}

func (bm *BalanceManager) GetBalanceSummary() *msgpb.BalanceInfo {
	balanceInfo := &msgpb.BalanceInfo{
		PlayerBalances: make([]*msgpb.PlayerBalance, 0),
	}
	for name, balance := range bm.balance {
		balanceInfo.PlayerBalances = append(balanceInfo.PlayerBalances, &msgpb.PlayerBalance{
			PlayerName: name,
			Balance:    int32(balance) + int32(bm.inUse[name]),
		})
	}
	return balanceInfo
}

func (bm *BalanceManager) Reset(name string) {
	bm.balance[name] = 0
}

func (bm *BalanceManager) ResetAll() {
	bm.balance = make(map[string]int)
}
