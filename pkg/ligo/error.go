package ligo

type LigoError string

func (le LigoError) Error() string {
	return string(le)
}
