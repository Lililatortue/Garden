

export class FileNode {
    /**
     * constructor for the FileNode class
     * @param {number} id 
     * @param {string} filename 
     * @param {string} signature
     * @param {string} path
     * @param {string} content
     */
    constructor(id, filename, signature, path, content = "") {
        this.id = id;
        this.signature = signature;
        this.filename = filename;
        this.path = path; // path to the file in the repository
        this.content = content; // contents of the file
    }

    static fromJSON(json) {
        return new FileNode(json.id, json.filename, json.signature, json.path, json.content);
    }

    static async fetchFileData(fileName, folderID) {
        try {
            const response = await fetch(`/api/v1/file/${fileName}?from=${folderID}`);
            if (!response.ok) {
                throw new Error("File not found");
            }
            const data = await response.json();
            console.log("File data fetched")
            console.log(data)
            return FileNode.fromJSON(data);
        } catch (error) {
            console.error("Error fetching file data:", error);
            alert("File not found");
            return null;
        }
    }

    static async fetchFileContent(fileID) {
        try {
            const response = await fetch(`/api/v1/file/content/${fileID}`);
            if (!response.ok) {
                throw new Error("File content not found");
            }
            const data = await response.json();
            return data.content;
        } catch (error) {
            console.error("Error fetching file content:", error);
            alert("File content not found");
            return null;
        }
    }

    static async fetchSubFiles(folderID) {
        try {
            const response = await fetch(`/api/v1/files?from=${folderID}`);
            if (!response.ok) {
                throw new Error("Files not found");
            }
            const data = await response.json();
            return data.map(file => FileNode.fromJSON(file));
        } catch (error) {
            console.error("Error fetching files:", error);
            return null;
        }
    }
}