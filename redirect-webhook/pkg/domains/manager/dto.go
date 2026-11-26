package manager

type ManagerHealthDB struct {
	Status string `json:"status"`
	Hello  string `json:"hello"`
}

type ManagerHealthDiskSpace struct {
	Status    string `json:"status"`
	Total     uint64 `json:"total"`
	Free      uint64 `json:"free"`
	Threshold uint64 `json:"threshold"`
}

type ManagerHealth struct {
	Status    string                 `json:"status"`
	DB        ManagerHealthDB        `json:"db"`
	DiskSpace ManagerHealthDiskSpace `json:"diskSpace"`
}

type ManagerInfoResponse struct {
	App *ManagerInfoResponseApp `json:"app"`

	Git *ManagerInfoResponseGitCommit `json:"git"`
}

type ManagerInfoResponseApp struct {
	Profile string `json:"profile"`

	Name string `json:"name"`

	Description string `json:"description"`

	Version string `json:"version"`
}

type ManagerInfoResponseGitCommit struct {
	Id string `json:"id"`

	Author string `json:"author"`

	Email string `json:"email"`

	Time string `json:"time"`
}
