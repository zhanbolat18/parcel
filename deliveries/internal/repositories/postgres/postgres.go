package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/zhanbolat18/parcel/deliveries/internal/entities"
	"github.com/zhanbolat18/parcel/deliveries/internal/repositories"
	"github.com/zhanbolat18/parcel/deliveries/internal/valueobjects"
	"time"
)

const DateFormat = time.RFC3339

type delivery struct {
	db *sqlx.DB
}

func NewDeliveryRepository(db *sqlx.DB) repositories.DeliveriesRepository {
	return &delivery{db: db}
}

type deliveryModel struct {
	Id          uint   `db:"id" json:"id"`
	Status      string `db:"status" json:"status"`
	Destination string `db:"destination" json:"destination"`
	RecipientId uint   `db:"recipient_id" json:"recipientId"`
	CourierId   *uint  `db:"courier_id" json:"courierId,omitempty"`
	CreatedAt   string `db:"created_at" json:"createdAt"`
	UpdatedAt   string `db:"updated_at" json:"updatedAt"`
}

func (d *delivery) GetAll(ctx context.Context) ([]*entities.Delivery, error) {
	dm := make([]deliveryModel, 0)
	q := "SELECT * FROM deliveries"
	err := d.db.SelectContext(ctx, &dm, q)
	if err != nil {
		return nil, err
	}
	dls := make([]*entities.Delivery, 0, len(dm))
	for _, model := range dm {
		model := model
		dls = append(dls, d.hydrateToEntity(&model))
	}
	return dls, nil
}

func (d *delivery) GetAllByCourier(ctx context.Context, courierId uint) ([]*entities.Delivery, error) {
	dm := make([]deliveryModel, 0)
	q := "SELECT * FROM deliveries WHERE courier_id=$1"
	err := d.db.SelectContext(ctx, &dm, q, courierId)
	if err != nil {
		return nil, err
	}
	dls := make([]*entities.Delivery, 0, len(dm))
	for _, model := range dm {
		model := model
		dls = append(dls, d.hydrateToEntity(&model))
	}
	return dls, nil
}

func (d *delivery) GetById(ctx context.Context, id uint) (*entities.Delivery, error) {
	dm := &deliveryModel{}
	q := "SELECT * FROM deliveries WHERE id=$1"
	err := d.db.GetContext(ctx, dm, q, id)
	if err != nil {
		return nil, err
	}
	return d.hydrateToEntity(dm), nil
}

func (d *delivery) Store(ctx context.Context, delivery *entities.Delivery) error {
	q := `INSERT INTO deliveries(status, destination, recipient_id, courier_id, created_at, updated_at) 
			VALUES($1, $2, $3, $4, $5, $6)
			RETURNING id;`
	dm := d.hydrateFromEntity(delivery)
	var id int
	stmt, err := d.db.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	err = stmt.QueryRowContext(ctx, dm.Status, dm.Destination, dm.RecipientId, dm.CourierId, dm.CreatedAt, dm.UpdatedAt).
		Scan(&id)
	if err != nil {
		return err
	}
	delivery.Id = uint(id)
	return nil
}

func (d *delivery) Update(ctx context.Context, delivery *entities.Delivery) error {
	q := `UPDATE deliveries SET 
			status=:status, 
			destination=:destination,
			recipient_id=:recipient_id,
			courier_id=:courier_id,
			created_at=:created_at,
			updated_at=:updated_at
		WHERE id=:id`
	_, err := d.db.NamedExecContext(ctx, q, d.hydrateFromEntity(delivery))
	return err
}

func (d *delivery) hydrateFromEntity(delivery *entities.Delivery) *deliveryModel {
	var c, u string

	if !delivery.CreatedAt.IsZero() {
		c = delivery.CreatedAt.Format(DateFormat)
	}
	if !delivery.UpdatedAt.IsZero() {
		u = delivery.UpdatedAt.Format(DateFormat)
	}

	return &deliveryModel{
		Id:          delivery.Id,
		Status:      string(delivery.Status),
		Destination: delivery.Destination,
		RecipientId: delivery.RecipientId,
		CourierId:   delivery.CourierId,
		CreatedAt:   c,
		UpdatedAt:   u,
	}
}

func (d *delivery) hydrateToEntity(model *deliveryModel) *entities.Delivery {
	var c, u time.Time
	if len(model.CreatedAt) > 0 {
		c, _ = time.Parse(DateFormat, model.CreatedAt)
	}
	if len(model.UpdatedAt) > 0 {
		u, _ = time.Parse(DateFormat, model.UpdatedAt)
	}

	return &entities.Delivery{
		Id:          model.Id,
		Status:      valueobjects.Status(model.Status),
		Destination: model.Destination,
		RecipientId: model.RecipientId,
		CourierId:   model.CourierId,
		CreatedAt:   c,
		UpdatedAt:   u,
	}
}
