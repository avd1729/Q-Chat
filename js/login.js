function login(user) {
    localStorage.setItem('currentUser', user);
    if (user === 'Alice') {
        window.location.href = 'alice/chat.html';
    } else if (user === 'Bob') {
        window.location.href = 'bob/chat.html';
    }
}