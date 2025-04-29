import { User } from "./User.js";
import { Branch } from "./Branch.js";

export class Repository {

    /**
     * constructor for the Repository class
     * @param {number} id 
     * @param {string} name 
     * @param {User} owner 
     * @param {Branch[]} branches 
     */
    constructor(id, name, user_id, branches = []) {
        this.id = id;
        this.name = name;
        this.user_id = user_id;
        this.branches = branches;
    }

    /**
     * Converts a JSON object to a Repository instance
     * @param {Object} json - JSON object to convert to Repository
     * @returns {Repository} A new Repository object created from the JSON data
     */
    static fromJSON(json) {
        if (json.branches == null) {
            json.branches = [];
        }
        const branches = json.branches.map(branch => Branch.fromJSON(branch));
        return new Repository(json.id, json.name, json.user_id, branches);
    }

    /**
     * Fetches repository data from the server
     * @param {string} repoName - The name of the repository to fetch
     * @param {number} userId - The ID of the user who owns the repository
     * @returns {Promise<Repository>} A promise that resolves to a Repository object
     */
    static async fetchRepositoryData(repoName, userId) {
        try {
            const response = await fetch(`/api/v1/repo/${repoName}?user_id=${userId}`);
            if (!response.ok) {
                throw new Error("Repository not found");
            }
            const data = await response.json();
            console.log("Repository data fetched")
            console.log(data)
            return Repository.fromJSON(data);
        } catch (error) {
            console.error("Error fetching repository data:", error);
            alert("Repository not found");
            return null;
        }
    }
}