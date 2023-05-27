package input

import "anio/input/localapp"

type Config struct {
	WebListener  any              `json:"webListener"`
	WebPollers   any              `json:"webPollers"`
	LocalPollers *localapp.Config `json:"localPollers"`
}
