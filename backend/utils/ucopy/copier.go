package ucopy

import "github.com/gotidy/copy"

var (
	copiers = copy.New(copy.Skip())
)

type Post func(src, dst any)

func CpPost(src, dst any, post Post) {
	Cp(src, dst)
	
	if post != nil {
		post(src, dst)
	}
}

func Cp(src, dst any) {
	copier := copiers.Get(dst, src)
	copier.Copy(dst, src)
}
