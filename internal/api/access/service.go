package access

import (
	"github.com/Timofey335/auth/internal/service"
	descAccess "github.com/Timofey335/auth/pkg/access_v1"
)

// AccessImplementation - структура содержащая заглушки GRPC методов (в случае если они еще не созданы) и
// объект (интерфейс) сервисного слоя
type AccessImplementation struct {
	descAccess.UnimplementedAccessV1Server
	accessService service.AccessService
}

// NewAccessImplementation - конструктор, который возвращает объект сервисного слоя
func NewAccessImplementation(accessService service.AccessService) *AccessImplementation {
	return &AccessImplementation{
		accessService: accessService,
	}
}
