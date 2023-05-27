package input

import "anio/input/localapp"

type Config struct {
	WebListener  any
	WebPollers   any
	LocalPollers *localapp.Config
}
