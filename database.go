package photoshare

import (
	"database/sql"
	"fmt"
	"github.com/coopernurse/gorp"
	"github.com/juju/errgo"
	_ "github.com/lib/pq" // PostgreSQL library
	"log"
	"os"
	"strings"
)

func dbConnect(user, pwd, name, host string) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s host=%s",
		user,
		name,
		pwd,
		host,
	))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func initDB(db *sql.DB, logSql bool) (*gorp.DbMap, error) {
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	if logSql {
		dbMap.TraceOn("[sql]", log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds))
	}

	dbMap.AddTableWithName(user{}, "users").SetKeys(true, "ID")
	dbMap.AddTableWithName(photo{}, "photos").SetKeys(true, "ID")
	dbMap.AddTableWithName(tag{}, "tags").SetKeys(true, "ID")

	return dbMap, nil
}

type dataMapper interface {
	createPhoto(*photo) error
	removePhoto(*photo) error
	updatePhoto(*photo) error
	updateTags(*photo) error

	createUser(*user) error
	updateUser(*user) error

	updateMany(...interface{}) error

	getPhoto(int64) (*photo, error)
	getPhotoDetail(int64, *user) (*photoDetail, error)
	getTagCounts() ([]tagCount, error)
	getPhotos(*page, string) (*photoList, error)
	getPhotosByOwnerID(*page, int64) (*photoList, error)
	searchPhotos(*page, string) (*photoList, error)

	isUserNameAvailable(user *user) (bool, error)
	isUserEmailAvailable(user *user) (bool, error)
	getActiveUser(userID int64) (*user, error)
	getUserByRecoveryCode(string) (*user, error)
	getUserByEmail(string) (*user, error)
	getUserByNameOrEmail(identifier string) (*user, error)
}

type defaultDataMapper struct {
	*gorp.DbMap
}

type transaction struct {
	*gorp.Transaction
}

func (t *transaction) updateTags(photo *photo) error {

	var (
		args    = []string{"$1"}
		params  = []interface{}{interface{}(photo.ID)}
		isEmpty = true
		counter = 1
	)
	for _, name := range photo.Tags {
		name = strings.TrimSpace(name)
		if name != "" {
			counter++
			args = append(args, fmt.Sprintf("$%d", counter))
			params = append(params, interface{}(strings.ToLower(name)))
			isEmpty = false
		}
	}

	if isEmpty && photo.ID != 0 {
		_, err := t.Exec("DELETE FROM photo_tags WHERE photo_id=$1", photo.ID)
		return errgo.Mask(err)
	}
	if _, err := t.Exec(fmt.Sprintf("SELECT add_tags(%s)", strings.Join(args, ",")), params...); err != nil {
		return errgo.Mask(err)
	}
	return nil

}

func newDataMapper(db *sql.DB, logSql bool) (dataMapper, error) {
	dbMap, err := initDB(db, logSql)
	if err != nil {
		return nil, err
	}
	return &defaultDataMapper{dbMap}, nil
}

func (d *defaultDataMapper) begin() (*transaction, error) {
	tx, err := d.Begin()
	if err != nil {
		return nil, err
	}
	return &transaction{tx}, nil
}

func (d *defaultDataMapper) createPhoto(photo *photo) error {
	t, err := d.begin()
	if err != nil {
		return errgo.Mask(err)
	}
	if err := t.Insert(photo); err != nil {
		return errgo.Mask(err)
	}
	if err := t.updateTags(photo); err != nil {
		t.Rollback()
		return errgo.Mask(err)
	}
	return errgo.Mask(t.Commit())
}

func (d *defaultDataMapper) createUser(user *user) error {
	return errgo.Mask(d.Insert(user))
}

func (d *defaultDataMapper) updatePhoto(photo *photo) error {
	if _, err := d.Update(photo); err != nil {
		return errgo.Mask(err)
	}
	return nil
}

func (d *defaultDataMapper) updateUser(user *user) error {
	if _, err := d.Update(user); err != nil {
		return errgo.Mask(err)
	}
	return nil
}

func (d *defaultDataMapper) removePhoto(photo *photo) error {
	if _, err := d.Delete(photo); err != nil {
		return errgo.Mask(err)
	}
	return nil
}

