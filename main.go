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
	ciphertext string
	sender     string
}

// BB84Protocol represents the complete QKD protocol
type BB84Protocol struct {
	alice          Participant
	bob            Participant
	numberOfBits   int
	sharedKey      []int
	quantumChannel []int // Simulated quantum channel
}

// SecureChannel represents the communication channel between Alice and Bob
type SecureChannel struct {
	sharedKey []int
	messages  []Message
}

// NewBB84Protocol creates a new instance of the BB84 protocol
func NewBB84Protocol(bits int) *BB84Protocol {
	return &BB84Protocol{
		alice:        Participant{name: "Alice"},
		bob:          Participant{name: "Bob"},
		numberOfBits: bits,
	}
}

// NewSecureChannel creates a new secure channel using the shared key
func NewSecureChannel(key []int) *SecureChannel {
	return &SecureChannel{
		sharedKey: key,
		messages:  make([]Message, 0),
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
	bb84.quantumChannel = make([]int, bb84.numberOfBits)

	// Simulate quantum state preparation and measurement
	for i := 0; i < bb84.numberOfBits; i++ {
		// If bases match, Bob measures the correct value
		if bb84.alice.bases[i] == bb84.bob.bases[i] {
			bb84.quantumChannel[i] = bb84.alice.bits[i]
		} else {
			// If bases don't match, Bob gets a random result
			num, _ := rand.Int(rand.Reader, big.NewInt(2))
			bb84.quantumChannel[i] = int(num.Int64())
		}
	}
}

// generateSharedKey creates the final shared key from matching bases
func (bb84 *BB84Protocol) generateSharedKey() {
	bb84.sharedKey = make([]int, 0)
	for i := 0; i < bb84.numberOfBits; i++ {
		if bb84.alice.bases[i] == bb84.bob.bases[i] {
			bb84.sharedKey = append(bb84.sharedKey, bb84.alice.bits[i])
		}
	}
}

// calculateErrorRate computes the quantum bit error rate
func (bb84 *BB84Protocol) calculateErrorRate() float64 {
	errors := 0
	total := 0

	for i := 0; i < bb84.numberOfBits; i++ {
		if bb84.alice.bases[i] == bb84.bob.bases[i] {
			total++
			if bb84.alice.bits[i] != bb84.quantumChannel[i] {
				errors++
			}
		}
	}

	if total == 0 {
		return 0
	}
	return float64(errors) / float64(total)
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
	keyBytes := convertKeyToBytes(sc.sharedKey)
	plaintextBytes := []byte(plaintext)

	cipherBytes := xorBytes(plaintextBytes, keyBytes)
	ciphertext := base64.StdEncoding.EncodeToString(cipherBytes)

	msg := &Message{
		ciphertext: ciphertext,
		sender:     sender,
	}

	sc.messages = append(sc.messages, *msg)
	return msg, nil
}

// DecryptMessage decrypts a message using the shared key
func (sc *SecureChannel) DecryptMessage(msg *Message) (string, error) {
	keyBytes := convertKeyToBytes(sc.sharedKey)

	cipherBytes, err := base64.StdEncoding.DecodeString(msg.ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %v", err)
	}

	plaintextBytes := xorBytes(cipherBytes, keyBytes)
	return string(plaintextBytes), nil
}

// RunProtocol executes the complete BB84 protocol
func (bb84 *BB84Protocol) RunProtocol() error {
	// Step 1: Alice generates random bits and bases
	if err := bb84.alice.generateRandomBits(bb84.numberOfBits); err != nil {
		return fmt.Errorf("alice bits generation failed: %v", err)
	}
	if err := bb84.alice.generateRandomBases(bb84.numberOfBits); err != nil {
		return fmt.Errorf("alice bases generation failed: %v", err)
	}

	// Step 2: Bob generates random measurement bases
	if err := bb84.bob.generateRandomBases(bb84.numberOfBits); err != nil {
		return fmt.Errorf("bob bases generation failed: %v", err)
	}

	// Step 3: Simulate quantum transmission
	bb84.simulateQuantumTransmission()

	// Step 4: Generate shared key from matching bases
	bb84.generateSharedKey()

	return nil
}

func main() {
	// Initialize protocol with 256 bits for key generation
	bb84 := NewBB84Protocol(256)

	// Run BB84 protocol
	if err := bb84.RunProtocol(); err != nil {
		fmt.Printf("Protocol failed: %v\n", err)
		return
	}

	// Print QKD results
	fmt.Printf("\n=== QKD Protocol Results ===\n")
	fmt.Printf("Initial number of bits: %d\n", bb84.numberOfBits)
	fmt.Printf("Final key length: %d\n", len(bb84.sharedKey))
	fmt.Printf("Error rate: %.2f%%\n", bb84.calculateErrorRate()*100)

	// Create secure channel
	secureChannel := NewSecureChannel(bb84.sharedKey)

	// Demonstrate message exchange
	fmt.Printf("\n=== Secure Message Exchange ===\n")

	// Alice sends a message
	aliceMsg := "Hello Bob! This is a secret message from Alice."
	encryptedMsg, err := secureChannel.EncryptMessage(aliceMsg, "Alice")
	if err != nil {
		fmt.Printf("Encryption failed: %v\n", err)
		return
	}

	fmt.Printf("\nAlice's original message: %s\n", aliceMsg)
	fmt.Printf("Encrypted message: %s\n", encryptedMsg.ciphertext)

	// Bob decrypts Alice's message
	decryptedMsg, err := secureChannel.DecryptMessage(encryptedMsg)
	if err != nil {
		fmt.Printf("Decryption failed: %v\n", err)
		return
	}
	fmt.Printf("Bob decrypted the message: %s\n", decryptedMsg)

	// Bob sends a reply
	bobMsg := "Hi Alice! I received your secret message successfully!"
	encryptedReply, err := secureChannel.EncryptMessage(bobMsg, "Bob")
	if err != nil {
		fmt.Printf("Encryption failed: %v\n", err)
		return
	}

	fmt.Printf("\nBob's original message: %s\n", bobMsg)
	fmt.Printf("Encrypted reply: %s\n", encryptedReply.ciphertext)

	// Alice decrypts Bob's reply
	decryptedReply, err := secureChannel.DecryptMessage(encryptedReply)
	if err != nil {
		fmt.Printf("Decryption failed: %v\n", err)
		return
	}
	fmt.Printf("Alice decrypted the reply: %s\n", decryptedReply)
}
