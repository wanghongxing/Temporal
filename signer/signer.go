package signer

import (
	"crypto/ecdsa"
	"errors"
	"io/ioutil"
	"math/big"

	"github.com/RTradeLtd/Temporal/utils"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

/*
This is used to generated signed messages that can be submitted to a smart contract in order to process a payment.
This module is only used for the "frontend" web GUI. Regular API use will incorporate a payment channel style contract
*/

type PaymentSigner struct {
	Key *ecdsa.PrivateKey
}

type SignedMessage struct {
	H             [32]byte       `json:"h"`
	R             [32]byte       `json:"r"`
	S             [32]byte       `json:"s"`
	V             uint8          `json:"v"`
	Address       common.Address `json:"address"`
	PaymentMethod uint8          `json:"payment_method"`
	PaymentNumber *big.Int       `json:"payment_number"`
	ChargeAmount  *big.Int       `json:"charge_amount"`
}

// GeneratePaymentSigner is used to generate our helper struct for signing payments
// keyFilePath is the path to a key as generated by geth
func GeneratePaymentSigner(keyFilePath, keyPass string) (*PaymentSigner, error) {
	fileBytes, err := ioutil.ReadFile(keyFilePath)
	if err != nil {
		return nil, err
	}
	pk, err := keystore.DecryptKey(fileBytes, keyPass)
	if err != nil {
		return nil, err
	}
	return &PaymentSigner{Key: pk.PrivateKey}, nil
}

func (ps *PaymentSigner) GenerateSignedPaymentMessagePrefixed(ethAddress common.Address, paymentMethod uint8, paymentNumber, chargeAmountInWei *big.Int) (*SignedMessage, error) {
	//  return keccak256(abi.encodePacked(msg.sender, _paymentNumber, _paymentMethod, _chargeAmountInWei));
	hashToSign := utils.SoliditySHA3(
		utils.Address(ethAddress),
		utils.Uint256(paymentNumber),
		utils.Uint8(paymentMethod),
		utils.Uint256(chargeAmountInWei),
	)
	hashPrefixed := utils.SoliditySHA3WithPrefix(hashToSign)
	sig, err := crypto.Sign(hashPrefixed, ps.Key)
	if err != nil {
		return nil, err
	}
	var h, r, s [32]byte
	for k := range hashPrefixed {
		h[k] = hashPrefixed[k]
	}
	if len(h) > 32 || len(h) < 32 {
		return nil, errors.New("failed to parse h")
	}
	for k := range sig[0:64] {
		if k < 32 {
			r[k] = sig[k]
		}
		if k >= 32 {
			s[k-32] = sig[k]
		}
	}
	if len(r) != len(s) && len(r) != 32 {
		return nil, errors.New("failed to parse R+S")
	}

	msg := &SignedMessage{
		H:             h,
		R:             r,
		S:             s,
		V:             uint8(sig[64]) + 27,
		Address:       ethAddress,
		PaymentMethod: paymentMethod,
		PaymentNumber: paymentNumber,
		ChargeAmount:  chargeAmountInWei,
	}

	// Here we do an off-chain validation to ensure that when validated on-chain the transaction won't rever
	// however for some reason, the data isn't validating on-chain
	pub := ps.Key.PublicKey
	compressedKey := crypto.CompressPubkey(&pub)
	valid := crypto.VerifySignature(compressedKey, hashPrefixed, sig[0:64])
	if !valid {
		return nil, errors.New("failed to validate signature off-chain")
	}
	return msg, nil
}
