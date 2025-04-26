package api

import (
	"encoding/json"
	"garden/data"
	"garden/data/sql"
	"garden/types"
	"log"
	"net/http"
	"strconv"
)

type Api struct {
	*http.Server
}

func (api *Api) Close() error {
	err := api.Handler.(*Mux).Close()
	if err != nil {
		return err
	}
	return api.Server.Close()
}

type Mux struct {
	*http.ServeMux
	repoManager data.Repo
}

func (mux *Mux) Close() error {
	return mux.repoManager.Close()
}

func NewMux() *Mux {
	db, err := sql.NewDBAccess()
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	mux := &Mux{
		ServeMux:    http.NewServeMux(),
		repoManager: *data.NewRepoWith(db),
	}

	mux.setRoutes()
	return mux
}

func NewApi(port string) *Api {
	return &Api{
		Server: &http.Server{
			Addr:    ":" + port,
			Handler: NewMux(),
		},
	}
}

func (api *Mux) setRoutes() {
	api.setPushedRoute()
}

func (api *Mux) setPushedRoute() {
	api.HandleFunc("POST /api/v1/{username}/{repo}/{branch}", func(w http.ResponseWriter, r *http.Request) {
		var (
			username   = r.PathValue("username")
			repoName   = r.PathValue("repo")
			branchName = r.PathValue("branch")
		)
		w.Header().Set("Content-Type", "application/json")
		log.Println("Pushed route")

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

		_, err = api.repoManager.UpdateBranchHead(branch)
		if err != nil {
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
