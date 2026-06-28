package ports

import "github.com/ruandg/microservices/shipping/internal/application/core/domain"

type DBPort interface {
	Save(*domain.Shipping) error
	Get(id string) (domain.Shipping, error)
}
