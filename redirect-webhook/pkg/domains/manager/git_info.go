package manager

import (
	"gopkg.in/src-d/go-git.v4"
	"os"
)

//func GetGitInfo() response.ManagerInfoResponse{
func GetGitInfo() (string, string, string, string) {

	path, _ := os.Getwd()
	r, _ := git.PlainOpen(path)
	ref, _ := r.Head()
	cIter, _ := r.Log(&git.LogOptions{From: ref.Hash(), All: false})
	commit, _ := cIter.Next()

	return commit.Hash.String(), commit.Author.String(), commit.Author.Email, commit.Author.When.String()

}
