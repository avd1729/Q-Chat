package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
)

// Basis represents measurement basis (Z = 0, X = 1)
type Basis int

const (
	ZBasis Basis = iota // Computational basis {|0⟩, |1⟩}
	XBasis              // Hadamard basis {|+⟩, |-⟩}
)

// Participant represents either Alice or Bob
type Participant struct {
	bits  []int
	bases []Basis
	name  string
}

// Message represents an encrypted message
type Message struct {
	Ciphertext string `json:"ciphertext"`
	Sender     string `json:"sender"`
}

// BB84Protocol represents the complete QKD protocol
type BB84Protocol struct {
	Alice          Participant
	Bob            Participant
	NumberOfBits   int
	SharedKey      []int
	QuantumChannel []int
	SecureChannel  *SecureChannel
}

// SecureChannel represents the communication channel between Alice and Bob
type SecureChannel struct {
	SharedKey []int
	Messages  []Message
}

// NewBB84Protocol creates a new instance of the BB84 protocol
func NewBB84Protocol(bits int) *BB84Protocol {
	return &BB84Protocol{
		Alice:        Participant{name: "Alice"},
		Bob:          Participant{name: "Bob"},
		NumberOfBits: bits,
	}
}

// generateRandomBits creates random classical bits
func (p *Participant) generateRandomBits(n int) error {
	p.bits = make([]int, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(2))
		if err != nil {
			return fmt.Errorf("failed to generate random bit: %v", err)
		}
		p.bits[i] = int(num.Int64())
	}
	return nil
}

// generateRandomBases creates random measurement bases
func (p *Participant) generateRandomBases(n int) error {
	p.bases = make([]Basis, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(2))
		if err != nil {
			return fmt.Errorf("failed to generate random basis: %v", err)
		}
		p.bases[i] = Basis(num.Int64())
	}
	return nil
}

// simulateQuantumTransmission simulates quantum state preparation and transmission
func (bb84 *BB84Protocol) simulateQuantumTransmission() {
	bb84.QuantumChannel = make([]int, bb84.NumberOfBits)
	for i := 0; i < bb84.NumberOfBits; i++ {
		if bb84.Alice.bases[i] == bb84.Bob.bases[i] {
			bb84.QuantumChannel[i] = bb84.Alice.bits[i]
		} else {
			num, _ := rand.Int(rand.Reader, big.NewInt(2))
			bb84.QuantumChannel[i] = int(num.Int64())
		}
	}
}

// generateSharedKey creates the final shared key from matching bases
func (bb84 *BB84Protocol) generateSharedKey() {
	bb84.SharedKey = make([]int, 0)
	for i := 0; i < bb84.NumberOfBits; i++ {
		if bb84.Alice.bases[i] == bb84.Bob.bases[i] {
			bb84.SharedKey = append(bb84.SharedKey, bb84.Alice.bits[i])
		}
	}
}

// NewSecureChannel creates a new SecureChannel using the shared key
func NewSecureChannel(sharedKey []int) *SecureChannel {
	return &SecureChannel{
		SharedKey: sharedKey,
		Messages:  make([]Message, 0),
	}
}

// RunProtocol executes the complete BB84 protocol and initializes the secure channel
func (bb84 *BB84Protocol) RunProtocol() error {
	if err := bb84.Alice.generateRandomBits(bb84.NumberOfBits); err != nil {
		return fmt.Errorf("alice bits generation failed: %v", err)
	}
	if err := bb84.Alice.generateRandomBases(bb84.NumberOfBits); err != nil {
		return fmt.Errorf("alice bases generation failed: %v", err)
	}
	if err := bb84.Bob.generateRandomBases(bb84.NumberOfBits); err != nil {
		return fmt.Errorf("bob bases generation failed: %v", err)
	}
	bb84.simulateQuantumTransmission()
	bb84.generateSharedKey()

	// Initialize the SecureChannel using the shared key
	bb84.SecureChannel = NewSecureChannel(bb84.SharedKey)
	return nil
}

// convertKeyToBytes converts bit array to byte array
func convertKeyToBytes(key []int) []byte {
	byteLen := (len(key) + 7) / 8
	bytes := make([]byte, byteLen)

	for i := 0; i < len(key); i++ {
		if key[i] == 1 {
			bytes[i/8] |= 1 << uint(7-i%8)
		}
	}
	return bytes
}

// xorBytes performs XOR operation on byte slices
func xorBytes(a, b []byte) []byte {
	result := make([]byte, len(a))
	for i := range a {
		result[i] = a[i] ^ b[i%len(b)]
	}
	return result
}

// EncryptMessage encrypts a message using the shared key
func (sc *SecureChannel) EncryptMessage(plaintext string, sender string) (*Message, error) {
	keyBytes := convertKeyToBytes(sc.SharedKey)
	plaintextBytes := []byte(plaintext)

	cipherBytes := xorBytes(plaintextBytes, keyBytes)
	ciphertext := base64.StdEncoding.EncodeToString(cipherBytes)

	msg := &Message{
		Ciphertext: ciphertext,
		Sender:     sender,
	}

	sc.Messages = append(sc.Messages, *msg)
	return msg, nil
}

// DecryptMessage decrypts a message using the shared key
func (sc *SecureChannel) DecryptMessage(msg *Message) (string, error) {
	keyBytes := convertKeyToBytes(sc.SharedKey)

	cipherBytes, err := base64.StdEncoding.DecodeString(msg.Ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %v", err)
	}

	plaintextBytes := xorBytes(cipherBytes, keyBytes)
	return string(plaintextBytes), nil
}
