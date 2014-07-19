package api

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"strings"
)

func InitDB(config *AppConfig) (*sqlx.DB, error) {
	return sqlx.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s host=%s",
		config.DBUser,
		config.DBName,
		config.DBPassword,
		config.DBHost,
	))
}

type DataStore interface {
	Begin() (Transaction, error)

	GetPhoto(int64) (*Photo, error)
	GetPhotoDetail(int64, *User) (*PhotoDetail, error)
	GetTagCounts() ([]TagCount, error)
	GetPhotos(*Page, string) (*PhotoList, error)
	GetPhotosByOwnerID(*Page, int64) (*PhotoList, error)
	SearchPhotos(*Page, string) (*PhotoList, error)

	IsUserNameAvailable(user *User) (bool, error)
	IsUserEmailAvailable(user *User) (bool, error)

	GetActiveUser(userID int64) (*User, error)
	GetUserByRecoveryCode(string) (*User, error)
	GetUserByEmail(string) (*User, error)
	GetUserByNameOrEmail(identifier string) (*User, error)
}

type Transaction interface {
	InsertPhoto(*Photo) error
	UpdatePhoto(*Photo) error
	DeletePhoto(*Photo) error
	UpdateTags(*Photo) error
	InsertUser(user *User) error
	UpdateUser(user *User) error
	Commit() error
	Rollback() error
}

type defaultDataStore struct {
	*sqlx.DB
}

type defaultTransaction struct {
	*sqlx.Tx
}

func NewDataStore(db *sqlx.DB) DataStore {
	return &defaultDataStore{db}
}

func (tx *defaultTransaction) DeletePhoto(photo *Photo) error {
	_, err := tx.Exec("DELETE FROM photos WHERE id=$1", photo.ID)
	return err
}

func (tx *defaultTransaction) UpdatePhoto(photo *Photo) error {
	_, err := tx.Exec(
		"UPDATE photos SET title=$1, up_votes=$2, down_votes=$3 WHERE id=$4",
		photo.Title,
		photo.UpVotes,
		photo.DownVotes,
		photo.ID,
	)
	return err
}

func (tx *defaultTransaction) InsertPhoto(photo *Photo) error {
	stmt, err := tx.PrepareNamed(
		"INSERT INTO photos (owner_id, title, photo) " +
			"VALUES(:owner_id, :title, :photo) RETURNING id, created_at",
	)
	if err != nil {
		return err
	}
	return stmt.QueryRowx(photo).Scan(&photo.ID, &photo.CreatedAt)
}

func (tx *defaultTransaction) UpdateTags(photo *Photo) error {

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
		_, err := tx.Exec("DELETE FROM photo_tags WHERE photo_id=$1", photo.ID)
		return err
	}
	_, err := tx.Exec(fmt.Sprintf("SELECT add_tags(%s)", strings.Join(args, ",")), params...)
	return err

}

func (ds *defaultDataStore) Begin() (Transaction, error) {
	tx, err := ds.Beginx()
	if err != nil {
		return nil, err
	}
	return &defaultTransaction{tx}, nil
}

func (ds *defaultDataStore) GetPhoto(photoID int64) (*Photo, error) {

	photo := &Photo{}

	if photoID == 0 {
		return photo, sql.ErrNoRows
	}
	if err := ds.Get(photo, "SELECT * FROM photos WHERE id=$1", photoID); err != nil {
		return photo, err
	}
	return photo, nil
}

func (ds *defaultDataStore) GetPhotoDetail(photoID int64, user *User) (*PhotoDetail, error) {

	photo := &PhotoDetail{}

	if photoID == 0 {
		return photo, sql.ErrNoRows
	}

	q := "SELECT p.*, u.name AS owner_name " +
		"FROM photos p JOIN users u ON u.id = p.owner_id " +
		"WHERE p.id=$1"

	if err := ds.Get(photo, q, photoID); err != nil {
		return photo, err
	}

	var tags []Tag

	if err := ds.Select(&tags,
		"SELECT t.* FROM tags t JOIN photo_tags pt ON pt.tag_id=t.id "+
			"WHERE pt.photo_id=$1", photo.ID); err != nil {
		return photo, err
	}
	for _, tag := range tags {
		photo.Tags = append(photo.Tags, tag.Name)
	}

	photo.Permissions = &Permissions{
		photo.CanEdit(user),
		photo.CanDelete(user),
		photo.CanVote(user),
	}
	return photo, nil

}

func (ds *defaultDataStore) GetPhotosByOwnerID(page *Page, ownerID int64) (*PhotoList, error) {

	var (
		photos []Photo
		err    error
		total  int64
	)
	if ownerID == 0 {
		return nil, nil
	}

	row := ds.QueryRow("SELECT COUNT(id) FROM photos WHERE owner_id=$1", ownerID)
	if err := row.Scan(&total); err != nil {
		return nil, err
	}

	if err = ds.Select(&photos,
		"SELECT * FROM photos WHERE owner_id = $1"+
			"ORDER BY (up_votes - down_votes) DESC, created_at DESC LIMIT $2 OFFSET $3",
		ownerID, page.Size, page.Offset); err != nil {
		return nil, err
	}
	return NewPhotoList(photos, total, page.Index), nil
}

