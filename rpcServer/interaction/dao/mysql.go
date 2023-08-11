package dao

import (
	"fmt"
	commysql "github.com/douyin/common/mysql"
	_ "github.com/douyin/kitex_gen/interaction"
	"github.com/douyin/kitex_gen/model"
	"github.com/douyin/models"
	"gorm.io/gorm"
)

type mysql struct {
	cli *gorm.DB
}

func NewMysql() *mysql {
	conn, _ := commysql.NewMysqlConn()
	return &mysql{cli: conn}
}

func (c *mysql) SearchFavoriteExist(m *models.FavoriteVideoRelation) (exist bool, err error) {
	result := c.cli.Where("user_id = ? and video_id = ?", m.UserID, m.VideoID).First(m)
	if result.Error != nil {
		// TODO : log err
		fmt.Println(result.Error)

	}
	return result.RowsAffected > 0, result.Error
}

func (c *mysql) InsertFavorite(m *models.FavoriteVideoRelation) (rows int64, err error) {

	result := c.cli.Create(m)
	if result.Error != nil {
		// TODO : log err
		fmt.Println(result.Error)
	}
	return result.RowsAffected, result.Error

}

func (c *mysql) CancelFavorite(m *models.FavoriteVideoRelation) (rows int64, err error) {

	result := c.cli.Where("user_id = ? and video_id = ?", m.UserID, m.VideoID).Delete(m)
	return result.RowsAffected, result.Error

}
func (c *mysql) SearchVideoListById(id int64) (videoList []*models.Video, err error) {
	result := c.cli.Raw("SELECT * FROM videos WHERE id in (SELECT video_id from user_favorite_videos WHERE user_id = ?)", id)
	var t []*models.Video
	result.Scan(&t)
	return t, nil
}

func (c *mysql) SearchUserById(id int64) (user *model.User, err error) {
	result := c.cli.Raw("SELECT * FROM users WHERE id = ?", id)
	var t model.User
	result.Scan(&t)
	return &t, err
}

// 评论

func (c *mysql) InsertComment(m *models.Comment) (rows int64, err error) {
	result := c.cli.Create(m)
	if result.Error != nil {
		// TODO : log err
		fmt.Println(result.Error)
	}
	return result.RowsAffected, result.Error
}

func (c *mysql) DeleteComment(m *models.Comment) (rows int64, err error) {
	result := c.cli.Where("id = ? and video_id = ?", m.ID, m.VideoID).Delete(m)
	return result.RowsAffected, result.Error
}

func (c *mysql) SearchCommentListSort(videoId int64) (videoList []*models.Comment, err error) {
	//result := c.mysql.Raw("SELECT * from comments c join users u on c.user_id = u.id  WHERE video_id = 3 ORDER BY c.id DESC")
	result := c.cli.Raw("SELECT * from comments where video_id = ? ORDER BY id DESC", videoId)
	var t []*models.Comment
	result.Scan(&t)
	return t, nil
}
