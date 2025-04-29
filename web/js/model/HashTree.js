import { FolderNode } from "./FolderNode.js";

export class HashTree {
    /**
     * constructor for the HashTree class
     * @param {FolderNode} folder_node 
     */
    constructor(folder_node) {

        /** @type {FolderNode} */
        this.folder_node = folder_node; // root folder node 
    }

    /**
     * Creates a HashTree object from a JSON object.
     * @param {Object} json - The JSON object to convert.
     * @param {number} json.id - The ID of the folder node.
     * @return {HashTree} The HashTree object.
     */
    static fromJSON(json) {
        const folder_node = FolderNode.fromJSON(json);
        return new HashTree(folder_node);
    }

    static async fetchHashTreeData(folderID) {
        try {
            const response = await fetch(`/api/v1/hashtree/${folderID}`);
            if (!response.ok) {
                throw new Error("HashTree not found");
            }
            const data = await response.json();
            console.log("HashTree data fetched")
            console.log(data)
            return HashTree.fromJSON(data);
        } catch (error) {
            console.error("Error fetching HashTree data:", error);
            alert("HashTree not found");
            return null;
        }
    }
}