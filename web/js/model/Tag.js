import { HashTree } from "./HashTree.js";

export class Tag {
    /**
     * constructor for the Tag class
     * @param {number} id 
     * @param {string} name 
     * @param {string} signature 
     * @param {string} message 
     * @param {string} timestamp 
     * @param {HashTree} tree 
     * @param {Tag} parent 
     */
    constructor(id, name, signature, message, timestamp, tree, parent = null) {
        this.id = id;
        this.parent = parent; // parent tag
        this.name = name;
        this.signature = signature; 
        this.message = message; // tag message
        this.timestamp = timestamp; // tag timestamp
        this.tree = tree; // tag tree
    }

    /**
     * Converts a JSON object to a Tag object
     * @param {Object} json - The JSON object to convert
     * @param {number} json.id - The ID of the tag
     * @param {string} json.name - The name of the tag
     * @param {string} json.signature - The signature of the tag
     * @param {string} json.message - The message of the tag
     * @param {string} json.timestamp - The timestamp of the tag
     * @param {HashTree} json.tree - The tree of the tag
     * @param {Tag?} json.parent - The parent tag of the tag
     * @returns {Tag} - The Tag object
     */
    static fromJSON(json) {
        const tree = json.tree ? HashTree.fromJSON(json.tree) : null;
        const parent = json.parent ? Tag.fromJSON(json.parent) : null;
        return new Tag(json.id, json.name, json.signature, json.message, json.timestamp, tree, parent);
    }

    /**
     * Fetches the tag data from the server
     * @param {number} tagID - The ID of the tag to fetch
     * @returns {Tag} - The Tag object or null if not found
     */
    static async fetchTagData(tagID) {
        try {
            const response = await fetch(`/api/v1/tag/${tagID}`);
            if (!response.ok) {
                throw new Error("Tag not found");
            }
            const data = await response.json();
            console.log("Tag data fetched")
            console.log(data)
            return Tag.fromJSON(data);
        } catch (error) {
            console.error("Error fetching tag data:", error);
            alert("Tag not found");
            return null;
        }
    }
}