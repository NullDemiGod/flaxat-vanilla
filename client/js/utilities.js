function validateUsername(username) {
    const regex = /^[A-Za-z0-9_]{6,20}$/;

    return regex.test(username)
}

function validateEmail(email) {
    const regex = /^[A-Za-z0-9_.]{3,}@[A-Za-z0-9.-]{3,}\.[A-Za-z0-9]{2,}$/;
    
    return regex.test(email)
}
