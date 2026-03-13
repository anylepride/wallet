package wallet

import (
	"github.com/shopspring/decimal"
)

type Wallet struct {
	id      string
	balance decimal.Decimal
}

func NewWallet(id string) *Wallet {
	return &Wallet{id: id}
}

func (w *Wallet) GetBalance() float64 {
	return w.balance.InexactFloat64()
}

func (w *Wallet) GetWalletId() string {
	return w.id
}

func (w *Wallet) AddBalance(amount float64) {
	w.balance.Add(decimal.NewFromFloat(amount))
}

func (w *Wallet) SubBalance(amount float64) {
	w.balance.Sub(decimal.NewFromFloat(amount))
}
