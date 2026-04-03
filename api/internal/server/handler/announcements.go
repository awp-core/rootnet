package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/cortexia/rootnet/api/internal/db/gen"
)

// Announcement represents a protocol announcement/message
type Announcement struct {
	ID        int64            `json:"id"`
	ChainID   int64            `json:"chainId"`
	Title     string           `json:"title"`
	Content   string           `json:"content"`
	Category  string           `json:"category"`
	Priority  int              `json:"priority"` // 0=info, 1=warning, 2=critical
	Active    bool             `json:"active"`
	CreatedAt time.Time        `json:"createdAt"`
	ExpiresAt *time.Time       `json:"expiresAt,omitempty"`
	Metadata  *json.RawMessage `json:"metadata,omitempty"`
}

// AnnouncementCreate is the request body for creating an announcement
type AnnouncementCreate struct {
	ChainID   int64            `json:"chainId"`
	Title     string           `json:"title"`
	Content   string           `json:"content"`
	Category  string           `json:"category"`
	Priority  int              `json:"priority"`
	ExpiresAt *time.Time       `json:"expiresAt,omitempty"`
	Metadata  *json.RawMessage `json:"metadata,omitempty"`
}

// listAnnouncements returns active, non-expired announcements
// GET /api/announcements?chainId=&category=&limit=&offset=
func listAnnouncements(db gen.DBTX) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chainID, _ := strconv.ParseInt(r.URL.Query().Get("chainId"), 10, 64)
		category := r.URL.Query().Get("category")
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		if limit <= 0 || limit > 100 {
			limit = 20
		}

		query := `SELECT id, chain_id, title, content, category, priority, active, created_at, expires_at, metadata
			FROM announcements
			WHERE active = TRUE AND (expires_at IS NULL OR expires_at > NOW())`
		args := []interface{}{}
		argIdx := 1

		if chainID > 0 {
			query += ` AND (chain_id = 0 OR chain_id = $` + strconv.Itoa(argIdx) + `)`
			args = append(args, chainID)
			argIdx++
		}
		if category != "" {
			query += ` AND category = $` + strconv.Itoa(argIdx)
			args = append(args, category)
			argIdx++
		}

		query += ` ORDER BY priority DESC, created_at DESC LIMIT $` + strconv.Itoa(argIdx) + ` OFFSET $` + strconv.Itoa(argIdx+1)
		args = append(args, limit, offset)

		rows, err := db.Query(r.Context(), query, args...)
		if err != nil {
			annWriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "query failed"})
			return
		}
		defer rows.Close()

		results := []Announcement{}
		for rows.Next() {
			var a Announcement
			if err := rows.Scan(&a.ID, &a.ChainID, &a.Title, &a.Content, &a.Category, &a.Priority, &a.Active, &a.CreatedAt, &a.ExpiresAt, &a.Metadata); err != nil {
				continue
			}
			results = append(results, a)
		}

		annWriteJSON(w, http.StatusOK, results)
	}
}

// getAnnouncement returns a single announcement by ID
// GET /api/announcements/{id}
func getAnnouncement(db gen.DBTX) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			annWriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
			return
		}

		var a Announcement
		err = db.QueryRow(r.Context(),
			`SELECT id, chain_id, title, content, category, priority, active, created_at, expires_at, metadata
			FROM announcements WHERE id = $1`, id).
			Scan(&a.ID, &a.ChainID, &a.Title, &a.Content, &a.Category, &a.Priority, &a.Active, &a.CreatedAt, &a.ExpiresAt, &a.Metadata)
		if err != nil {
			annWriteJSON(w, http.StatusNotFound, map[string]string{"error": "announcement not found"})
			return
		}

		annWriteJSON(w, http.StatusOK, a)
	}
}

// getLLMContext returns announcements formatted for LLM consumption
// GET /api/announcements/llm-context?chainId=
func getLLMContext(db gen.DBTX) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chainID, _ := strconv.ParseInt(r.URL.Query().Get("chainId"), 10, 64)

		query := `SELECT title, content, category, priority, metadata
			FROM announcements
			WHERE active = TRUE AND (expires_at IS NULL OR expires_at > NOW())`
		args := []interface{}{}

		if chainID > 0 {
			query += ` AND (chain_id = 0 OR chain_id = $1)`
			args = append(args, chainID)
		}
		query += ` ORDER BY priority DESC, created_at DESC LIMIT 20`

		rows, err := db.Query(r.Context(), query, args...)
		if err != nil {
			annWriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "query failed"})
			return
		}
		defer rows.Close()

		type llmItem struct {
			Title    string           `json:"title"`
			Content  string           `json:"content"`
			Category string           `json:"category"`
			Priority string           `json:"priority"`
			Metadata *json.RawMessage `json:"metadata,omitempty"`
		}

		items := []llmItem{}
		for rows.Next() {
			var item llmItem
			var pri int
			var meta *json.RawMessage
			if err := rows.Scan(&item.Title, &item.Content, &item.Category, &pri, &meta); err != nil {
				continue
			}
			switch pri {
			case 0:
				item.Priority = "info"
			case 1:
				item.Priority = "warning"
			case 2:
				item.Priority = "critical"
			}
			item.Metadata = meta
			items = append(items, item)
		}

		annWriteJSON(w, http.StatusOK, map[string]interface{}{
			"protocol":      "AWP (Agent Working Protocol)",
			"announcements": items,
			"count":         len(items),
			"timestamp":     time.Now().UTC().Format(time.RFC3339),
		})
	}
}

