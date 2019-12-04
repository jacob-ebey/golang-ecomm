package schema

import (
	"github.com/graphql-go/graphql"

	core "github.com/jacob-ebey/graphql-core"
)

var UploadScalar = graphql.NewScalar(graphql.ScalarConfig{
	Name: "Upload",
	ParseValue: func(value interface{}) interface{} {
		if v, ok := value.(*core.MultipartFile); ok {
			return v
		}

		return nil
	},
})
