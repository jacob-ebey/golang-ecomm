package runtime

import (
	"fmt"

	core "github.com/jacob-ebey/graphql-core"
)

func PrintError(err error) {
	fmt.Println(err)

	wrapped, ok := err.(*core.WrappedError)
	if ok && wrapped != nil && wrapped.InternalError != nil {
		PrintError(wrapped.InternalError)
	}
}
