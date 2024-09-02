package user

import (
	"github.com/Timofey335/auth/internal/service"
	desc "github.com/Timofey335/auth/pkg/auth_v1"
)

// Implementation - структура содержащая заглушки GRPC методов (в случае если они еще не созданы) и
// объект (интерфейс) сервисного слоя
type Implementation struct {
	desc.UnimplementedAuthV1Server
	userService service.UserService
}

// NewImplementation - конструктор, который возвращает объект сервисного слоя
func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{
		userService: userService,
	}
}
