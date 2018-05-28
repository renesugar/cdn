package app

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/subutai-io/cdn/db"
)

type SearchRequest struct {
	fileID    string // files' UUID (or MD5)
	owner     string // files' owner username
	name      string // files' name within CDN
	repo      string // files' repository - either "apt", "raw", or "template"
	version   string // files' version
	tags      string // files' tags in format: "tag1,tag2,tag3"
	token     string // user's token
	operation string // operation type requested
}

// ParseRequest takes HTTP request and converts it into Request struct
func (r *SearchRequest) ParseRequest(req *http.Request) (err error) {
	r.fileID = req.URL.Query().Get("id")
	r.name = req.URL.Query().Get("name")
	r.owner = req.URL.Query().Get("owner")
	r.repo = strings.Split(req.RequestURI, "/")[3] // Splitting /kurjun/rest/repo/func into ["", "kurjun", "rest", "repo" (index: 3), "func"]
	r.version = req.URL.Query().Get("version")
	r.tags = req.URL.Query().Get("tags")
	r.token = req.URL.Query().Get("token")
	return
}

// BuildQuery constructs the query out of the existing parameters in SearchRequest
func (r *SearchRequest) BuildQuery() (query map[string]string) {
	if r.fileID != "" {
		query["fileID"] = r.fileID
	}
	if r.owner != "" {
		query["owner"] = r.owner
	}
	if r.name != "" {
		query["name"] = r.name
	}
	if r.repo != "" {
		query["repo"] = r.repo
	}
	if r.version != "" {
		query["version"] = r.version
	}
	if r.tags != "" {
		query["tags"] = r.tags
	}
	return
}

// SearchResult is a struct which return after search in db by parameters of SearchRequest
type SearchResult struct {
	fileID       string `json:"id,omitempty"`
	owner        string `json:",omitempty"`
	name         string `json:",omitempty"`
	filename     string `json",omitempty"`
	repo         string `json:"type,omitempty"`
	version      string `json:",omitempty"`
	scope        string `json:",omitempty"`
	md5          string `json:",omitempty"`
	sha256       string `json:",omitempty"`
	size         int    `json:",omitempty"`
	tags         string `json:",omitempty"`
	date         string `json:"upload-date-formatted,omitempty"`
	timestamp    string `json:"upload-date-timestamp,omitempty"`
	description  string `json:",omitempty"`
	architecture string `json:",omitempty"`
	parent       string `json:",omitempty"`
	pversion     string `json:"parent-version,omitempty"`
	powner       string `json:"parent-owner,omitempty"`
	prefsize     string `json:",omitempty"`
}

// BuildResult is make SearchResult struct from map of values
func BuildResult(res map[string]string) (searchRes SearchResult) {
	for k, v := range res {
		if k == "fileID" {
			searchRes.fileID = v
		}
		if k == "owner" {
			searchRes.owner = v
		}
		if k == "name" {
			searchRes.name = v
		}
		if k == "filename" {
			searchRes.filename = v
		}
		if k == "repo" {
			searchRes.repo = v
		}
		if k == "version" {
			searchRes.version = v
		}
		if k == "scope" {
			searchRes.scope = v
		}
		if k == "md5" {
			searchRes.md5 = v
		}
		if k == "sha256" {
			searchRes.sha256 = v
		}
		if k == "size" {
			sz, _ := strconv.Atoi(v)
			searchRes.size = sz
		}
		if k == "tags" {
			searchRes.tags = v
		}
		if k == "date" {
			searchRes.date = v
		}
		if k == "timestamp" {
			searchRes.timestamp = v
		}
		if k == "description" {
			searchRes.description = v
		}
		if k == "architecture" {
			searchRes.architecture = v
		}
		if k == "parent" {
			searchRes.parent = v
		}
		if k == "pversion" {
			searchRes.pversion = v
		}
		if k == "powner" {
			searchRes.powner = v
		}
		if k == "prefsize" {
			searchRes.prefsize = v
		}
	}
	return searchRes
}

func Retrieve(request SearchRequest) []SearchResult {
	query := request.BuildQuery()
	results, err := Search(query)
	if err != nil {

	}
	if request.operation == "info" {

	} else if request.operation == "list" {

	}
	return results
}

func GetFileInfo(id string) (info map[string]string, err error) {
	info["fileID"] = id
	err = db.DB.View(func(tx *bolt.Tx) error {
		file := tx.Bucket(db.MyBucket).Bucket([]byte(id))
		if file == nil {
			return fmt.Errorf("file %s not found", id)
		}
		owner := file.Bucket([]byte("owner"))
		key, _ := owner.Cursor().First()
		info["owner"] = string(key)
		info["name"] = string(file.Get([]byte("name")))
		repo := file.Bucket([]byte("type"))
		if repo != nil {
			key, _ = repo.Cursor().First()
			info["repo"] = string(key)
		}
		if len(info["repo"]) == 0 {
			return fmt.Errorf("couldn't find repo for file %s", id)
		}
		info["version"] = string(file.Get([]byte("version")))
		info["tags"] = string(file.Get([]byte("tag")))
		info["date"] = string(file.Get([]byte("date")))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return
}

func MatchQuery(file, query map[string]string) bool {
	for key, value := range query {
		if file[key] != value {
			return false
		}
	}
	return true
}

// Search return list of files with parameters like query
func Search(query map[string]string) (list []SearchResult, err error) {
	var sr SearchResult
	db.DB.View(func(tx *bolt.Tx) error {
		files := tx.Bucket(db.MyBucket)
		files.ForEach(func(k, v []byte) error {
			file, err := GetFileInfo(string(k))
			if err != nil {
				return err
			}
			if MatchQuery(file, query) {
				sr = BuildResult(file)
				list = append(list, sr)
			}
			return nil
		})
		return nil
	})
	return list, nil
}
