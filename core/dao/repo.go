package dao

import "sync"

var (
	repoPool = sync.Pool{
		New: func() interface{} {
			return &Repository{}
		},
	}
)

func Repo(r string) *Repository {
	repo := repoPool.Get().(*Repository)
	repo.Collection = r
	repoPool.Put(repo)
	return repo
}