func (d *defaultDataMapper) updateTags(photo *photo) error {
	tx, err := d.begin()
	if err != nil {
		return errgo.Mask(err)
	}
	if err := tx.updateTags(photo); err != nil {
		tx.Rollback()
		return err
	}
	return errgo.Mask(tx.Commit())
}

func (d *defaultDataMapper) updateMany(items ...interface{}) error {
	tx, err := d.begin()
	if err != nil {
		return errgo.Mask(err)
	}
	for _, item := range items {
		if _, err := tx.Update(item); err != nil {
			tx.Rollback()
			return errgo.Mask(err)
		}
	}
	return errgo.Mask(tx.Commit())
}

func (d *defaultDataMapper) getPhoto(photoID int64) (*photo, error) {

	p := &photo{}

	if photoID == 0 {
		return p, sql.ErrNoRows
	}

	obj, err := d.Get(p, photoID)
	if err != nil {
		return p, errgo.Mask(err)
	}
	if obj == nil {
		return p, sql.ErrNoRows
	}
	return obj.(*photo), nil
}

func (d *defaultDataMapper) getPhotoDetail(photoID int64, user *user) (*photoDetail, error) {

	photo := &photoDetail{}

	if photoID == 0 {
		return photo, sql.ErrNoRows
	}

	q := "SELECT p.*, u.name AS owner_name " +
		"FROM photos p JOIN users u ON u.id = p.owner_id " +
		"WHERE p.id=$1"

	if err := d.SelectOne(photo, q, photoID); err != nil {
		return photo, errgo.Mask(err)
	}

	var tags []tag

	if _, err := d.Select(&tags,
		"SELECT t.* FROM tags t JOIN photo_tags pt ON pt.tag_id=t.id "+
			"WHERE pt.photo_id=$1", photo.ID); err != nil {
		return photo, errgo.Mask(err)
	}
	for _, tag := range tags {
		photo.Tags = append(photo.Tags, tag.Name)
	}

	photo.Permissions = &permissions{
		photo.canEdit(user),
		photo.canDelete(user),
		photo.canVote(user),
	}
	return photo, nil

}

func (d *defaultDataMapper) getPhotosByOwnerID(page *page, ownerID int64) (*photoList, error) {
	var (
		photos []photo
		err    error
		total  int64
	)

	if ownerID == 0 {
		return nil, sql.ErrNoRows
	}
	if total, err = d.SelectInt("SELECT COUNT(id) FROM photos WHERE owner_id=$1", ownerID); err != nil {
		return nil, errgo.Mask(err)
	}

	if _, err = d.Select(&photos,
		"SELECT * FROM photos WHERE owner_id = $1"+
			"ORDER BY (up_votes - down_votes) DESC, created_at DESC LIMIT $2 OFFSET $3",
		ownerID, page.size, page.offset); err != nil {
		return nil, errgo.Mask(err)
	}
	return newPhotoList(photos, total, page.index), nil

}

func (d *defaultDataMapper) searchPhotos(page *page, q string) (*photoList, error) {

	var (
		clauses []string
		params  []interface{}
		err     error
		photos  []photo
		total   int64
	)

	if q == "" {
		return nil, nil
	}

	for num, word := range strings.Split(q, " ") {
		word = strings.TrimSpace(word)
		if word == "" || num > 6 {
			break
		}

		num++

		if strings.HasPrefix(word, "@") {
			word = word[1:]
			clauses = append(clauses, fmt.Sprintf(
				"SELECT p.* FROM photos p "+
					"INNER JOIN users u ON u.id = p.owner_id  "+
					"WHERE UPPER(u.name::text) = UPPER($%d)", num))
		} else if strings.HasPrefix(word, "#") {
			word = word[1:]
			clauses = append(clauses, fmt.Sprintf(
				"SELECT p.* FROM photos p "+
					"INNER JOIN photo_tags pt ON pt.photo_id = p.id "+
					"INNER JOIN tags t ON pt.tag_id=t.id "+
					"WHERE UPPER(t.name::text) = UPPER($%d)", num))
		} else {
			word = "%" + word + "%"
			clauses = append(clauses, fmt.Sprintf(
				"SELECT DISTINCT p.* FROM photos p "+
					"INNER JOIN users u ON u.id = p.owner_id  "+
					"LEFT JOIN photo_tags pt ON pt.photo_id = p.id "+
					"LEFT JOIN tags t ON pt.tag_id=t.id "+
					"WHERE UPPER(p.title::text) LIKE UPPER($%d) OR "+
					"UPPER(u.name::text) LIKE UPPER($%d) OR t.name LIKE $%d",
				num, num, num))
		}

		params = append(params, interface{}(word))
	}

	clausesSql := strings.Join(clauses, " INTERSECT ")

	countSql := fmt.Sprintf("SELECT COUNT(id) FROM (%s) q", clausesSql)

	if total, err = d.SelectInt(countSql, params...); err != nil {
		return nil, errgo.Mask(err)
	}

	numParams := len(params)

	sql := fmt.Sprintf("SELECT * FROM (%s) q ORDER BY (up_votes - down_votes) DESC, created_at DESC LIMIT $%d OFFSET $%d",
		clausesSql, numParams+1, numParams+2)

	params = append(params, interface{}(page.size))
	params = append(params, interface{}(page.offset))

	if _, err = d.Select(&photos, sql, params...); err != nil {
		return nil, errgo.Mask(err)
	}
	return newPhotoList(photos, total, page.index), nil
}

