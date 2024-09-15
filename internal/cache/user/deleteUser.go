package user

import (
	"context"
	"strconv"
)

// DeleteUser - удаляет клиента из кэша
func (c *cacheImplementation) DeleteUser(ctx context.Context, id int64) error {
	idStr := strconv.Itoa(int(id))
	err := c.cacheClient.DeleteHashSet(ctx, idStr)
	if err != nil {
		return err
	}

	return nil
}
