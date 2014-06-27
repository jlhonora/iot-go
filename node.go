package main

import (
	"encoding/json"
	"io"
)

type Node struct {
    Id			int		`json:"id"`
    Identifier	string	`json:"identifier"`
    Name		string	`json:"name"`
    CreatedAt   string	`json:"created_at"`
    UpdatedAt   string	`json:"updated_at"`
}

func Decode(r io.Reader) (x *Node, err error) {
    x = new(Node)
    err = json.NewDecoder(r).Decode(x)
    return
}
