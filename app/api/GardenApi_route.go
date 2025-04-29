package api

import (
	"encoding/json"
	"garden/types"
	"log"
	"net/http"
	"strconv"
)

func (api *GardenApi) setRoutes() {
	api.setPushedRoute()
	api.setTestRoutes()
	api.setGetUsersRoute()
	api.setGetRepositoriesRoute()
	api.setGetBrancheRoute()
	api.setGetReposBranchesRoute()
	api.setGetTagsRoute()
}

func (api *GardenApi) setTestRoutes() {

	test := func(w http.ResponseWriter, r *http.Request) {
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

		w.WriteHeader(http.StatusOK)
	}

	api.HandleFunc("GET /api/test", test)
	api.HandleFunc("POST /api/test", test)
}

func (api *GardenApi) setPushedRoute() {
	api.HandleFunc("POST /api/v1/{username}/{repo}/{branch}", func(w http.ResponseWriter, r *http.Request) {
		var (
			username   = r.PathValue("username")
			repoName   = r.PathValue("repo")
			branchName = r.PathValue("branch")
		)
		w.Header().Set("Content-Type", "application/json")
		log.Println("Pushed route called for user", username, " repo", repoName, " branch", branchName)

		var tag types.GardenTag
		if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
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
				Head: &tag,
			}
			_, err := api.repoManager.AddBranch(branch, repo.ID)
			if err != nil {
				log.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"error": "Error adding non-existing branch to database",
				})
				return
			}
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
		tag.Parent = branch.Head
		branch.Head = &tag
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
	})
}

func (api *GardenApi) setGetUsersRoute() {
	api.HandleFunc("GET /api/v1/user/{username}", func(w http.ResponseWriter, r *http.Request) {
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
	})
}

func (api *GardenApi) setGetRepositoriesRoute() {
	api.HandleFunc("GET /api/v1/repo/{repoName}", func(w http.ResponseWriter, r *http.Request) {
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
	})
}

func (api *GardenApi) setGetBrancheRoute() {
	api.HandleFunc("GET /api/v1/branch/{branchname}", func(w http.ResponseWriter, r *http.Request) {
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

	})
}

func (api *GardenApi) setGetReposBranchesRoute() {
	api.HandleFunc("GET /api/v1/branches", func(w http.ResponseWriter, r *http.Request) {
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
	})
}

func (api *GardenApi) setGetTagsRoute() {
	api.HandleFunc("GET /api/v1/tag/{id}", func(w http.ResponseWriter, r *http.Request) {
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
	})
}
