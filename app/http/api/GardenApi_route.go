package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"garden/types"
	"log"
	"net/http"
	"strconv"
)

func (api *GardenApi) setRoutes() {
	api.HandleFunc("GET /api/test", api.handleTest)
	api.HandleFunc("POST /api/test", api.handleTest)
	api.HandleFunc("POST /api/v1/push/{username}/{repo}/{branch}", api.handlePush)
	api.HandleFunc("GET /api/v1/user/{username}", api.handleGetUser)
	api.HandleFunc("GET /api/v1/repo/{repoName}", api.handleGetRepository)
	api.HandleFunc("GET /api/v1/branch/{branchname}", api.handleGetBranch)
	api.HandleFunc("GET /api/v1/branches", api.handleGetBranchesFromRepo)
	api.HandleFunc("GET /api/v1/tag/{id}", api.handleGetTag)
	api.HandleFunc("GET /api/v1/hashtree/{id}", api.handleGetHashtree)
	api.HandleFunc("GET /api/v1/file/{id}", api.handleGetFile)
}

func (api *GardenApi) handleTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Write requestbody to response body
	err := r.Write(w)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Error writing response",
		})
		return
	}
}

func (api *GardenApi) handlePush(w http.ResponseWriter, r *http.Request) {
	var (
		username   = r.PathValue("username")
		repoName   = r.PathValue("repo")
		branchName = r.PathValue("branch")
	)
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"message": "Error writing response",
				"error":   r.(error).Error(),
			})
		}
	}()
	w.Header().Set("Content-Type", "application/json")
	log.Println("Pushed route called for user", username, " repo", repoName, " branch", branchName)

	var tag types.GardenTag
	if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": "Error decoding request body",
			"error":   err.Error(),
			"ok":      false,
		})
		return
	}
	marshed, _ := json.MarshalIndent(tag, "\n", "\t")
	log.Printf("Tag: %s", marshed)

	user, err := api.repoManager.ReadUserByUsername(username)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Error reading user from database",
		})
		return
	}

	repo := user.GetRepository(repoName)
	if repo == nil {
		log.Println("Repository not found")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Repository not found. Make sure you have the correct username and repository name.",
		})
		return
	}

	branch := repo.GetBranch(branchName)
	if branch == nil {
		branch = &types.Branch{
			Name: branchName,
			Head: tag,
		}
		branchID, err := api.repoManager.AddBranch(branch, repo.ID)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "Error adding non-existing branch to database",
			})
			return
		}
		branch.ID = branchID
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message": "Branch created",
			"user":    username,
			"repo":    repoName,
			"branch":  branchName,
			"id":      strconv.FormatInt(branch.ID, 10),
		})
		return
	}

	if branch.Head.Signature != tag.Parent.Signature {
		log.Println("Tag signature does not match branch head")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": "Tag signature does not match branch head. Make sure you have the correct tag signature.",
			"got":   tag.Signature,
			"head":  branch.Head,
		})
		return
	}
	tag.Parent = &branch.Head
	branch.Head = tag
	tagId, err := api.repoManager.AddTag(&tag)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Error adding tag to database",
		})
		return
	}

	if err := api.repoManager.UpdateBranchHead(branch); err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Error updating branch head",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Tag pushed",
		"tag":     tag.Signature,
		"user":    username,
		"repo":    repoName,
		"branch":  branchName,
		"id":      strconv.FormatInt(tagId, 10),
	})
}

func (api *GardenApi) handleGetUser(w http.ResponseWriter, r *http.Request) {
	var (
		username = r.PathValue("username")
	)
	w.Header().Set("Content-Type", "application/json")
	log.Println("Get users route called for user", username)

	user, err := api.repoManager.ReadUserByUsername(username)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Error reading user from database",
		})
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Error writing response",
		})
		return
	}

	log.Printf("User found: %#v\n", user)
}

func (api *GardenApi) handleGetRepository(w http.ResponseWriter, r *http.Request) {
	var (
		userId   = r.URL.Query().Get("user_id")
		repoName = r.PathValue("repoName")
		repo     *types.Repository
	)

	w.Header().Set("Content-Type", "application/json")
	log.Println("Get repositories route called for repo", repoName, "from user", userId)

	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Missing user id in query",
		})
		return
	}
	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid user id",
		})
		return
	}

	repo, err = api.repoManager.ReadRepository(repoName, id)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Error reading repository from database",
		})
	}

	if err := json.NewEncoder(w).Encode(repo); err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Error writing response",
		})
		return
	}
	log.Printf("Repository found: %#v\n", repo)
}

