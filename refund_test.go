package mango

import (
	"strings"
	"testing"
)

func TestRefund(test *testing.T) {
	serv := newTestService(test)
	user := createTestUser(serv)
	if _, err := user.Save(); err != nil {
		test.Fatal("Unable to store user:", err)
	}
	wallet := createTestWallet(test, serv, user)

	payin := createTestDirectDebitWebPayIn(test, serv, user, EUR100, EUR0, wallet)
	_, _, err := payin.Refund()
	if err != nil {
		// TODO: obtain succeeded transaction to test refund
		// Now ignore 'must have a SUCCEEDED Status' error
		if !strings.Contains(err.Error(), "original transaction must have a SUCCEEDED Status") {
			test.Fatal("Unable to refund a payin:", err)
		}
	}

	// if refund.RefundReasonType != RefundReasonOwnerDoNotMatchBankaccount {
	// 	test.Fatal("Invalid refund reason type ", refund.RefundReasonType)
	// }
}
