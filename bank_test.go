package mango

import (
	"testing"
)

const (
	testIBAN = "FR3020041010124530725S03383"
	testBIC  = "CRLYFRPP"
)

func TestBankAccount_Save(test *testing.T) {
	serv := newTestService(test)
	user := createTestUser(serv)
	if _, err := user.Save(); err != nil {
		test.Fatal("Unable to store user:", err)
	}

	account := createTestBankAccount(test, serv, user)

	if account.IBAN != testIBAN {
		test.Fatalf("Invalid IBAN is saved bank account: got %s, should be %s", account.IBAN, testIBAN)
	}
}

func createTestBankAccount(test *testing.T, serv *MangoPay, user *NaturalUser) *BankAccount {
	acc, err := serv.NewBankAccount(user, user.FirstName, "one great place", IBAN)
	if err != nil {
		test.Fatal("Unable to create BankAccount", err.Error())
	}
	acc.IBAN, acc.BIC = testIBAN, testBIC
	if err := acc.Save(); err != nil {
		test.Fatal("Unable to save BankAccount", err.Error())
	}

	return acc
}
