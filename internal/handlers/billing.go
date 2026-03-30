package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gated-community-api/internal/domain"
	"gated-community-api/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Payment struct {
	ID            string    `json:"id"`
	InvoiceID     string    `json:"invoice_id"`
	CommunityID   string    `json:"community_id"`
	Amount        float64   `json:"amount"`
	PaymentMethod string    `json:"payment_method"` // 'cash', 'transfer', 'card'
	TransactionID string    `json:"transaction_id"`
	PaymentDate   time.Time `json:"payment_date"`
}

type Invoice struct {
	ID            string    `json:"id"`
	UnitID        string    `json:"unit_id"`
	UnitNumber    string    `json:"unit_number,omitempty"`
	Amount        float64   `json:"amount"`
	DueDate       time.Time `json:"due_date"`
	Status        string    `json:"status"`
	BillingPeriod string    `json:"billing_period"`
	PaidAmount    float64   `json:"paid_amount"`
	Description   string    `json:"description"`
}

func HandleCreateInvoice(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		communityID, ok := r.Context().Value(middleware.CommunityContextKey).(string)
		if !ok {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Missing community context"}, nil)
			return
		}
		var inv Invoice
		if err := json.NewDecoder(r.Body).Decode(&inv); err != nil {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Invalid input"}, nil)
			return
		}

		query := `
			INSERT INTO invoices (community_id, unit_id, amount, due_date, status, billing_period, description)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id`

		err := db.QueryRow(r.Context(), query,
			communityID, inv.UnitID, inv.Amount, inv.DueDate, "UNPAID", inv.BillingPeriod, inv.Description,
		).Scan(&inv.ID)

		if err != nil {
			domain.WriteJSON(w, http.StatusInternalServerError, domain.Envelope{Success: false, Message: "Failed to create invoice"}, nil)
			return
		}

		domain.WriteJSON(w, http.StatusCreated, domain.Envelope{Success: true, Data: inv}, nil)
	}
}

func HandleCreatePayment(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		communityID, ok := r.Context().Value(middleware.CommunityContextKey).(string)
		if !ok {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Missing community context"}, nil)
			return
		}
		var p Payment
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Invalid input"}, nil)
			return
		}

		tx, err := db.Begin(r.Context())
		if err != nil {
			domain.WriteJSON(w, http.StatusInternalServerError, domain.Envelope{Success: false, Message: "Failed to start transaction"}, nil)
			return
		}
		defer tx.Rollback(r.Context())

		// 1. Record the payment
		queryPayment := `
			INSERT INTO payments (community_id, invoice_id, amount, payment_method, transaction_id)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, payment_date`

		err = tx.QueryRow(r.Context(), queryPayment,
			communityID, p.InvoiceID, p.Amount, p.PaymentMethod, p.TransactionID,
		).Scan(&p.ID, &p.PaymentDate)
		if err != nil {
			domain.WriteJSON(w, http.StatusInternalServerError, domain.Envelope{Success: false, Message: "Failed to record payment"}, nil)
			return
		}

		// 2. Update invoice balance and status
		queryInvoice := `
			UPDATE invoices 
			SET paid_amount = paid_amount + $1,
			    status = CASE 
					WHEN (paid_amount + $1) >= amount THEN 'PAID'
					WHEN (paid_amount + $1) > 0 THEN 'PARTIAL'
					ELSE 'UNPAID'
				END
			WHERE id = $2 AND community_id = $3`

		_, err = tx.Exec(r.Context(), queryInvoice, p.Amount, p.InvoiceID, communityID)
		if err != nil {
			domain.WriteJSON(w, http.StatusInternalServerError, domain.Envelope{Success: false, Message: "Failed to update invoice"}, nil)
			return
		}

		if err := tx.Commit(r.Context()); err != nil {
			domain.WriteJSON(w, http.StatusInternalServerError, domain.Envelope{Success: false, Message: "Failed to commit payment"}, nil)
			return
		}

		// Trigger Notification
		var unitUserID string
		err = db.QueryRow(r.Context(),
			`SELECT r.user_id 
			 FROM residents r 
			 JOIN invoices i ON r.unit_id = i.unit_id 
			 WHERE i.id = $1 AND r.is_primary_contact = true AND r.user_id IS NOT NULL`,
			p.InvoiceID).Scan(&unitUserID)

		if err == nil {
			CreateNotification(r.Context(), db, unitUserID, communityID,
				"Payment Received", "A payment of $"+fmt.Sprintf("%.2f", p.Amount)+" has been recorded.", "/billing/invoices")
		}

		domain.WriteJSON(w, http.StatusCreated, domain.Envelope{Success: true, Data: p}, nil)
	}
}

