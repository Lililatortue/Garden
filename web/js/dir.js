import { User } from "./model/User.js";
import { Repository } from "./model/Repo.js";
import { Branch } from "./model/Branch.js";
import { Tag } from "./model/Tag.js";
import { HashTree } from "./model/HashTree.js";
import { FolderNode } from "./model/FolderNode.js";
import { FileNode } from "./model/FileNode.js";

/** @type {User|null} */
var user = null
if (User.isLoggedIn()) {
    user = User.getUserFromSessionStorage()
    if (user === null || user === undefined) {
        console.error("User not found")
        window.location.href = "/login.html"
    }
} else {
    window.location.href = "/login.html"
}
console.log(user)

// Get URL parameters
const urlParams = new URLSearchParams(window.location.search)
const repoName = urlParams.get("repo")
const branchName = urlParams.get("branch") ?? "main"
const folderID = urlParams.get("folder")
const parentID = urlParams.get("parent")

// Get Element references
const usernameTitle = document.getElementById("user-name")
const repoTitle = document.getElementById("repo-name")
const repoBranches = document.getElementById("repo-branches")
const repoContent = document.getElementById("repo-content")
const codeLink = document.getElementById("code-link")


usernameTitle.innerText = user?.name
repoTitle.href = `/repo.html?username=${user?.name}&repo=${repoName}`
codeLink.href = `/repo.html?username=${user?.name}&repo=${repoName}`

const repo = await Repository.fetchRepositoryData(repoName, user?.id)
.then(repo => {
    if (!user.repositories.find(r => r.name === repo.name)) {
        user.repositories.push(repo)
        user.setUserToSessionStorage()
    }
    return repo
})
.then(repo => {
    repoTitle.innerText = repo.name
    return repo
})
.catch(err => {
    console.error("Error fetching repository data:", err)
    window.location.href = "/home.html"
})

/** @type {Branch} */
var currBranch = null
repo.branches = await Branch.fetchBranches(repo.id)
.then(branches => {
    repoBranches.innerHTML = ""
    branches.forEach(branch => {
        console.log(branch)
        const branchElement = document.createElement("option")
        branchElement.classList.add("dropdown-item")
        branchElement.value = branch.name
        branchElement.innerText = branch.name
        if (branch.name === branchName) {
            branchElement.selected = true
            currBranch = branch
        }
        branchElement.addEventListener("click", () => {
            window.location.href = `dir.html?repo=${repoName}&branch=${branch.name}&folder=${folderName}&parent=${parentID}`
        })
        branchElement.addEventListener("change", () => {
            window.location.href = `dir.html?repo=${repoName}&branch=${branch.name}&folder=${folderName}&parent=${parentID}`
        })        
        repoBranches.appendChild(branchElement)
    })
})
.catch(err => {
    console.error("Error fetching branches:", err)
    window.location.href = "/home.html"
})

Tag.fetchTagData(currBranch.head.id)
.then(tag => {
    console.log(tag)
    return tag
})
.catch(err => {
    console.error("Error fetching tag data:", err)
    return null
})

HashTree.fetchHashTreeData(folderID)
.then(tree => {
    console.log("HashTree data fetched")
    console.log(tree)
    tree.folder_node.contents.subfolders.forEach(folder => {
        console.log("subfolder: " + folder.filename)
        let fileListItem = document.createElement("div")
        fileListItem.className = "file-list-item folder"
        let folderIcon = document.createElement("span")       
        folderIcon.className = "folder-icon" 
        let svg = document.createElement("svg")
        svg.className = "octicon octicon-file-directory"
        svg.setAttribute("viewBox", "0 0 16 16")
        folderIcon.appendChild(svg)
        fileListItem.appendChild(folderIcon)
        let fileName = document.createElement("span")
        fileName.className = "file-name"
        fileName.innerText = folder.filename
        fileListItem.appendChild(fileName)
        fileListItem.addEventListener("click", () => {
            window.location.href = `/dir.html?repo=${repoName}&branch=${branchName}&folder=${folder.id}&parent=${tree.folder_node.id}`
        })
        repoContent.appendChild(fileListItem)
    })
    tree.folder_node.contents.subfiles.forEach(file => {
        console.log("file: " + file.filename)
        let fileListItem = document.createElement("div")
        fileListItem.className = "file-list-item file"
        let fileIcon = document.createElement("span")       
        fileIcon.className = "file-icon" 
        let svg = document.createElement("svg")
        svg.className = "octicon octicon-file"
        svg.setAttribute("viewBox", "0 0 16 16")
        fileIcon.appendChild(svg)
        fileListItem.appendChild(fileIcon)
        let fileName = document.createElement("span")
        fileName.className = "file-name"
        fileName.innerText = file.filename
        fileListItem.appendChild(fileName)
        repoContent.appendChild(fileListItem)
    })

    return tree
})
.catch(error => {
    console.error("Error fetching HashTree data:", error);
    return null;
})


