package repository

import (
	"database/sql"
	"frappuchino/internal/apperrors"
	"log/slog"
)

func checkRowsAffected(result sql.Result, id interface{}) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("Repository error from check row affected: failed to get rows affected", "id", id, "error", err)
		return err
	}
	if rowsAffected == 0 {
		return apperrors.ErrNotExistConflict
	}
	slog.Info("Repository info: rows affected", "id", id, "rows affected", rowsAffected)
	return nil
}
