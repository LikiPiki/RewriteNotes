package postgres

import (
	"github.com/likipiki/RewriteNotes/app"
)

func Install(services ...app.UserService) {
	for _, s := range services {
		s.Install()
	}
}
