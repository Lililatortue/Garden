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
} else {
    window.location.href = "/login.html"
}
console.log(user)

const urlParams = new URLSearchParams(window.location.search)
const repoName = urlParams.get("repo")
const branchName = urlParams.get("branch") ?? "main"

const repo = await Repository.fetchRepositoryData(repoName, user.id)
if (repo === null || repo === undefined) {
    alert("Repository not found")
    //window.location.href = "/home.html"
}

const branch = await Branch.fetchBranchData(branchName, repo.id)
if (branch === null || branch === undefined) {
    alert("Branch not found")
    //window.location.href = "/home.html"
}

const usernameTitle = document.getElementById("user-name")
const repoTitle = document.getElementById("repo-name")
const repoBranches = document.getElementById("repo-branches")
const repoContent = document.getElementById("repo-content")


const tag = await Tag.fetchTagData(branch.tag_id)
if (tag === null || tag === undefined) {
    alert("Tag not found")
    //window.location.href = "/home.html"
}

const tree = await HashTree.fetchHashTreeData(tag.tree.folder_node.id)
if (tree === null || tree === undefined) {
    alert("HashTree not found")
    //window.location.href = "/home.html"
}



const branches = await Branch.fetchBranches(repo.id)
if (branches === null || branches === undefined) {
    alert("Branches not found")
    //window.location.href = "/home.html"
}
if (branches.length === 0) {
    alert("No branches found")
    //window.location.href = "/home.html"
}
repoBranches.innerHTML = "" // clear the branches dropdown
branches.map(branch => branch.name)
    .forEach(name => {
        let opt = document.createElement("option")
        opt.value = name
        opt.innerText = name
        if (name === branchName) {
            opt.selected = true
        }
        repoBranches.appendChild()
    })



usernameTitle.innerText = user.name
repoTitle.innerText = repo.name


