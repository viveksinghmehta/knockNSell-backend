package routes

import db "knockNSell/db/gen"

type Server struct {
	q *db.Queries
}

func NewServer(q *db.Queries) *Server {
	return &Server{q: q}
}
