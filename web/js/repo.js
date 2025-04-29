import markdownit from 'https://cdn.jsdelivr.net/npm/markdown-it@14.1.0/+esm'

import { User } from "./model/User.js";
import { Repository } from "./model/Repo.js";
import { Branch } from "./model/Branch.js";
import { Tag } from "./model/Tag.js";
import { HashTree } from "./model/HashTree.js";
import { FolderNode } from "./model/FolderNode.js";
import { FileNode } from "./model/FileNode.js";

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

const urlParams = new URLSearchParams(window.location.search)
const repoName = urlParams.get("repo")
const branchName = urlParams.get("branch") ?? "main"

const repo = await Repository.fetchRepositoryData(repoName, user?.id)
if (repo === null || repo === undefined) {
    console.error("Repository not found")
    console.info(repo)
    //window.location.href = "/home.html"
}

const branch = await Branch.fetchBranchData(branchName, repo.id)
if (branch === null || branch === undefined) {
    console.error("Branch not found")
    console.info(branch)
    //window.location.href = "/home.html"
}

const usernameTitle = document.getElementById("user-name")
const repoTitle = document.getElementById("repo-name")
const repoBranches = document.getElementById("repo-branches")
const repoContent = document.getElementById("repo-content")
const readmeContainer = document.getElementById("readme-container")

usernameTitle.innerText = user.name
repoTitle.innerText = repo.name

const tag = await Tag.fetchTagData(branch.head.id)
if (tag === null || tag === undefined) {
    console.error("Tag not found")
    console.info(tag)
    // window.location.href = "/home.html"
}

const tree = await HashTree.fetchHashTreeData(tag.tree.folder_node.id)
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
            window.location.href = `/dir.html?repo=${repoName}&branch=${branchName}&folder=${folder.filename}&parent=${tree.folder_node.id}`
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
.then(tree => {
    const readmeFile = tree.folder_node.contents.subfiles.find(file => file.filename === "README.md")
    if (readmeFile) {
        console.log("README.md found")
        const md = markdownit()
        const renderedReadme = md.render(readmeFile.content)
        readmeContainer.innerHTML = renderedReadme
    } else {
        console.log("README.md not found")
        readmeContainer.innerHTML = "<p>No README.md found</p>"
    }
})
.catch(error => {
    console.error("Error fetching HashTree data:", error);
    return null;
})
if (tree === null || tree === undefined) {
    console.error("HashTree not found")
    console.info(tree)
    // window.location.href = "/home.html"
}



const branches = await Branch.fetchBranches(repo.id)
if (branches === null || branches === undefined) {
    console.error("Branches not found")
    console.info(branches)
    // window.location.href = "/home.html"
}
if (branches.length === 0) {
    console.error("No branches found")
    console.info(branches)
    // window.location.href = "/home.html"
}
repoBranches.innerHTML = "" // clear the branches dropdown
branches.map(branch => branch.name)
    .forEach(name => {
        console.log("adding branch: " + name)
        let opt = document.createElement("option")
        opt.value = name
        opt.innerText = name
        if (name === branchName) {
            opt.selected = true
        }
        repoBranches.appendChild(opt)
    })





