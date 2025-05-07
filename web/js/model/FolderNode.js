import { FileNode } from "./FileNode.js";

export class FolderNode {
    /**
     * constructor for the FolderNode class
     * @param {number} id 
     * @param {string} filename 
     * @param {string} signature
     * @param {string} path
     * @param {FolderNode[]} subfolders
     * @param {FileNode[]} subfiles
     */
    constructor(id, filename, signature, path, subfolders = [], subfiles = []) {
        this.id = id;
        this.signature = signature;
        this.filename = filename;
        this.path = path; 
        this.contents = {
            subfolders: subfolders, // array of FolderNode objects
            subfiles: subfiles, // array of FileNode objects
        }
    }

    static fromJSON(json) {
        const subfolders = json.subfolders?.map(folder => FolderNode.fromJSON(folder)) ?? [];
        const subfiles = json.subfiles?.map(file => FileNode.fromJSON(file)) ?? [];
        return new FolderNode(json.id, json.filename, json.signature, json.path, subfolders, subfiles);
    }

    static async fetchFolderData(folderID) {
        try {
            const response = await fetch(`/api/v1/folder/${folderID}`);
            if (!response.ok) {
                throw new Error("Folder not found");
            }
            const data = await response.json();
            console.log("Folder data fetched")
            console.log(data)
            return FolderNode.fromJSON(data);
        } catch (error) {
            console.error("Error fetching folder data:", error);
            return null;
        }
    }

    static async fetchSubFolders(folderID) {
        try {
            const response = await fetch(`/api/v1/folders?from=${folderID}`);
            if (!response.ok) {
                throw new Error("Folders not found");
            }
            const data = await response.json();
            return data.map(folder => FolderNode.fromJSON(folder));
        } catch (error) {
            console.error("Error fetching folders:", error);
            alert("Folders not found");
            return null;
        }
    }
}