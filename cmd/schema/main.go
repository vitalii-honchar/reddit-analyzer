package main

import (
	"fmt"
	"reddit-analyzer/internal/lib/domain"

	"github.com/vitalii-honchar/go-agent/pkg/goagent/schema"
)

func main() {
	s := &domain.SearchResult{}

	res, err := schema.GenerateSchemaStr(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
