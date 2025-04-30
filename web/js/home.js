import { User } from "./model/User.js"
import { Repository } from "./model/Repo.js"
import { Branch } from "./model/Branch.js"
import { Tag } from "./model/Tag.js"
import { HashTree } from "./model/HashTree.js"
import { FolderNode } from "./model/FolderNode.js"
import { FileNode } from "./model/FileNode.js"

const myRepository = document.getElementById("my-repositories")
const usernameTitle = document.getElementById("user-name")
const usernameLinks = document.getElementsByClassName("username-link")

/** @type {User} */
var user = null
if (User.isLoggedIn()) {
    user = User.getUserFromSessionStorage()
} else {
    user = await User.fetchUserData(username)
    if (user) {
        user.setUserToLocalStorage()
    } else {
        alert("User not found")
        window.location.href = "/login.html"
    }
}

usernameTitle.innerText = user.name.charAt(0).toUpperCase() + user.name.slice(1)
for (let usernameLink of usernameLinks) {
    usernameLink.href = `./home.html`
}

if (user.repositories.length === 0) {
    console.log("No repositories found for this user.");
    let li = document.createElement("li");
    li.innerText = "No repositories found";
    myRepository.append(li);
} else {
    myRepository.innerHTML = ""; // Clear previous content
    user.repositories.forEach(repo => {
    console.log("adding repo: " + repo.name)
    console.info(repo)
    let anchor = document.createElement("a");
    anchor.href = "repo.html?username=" + user.name + "&repo=" + repo.name;
    anchor.innerText = repo.name;
    let li = document.createElement("li");
    li.appendChild(anchor);
    myRepository.append(li);
    })
}



