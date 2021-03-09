// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Custom error returned in case of failed transaction.
type ErrTransferFailed struct {
	transferId string
	msg        string
}

func (e *ErrTransferFailed) Error() string {
	return fmt.Sprintf("transfer %s failed: %s ", e.transferId, e.msg)
}

// List of transactions.
type TransferList []*Transfer

// Transfer hold details about relocating e-money from a wallet
// to another one.
//
// See http://docs.mangopay.com/api-references/transfers/.
type Transfer struct {
	ProcessReply
	AuthorId         string
	CreditedUserId   string
	DebitedFunds     Money
	Fees             Money
	DebitedWalletId  string
	CreditedWalletId string
	CreditedFunds    Money
	service          *MangoPay
}

// List of transactions.
type TransactionList []*Transaction

// Transfer hold details about relocating e-money from a wallet
// to another one.
//
// See http://docs.mangopay.com/api-references/transfers/.
type Transaction struct {
	ProcessReply
	AuthorId         string
	CreditedUserId   string
	DebitedFunds     Money
	Fees             Money
	DebitedWalletId  string
	CreditedWalletId string
	CreditedFunds    Money
	Type             string
	service          *MangoPay
}

func (t *Transfer) String() string {
	return struct2string(t)
}

func (t *Transaction) String() string {
	return struct2string(t)
}

// NewTransfer creates a new tranfer (or transaction).
func (m *MangoPay) NewTransfer(author Consumer, amount Money, fees Money, from, to *Wallet) (*Transfer, error) {
	msg := "new tranfer: "
	if author == nil {
		return nil, errors.New(msg + "nil author")
	}
	if from == nil {
		return nil, errors.New(msg + "nil source wallet")
	}
	if to == nil {
		return nil, errors.New(msg + "nil dest wallet")
	}
	if from.Id == "" {
		return nil, errors.New(msg + "source wallet has empty Id")
	}
	if to.Id == "" {
		return nil, errors.New(msg + "dest wallet has empty Id")
	}
	id := consumerId(author)
	if id == "" {
		return nil, errors.New(msg + "author has empty Id")
	}
	t := &Transfer{
		AuthorId:         id,
		DebitedFunds:     amount,
		Fees:             fees,
		DebitedWalletId:  from.Id,
		CreditedWalletId: to.Id,
		ProcessReply:     ProcessReply{},
	}
	t.service = m
	return t, nil
}

func (t *Transfer) Refund() (*Refund, *RateLimitInfo, error) {
	r := &Refund{
		ProcessReply: ProcessReply{},
		transfer:     t,
		kind:         transferRefund,
	}

	rateLimitInfo, err := r.save()
	if err != nil {
		return nil, nil, err
	}

	return r, rateLimitInfo, nil
}

// Save sends an HTTP query to create a transfer. Upon successful creation,
// it may return an ErrTransferFailed error if the transaction has been
// rejected (unsufficient wallet balance for example).
func (t *Transfer) Save() (*RateLimitInfo, error) {
	data := JsonObject{}
	j, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(j, &data); err != nil {
		return nil, err
	}

	// Force float64 to int conversion after unmarshalling.
	for _, field := range []string{"CreationDate", "ExecutionDate"} {
		data[field] = int(data[field].(float64))
	}

	// Fields not allowed when creating a tranfer.
	for _, field := range []string{"Id", "CreationDate", "ExecutionDate", "CreditedFunds", "CreditedUserId", "ResultCode", "ResultMessage", "Status"} {
		delete(data, field)
	}

	tr, rateLimitInfo, err := t.service.anyRequest(new(Transfer), actionCreateTransfer, data)
	if err != nil {
		return nil, err
	}
	serv := t.service
	*t = *(tr.(*Transfer))
	t.service = serv

	if t.Status == "FAILED" {
		return nil, &ErrTransferFailed{t.Id, t.ResultMessage}
	}
	return rateLimitInfo, nil
}

// Transfer finds a transaction by id.
func (m *MangoPay) Transfer(id string) (*Transfer, *RateLimitInfo, error) {
	w, rateLimitInfo, err := m.anyRequest(new(Transfer), actionFetchTransfer, JsonObject{"Id": id})
	if err != nil {
		return nil, nil, err
	}
	return w.(*Transfer), rateLimitInfo, nil
}

// Transfer finds all user's transactions. Provided for convenience.
func (m *MangoPay) Transfers(user Consumer) (TransferList, *RateLimitInfo, error) {
	return m.transfers(user)
}

func (m *MangoPay) transfers(u Consumer) (TransferList, *RateLimitInfo, error) {
	id := consumerId(u)
	if id == "" {
		return nil, nil, errors.New("user has empty Id")
	}
	trs, rateLimitInfo, err := m.anyRequest(new(TransferList), actionFetchUserTransfers, JsonObject{"Id": id})
	if err != nil {
		return nil, nil, err
	}
	return *(trs.(*TransferList)), rateLimitInfo, nil
}

// Transfer finds all user's transactions. Provided for convenience.
func (m *MangoPay) Transactions(user Consumer) (TransactionList, *RateLimitInfo, error) {
	return m.transactions(user)
}

func (m *MangoPay) transactions(u Consumer) (TransactionList, *RateLimitInfo, error) {
	id := consumerId(u)
	if id == "" {
		return nil, nil, errors.New("user has empty Id")
	}
	trs, rateLimitInfo, err := m.anyRequest(new(TransactionList), actionFetchUserTransfers, JsonObject{"Id": id})
	if err != nil {
		return nil, nil, err
	}
	return *(trs.(*TransactionList)), rateLimitInfo, nil
}
