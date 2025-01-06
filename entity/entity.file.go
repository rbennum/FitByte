package entity

import "database/sql"

type File struct {
	FileId   string         `json:"fileid"`
	FileURI  sql.NullString `json:"fileuri"`
	FileName sql.NullString `json:"filename"`
}
