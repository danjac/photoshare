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

type dataStore struct {
	*gorp.DbMap
	photos photoDataManager
	users  userDataManager
}

func newDataStore(dbMap *gorp.DbMap) *dataStore {
	ds := &dataStore{DbMap: dbMap}
	ds.photos = &defaultPhotoDataManager{dbMap}
	ds.users = &defaultUserDataManager{dbMap}
	return ds
}

type photoDataManager interface {
	create(*photo) error
	update(*photo) error
	remove(*photo) error
	get(int64) (*photo, error)
	getDetail(int64, *user) (*photoDetail, error)
	getTagCounts() ([]tagCount, error)
	all(*page, string) (*photoList, error)
	byOwnerID(*page, int64) (*photoList, error)
	search(*page, string) (*photoList, error)
	updateVotes(*photo, *user) error
	updateTags(*photo) error
}

type defaultPhotoDataManager struct {
	dbMap *gorp.DbMap
}

func (m *defaultPhotoDataManager) remove(photo *photo) error {
	_, err := m.dbMap.Delete(photo)
	return errgo.Mask(err)
}

func (m *defaultPhotoDataManager) update(photo *photo) error {
	_, err := m.dbMap.Update(photo)
	return errgo.Mask(err)
}

func (m *defaultPhotoDataManager) create(photo *photo) error {
	t, err := m.dbMap.Begin()
	if err != nil {
		return errgo.Mask(err)
	}
	if err := t.Insert(photo); err != nil {
		t.Rollback()
		return errgo.Mask(err)
	}
	if err := m.txUpdateTags(t, photo); err != nil {
		return err
	}
	return errgo.Mask(t.Commit())
}

func (m *defaultPhotoDataManager) updateVotes(photo *photo, user *user) error {

	t, err := m.dbMap.Begin()
	if err != nil {
		return err
	}

	if _, err = t.Update(photo); err != nil {
		t.Rollback()
		return errgo.Mask(err)
	}

	user.registerVote(photo.ID)

	if _, err = t.Update(user); err != nil {
		t.Rollback()
		return errgo.Mask(err)
	}

	return errgo.Mask(t.Commit())

}

func (m *defaultPhotoDataManager) updateTags(photo *photo) error {
	t, err := m.dbMap.Begin()
	if err != nil {
		return err
	}

	if err := m.txUpdateTags(t, photo); err != nil {
		return err
	}
	return errgo.Mask(t.Commit())
}

// handles delete/update of tags within a transaction
func (m *defaultPhotoDataManager) txUpdateTags(t *gorp.Transaction, photo *photo) error {

	var (
		args    = []string{"$1"}
		params  = []interface{}{interface{}(photo.ID)}
		isEmpty = true
	)
	for i, name := range photo.Tags {
		name = strings.TrimSpace(name)
		if name != "" {
			args = append(args, fmt.Sprintf("$%d", i+2))
			params = append(params, interface{}(strings.ToLower(name)))
			isEmpty = false
		}
	}

	if isEmpty && photo.ID != 0 {
		_, err := t.Exec("delete FROM photo_tags WHERE photo_id=$1", photo.ID)
		t.Rollback()
		return errgo.Mask(err)
	}
	if _, err := t.Exec(fmt.Sprintf("SELECT add_tags(%s)", strings.Join(args, ",")), params...); err != nil {
		t.Rollback()
		return errgo.Mask(err)
	}
	return nil

}

func (m *defaultPhotoDataManager) get(photoID int64) (*photo, error) {

	p := &photo{}

	if photoID == 0 {
		return p, sql.ErrNoRows
	}

	obj, err := m.dbMap.Get(p, photoID)
	if err != nil {
		return p, errgo.Mask(err)
	}
	if obj == nil {
		return p, sql.ErrNoRows
	}
	return obj.(*photo), nil
}

