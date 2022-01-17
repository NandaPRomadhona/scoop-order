package src

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"
)

func TestGenerateInvoiceNumber(t *testing.T) {
	type args struct {
		userID int32
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{name: "Generate correct data", args: args{userID: 1234}, want: 123},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateInvoiceNumber(tt.args.userID)
			fmt.Println("got: ", got)
			if got == tt.want {
				t.Errorf("GenerateInvoiceNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateSignature(t *testing.T) {
	type args struct {
		userID           int32
		offerID          []int32
		paymentGatewayID int32
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Valid test", args{1855406, []int32{43103}, 29}, "c38376514cf6b5a264a10869e0de85d85f6906a4c9337bf2a3e267b2cb970527"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateSignature(tt.args.userID, tt.args.offerID, tt.args.paymentGatewayID); !reflect.DeepEqual(hex.EncodeToString(got[:]), tt.want) {
				t.Errorf("GenerateSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}