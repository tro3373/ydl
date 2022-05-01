package worker

type Ctx struct {
	WorkDir   string
	QueueDir  string
	LibDir    string
	DoneDir   string
	YoutubeDl string
}

func NewCtx(work, queue, lib, done string) Ctx {
	return Ctx{
		WorkDir:   work,
		QueueDir:  queue,
		LibDir:    lib,
		DoneDir:   done,
		YoutubeDl: "youtube-dl",
	}
}
