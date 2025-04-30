import { User } from "./model/User.js";

if (User.isLoggedIn()) {
    window.location.href = "/home.html";
}

const form = document.querySelector('form');

form.addEventListener('submit', async function(event) {
    event.preventDefault(); // Prevent the default form submission

    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    const user = await User.fetchUserData(username);
    console.log(user)
    if (user) {
        if (user.password === password) {
            user.setUserToSessionStorage();
            window.location.href = "/home.html";
        } else {
            alert("Incorrect password. Please try again.");
        }
    } else {
        alert("User not found. Please check your username.");
    }
});