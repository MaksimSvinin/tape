package entity

type ReadFileQuery struct {
	Numbers []int  `query:"numbers"`
	Path    string `query:"path"`
}
