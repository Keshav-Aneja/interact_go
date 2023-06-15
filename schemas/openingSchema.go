package schemas

import "github.com/lib/pq"

type OpeningCreateScheam struct { // from request
	Title       string         `json:"title" validate:"required,max=50"`
	Description string         `json:"description" validate:"required,max=500"`
	Tags        pq.StringArray `json:"tags"`
}

type OpeningEditScheam struct { // from request
	Description string         `json:"description" validate:"max=500"`
	Tags        pq.StringArray `json:"tags"`
}
