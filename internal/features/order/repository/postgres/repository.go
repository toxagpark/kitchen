package order_pg_repository

import core_postgres_pool "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/repository/postgres/pool"

type Repository struct {
	pool core_postgres_pool.Pool
}

func NewRepository(pool core_postgres_pool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}
