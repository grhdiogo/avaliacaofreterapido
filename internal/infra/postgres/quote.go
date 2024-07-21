package postgres

import (
	"avaliacaofreterapido/internal/domain/quote"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
)

type quoteRepImpl struct {
	ctx context.Context
	tx  pgx.Tx
}

// ========================================================================================================
// tb_quote_volume
// ========================================================================================================

func (r *quoteRepImpl) storeVolume(quoteID string, e quote.Volume) error {
	// create sql
	sqlText := `INSERT INTO tb_quote_volume (
		quote_id,
		category,
		amount,
		unitary_weight,
		price,
		sku,
		height,
		width,
		length
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`
	// save
	result, err := r.tx.Exec(r.ctx, sqlText,
		quoteID,
		e.Category,
		e.Amount,
		e.UnitaryWeight,
		e.Price,
		e.Sku,
		e.Height,
		e.Width,
		e.Length,
	)
	if err != nil {
		return err
	}
	//
	if result.RowsAffected() != 1 {
		return errors.New("failed to save quote volume")
	}
	// success
	return nil
}

// ========================================================================================================
// tb_quote
// ========================================================================================================

func (r *quoteRepImpl) scan(rows pgx.Rows) (*quote.Entity, error) {
	id := sql.NullString{}
	cpfCnpj := sql.NullString{}
	addressCep := sql.NullString{}
	rawResponse := sql.NullString{}
	rawRequest := sql.NullString{}
	// scan
	err := rows.Scan(
		&id,
		&cpfCnpj,
		&addressCep,
		&rawResponse,
		&rawRequest,
	)
	if err != nil {
		return nil, err
	}
	// entity
	entity := new(quote.Entity)
	// add to entity
	if id.Valid {
		entity.ID = id.String
	}
	if cpfCnpj.Valid {
		entity.CpfCnpj = cpfCnpj.String
	}
	if addressCep.Valid {
		entity.Address.Cep = addressCep.String
	}
	if rawResponse.Valid {
		entity.RawResponse = []byte(rawResponse.String)
	}
	if rawRequest.Valid {
		entity.RawRequest = []byte(rawRequest.String)
	}
	return entity, nil
}

func (r *quoteRepImpl) list(limit int) ([]quote.Entity, error) {
	// create sql
	sqlText := `
	SELECT
		id,
		cpf_cnpj,
		address_cep,
		raw_response,
		raw_request
		from tb_quote
		order by created_at desc
		`
	// add limit case exists
	vars := []any{}
	if limit >= 0 {
		// add to sql
		sqlText = fmt.Sprintf("%s limit $1", sqlText)
		// add to variables
		vars = append(vars, limit)
	}
	rows, err := r.tx.Query(r.ctx, sqlText, vars...)
	if err != nil {
		return nil, err
	}
	result := make([]quote.Entity, 0)
	for rows.Next() {
		e, err := r.scan(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, *e)
	}
	// scan
	return result, nil
}

func (r *quoteRepImpl) store(e quote.Entity) error {
	// create sql
	sqlText := `INSERT INTO tb_quote (
		id,
		cpf_cnpj,
		address_cep,
		raw_response,
		raw_request
		)
		VALUES ($1, $2, $3, $4, $5)
		`
	// save
	result, err := r.tx.Exec(r.ctx, sqlText,
		e.ID,
		newString(e.CpfCnpj),
		newString(e.Address.Cep),
		e.RawResponse,
		e.RawRequest,
	)
	if err != nil {
		return err
	}
	//
	if result.RowsAffected() != 1 {
		return errors.New("failed to save quote")
	}
	// success
	return nil
}

// ========================================================================================================
// exported
// ========================================================================================================

func (r *quoteRepImpl) Store(e quote.Entity) error {
	// store quote
	err := r.store(e)
	if err != nil {
		return err
	}
	// store volumes
	for _, v := range e.Volumes {
		err = r.storeVolume(e.ID, v)
		if err != nil {
			return err
		}
	}
	// success
	return nil
}

// TODO: Retornar listagem de volumes
func (r *quoteRepImpl) List(limit int) ([]quote.Entity, error) {
	return r.list(limit)
}

func NewQuoteRepository(ctx context.Context, tx pgx.Tx) quote.Repository {
	return &quoteRepImpl{
		ctx: ctx,
		tx:  tx,
	}
}
