package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gated-community-api/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SecurityLogRecord struct {
	ID         string     `json:"id"`
	PersonType string     `json:"person_type"`
	PersonName string     `json:"person_name"`
	UnitNumber string     `json:"unit_number,omitempty"`
	TimeIn     time.Time  `json:"time_in"`
	TimeOut    *time.Time `json:"time_out"`
	GuardName  string     `json:"guard_name"`
	Remarks    string     `json:"remarks"`
}

func HandleListSecurityLogs(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		communityID := r.Header.Get("X-Community-ID")
		if communityID == "" {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Community ID is required"}, nil)
			return
		}

		// Parse Query Params
		query := r.URL.Query()
		personType := query.Get("type")
		startDate := query.Get("start_date")
		endDate := query.Get("end_date")
		sortBy := query.Get("sort_by")
		order := strings.ToUpper(query.Get("order"))
		limit, _ := strconv.Atoi(query.Get("limit"))
		page, _ := strconv.Atoi(query.Get("page"))

		if limit == 0 {
			limit = 20
		}
		if page == 0 {
			page = 1
		}
		if order != "DESC" {
			order = "ASC"
		}
		if sortBy == "" {
			sortBy = "time_in"
		}

		offset := (page - 1) * limit

		// Dynamic SQL Construction
		whereClauses := []string{"sl.community_id = $1"}
		args := []any{communityID}
		argCount := 2

		if personType != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("sl.person_type = $%d", argCount))
			args = append(args, personType)
			argCount++
		}
		if startDate != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("sl.time_in >= $%d", argCount))
			args = append(args, startDate)
			argCount++
		}
		if endDate != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("sl.time_in <= $%d", argCount))
			args = append(args, endDate)
			argCount++
		}

		sql := fmt.Sprintf(`
			SELECT 
				sl.id, sl.person_type, sl.time_in, sl.time_out, sl.remarks,
				COALESCE(v.name, r.first_name || ' ' || r.last_name, vn.company_name) as person_name,
				u.unit_number,
				gu.first_name || ' ' || gu.last_name as guard_name
			FROM security_logs sl
			LEFT JOIN visitors v ON sl.visitor_id = v.id
			LEFT JOIN residents r ON sl.resident_id = r.id
			LEFT JOIN vendors vn ON sl.vendor_id = vn.id
			LEFT JOIN units u ON sl.unit_id = u.id
			LEFT JOIN users gu ON sl.guard_id = gu.id
			WHERE %s
			ORDER BY %s %s
			LIMIT $%d OFFSET $%d`,
			strings.Join(whereClauses, " AND "), sortBy, order, argCount, argCount+1)

		args = append(args, limit, offset)

		rows, err := db.Query(r.Context(), sql, args...)
		if err != nil {
			domain.WriteJSON(w, http.StatusInternalServerError, domain.Envelope{Success: false, Message: "Failed to fetch logs"}, nil)
			return
		}
		defer rows.Close()

		logs := []SecurityLogRecord{}
		for rows.Next() {
			var l SecurityLogRecord
			err := rows.Scan(
				&l.ID, &l.PersonType, &l.TimeIn, &l.TimeOut, &l.Remarks,
				&l.PersonName, &l.UnitNumber, &l.GuardName,
			)
			if err != nil {
				continue
			}
			logs = append(logs, l)
		}

		domain.WriteJSON(w, http.StatusOK, domain.Envelope{
			Success: true,
			Data:    logs,
		}, nil)
	}
}