func (api *GardenApi) handleGetBranch(w http.ResponseWriter, r *http.Request) {
	var (
		repoID     = r.URL.Query().Get("from")
		branchName = r.PathValue("branchname")
	)
	w.Header().Set("Content-Type", "application/json")
	log.Println("Get branches route called for branch", branchName, " from repo", repoID)

	if repoID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Missing repository id in query",
		})
		return
	}
	id, err := strconv.ParseInt(repoID, 10, 64)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid repository id",
		})
		return
	}

	branch, err := api.repoManager.ReadBranch(branchName, id)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Error reading branch from database",
		})
		return
	}
	if err := json.NewEncoder(w).Encode(branch); err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Error writing response",
		})
		return
	}

	log.Printf("Branch found: %#v\n", branch)

}

func (api *GardenApi) handleGetBranchesFromRepo(w http.ResponseWriter, r *http.Request) {
	var (
		repoID = r.URL.Query().Get("from")
	)
	w.Header().Set("Content-Type", "application/json")
	log.Println("Get all tags for repo", repoID, " route called")

	if repoID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Missing repository id in query",
		})
		return
	}
	id, err := strconv.ParseInt(repoID, 10, 64)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid repository id",
		})
	}

	branches, err := api.repoManager.ReadBranchesOf(id)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Error reading branches from database",
		})
	}

	if err := json.NewEncoder(w).Encode(branches); err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Error writing response",
		})
		return
	}

	log.Printf("Branches found: %#v\n", branches)
}

func (api *GardenApi) handleGetTag(w http.ResponseWriter, r *http.Request) {
	var (
		tagId = r.PathValue("id")
	)
	w.Header().Set("Content-Type", "application/json")
	log.Println("Get tags route called for tag", tagId)

	if tagId == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Missing tag id in query",
		})
		return
	}

	id, err := strconv.ParseInt(tagId, 10, 64)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid tag id",
		})
		return
	}

	tag, err := api.repoManager.ReadTagBy(id)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Error reading tag from database",
		})
		return
	}

	if err := json.NewEncoder(w).Encode(tag); err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Error writing response",
		})
		return
	}

	log.Printf("Tag found: %#v\n", tag)
}

func (api *GardenApi) handleGetHashtree(w http.ResponseWriter, r *http.Request) {
	var (
		tagIdStr = r.PathValue("id")
		tree     *types.HashTree
	)

	w.Header().Set("Content-Type", "application/json")
	log.Println("Get hashtree route called for tag", tagIdStr)

	if tagIdStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Missing tag id in query",
		})
		return
	}

	tagId, err := strconv.ParseInt(tagIdStr, 10, 64)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid tag id",
		})
	}

	tree, err = api.repoManager.ReadTree(tagId)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Error reading tag from database",
		})
		return
	}

	if err := json.NewEncoder(w).Encode(tree); err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Error writing response",
		})
	}

	log.Printf("Tag found: %#v\n", tree)
}

func (api *GardenApi) handleHeadFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}

func (api *GardenApi) handleGetFile(w http.ResponseWriter, r *http.Request) {
	var (
		fileId = r.PathValue("id")
	)
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in handleGetFile", r)
			_ = json.NewEncoder(w).Encode(r)
		}
	}()

	api.handleHeadFile(w, r)

	log.Println("Get file route called for file", fileId)

	if fileId == "" {
		w.WriteHeader(http.StatusBadRequest)
		panic(map[string]any{
			"message": "Missing file id in query",
			"file":    fileId,
			"error":   errors.New("missing file id in query"),
		})
	}

	id, err := strconv.ParseInt(fileId, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(map[string]any{
			"message": "Invalid file id",
			"file":    fileId,
			"error":   errors.New("invalid file id"),
		})
	}

	file, err := api.repoManager.ReadFileByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			panic(map[string]any{
				"message": "File not found",
				"file":    fileId,
				"error":   err,
			})
		}
		w.WriteHeader(http.StatusInternalServerError)
		panic(map[string]any{
			"message": "Error reading file from database",
			"file":    fileId,
			"error":   err,
		})
	}

	if err := json.NewEncoder(w).Encode(file); err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Error writing response",
		})
		return
	}
}