func HandleListPayments(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		communityID, ok := r.Context().Value(middleware.CommunityContextKey).(string)
		if !ok {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Missing community context"}, nil)
			return
		}
		invoiceID := r.URL.Query().Get("invoice_id")

		query := `SELECT id, invoice_id, amount, payment_method, transaction_id, payment_date 
		          FROM payments WHERE community_id = $1`
		args := []any{communityID}

		if invoiceID != "" {
			query += " AND invoice_id = $2"
			args = append(args, invoiceID)
		}
		query += " ORDER BY payment_date DESC"

		rows, err := db.Query(r.Context(), query, args...)
		if err != nil {
			domain.WriteJSON(w, http.StatusInternalServerError, domain.Envelope{Success: false, Message: "Failed to fetch payments"}, nil)
			return
		}
		defer rows.Close()

		payments := []Payment{}
		for rows.Next() {
			var p Payment
			if err := rows.Scan(&p.ID, &p.InvoiceID, &p.Amount, &p.PaymentMethod, &p.TransactionID, &p.PaymentDate); err != nil {
				continue
			}
			payments = append(payments, p)
		}
		domain.WriteJSON(w, http.StatusOK, domain.Envelope{Success: true, Data: payments}, nil)
	}
}

func HandleListInvoices(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		communityID, ok := r.Context().Value(middleware.CommunityContextKey).(string)
		if !ok {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Missing community context"}, nil)
			return
		}
		status := r.URL.Query().Get("status")

		query := `
			SELECT i.id, i.unit_id, u.unit_number, i.amount, i.due_date, i.status, i.billing_period, i.description
			FROM invoices i
			JOIN units u ON i.unit_id = u.id
			WHERE i.community_id = $1`

		args := []any{communityID}
		if status != "" {
			query += " AND i.status = $2"
			args = append(args, status)
		}
		query += " ORDER BY i.due_date DESC"

		rows, err := db.Query(r.Context(), query, args...)
		if err != nil {
			domain.WriteJSON(w, http.StatusInternalServerError, domain.Envelope{Success: false, Message: "Failed to fetch invoices"}, nil)
			return
		}
		defer rows.Close()

		invoices := []Invoice{}
		for rows.Next() {
			var i Invoice
			if err := rows.Scan(&i.ID, &i.UnitID, &i.UnitNumber, &i.Amount, &i.DueDate, &i.Status, &i.BillingPeriod, &i.Description); err != nil {
				continue
			}
			invoices = append(invoices, i)
		}
		domain.WriteJSON(w, http.StatusOK, domain.Envelope{Success: true, Data: invoices}, nil)
	}
}

func HandleGetOutstandingBalances(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		communityID, ok := r.Context().Value(middleware.CommunityContextKey).(string)
		if !ok {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "Missing community context"}, nil)
			return
		}

		query := `
			SELECT u.unit_number, SUM(i.amount) as total_due
			FROM invoices i
			JOIN units u ON i.unit_id = u.id
			WHERE i.community_id = $1 AND i.status IN ('UNPAID', 'PARTIAL', 'OVERDUE')
			GROUP BY u.unit_number
			ORDER BY total_due DESC`

		rows, err := db.Query(r.Context(), query, communityID)
		if err != nil {
			domain.WriteJSON(w, http.StatusInternalServerError, domain.Envelope{Success: false, Message: "Failed to calculate balances"}, nil)
			return
		}
		defer rows.Close()

		type Balance struct {
			UnitNumber string  `json:"unit_number"`
			TotalDue   float64 `json:"total_due"`
		}
		balances := []Balance{}
		for rows.Next() {
			var b Balance
			if err := rows.Scan(&b.UnitNumber, &b.TotalDue); err != nil {
				continue
			}
			balances = append(balances, b)
		}
		domain.WriteJSON(w, http.StatusOK, domain.Envelope{Success: true, Data: balances}, nil)
	}
}

func HandleGetUnitLedger(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		unitID := chi.URLParam(r, "unitID")
		query := `SELECT id, amount, due_date, status, billing_period, description FROM invoices WHERE unit_id = $1 ORDER BY due_date DESC`

		rows, err := db.Query(r.Context(), query, unitID)
		if err != nil {
			domain.WriteJSON(w, http.StatusInternalServerError, domain.Envelope{Success: false, Message: "Failed to fetch ledger"}, nil)
			return
		}
		defer rows.Close()

		invoices := []Invoice{}
		for rows.Next() {
			var i Invoice
			rows.Scan(&i.ID, &i.Amount, &i.DueDate, &i.Status, &i.BillingPeriod, &i.Description)
			invoices = append(invoices, i)
		}
		domain.WriteJSON(w, http.StatusOK, domain.Envelope{Success: true, Data: invoices}, nil)
	}
}