func (m *defaultPhotoDataManager) getDetail(photoID int64, user *user) (*photoDetail, error) {

	photo := &photoDetail{}

	if photoID == 0 {
		return photo, sql.ErrNoRows
	}

	q := "SELECT p.*, u.name AS owner_name " +
		"FROM photos p JOIN users u ON u.id = p.owner_id " +
		"WHERE p.id=$1"

	if err := m.dbMap.SelectOne(photo, q, photoID); err != nil {
		return photo, errgo.Mask(err)
	}

	var tags []tag

	if _, err := m.dbMap.Select(&tags,
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

func (m *defaultPhotoDataManager) byOwnerID(page *page, ownerID int64) (*photoList, error) {
	var (
		photos []photo
		err    error
		total  int64
	)

	if ownerID == 0 {
		return nil, sql.ErrNoRows
	}
	if total, err = m.dbMap.SelectInt("SELECT COUNT(id) FROM photos WHERE owner_id=$1", ownerID); err != nil {
		return nil, errgo.Mask(err)
	}

	if _, err = m.dbMap.Select(&photos,
		"SELECT * FROM photos WHERE owner_id = $1"+
			"ORDER BY (up_votes - down_votes) DESC, created_at DESC LIMIT $2 OFFSET $3",
		ownerID, page.size, page.offset); err != nil {
		return nil, errgo.Mask(err)
	}
	return newPhotoList(photos, total, page.index), nil

}

func (m *defaultPhotoDataManager) search(page *page, q string) (*photoList, error) {

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

	if total, err = m.dbMap.SelectInt(countSql, params...); err != nil {
		return nil, errgo.Mask(err)
	}

	numParams := len(params)

	sql := fmt.Sprintf("SELECT * FROM (%s) q ORDER BY (up_votes - down_votes) DESC, created_at DESC LIMIT $%d OFFSET $%d",
		clausesSql, numParams+1, numParams+2)

	params = append(params, interface{}(page.size))
	params = append(params, interface{}(page.offset))

	if _, err = m.dbMap.Select(&photos, sql, params...); err != nil {
		return nil, errgo.Mask(err)
	}
	return newPhotoList(photos, total, page.index), nil
}

func (m *defaultPhotoDataManager) all(page *page, orderBy string) (*photoList, error) {

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

	if total, err = m.dbMap.SelectInt("SELECT COUNT(id) FROM photos"); err != nil {
		return nil, errgo.Mask(err)
	}

	if _, err = m.dbMap.Select(&photos,
		"SELECT * FROM photos "+
			"ORDER BY "+orderBy+" DESC LIMIT $1 OFFSET $2", page.size, page.offset); err != nil {
		return nil, errgo.Mask(err)
	}
	return newPhotoList(photos, total, page.index), nil
}

func (m *defaultPhotoDataManager) getTagCounts() ([]tagCount, error) {
	var tags []tagCount
	if _, err := m.dbMap.Select(&tags, "SELECT name, photo, num_photos FROM tag_counts"); err != nil {
		return tags, errgo.Mask(err)
	}
	return tags, nil
}

type userDataManager interface {
	create(user *user) error
	update(user *user) error
	isNameAvailable(user *user) (bool, error)
	isEmailAvailable(user *user) (bool, error)
	getActive(userID int64) (*user, error)
	getByRecoveryCode(string) (*user, error)
	getByEmail(string) (*user, error)
	getByNameOrEmail(identifier string) (*user, error)
}

type defaultUserDataManager struct {
	dbMap *gorp.DbMap
}

func (m *defaultUserDataManager) create(user *user) error {
	return errgo.Mask(m.dbMap.Insert(user))
}

func (m *defaultUserDataManager) update(user *user) error {
	_, err := m.dbMap.Update(user)
	return errgo.Mask(err)
}

func (m *defaultUserDataManager) isNameAvailable(user *user) (bool, error) {
	var (
		num int64
		err error
	)
	q := "SELECT COUNT(id) FROM users WHERE name=$1"
	if user.ID == 0 {
		num, err = m.dbMap.SelectInt(q, user.Name)
	} else {
		q += " AND id != $2"
		num, err = m.dbMap.SelectInt(q, user.Name, user.ID)
	}
	if err != nil {
		return false, errgo.Mask(err)
	}
	return num == 0, nil
}

func (m *defaultUserDataManager) isEmailAvailable(user *user) (bool, error) {
	var (
		num int64
		err error
	)
	q := "SELECT COUNT(id) FROM users WHERE email=$1"
	if user.ID == 0 {
		num, err = m.dbMap.SelectInt(q, user.Email)
	} else {
		q += " AND id != $2"
		num, err = m.dbMap.SelectInt(q, user.Email, user.ID)
	}
	if err != nil {
		return false, errgo.Mask(err)
	}
	return num == 0, nil
}

func (m *defaultUserDataManager) getActive(userID int64) (*user, error) {

	user := &user{}
	if err := m.dbMap.SelectOne(user, "SELECT * FROM users WHERE active=$1 AND id=$2", true, userID); err != nil {
		return user, errgo.Mask(err)
	}
	return user, nil

}

func (m *defaultUserDataManager) getByRecoveryCode(code string) (*user, error) {

	user := &user{}
	if code == "" {
		return user, sql.ErrNoRows
	}
	if err := m.dbMap.SelectOne(user, "SELECT * FROM users WHERE active=$1 AND recovery_code=$2", true, code); err != nil {
		return user, errgo.Mask(err)
	}
	return user, nil

}
func (m *defaultUserDataManager) getByEmail(email string) (*user, error) {
	user := &user{}
	if err := m.dbMap.SelectOne(user, "SELECT * FROM users WHERE active=$1 AND email=$2", true, email); err != nil {
		return user, errgo.Mask(err)
	}
	return user, nil
}

func (m *defaultUserDataManager) getByNameOrEmail(identifier string) (*user, error) {
	user := &user{}

	if err := m.dbMap.SelectOne(user, "SELECT * FROM users WHERE active=$1 AND (email=$2 OR name=$2)", true, identifier); err != nil {
		return user, errgo.Mask(err)
	}

	return user, nil
}
