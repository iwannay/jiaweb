package jiaweb

import (
	"github.com/iwannay/jiaweb/config"

	"github.com/iwannay/jiaweb/express"
)

type JiaWeb struct {
	Express *express.Express
	Config  *config.Config
}

func ListenAndServe() {

}

func New() {
	return &JiaWeb{}
}
