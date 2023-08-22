package dao

import (
	"fmt"
	commysql "github.com/douyin/common/mysql"
	_ "github.com/douyin/kitex_gen/interaction"
	"github.com/douyin/kitex_gen/model"
	"github.com/douyin/models"
	"gorm.io/gorm"
	"log"
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
		log.Println(err)
	}
	return result.RowsAffected > 0, result.Error
}

func (c *mysql) InsertFavorite(m *models.FavoriteVideoRelation) (rows int64, err error) {
	result := c.cli.Create(m)
	if result.Error != nil {
		log.Println(err)
	}
	return result.RowsAffected, result.Error

}

func (c *mysql) CancelFavorite(m *models.FavoriteVideoRelation) (rows int64, err error) {
	result := c.cli.Where("user_id = ? and video_id = ?", m.UserID, m.VideoID).Delete(m)
	if result.Error != nil {
		log.Println(err)
	}
	return result.RowsAffected, result.Error
}

func (c *mysql) VideoFavoriteCountIncr(videoId int64, num int64) (rows int64, err error) {
	result := c.cli.Exec("UPDATE videos SET favorite_count = favorite_count + ? WHERE id = ?;", num, videoId)
	if result.Error != nil {
		log.Println(err)
	}
	return result.RowsAffected, result.Error
}

func (c *mysql) VideoCommentCountIncr(videoId int64, num int64) (rows int64, err error) {
	result := c.cli.Exec("UPDATE videos SET comment_count = comment_count + ? WHERE id = ?;", num, videoId)
	if result.Error != nil {
		log.Println(err)
	}
	return result.RowsAffected, result.Error
}

func (c *mysql) UserFavoriteCountIncr(userId int64, num int64) (rows int64, err error) {
	result := c.cli.Exec("UPDATE users SET favorite_count = favorite_count + ? WHERE id = ?;", num, userId)
	if result.Error != nil {
		log.Println(err)
	}
	return result.RowsAffected, result.Error
}

func (c *mysql) UserTotalFavoritedCountIncr(userId int64, num int64) (rows int64, err error) {
	result := c.cli.Exec("UPDATE users SET total_favorited = total_favorited + ? WHERE id = ?;", num, userId)
	if result.Error != nil {
		log.Println(err)
	}
	return result.RowsAffected, result.Error
}

// SearchFavoriteVideoIds 根据userId查询喜欢视频ids列表
func (c *mysql) SearchFavoriteVideoIds(userId int64) (favoriteVideoIds []int64, err error) {
	result := c.cli.Raw("SELECT video_id from user_favorite_videos WHERE user_id = ?", userId)
	var t []int64
	result.Scan(&t)
	return t, result.Error
}

// SearchAuthorIdsByVideoId 根据视频ids查询author列表
func (c *mysql) SearchAuthorIdsByVideoId(id int64) (authorId int64, err error) {
	result := c.cli.Raw("SELECT author_id from videos WHERE id = ?", id)
	var t int64
	result.Scan(&t)
	return t, result.Error
}

// SearchAuthorIdsByVideoIds 根据视频ids列表查询author列表
func (c *mysql) SearchAuthorIdsByVideoIds(ids int64) (authorIds []int64, err error) {
	result := c.cli.Raw("SELECT author_id from videos WHERE id in ?", ids)
	var t []int64
	result.Scan(&t)
	return t, result.Error
}

//func (c *mysql) SearchFavoriteVideoList(ids []int64) (favoriteVideoIds []int64, err error) {
//	result := c.cli.Raw("SELECT video_id from videos WHERE id in ?", ids)
//	var t []int64
//	result.Scan(&t)
//	return t, result.Error
//}

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

func (c *mysql) SearchUserByIds(authorIds []int64, userId int64) (userList []*model.User, err error) {
	result := c.cli.Raw("SELECT u.*, if(r.to_user_id is NULL,0, 1) as is_follow "+
		"FROM users u left join relations r on u.id = r.user_id and r.to_user_id = ? "+
		"WHERE u.id in ?", userId, authorIds)
	var t []*model.User
	result.Scan(&t)
	return t, result.Error
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
