package config

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

// MidtransClient adalah instance untuk transaksi Midtrans
var MidtransClient snap.Client

// InitMidtrans menginisialisasi Midtrans dengan server key
func InitMidtrans(serverKey string) {
	MidtransClient.New(serverKey, midtrans.Sandbox)
}
