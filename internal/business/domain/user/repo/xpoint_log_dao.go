package repo

import (
	"errors"
	"gim/internal/business/domain/user/model"
	"gim/pkg/db"
	"gim/pkg/gerrors"
	"github.com/jinzhu/gorm"
)

type xpointLogDao struct {
}

var XpointLogDao = new(xpointLogDao)

// GetLogsByUserID 获取用户积分变动日志
func (d *xpointLogDao) GetLogsByUserID(userId int64, limit int) ([]model.XPointLog, error) {
	var logs []model.XPointLog
	err := db.DB.Where("user_id = ?", userId).Order("create_time desc").Limit(limit).Find(&logs).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gerrors.WrapError(err)
	}
	return logs, nil
}

// AddLog 插入积分变动日志
func (d *xpointLogDao) AddLog(log model.XPointLog) (int64, error) {
	err := db.DB.Create(&log).Error
	return int64(log.ID), gerrors.WrapError(err)
}