// adminCreateAnnouncement creates a new announcement (admin only)
// POST /api/admin/announcements
func adminCreateAnnouncement(db gen.DBTX) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AnnouncementCreate
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 65536)).Decode(&req); err != nil {
			annWriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}
		if req.Title == "" || req.Content == "" {
			annWriteJSON(w, http.StatusBadRequest, map[string]string{"error": "title and content required"})
			return
		}
		if req.Category == "" {
			req.Category = "general"
		}

		var id int64
		err := db.QueryRow(r.Context(),
			`INSERT INTO announcements (chain_id, title, content, category, priority, expires_at, metadata)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			req.ChainID, req.Title, req.Content, req.Category, req.Priority, req.ExpiresAt, req.Metadata).
			Scan(&id)
		if err != nil {
			annWriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "insert failed"})
			return
		}

		annWriteJSON(w, http.StatusCreated, map[string]interface{}{"id": id})
	}
}

// adminUpdateAnnouncement updates an announcement (admin only)
// PUT /api/admin/announcements/{id}
func adminUpdateAnnouncement(db gen.DBTX) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			annWriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
			return
		}

		var req AnnouncementCreate
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 65536)).Decode(&req); err != nil {
			annWriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}

		_, err = db.Exec(r.Context(),
			`UPDATE announcements SET chain_id=$1, title=$2, content=$3, category=$4, priority=$5, expires_at=$6, metadata=$7 WHERE id=$8`,
			req.ChainID, req.Title, req.Content, req.Category, req.Priority, req.ExpiresAt, req.Metadata, id)
		if err != nil {
			annWriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "update failed"})
			return
		}

		annWriteJSON(w, http.StatusOK, map[string]string{"status": "updated"})
	}
}

// adminDeleteAnnouncement deactivates an announcement (admin only)
// DELETE /api/admin/announcements/{id}
func adminDeleteAnnouncement(db gen.DBTX) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			annWriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
			return
		}

		_, err = db.Exec(r.Context(), `UPDATE announcements SET active = FALSE WHERE id = $1`, id)
		if err != nil {
			annWriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "delete failed"})
			return
		}

		annWriteJSON(w, http.StatusOK, map[string]string{"status": "deactivated"})
	}
}

// writeJSON helper (reuses handler pattern)
func annWriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// RegisterAnnouncementRoutes registers public and admin announcement routes
func RegisterAnnouncementRoutes(r chi.Router, db gen.DBTX, adminAuth func(http.Handler) http.Handler) {
	// Public endpoints
	r.Route("/api/announcements", func(r chi.Router) {
		r.Get("/", listAnnouncements(db))
		r.Get("/llm-context", getLLMContext(db))
		r.Get("/{id}", getAnnouncement(db))
	})

	// Admin endpoints (behind auth middleware)
	r.Route("/api/admin/announcements", func(r chi.Router) {
		r.Use(adminAuth)
		r.Post("/", adminCreateAnnouncement(db))
		r.Put("/{id}", adminUpdateAnnouncement(db))
		r.Delete("/{id}", adminDeleteAnnouncement(db))
	})
}

// svcListAnnouncements is the RPC-compatible service method
func svcListAnnouncements(ctx context.Context, db gen.DBTX, chainID int64, category string, limit, offset int) ([]Announcement, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	query := `SELECT id, chain_id, title, content, category, priority, active, created_at, expires_at, metadata
		FROM announcements
		WHERE active = TRUE AND (expires_at IS NULL OR expires_at > NOW())`
	args := []interface{}{}
	argIdx := 1

	if chainID > 0 {
		query += ` AND (chain_id = 0 OR chain_id = $` + strconv.Itoa(argIdx) + `)`
		args = append(args, chainID)
		argIdx++
	}
	if category != "" {
		query += ` AND category = $` + strconv.Itoa(argIdx)
		args = append(args, category)
		argIdx++
	}
	query += ` ORDER BY priority DESC, created_at DESC LIMIT $` + strconv.Itoa(argIdx) + ` OFFSET $` + strconv.Itoa(argIdx+1)
	args = append(args, limit, offset)

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []Announcement{}
	for rows.Next() {
		var a Announcement
		if err := rows.Scan(&a.ID, &a.ChainID, &a.Title, &a.Content, &a.Category, &a.Priority, &a.Active, &a.CreatedAt, &a.ExpiresAt, &a.Metadata); err != nil {
			continue
		}
		results = append(results, a)
	}
	return results, nil
}
