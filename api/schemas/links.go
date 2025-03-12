package schemas

type Link struct {
	Self string `json:"self"`
}

type DocumentLink struct {
	Self  string `json:"self"`
	First string `json:"first,omitempty"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
}
