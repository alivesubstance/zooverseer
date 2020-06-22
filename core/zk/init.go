package zk

var (
	Repo        = Repository{}
	CachingRepo = CachingRepository{Repo: Repo}
)

func Close() {
	ConnCache.InvalidateAll()
	CachingRepo.Close()
}
