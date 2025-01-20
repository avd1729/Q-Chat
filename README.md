# Q-Chat

Q-Chat is a one-to-one messaging platform that implements BB84 quantum protocol for secure key exchange between users (Alice and Bob) before message encryption and decryption.

## Technology Stack

### Backend
- Go programming language
- Gin web framework

### Frontend
- HTML
- CSS 
- Vanilla JavaScript

## Core Features

- Implementation of BB84 quantum key distribution protocol for secure key exchange between Alice and Bob
- Message encryption and decryption using the shared key
- Simple user interface for message exchange

## BB84 Protocol Implementation

The platform uses BB84 quantum protocol for key distribution:
1. Alice generates random bits and random basis (rectilinear or diagonal)
2. Alice sends qubits to Bob based on these bits and basis
3. Bob measures received qubits using randomly chosen basis
4. Alice and Bob publicly compare their basis choices
5. They keep only the bits where they used the same basis
6. These matching bits become their shared secret key


## Setup and Installation

1. Clone the repository
2. Set up the Go backend
3. Start the frontend

## Usage

1. Alice and Bob connect to the platform
2. They perform key exchange using BB84 protocol
3. Once they have a shared key, they can send encrypted messages
