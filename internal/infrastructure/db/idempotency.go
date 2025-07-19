package db

import (
	"context"

	dbsqlc "codematic/internal/infrastructure/db/sqlc"
	"codematic/internal/shared/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// GetIdempotencyRecord wraps the sqlc method and converts string IDs
func (db *DBConn) GetIdempotencyRecord(ctx context.Context, tenantID, key,
	endpoint, requestHash string) (
	record *dbsqlc.IdempotencyKey, found bool, err error) {
	tid, err := utils.StringToPgUUID(tenantID)
	if err != nil {
		return nil, false, err
	}
	params := dbsqlc.GetIdempotencyRecordParams{
		TenantID:       tid,
		IdempotencyKey: key,
		Endpoint:       endpoint,
		RequestHash:    requestHash,
	}
	rec, err := db.Queries.GetIdempotencyRecord(ctx, params)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &rec, true, nil
}

// SaveIdempotencyRecord wraps the sqlc method and converts string IDs
func (db *DBConn) SaveIdempotencyRecord(ctx context.Context, tenantID, userID,
	key, endpoint, requestHash string, responseBody []byte, statusCode int) error {
	tid, err := utils.StringToPgUUID(tenantID)
	if err != nil {
		return err
	}
	var uid pgtype.UUID
	if userID != "" {
		u, err := utils.StringToPgUUID(userID)
		if err == nil {
			uid = u
		}
	}
	id := uuid.New()
	params := dbsqlc.SaveIdempotencyRecordParams{
		ID:             pgtype.UUID{Bytes: id, Valid: true},
		TenantID:       tid,
		UserID:         uid,
		IdempotencyKey: key,
		Endpoint:       endpoint,
		RequestHash:    requestHash,
		ResponseBody:   responseBody,
		StatusCode:     pgtype.Int4{Int32: int32(statusCode), Valid: true},
	}
	_, err = db.Queries.SaveIdempotencyRecord(ctx, params)
	return err
}