func (ds *defaultDataStore) SearchPhotos(page *Page, q string) (*PhotoList, error) {

	var (
		clauses []string
		params  []interface{}
		err     error
		photos  []Photo
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

		num += 1

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

	row := ds.QueryRow(countSql, params...)
	if err := row.Scan(&total); err != nil {
		return nil, err
	}

	numParams := len(params)

	sql := fmt.Sprintf(
		"SELECT * FROM (%s) q ORDER BY (up_votes - down_votes) DESC, created_at DESC LIMIT $%d OFFSET $%d",
		clausesSql, numParams+1, numParams+2)

	params = append(params, interface{}(page.Size))
	params = append(params, interface{}(page.Offset))

	if err = ds.Select(&photos, sql, params...); err != nil {
		return nil, err
	}
	return NewPhotoList(photos, total, page.Index), nil
}

func (ds *defaultDataStore) GetPhotos(page *Page, orderBy string) (*PhotoList, error) {

	var (
		total  int64
		photos []Photo
		err    error
	)
	if orderBy == "votes" {
		orderBy = "(up_votes - down_votes)"
	} else {
		orderBy = "created_at"
	}

	row := ds.QueryRow("SELECT COUNT(id) FROM photos")
	if err := row.Scan(&total); err != nil {
		return nil, err
	}

	if err = ds.Select(&photos,
		"SELECT * FROM photos "+
			"ORDER BY "+orderBy+" DESC LIMIT $1 OFFSET $2", page.Size, page.Offset); err != nil {
		return nil, err
	}
	return NewPhotoList(photos, total, page.Index), nil
}

func (ds *defaultDataStore) GetTagCounts() ([]TagCount, error) {
	var tags []TagCount
	if err := ds.Select(&tags, "SELECT name, photo, num_photos FROM tag_counts"); err != nil {
		return tags, err
	}
	return tags, nil
}

func (tx *defaultTransaction) InsertUser(user *User) error {
	stmt, err := tx.PrepareNamed(
		"INSERT INTO users (name, email, password, admin) " +
			"VALUES(:name, :email, :password, :admin) RETURNING id, created_at",
	)
	if err != nil {
		return err
	}
	return stmt.QueryRowx(user).Scan(&user.ID, &user.CreatedAt)
}

func (tx *defaultTransaction) UpdateUser(user *User) error {
	_, err := tx.Exec(
		"UPDATE users SET name=$1, email=$2, password=$3, is_admin=$4, "+
			"is_active=$5 recovery_code=$5 WHERE id=$6",
		user.Name,
		user.Email,
		user.Password,
		user.IsAdmin,
		user.IsActive,
		user.RecoveryCode,
		user.ID,
	)
	return err
}

func (ds *defaultDataStore) IsUserNameAvailable(user *User) (bool, error) {
	var (
		num int64
		err error
	)
	q := "SELECT COUNT(id) FROM users WHERE name=$1"
	if user.ID == 0 {
		err = ds.Get(&num, q, user.Name)
	} else {
		q += " AND id != $2"
		err = ds.Get(&num, q, user.Name, user.ID)
	}
	if err != nil {
		return false, err
	}
	return num == 0, nil
}

func (ds *defaultDataStore) IsUserEmailAvailable(user *User) (bool, error) {
	var (
		num int64
		err error
	)
	q := "SELECT COUNT(id) FROM users WHERE email=$1"
	if user.ID == 0 {
		err = ds.Get(&num, q, user.Email)
	} else {
		q += " AND id != $2"
		err = ds.Get(&num, q, user.Email, user.ID)
	}
	if err != nil {
		return false, err
	}
	return num == 0, nil
}

func (ds *defaultDataStore) GetActiveUser(userID int64) (*User, error) {

	user := &User{}
	if err := ds.Get(user, "SELECT * FROM users WHERE active=$1 AND id=$2", true, userID); err != nil {
		return user, err
	}
	return user, nil

}

func (ds *defaultDataStore) GetUserByRecoveryCode(code string) (*User, error) {

	user := &User{}
	if code == "" {
		return user, sql.ErrNoRows
	}
	if err := ds.Get(user, "SELECT * FROM users WHERE active=$1 AND recovery_code=$2", true, code); err != nil {
		return user, err
	}
	return user, nil

}
func (ds *defaultDataStore) GetUserByEmail(email string) (*User, error) {
	user := &User{}
	if err := ds.Get(user, "SELECT * FROM users WHERE active=$1 AND email=$2", true, email); err != nil {
		return user, err
	}
	return user, nil
}

func (ds *defaultDataStore) GetUserByNameOrEmail(identifier string) (*User, error) {
	user := &User{}

	if err := ds.Get(user, "SELECT * FROM users WHERE active=$1 AND (email=$2 OR name=$2)", true, identifier); err != nil {
		return user, err
	}

	return user, nil
}
