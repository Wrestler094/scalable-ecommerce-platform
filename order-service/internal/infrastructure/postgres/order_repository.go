package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Wrestler094/scalable-ecommerce-platform/order-service/internal/domain"
	"github.com/Wrestler094/scalable-ecommerce-platform/order-service/internal/infrastructure/postgres/dao"
)

var _ domain.OrderRepository = (*OrderRepository)(nil)
var _ domain.OrderPaymentRepository = (*OrderRepository)(nil)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, order domain.Order) error {
	const op = "orderRepository.Create"

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer func() { _ = tx.Rollback() }()

	var orderID int64
	err = tx.QueryRowxContext(ctx, `
		INSERT INTO orders (uuid, user_id, status, total_amount, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, order.UUID, order.UserID, order.Status, order.TotalAmount, order.CreatedAt).Scan(&orderID)
	if err != nil {
		return fmt.Errorf("%s: failed to insert order: %w", op, err)
	}

	items := dao.ToDBOrderItems(orderID, order.Items)
	_, err = tx.NamedExecContext(ctx, `
		INSERT INTO order_items (order_id, product_id, quantity, price)
		VALUES (:order_id, :product_id, :quantity, :price)
	`, items)
	if err != nil {
		return fmt.Errorf("%s: failed to insert order items: %w", op, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}

func (r *OrderRepository) FindByUUID(ctx context.Context, uuid string) (domain.Order, error) {
	const op = "orderRepository.FindByUUID"

	var rows []struct {
		dao.DBOrder
		dao.DBOrderItem
		HasItem bool `db:"has_item"`
	}

	err := r.db.SelectContext(ctx, &rows, `
		SELECT 
			o.id, o.uuid, o.user_id, o.status, o.total_amount, o.created_at,
			oi.order_id, oi.product_id, oi.quantity, oi.price,
			oi.order_id IS NOT NULL as has_item
		FROM orders o
		LEFT JOIN order_items oi ON o.id = oi.order_id
		WHERE o.uuid = $1
	`, uuid)
	if err != nil {
		return domain.Order{}, fmt.Errorf("%s: failed to fetch order with items: %w", op, err)
	}
	if len(rows) == 0 {
		return domain.Order{}, fmt.Errorf("%s: failed to find order: %w", op, domain.ErrOrderNotFound)
	}

	order := domain.Order{
		UUID:        rows[0].UUID,
		UserID:      rows[0].UserID,
		Status:      rows[0].Status,
		TotalAmount: rows[0].TotalAmount,
		CreatedAt:   rows[0].CreatedAt,
		Items:       make([]domain.OrderItem, 0, len(rows)),
	}

	for _, row := range rows {
		if row.HasItem {
			order.Items = append(order.Items, dao.ToDomainOrderItem(row.DBOrderItem))
		}
	}

	return order, nil
}

func (r *OrderRepository) FindByUserID(ctx context.Context, userID int64) ([]domain.Order, error) {
	const op = "orderRepository.FindByUserID"

	var rows []struct {
		dao.DBOrder
		dao.DBOrderItem
		HasItem bool `db:"has_item"`
	}

	err := r.db.SelectContext(ctx, &rows, `
		SELECT 
			o.id, o.uuid, o.user_id, o.status, o.total_amount, o.created_at,
			oi.order_id, oi.product_id, oi.quantity, oi.price,
			oi.order_id IS NOT NULL as has_item
		FROM orders o
		LEFT JOIN order_items oi ON o.id = oi.order_id
		WHERE o.user_id = $1
		ORDER BY o.created_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to fetch user orders: %w", op, err)
	}

	orderMap := make(map[string]*domain.Order)
	for _, row := range rows {
		order, exists := orderMap[row.DBOrder.UUID]
		if !exists {
			order = &domain.Order{
				UUID:        row.UUID,
				UserID:      row.UserID,
				Status:      row.Status,
				TotalAmount: row.TotalAmount,
				CreatedAt:   row.CreatedAt,
				Items:       []domain.OrderItem{},
			}
			orderMap[row.UUID] = order
		}

		if row.HasItem {
			order.Items = append(order.Items, dao.ToDomainOrderItem(row.DBOrderItem))
		}
	}

	orders := make([]domain.Order, 0, len(orderMap))
	for _, order := range orderMap {
		orders = append(orders, *order)
	}

	return orders, nil
}

func (r *OrderRepository) MarkAsPaid(ctx context.Context, uuid string) error {
	const op = "orderRepository.MarkAsPaid"

	res, err := r.db.ExecContext(ctx, `
		UPDATE orders SET status = 'paid'
		WHERE uuid = $1
	`, uuid)
	if err != nil {
		return fmt.Errorf("%s: failed to mark order as paid: %w", op, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get affected rows: %w", op, err)
	}

	if count == 0 {
		return fmt.Errorf("%s: failed to mark order as paid: %w", op, domain.ErrOrderNotFound)
	}

	return nil
}
