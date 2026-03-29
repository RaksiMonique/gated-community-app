package handlers

import (
	"net/http"
	"time"

	"gated-community-api/internal/domain"
	"gated-community-api/internal/middleware"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardData struct {
	Metrics struct {
		OccupancyRate         float64 `json:"occupancy_rate"`
		PaymentCollectionRate float64 `json:"payment_collection_rate"`
		OpenMaintenanceCount  int     `json:"open_maintenance_count"`
		OutstandingBalance    float64 `json:"outstanding_balance"`
		VisitorVolume7D       int     `json:"visitor_volume_7d"`
	} `json:"metrics"`
	RecentActivity struct {
		Announcements []domain.Announcement `json:"announcements"`
		SecurityLogs  []SecurityLogRecord   `json:"security_logs"`
	} `json:"recent_activity"`
}

func HandleGetAdminDashboard(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		communityID, ok := r.Context().Value(middleware.CommunityContextKey).(string)
		if !ok {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Missing community context"}, nil)
			return
		}
		currentMonth := time.Now().Format("2006-01")

		var data DashboardData

		// 1. Occupancy Rate: (Units with active residents / Total Units)
		occupancyQuery := `
			SELECT 
				(SELECT COUNT(DISTINCT unit_id) FROM residents WHERE community_id = $1 AND status = 'ACTIVE')::float / 
				NULLIF((SELECT COUNT(*) FROM units WHERE community_id = $1), 0) * 100`
		db.QueryRow(r.Context(), occupancyQuery, communityID).Scan(&data.Metrics.OccupancyRate)

		// 2. Payment Collection Rate (Current Month)
		paymentQuery := `
			SELECT 
				COALESCE(SUM(paid_amount), 0) / NULLIF(SUM(amount), 0) * 100
			FROM invoices 
			WHERE community_id = $1 AND billing_period = $2`
		db.QueryRow(r.Context(), paymentQuery, communityID, currentMonth).Scan(&data.Metrics.PaymentCollectionRate)

		// 3. Open Maintenance Requests
		maintenanceQuery := `SELECT COUNT(*) FROM maintenance_requests WHERE community_id = $1 AND status IN ('OPEN', 'IN_PROGRESS')`
		db.QueryRow(r.Context(), maintenanceQuery, communityID).Scan(&data.Metrics.OpenMaintenanceCount)

		// 4. Total Outstanding Balance
		balanceQuery := `SELECT COALESCE(SUM(amount - paid_amount), 0) FROM invoices WHERE community_id = $1 AND status != 'PAID'`
		db.QueryRow(r.Context(), balanceQuery, communityID).Scan(&data.Metrics.OutstandingBalance)

		// 5. Visitor Volume (Last 7 Days)
		visitorQuery := `SELECT COUNT(*) FROM visitors WHERE community_id = $1 AND created_at >= NOW() - INTERVAL '7 days'`
		db.QueryRow(r.Context(), visitorQuery, communityID).Scan(&data.Metrics.VisitorVolume7D)

		// 6. Recent Announcements (Top 3)
		annRows, _ := db.Query(r.Context(),
			"SELECT id, title, created_at, priority FROM announcements WHERE community_id = $1 ORDER BY created_at DESC LIMIT 3", communityID)
		defer annRows.Close()
		for annRows.Next() {
			var a domain.Announcement
			annRows.Scan(&a.ID, &a.Title, &a.CreatedAt, &a.Priority)
			data.RecentActivity.Announcements = append(data.RecentActivity.Announcements, a)
		}

		// 7. Recent Security Logs (Top 5)
		logSql := `
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
			WHERE sl.community_id = $1
			ORDER BY sl.time_in DESC LIMIT 5`

		logRows, _ := db.Query(r.Context(), logSql, communityID)
		defer logRows.Close()
		for logRows.Next() {
			var l SecurityLogRecord
			logRows.Scan(&l.ID, &l.PersonType, &l.TimeIn, &l.TimeOut, &l.Remarks, &l.PersonName, &l.UnitNumber, &l.GuardName)
			data.RecentActivity.SecurityLogs = append(data.RecentActivity.SecurityLogs, l)
		}

		domain.WriteJSON(w, http.StatusOK, domain.Envelope{Success: true, Data: data}, nil)
	}
}
