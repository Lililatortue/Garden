import { Repository } from "./Repo.js";

export class User {

    /**
     * constructor for the User class
     * @param {number} id
     * @param {string} name
     * @param {string} email
     * @param {string} password
     * @param {Repository[]} repositories
     */
    constructor(id, name, email, password, repositories = []) {
        this.id = id;
        this.name = name;
        this.email = email;
        this.password = password;
        this.repositories = repositories;
    }

    /**
     * Fetches user data from the server using the provided username.
     * If the user is found, it returns a User object. Otherwise, it returns null.
     * @param {string} username 
     * @returns {Promise<User>} a promise that resolves to a User object
     */
    static async fetchUserData(username) {
        try {
            const response = await fetch(`/api/v1/user/${username}`);
            if (!response.ok) {
                throw new Error("User not found");
            }
            const data = await response.json();
            console.log("User data fetched")
            console.log(data)
            return User.fromJSON(data);
        } catch (error) {
            console.info(error)
            console.error("Error fetching user data:", error);
            alert("User not found");
            return null;
        }
    }

    /**
     * Creates a User object from a JSON object.
     * @param {Object} json - The JSON object to parse.
     * @return {User} The User object created from the JSON object.
     */
    static fromJSON(json) {
        const repositories = json.repositories.map(repo => Repository.fromJSON(repo));
        return new User(json.id, json.name, json.email, json.password, repositories);
    }

    /**
     * This function retrieves the user object from session storage and parses it into a User object.
     * @returns {User} the user object from session storage or null if not found
     */
    static getUserFromSessionStorage() {
        const userData = sessionStorage.getItem("user");
        if (userData) {
            const parsedData = JSON.parse(userData);
            return User.fromJSON(parsedData);
        }
        return null;
    }

    /**
     * This function saves the user object to session storage as json string.
     */
    setUserToSessionStorage() {
        sessionStorage.setItem("user", JSON.stringify(this));
    }

    /**
     * This function clears the user object from session storage.
     */
    static clearUserFromSessionStorage() {
        sessionStorage.removeItem("user");
    }

    /**
     * This function checks if the user is logged in by checking if the user object is present in session storage.
     * @returns {boolean} true if the user is logged in, false otherwise
     */
    static isLoggedIn() {
        return sessionStorage.getItem("user") !== null && sessionStorage.getItem("user") !== undefined;
    }

    /**
     * This function retrieves the user object from local storage and parses it into a User object.
     * @returns {User} the user object from local storage or null if not found
     */
    static logout() {
        User.clearUserFromLocalStorage();
        window.location.href = "/login.html"
    }

}