func (d *defaultDataMapper) getPhotos(page *page, orderBy string) (*photoList, error) {

	var (
		total  int64
		photos []photo
		err    error
	)
	if orderBy == "votes" {
		orderBy = "(up_votes - down_votes)"
	} else {
		orderBy = "created_at"
	}

	if total, err = d.SelectInt("SELECT COUNT(id) FROM photos"); err != nil {
		return nil, errgo.Mask(err)
	}

	if _, err = d.Select(&photos,
		"SELECT * FROM photos "+
			"ORDER BY "+orderBy+" DESC LIMIT $1 OFFSET $2", page.size, page.offset); err != nil {
		return nil, errgo.Mask(err)
	}
	return newPhotoList(photos, total, page.index), nil
}

func (d *defaultDataMapper) getTagCounts() ([]tagCount, error) {
	var tags []tagCount
	if _, err := d.Select(&tags, "SELECT name, photo, num_photos FROM tag_counts"); err != nil {
		return tags, errgo.Mask(err)
	}
	return tags, nil
}

func (d *defaultDataMapper) isUserNameAvailable(user *user) (bool, error) {
	var (
		num int64
		err error
	)
	q := "SELECT COUNT(id) FROM users WHERE name=$1"
	if user.ID == 0 {
		num, err = d.SelectInt(q, user.Name)
	} else {
		q += " AND id != $2"
		num, err = d.SelectInt(q, user.Name, user.ID)
	}
	if err != nil {
		return false, errgo.Mask(err)
	}
	return num == 0, nil
}

func (d *defaultDataMapper) isUserEmailAvailable(user *user) (bool, error) {
	var (
		num int64
		err error
	)
	q := "SELECT COUNT(id) FROM users WHERE email=$1"
	if user.ID == 0 {
		num, err = d.SelectInt(q, user.Email)
	} else {
		q += " AND id != $2"
		num, err = d.SelectInt(q, user.Email, user.ID)
	}
	if err != nil {
		return false, errgo.Mask(err)
	}
	return num == 0, nil
}

func (d *defaultDataMapper) getActiveUser(userID int64) (*user, error) {

	user := &user{}
	if err := d.SelectOne(user, "SELECT * FROM users WHERE active=$1 AND id=$2", true, userID); err != nil {
		return user, errgo.Mask(err)
	}
	return user, nil

}

func (d *defaultDataMapper) getUserByRecoveryCode(code string) (*user, error) {

	user := &user{}
	if code == "" {
		return user, sql.ErrNoRows
	}
	if err := d.SelectOne(user, "SELECT * FROM users WHERE active=$1 AND recovery_code=$2", true, code); err != nil {
		return user, errgo.Mask(err)
	}
	return user, nil

}
func (d *defaultDataMapper) getUserByEmail(email string) (*user, error) {
	user := &user{}
	if err := d.SelectOne(user, "SELECT * FROM users WHERE active=$1 AND email=$2", true, email); err != nil {
		return user, errgo.Mask(err)
	}
	return user, nil
}

func (d *defaultDataMapper) getUserByNameOrEmail(identifier string) (*user, error) {
	user := &user{}

	if err := d.SelectOne(user, "SELECT * FROM users WHERE active=$1 AND (email=$2 OR name=$2)", true, identifier); err != nil {
		return user, errgo.Mask(err)
	}

	return user, nil
}
