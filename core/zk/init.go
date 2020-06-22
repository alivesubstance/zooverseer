package zk

var (
	Repo        = Repository{}
	CachingRepo = CachingRepository{Repo: Repo}
)

func Reset() {
	ConnCache.InvalidateAll()
	CachingRepo.Close()
}
