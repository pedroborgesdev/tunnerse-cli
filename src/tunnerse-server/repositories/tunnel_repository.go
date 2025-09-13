package repositories

import (
	"database/sql"
	"time"
	"tunnerse-server/database"
	"tunnerse-server/logger"
	"tunnerse-server/models"
)

type TunnelRepository struct {
	DB *database.Database
}

func NewTunnelRepository(db *database.Database) *TunnelRepository {
	return &TunnelRepository{DB: db}
}

// Create cria ou inicializa registros de Tunnel e Info
func (r *TunnelRepository) Create(tunnel *models.Tunnel, info *models.Info) error {
	tx, err := r.DB.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT OR IGNORE INTO Tunnel (ID, Port, Url, Domain, Active, CreatedAt)
		VALUES (?, ?, ?, ?, ?, ?)`,
		tunnel.ID, tunnel.Port, tunnel.Url, tunnel.Domain, tunnel.Active, tunnel.CreatedAt,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT OR IGNORE INTO Info (ID, Requests, Healthchecks, Warns, Errors)
		VALUES (?, ?, ?, ?, ?)`,
		info.ID, info.Requests, info.Healthchecks, info.Warns, info.Errors,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *TunnelRepository) GetPID(tunnelID string) (int, error) {
	var pid int
	err := r.DB.DB.QueryRow(`SELECT Pid FROM Info WHERE ID = ?`, tunnelID).Scan(&pid)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return pid, nil
}

func (r *TunnelRepository) GetTunnel(id string) (*models.Tunnel, error) {
	var t models.Tunnel
	err := r.DB.DB.QueryRow(`
		SELECT ID, Port, Url, Domain, Active
		FROM Tunnel WHERE ID = ?`, id).Scan(&t.ID, &t.Port, &t.Url, &t.Domain, &t.Active)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TunnelRepository) UpdateTunnelStatus(tunnelID string, active bool) error {
	_, err := r.DB.DB.Exec(`UPDATE Tunnel SET Active = ? WHERE ID = ?`, active, tunnelID)
	return err
}

func (r *TunnelRepository) UpdateRequestCount(id string) {
	go func() {
		tx, err := r.DB.DB.Begin()
		if err != nil {
			logger.Log("ERROR", "failed to begin transaction", []logger.LogDetail{{Key: "error", Value: err.Error()}})
			return
		}
		defer tx.Rollback()

		res, err := tx.Exec(`UPDATE Info SET Requests = Requests + 1 WHERE ID = ?`, id)
		if err != nil {
			logger.Log("ERROR", "failed to increment requests", []logger.LogDetail{{Key: "error", Value: err.Error()}})
			return
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			_, err = tx.Exec(`INSERT INTO Info (ID, Requests, Pid, Healthcheck, Errors) VALUES (?, 1, 0, 0, 0)`, id)
			if err != nil {
				logger.Log("ERROR", "failed to insert new info record", []logger.LogDetail{{Key: "error", Value: err.Error()}})
				return
			}
		}

		if err := tx.Commit(); err != nil {
			logger.Log("ERROR", "failed to commit transaction", []logger.LogDetail{{Key: "error", Value: err.Error()}})
			return
		}
	}()
}

func (r *TunnelRepository) UpdateHealthcheckCount(id string) {
	go func() {
		tx, err := r.DB.DB.Begin()
		if err != nil {
			logger.Log("ERROR", "failed to begin transaction", []logger.LogDetail{{Key: "error", Value: err.Error()}})
			return
		}
		defer tx.Rollback()

		res, err := tx.Exec(`UPDATE Info SET Healthchecks = Healthchecks + 1 WHERE ID = ?`, id)
		if err != nil {
			logger.Log("ERROR", "failed to increment healthchecks", []logger.LogDetail{{Key: "error", Value: err.Error()}})
			return
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			_, err = tx.Exec(`INSERT INTO Info (ID, Healthchecks, Pid, Errors) VALUES (?, 1, 0, 0)`, id)
			if err != nil {
				logger.Log("ERROR", "failed to insert new info record", []logger.LogDetail{{Key: "error", Value: err.Error()}})
				return
			}
		}

		if err := tx.Commit(); err != nil {
			logger.Log("ERROR", "failed to commit transaction", []logger.LogDetail{{Key: "error", Value: err.Error()}})
			return
		}
	}()
}

