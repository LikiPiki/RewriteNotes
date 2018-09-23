package postgres

type Servicer interface {
	Install()
}

func Install(services ...Servicer) {
	for _, s := range services {
		s.Install()
	}
}
