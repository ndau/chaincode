package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToBytes(t *testing.T) {
	bcheck(t, toBytes(1), "0100000000000000")
	bcheck(t, toBytes(-1), "FFFFFFFFFFFFFFFF")
	bcheck(t, toBytes(0x1122334455667788), "8877665544332211")
}

func TestToBytesU(t *testing.T) {
	bcheck(t, toBytesU(1), "0100000000000000")
	bcheck(t, toBytesU(0xFFFFFFFFFFFFFFFF), "FFFFFFFFFFFFFFFF")
	bcheck(t, toBytes(0x1122334455667788), "8877665544332211")
}

func TestZero(t *testing.T) {
	op, err := newPushOpcode("0")
	assert.Nil(t, err)
	b := op.bytes()
	bcheck(t, b, "20")
}

func TestOne(t *testing.T) {
	op, err := newPushOpcode("1")
	assert.Nil(t, err)
	b := op.bytes()
	bcheck(t, b, "2A")
}

func TestNegOne(t *testing.T) {
	op, err := newPushOpcode("-1")
	assert.Nil(t, err)
	b := op.bytes()
	bcheck(t, b, "2B")
}

func TestPushLobyte(t *testing.T) {
	op, err := newPushOpcode("17")
	assert.Nil(t, err)
	b := op.bytes()
	bcheck(t, b, "2111")
}

func TestPushHibyte(t *testing.T) {
	op, err := newPushOpcode("192")
	assert.Nil(t, err)
	b := op.bytes()
	bcheck(t, b, "21C0")
}

func TestPushNeg2(t *testing.T) {
	op, err := newPushOpcode("-2")
	assert.Nil(t, err)
	b := op.bytes()
	bcheck(t, b, "21FE")
}

func TestPushNeg200(t *testing.T) {
	op, err := newPushOpcode("-207")
	assert.Nil(t, err)
	b := op.bytes()
	bcheck(t, b, "2131")
}

func TestPush2Bytes(t *testing.T) {
	op, err := newPushOpcode("0x3478")
	assert.Nil(t, err)
	b := op.bytes()
	bcheck(t, b, "227834")
}

func TestPush3Bytes(t *testing.T) {
	op, err := newPushOpcode("0x125678")
	assert.Nil(t, err)
	b := op.bytes()
	bcheck(t, b, "23785612")
}

func TestPush4Bytes(t *testing.T) {
	op, err := newPushOpcode("0x12345678")
	assert.Nil(t, err)
	b := op.bytes()
	bcheck(t, b, "2478563412")
}

func TestPush5Bytes(t *testing.T) {
	op, err := newPushOpcode("0x123456780A")
	assert.Nil(t, err)
	b := op.bytes()
	bcheck(t, b, "250A78563412")
}

func TestPush6Bytes(t *testing.T) {
	op, err := newPushOpcode("0x12345678AA55")
	assert.Nil(t, err)
	b := op.bytes()
	bcheck(t, b, "2655AA78563412")
}

func TestPush7Bytes(t *testing.T) {
	op, err := newPushOpcode("0x12345678CAFEDA")
	assert.Nil(t, err)
	b := op.bytes()
	bcheck(t, b, "27DAFECA78563412")
}

func TestPush8BytesPositive(t *testing.T) {
	op, err := newPushOpcode("0x1BAD1DEACAFEBABE")
	assert.Nil(t, err)
	b := op.bytes()
	bcheck(t, b, "28BEBAFECAEA1DAD1B")
}

func TestPush64(t *testing.T) {
	op, err := newPush64("0xFFEEDDCCBBAA0011")
	assert.Nil(t, err)
	b := op.bytes()
	bcheck(t, b, "291100AABBCCDDEEFF")
}

func TestPushTimestamp(t *testing.T) {
	op, err := newPushTimestamp("2018-07-18T20:00:58Z")
	assert.Nil(t, err)
	b := op.bytes()
	bcheck(t, b, "2C801292DB9F0F0000")
}
