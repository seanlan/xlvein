//generated by lazy
//author: seanlan

package dao

import (
	"context"
	"github.com/seanlan/xlvein/app/dao/sqlmodel"
	"gorm.io/gorm/clause"
)

func CountConversation(ctx context.Context, expr clause.Expression) (totalRows int64, err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.Conversation{})
	if expr != nil {
		db = db.Where(expr)
	}
	db.Count(&totalRows)
	return totalRows, nil
}

func SumConversation(ctx context.Context, sumField sqlmodel.FieldBase, expr clause.Expression) (sum int64, err error) {
	var sumValue = struct {
		N int64 `json:"n"`
	}{}
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.Conversation{})
	if expr != nil {
		db = db.Where(expr)
	}
	err = db.Select("sum(" + sumField.FieldName + ") as n").Scan(&sumValue).Error
	return sumValue.N, err
}

func FetchAllConversation(ctx context.Context, records interface{}, expr clause.Expression, page, pagesize int, orders ...string) (err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.Conversation{})
	if expr != nil {
		db = db.Where(expr)
	}
	if page > 0 {
		offset := (page - 1) * pagesize
		db = db.Offset(offset).Limit(pagesize)
	} else {
		db = db.Limit(pagesize)
	}
	for _, order := range orders {
		db = db.Order(order)
	}
	if err = db.Find(records).Error; err != nil {
		err = ErrNotFound
		return err
	}
	return nil
}

func FetchConversation(ctx context.Context, record interface{}, expr clause.Expression, orders ...string) (err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.Conversation{})
	if expr != nil {
		db = db.Where(expr)
	}
	for _, order := range orders {
		db = db.Order(order)
	}
	if err = db.First(record).Error; err != nil {
		err = ErrNotFound
		return err
	}
	return nil
}

func AddConversation(ctx context.Context, d *sqlmodel.Conversation) (RowsAffected int64, err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.Conversation{}).Create(d)
	if err = db.Error; err != nil {
		return -1, ErrInsertFailed
	}
	return db.RowsAffected, nil
}

func AddsConversation(ctx context.Context, d *[]sqlmodel.Conversation) (RowsAffected int64, err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.Conversation{}).Create(d)
	if err = db.Error; err != nil {
		return -1, ErrInsertFailed
	}
	return db.RowsAffected, nil
}

func UpdateConversation(ctx context.Context, updated *sqlmodel.Conversation) (RowsAffected int64, err error) {
	db := GetDB(ctx).WithContext(ctx).Save(updated)
	if err = db.Error; err != nil {
		return -1, ErrUpdateFailed
	}
	return db.RowsAffected, nil
}

func UpdatesConversation(ctx context.Context, expr clause.Expression, updated map[string]interface{}) (RowsAffected int64, err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.Conversation{})
	if expr != nil {
		db = db.Where(expr)
	}
	db = db.Updates(updated)
	if err = db.Error; err != nil {
		return -1, err
	}
	return db.RowsAffected, nil
}

func UpsertConversation(ctx context.Context, d *sqlmodel.Conversation, upsert map[string]interface{}) (RowsAffected int64, err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.Conversation{}).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(upsert),
	}).Create(d)
	if err = db.Error; err != nil {
		return -1, ErrInsertFailed
	}
	return db.RowsAffected, nil
}

func DeleteConversation(ctx context.Context, expr clause.Expression) (rowsAffected int64, err error) {
	db := GetDB(ctx).WithContext(ctx).Model(&sqlmodel.Conversation{})
	if expr != nil {
		db = db.Where(expr)
	}
	db = db.Delete(sqlmodel.Conversation{})
	if err = db.Error; err != nil {
		return -1, err
	}
	return db.RowsAffected, nil
}
