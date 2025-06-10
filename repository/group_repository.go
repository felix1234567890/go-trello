package repository

import (
	"gorm.io/gorm"
	"felix1234567890/go-trello/models"
)

type GroupRepository struct {
	DB *gorm.DB
}

func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{DB: db}
}

func (r *GroupRepository) CreateGroup(group *models.Group) error {
	return r.DB.Create(group).Error
}

func (r *GroupRepository) GetGroupByID(id uint) (*models.Group, error) {
	var group models.Group
	err := r.DB.Preload("Users").First(&group, id).Error
	return &group, err
}

func (r *GroupRepository) GetAllGroups() ([]models.Group, error) {
	var groups []models.Group
	err := r.DB.Preload("Users").Find(&groups).Error
	return groups, err
}

func (r *GroupRepository) UpdateGroup(group *models.Group) error {
	return r.DB.Save(group).Error
}

func (r *GroupRepository) DeleteGroup(id uint) error {
	return r.DB.Delete(&models.Group{}, id).Error
}

func (r *GroupRepository) AddUserToGroup(groupID, userID uint) error {
	var group models.Group
	if err := r.DB.First(&group, groupID).Error; err != nil {
		return err
	}
	var user models.User
	if err := r.DB.First(&user, userID).Error; err != nil {
		return err
	}
	return r.DB.Model(&group).Association("Users").Append(&user)
}

func (r *GroupRepository) RemoveUserFromGroup(groupID, userID uint) error {
	var group models.Group
	if err := r.DB.First(&group, groupID).Error; err != nil {
		return err
	}
	var user models.User
	if err := r.DB.First(&user, userID).Error; err != nil {
		return err
	}
	return r.DB.Model(&group).Association("Users").Delete(&user)
}
