package src

import (
	"fmt"
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
