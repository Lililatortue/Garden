import { Tag } from "./Tag.js";

export class Branch {
    /**
     * constructor for the Branch class
     * @param {number} id 
     * @param {string} name 
     * @param {Tag} head 
     */
    constructor(id, name, head) {
        this.id = id;
        this.name = name;
        this.head = head; // commit id
    }

    /**
     * Converts a JSON object to a Branch object
     * @param {Object} json 
     * @returns {Branch} A Branch object from the JSON object
     */
    static fromJSON(json) {
        const head = json.head ? Tag.fromJSON(json.head) : null;
        return new Branch(json.id, json.name, head);
    }

    /**
     * Fetches branch data from the server
     * @param {string} branchName
     * @param {number} repoID
     * @returns {Promise<Branch>} A promise that resolves to a Branch object
     */
    static async fetchBranchData(branchName, repoID) {
        try {
            const response = await fetch(`/api/v1/branch/${branchName}?from=${repoID}`);
            if (!response.ok) {
                throw new Error("Branch not found");
            }
            const data = await response.json();
            console.log("Branch data fetched")
            console.log(data)
            return Branch.fromJSON(data);
        } catch (error) {
            console.error("Error fetching branch data:", error);
            alert("Branch not found");
            return null;
        }
    }

    /**
     * Fetches all branches from the server for a given repository
     * @param {number} repoID
     * @returns {Promise<Branch[]>} A promise that resolves to an array of Branch objects
     */
    static async fetchBranches(repoID) {
        try {
            const response = await fetch(`/api/v1/branches?from=${repoID}`);
            if (!response.ok) {
                throw new Error("Branches not found");
            }
            const data = await response.json();
            return data.map(branch => Branch.fromJSON(branch));
        } catch (error) {
            console.error("Error fetching branches:", error);
            alert("Branches not found");
            return null;
        }
    }
}