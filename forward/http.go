package forward

import (
	"bufio"
	"net/http"
)

func Handler(r *http.Request, f func(bs []byte)) {
	defer r.Body.Close()
	reader := bufio.NewReader(r.Body)
	for {
		bs, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		if f != nil {
			f(bs)
		}
	}
}
