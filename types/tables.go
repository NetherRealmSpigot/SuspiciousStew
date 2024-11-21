package types

import "github.com/jackc/pgx/v5/pgtype"

type IpInfoResponse struct {
	Id int64 `json:"id"`
}

type SessionIdResponse struct {
	Id int64 `json:"id"`
}

type PlayerInfoResponse struct {
	UUID    pgtype.UUID `json:"uuid"`
	Name    string      `json:"name"`
	Version int         `json:"version"`
}
