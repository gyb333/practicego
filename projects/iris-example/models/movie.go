package models

import "iris-example/datamodels"

type Movie struct {
	datamodels.Movie
}

func (m Movie) Validate() (Movie, error) {
	/* do some checks and return an error if that Movie is not valid */
	return m,nil
}