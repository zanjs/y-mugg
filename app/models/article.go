package models

import (
	"fmt"
	"time"

	"github.com/nestgo/log"
	"github.com/zanjs/y-mugg/db"
)

type (
	ArticlePage struct {
		Data []Article `json:"data"`
		Page PageModel `json:"page"`
	}
	// Article is
	Article struct {
		BaseModel
		User    User   `gorm:"ForeignKey:UserId;AssociationForeignKey:Id" json:"user"`
		UserID  int    `json:"user_id" gorm:"type:integer(11)"`
		Title   string `json:"title" gorm:"type:varchar(100)"`
		Content string `json:"content" gorm:"type:text"`
	}
)

func CreateArticle(m *Article) error {
	var err error

	m.CreatedAt = time.Now()
	tx := gorm.MysqlConn().Begin()
	if err = tx.Create(&m).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return err
}

func (m *Article) UpdateArticle(data *Article) error {
	var err error

	m.UpdatedAt = time.Now()
	m.Title = data.Title
	m.Content = data.Content

	tx := gorm.MysqlConn().Begin()
	if err = tx.Save(&m).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return err
}

func (m Article) DeleteArticle() error {
	var err error
	tx := gorm.MysqlConn().Begin()
	if err = tx.Delete(&m).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return err
}

func GetArticleById(id uint64) (Article, error) {
	var (
		article Article
		err     error
	)

	tx := gorm.MysqlConn().Begin()
	if err = tx.Last(&article, id).Error; err != nil {
		tx.Rollback()
		return article, err
	}
	tx.Commit()

	return article, err
}

func GetArticles(p QueryParams) (ArticlePage, error) {
	var (
		articles []Article
		pageData ArticlePage
		err      error
	)

	pageData.Page.Limit = p.Limit
	pageData.Page.Offset = p.Offset
	// pageData.Page.Limit = 2
	// pageData.Page.Offset = 2

	tx := gorm.MysqlConn().Begin()

	// err = tx.Find(&articles).Count(&pageData.Page.Count).Error

	// if err != nil {
	// 	return pageData, err
	// }
	timeLayout := "2006-01-02 15:04:05"

	if p.EndTime == "" {
		p.EndTime = "2099-01-01 00:00:00"
		fmt.Println("endTime 为空")
	}

	startTime, _ := time.Parse(timeLayout, p.StartTime)
	endTime, _ := time.Parse(timeLayout, p.EndTime)

	if err = tx.Where("created_at BETWEEN ? AND ?", startTime, endTime).Preload("User").Order("id desc").Offset(pageData.Page.Offset).Limit(pageData.Page.Limit).Find(&articles).Count(&pageData.Page.Count).Error; err != nil {
		tx.Rollback()
		return pageData, err
	}

	tx.Commit()

	pageData.Data = articles

	return pageData, err
}

func GetArticlesFor() ([]Article, error) {
	var (
		articles []Article

		err error
	)

	tx := gorm.MysqlConn().Begin()

	if err = tx.Find(&articles).Error; err != nil {
		tx.Rollback()
		return articles, err
	}

	// select articles.title,users.username from articles inner join users on articles.user_id = users.id

	for key, article := range articles {
		if err := tx.Model(&article).Related(&article.User).Error; err != nil {
			log.Debugf("articles user related error: %v", err)
		}
		articles[key] = article
	}

	tx.Commit()

	return articles, err
}
