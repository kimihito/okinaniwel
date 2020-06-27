package model

type Dog struct {
	Id         string `json:"id"`
	Breed      string `json:"bread"`
	Place      string `json:"place"`
	Image      string `json:"image"`
	Note       string `json:"note"`
	AddedOn    string `json:"addedOn"`
	DeadlineOn string `json:"deadlineOn"`
	Color      string `json:"color"`
	Sex        string `json:"sex"`
	Size       string `json:"size"`
	Collar     string `json:"collar"`
	URL        string `json:"url"`
}
