package zk

var (
	Repo        = Repository{}
	CachingRepo = CachingRepository{Repo: Repo}
)
