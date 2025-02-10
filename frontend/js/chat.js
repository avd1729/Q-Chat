const currentUser = localStorage.getItem('currentUser');
const currentPath = window.location.pathname;

// Security check for correct path
if ((currentUser === 'Alice' && !currentPath.includes('/alice/')) || 
    (currentUser === 'Bob' && !currentPath.includes('/bob/'))) {
    window.location.href = '../index.html';
}

document.getElementById('currentUser').textContent = currentUser;

let protocolInitialized = false;
const initializeBtn = document.getElementById('initializeBtn');
const statusIndicator = document.getElementById('statusIndicator');

// Check initialization status on load
checkInitializationStatus();

// Function to check if protocol is already initialized
async function checkInitializationStatus() {
    try {
        const response = await fetch('http://localhost:8080/messages');
        const data = await response.json();

        if (data.messages && data.messages.length > 0) {
            protocolInitialized = true;
            updateInitializeButton();
        }
    } catch (error) {
        console.error('Failed to check initialization status:', error);
    }
}

function updateInitializeButton() {
    if (protocolInitialized) {
        initializeBtn.classList.add('initialized');
        initializeBtn.textContent = 'Protocol Initialized';
        statusIndicator.textContent = '● Protocol Active';
        statusIndicator.classList.add('status-initialized');
    }
}

async function initializeProtocol() {
    try {
        const response = await fetch('http://localhost:8080/initialize', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ bits: 256 })
        });

        const data = await response.json();
        if (data.message === 'Protocol initialized successfully') {
            protocolInitialized = true;
            updateInitializeButton();
        }
    } catch (error) {
        alert('Failed to initialize protocol:', error);
    }
}

function displayMessage(message, sent) {
    const messagesContainer = document.getElementById('messages');
    const messageElement = document.createElement('div');
    messageElement.className = `message ${sent ? 'sent' : 'received'}`;
    messageElement.textContent = message;
    messagesContainer.appendChild(messageElement);
    messagesContainer.scrollTop = messagesContainer.scrollHeight; // Auto-scroll to the latest message
}

async function decryptMessage(ciphertext) {
    try {
        const decryptResponse = await fetch('http://localhost:8080/decrypt', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ ciphertext: ciphertext })
        });

        const decryptData = await decryptResponse.json();
        return decryptData.plaintext;
    } catch (error) {
        console.error('Failed to decrypt message:', error);
        throw error;
    }
}

// Add a variable to track the last message count
let lastMessageCount = 0;
let lastProcessedMessageId = -1; // Add this to track processed messages

async function sendMessage() {
    if (!protocolInitialized) {
        alert('Please initialize the protocol first!');
        return;
    }

    const messageInput = document.getElementById('messageInput');
    const message = messageInput.value.trim();

    if (message) {
        try {
            const response = await fetch('http://localhost:8080/encrypt', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    plaintext: message,
                    sender: currentUser
                })
            });

            const data = await response.json();
            if (data.ciphertext) {
                messageInput.value = ''; // Clear the input field
                // Remove the manual pollMessages call - let the interval handle it
            }
        } catch (error) {
            alert('Failed to send message:', error);
        }
    }
}

async function pollMessages() {
    try {
        const response = await fetch('http://localhost:8080/messages');
        const data = await response.json();

        if (data.messages && data.messages.length > 0 && !protocolInitialized) {
            protocolInitialized = true;
            updateInitializeButton();
        }

        // Only update if we have new messages
        if (!data.messages || data.messages.length === lastMessageCount) {
            return; // No new messages, skip update
        }

        const messagesContainer = document.getElementById('messages');
        
        // If this is the first load or we're missing messages, do a full rebuild
        if (messagesContainer.children.length === 0) {
            messagesContainer.innerHTML = '';
            
            // Process all messages in order
            for (let i = 0; i < data.messages.length; i++) {
                const msg = data.messages[i];
                let plaintext;
                
                if (msg.sender === currentUser) {
                    plaintext = await decryptMessage(msg.ciphertext);
                    displayMessage(plaintext, true);
                } else {
                    try {
                        plaintext = await decryptMessage(msg.ciphertext);
                        displayMessage(plaintext, false);
                    } catch (error) {
                        console.error('Failed to decrypt message:', error);
                    }
                }
            }
        } else if (data.messages.length > lastMessageCount) {
            // Only process new messages
            for (let i = lastMessageCount; i < data.messages.length; i++) {
                const msg = data.messages[i];
                let plaintext;
                
                if (msg.sender === currentUser) {
                    plaintext = await decryptMessage(msg.ciphertext);
                    displayMessage(plaintext, true);
                } else {
                    try {
                        plaintext = await decryptMessage(msg.ciphertext);
                        displayMessage(plaintext, false);
                    } catch (error) {
                        console.error('Failed to decrypt message:', error);
                    }
                }
            }
        }

        lastMessageCount = data.messages.length;
    } catch (error) {
        console.error('Failed to fetch messages:', error);
    }
}

initializeBtn.addEventListener('click', initializeProtocol);
setInterval(pollMessages, 1000);

document.getElementById('messageInput').addEventListener('keypress', function(e) {
    if (e.key === 'Enter') {
        sendMessage();
    }
});