func (r *TunnelRepository) UpdateWarnCount(id string) {
	go func() {
		tx, err := r.DB.DB.Begin()
		if err != nil {
			logger.Log("ERROR", "failed to begin transaction", []logger.LogDetail{{Key: "error", Value: err.Error()}})
			return
		}
		defer tx.Rollback()

		res, err := tx.Exec(`UPDATE Info SET Warns = Warns + 1 WHERE ID = ?`, id)
		if err != nil {
			logger.Log("ERROR", "failed to increment warns", []logger.LogDetail{{Key: "error", Value: err.Error()}})
			return
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			_, err = tx.Exec(`INSERT INTO Info (ID, Warns, Pid, Errors) VALUES (?, 1, 0, 0)`, id)
			if err != nil {
				logger.Log("ERROR", "failed to insert new info record", []logger.LogDetail{{Key: "error", Value: err.Error()}})
				return
			}
		}

		if err := tx.Commit(); err != nil {
			logger.Log("ERROR", "failed to commit transaction", []logger.LogDetail{{Key: "error", Value: err.Error()}})
			return
		}
	}()
}

func (r *TunnelRepository) UpdateErrorCount(id string) {
	go func() {
		tx, err := r.DB.DB.Begin()
		if err != nil {
			logger.Log("ERROR", "failed to begin transaction", []logger.LogDetail{{Key: "error", Value: err.Error()}})
			return
		}
		defer tx.Rollback()

		res, err := tx.Exec(`UPDATE Info SET Errors = Errors + 1 WHERE ID = ?`, id)
		if err != nil {
			logger.Log("ERROR", "failed to increment errors", []logger.LogDetail{{Key: "error", Value: err.Error()}})
			return
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			_, err = tx.Exec(`INSERT INTO Info (ID, Errors, Pid) VALUES (?, 1, 0)`, id)
			if err != nil {
				logger.Log("ERROR", "failed to insert new info record", []logger.LogDetail{{Key: "error", Value: err.Error()}})
				return
			}
		}

		if err := tx.Commit(); err != nil {
			logger.Log("ERROR", "failed to commit transaction", []logger.LogDetail{{Key: "error", Value: err.Error()}})
			return
		}
	}()
}

func (r *TunnelRepository) ListTunnels() ([]*models.Tunnel, error) {
	rows, err := r.DB.DB.Query(`SELECT ID, Port, Url, Domain, Active FROM Tunnel`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tunnels []*models.Tunnel
	for rows.Next() {
		var t models.Tunnel
		if err := rows.Scan(&t.ID, &t.Port, &t.Url, &t.Domain, &t.Active); err != nil {
			return nil, err
		}
		tunnels = append(tunnels, &t)
	}
	return tunnels, nil
}

func (r *TunnelRepository) InfoTunnel(tunnelID string) (*struct {
	ID           string
	Port         int
	Url          string
	Domain       string
	Active       bool
	CreatedAt    time.Time
	Pid          int
	Requests     int
	Healthchecks int
	Warns        int
	Errors       int
}, error) {

	ti := &struct {
		ID           string
		Port         int
		Url          string
		Domain       string
		Active       bool
		CreatedAt    time.Time
		Pid          int
		Requests     int
		Healthchecks int
		Warns        int
		Errors       int
	}{}

	// Ler CreatedAt como string primeiro
	var createdAtStr string

	row := r.DB.DB.QueryRow(`
		SELECT t.ID, t.Port, t.Url, t.Domain, t.Active, t.CreatedAt,
		       i.Pid, i.Requests, i.Healthchecks, i.Warns, i.Errors
		FROM Tunnel t
		LEFT JOIN Info i ON t.ID = i.ID
		WHERE t.ID = ?
	`, tunnelID)

	err := row.Scan(
		&ti.ID, &ti.Port, &ti.Url, &ti.Domain, &ti.Active, &createdAtStr,
		&ti.Pid, &ti.Requests, &ti.Healthchecks, &ti.Warns, &ti.Errors,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Converter string para time.Time usando RFC3339
	ti.CreatedAt, err = time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		// Caso falhe, loga e deixa zero
		logger.Log("ERROR", "failed to parse CreatedAt", []logger.LogDetail{
			{Key: "error", Value: err.Error()},
			{Key: "value", Value: createdAtStr},
		})
		ti.CreatedAt = time.Time{}
	}

	return ti, nil
}
