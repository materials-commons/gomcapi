package urlpath

import (
	"errors"
	"fmt"
	"net/url"
	"path"
)

func Join(weburl string, paths ...string) (string, error) {
	u, err := url.Parse(weburl)
	if err != nil {
		return "", errors.New("invalid url")
	}
	allPaths := append([]string{u.Path}, paths...)
	u.Path = path.Join(allPaths...)
	return fmt.Sprintf("%s", u), nil
}
