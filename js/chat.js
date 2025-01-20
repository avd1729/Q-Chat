// Add this at the beginning of chat.js
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
        console.error('Failed to check initialization status');
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
            body: JSON.stringify({ bits: '1010' })
        });
        const data = await response.json();
        protocolInitialized = true;
        updateInitializeButton();
    } catch (error) {
        alert('Failed to initialize protocol');
    }
}

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
            messageInput.value = '';
        } catch (error) {
            alert('Failed to send message');
        }
    }
}

function displayMessage(message, sent) {
    const messagesContainer = document.getElementById('messages');
    const messageElement = document.createElement('div');
    messageElement.className = `message ${sent ? 'sent' : 'received'}`;
    messageElement.textContent = message;
    messagesContainer.appendChild(messageElement);
    messagesContainer.scrollTop = messagesContainer.scrollHeight;
}

async function pollMessages() {
    try {
        const response = await fetch('http://localhost:8080/messages');
        const data = await response.json();
        
        if (data.messages && data.messages.length > 0 && !protocolInitialized) {
            protocolInitialized = true;
            updateInitializeButton();
        }
        
        const messages = data.messages;
        document.getElementById('messages').innerHTML = '';
        messages.forEach(msg => {
            displayMessage(msg.message, msg.sender === currentUser);
        });
    } catch (error) {
        console.error('Failed to fetch messages');
    }
}

initializeBtn.addEventListener('click', initializeProtocol);
setInterval(pollMessages, 1000);

document.getElementById('messageInput').addEventListener('keypress', function(e) {
    if (e.key === 'Enter') {
        sendMessage();
    }
});