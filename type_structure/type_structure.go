package type_structure

type Article struct {
	Id             string `json:"Id"`
	Title          string `json:"Title"`
	Desciption     string `json:"desc"`
	ContentArtcile string `json:"content"`
}

type Users struct {
	Username string `json:"username"`
	Pass     string `json:"pass"`
}
