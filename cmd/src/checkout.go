package src

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func GenerateInvoiceNumber(userID int32) int {
	var invNum [14]string
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 14; i++ {
		invNum[i] = strconv.Itoa(rand.Intn(100))
	}
	newInvNum := strings.Join(invNum[:], "")
	invoiceNumber, _ := strconv.Atoi(newInvNum[1:8])
	newInvNum = fmt.Sprint(invoiceNumber)+fmt.Sprint(userID)
	invoiceNumber, _ = strconv.Atoi(newInvNum)
	return invoiceNumber
}

func GenerateSignature(userID int32, offerID []int32, paymentGatewayID int32) [32]byte{
	aOfferID := strings.Trim(strings.Replace(fmt.Sprint(offerID), " ", ",", -1), "[]")
	signature := fmt.Sprint(userID) + "|" + aOfferID + "|" + fmt.Sprint(paymentGatewayID)
	// calculate SHA-256 hash of the input
	signatureHash := sha256.Sum256([]byte(signature))

	return signatureHash
}
