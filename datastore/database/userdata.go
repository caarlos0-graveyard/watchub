package database

import (
	"encoding/json"

	"github.com/caarlos0/watchub/shared/model"
	"github.com/hashicorp/go-multierror"
	"github.com/jmoiron/sqlx"
)

// Userdatastore in database
type Userdatastore struct {
	*sqlx.DB
}

// NewUserdatastore datastore
func NewUserdatastore(db *sqlx.DB) *Userdatastore {
	return &Userdatastore{db}
}

func (db *Userdatastore) Disable(userID int64) error {
	_, err := db.Exec(`
		UPDATE tokens
		SET disabled = true,
			updated_at = current_timestamp
		WHERE id = $1
	`, userID)
	return err
}

// GetFollowers of a given userID
func (db *Userdatastore) GetFollowers(userID int64) ([]string, error) {
	var logins []string
	err := db.Select(&logins, `
		SELECT login FROM users u
		JOIN user_followers uf ON
			uf.user_id = $1 AND
			u.id = uf.follower_id
	`, userID)
	return logins, err
}

// SaveFollowers for a given userID
func (db *Userdatastore) SaveFollowers(userID int64, followers []string) error {
	return db.WithTx(func(tx *sqlx.Tx) error {
		for _, login := range followers {
			var followerID int64
			if err := tx.Select(&followerID, `
				INSERT INTO users(login)
				VALUES($1)
				ON CONFLICT(login) DO
					UPDATE SET login = $1, updated_at = current_timestamp
				RETURNING id
			`, login); err != nil {
				return err
			}

			if _, err := tx.Exec(`
				INSERT INTO user_followers(user_id, follower_id)
				VALUES ($1, $2)
				ON CONFLICT DO NOTHING
			`, userID, followerID); err != nil {
				return err
			}
		}
		return nil
	})
}

// GetStars of a given userID
func (db *Userdatastore) GetStars(userID int64) (result []model.Star, err error) {
	var stars json.RawMessage
	if err := db.QueryRow(`
		SELECT stars
		FROM tokens
		WHERE id = $1
	`, userID).Scan(&stars); err != nil {
		return result, err
	}
	return result, json.Unmarshal(stars, &result)
}

// SaveStars for a given userID
func (db *Userdatastore) SaveStars(userID int64, stars []model.Star) error {
	return db.WithTx(func(tx *sqlx.Tx) error {
		for _, star := range stars {
			var repoID int64
			// TODO: what if multiple users have the same repository? e.g. from orgs
			if err := tx.Select(&repoID, `
				INSERT INTO repositories(name, user_id)
				VALUES($1, $2)
				ON CONFLICT DO UPDATE
					SET name = $1,
						updated_at = current_timestamp
				RETURNING id
			`, star.RepoName, userID); err != nil {
				return err
			}

			for _, login := range star.Stargazers {
				var userID int64
				if err := tx.Select(&userID, `
					INSERT INTO users(login)
					VALUES($1)
					ON CONFLICT DO UPDATE
						SET login = $1,
							updated_at = current_timestamp
					RETURNING id
				`, login); err != nil {
					return err
				}

				if _, err := tx.Exec(`
					INSERT INTO starred_repositories(repository_id, stargazer_id)
					VALUES ($1, $2)
					ON CONFLICT DO NOTHING
				`, repoID, userID); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

const followerCountQuery = `
	SELECT count(id)
	FROM user_followers
	WHERE user_id = $1
`

// FollowerCount returns the amount of followers stored for a given userID
func (db *Userdatastore) FollowerCount(userID int64) (count int, err error) {
	err = db.QueryRow(followerCountQuery, userID).Scan(&count)
	return
}

const starCountQuery = `
	SELECT count(sr.id)
	FROM starred_repositories sr
	JOIN repositories r ON
		r.user_id = $1 AND
		sr.repository_id = r.id
`

// StarCount returns the amount of stargazers of all the user's repositories
func (db *Userdatastore) StarCount(userID int64) (count int, err error) {
	err = db.QueryRow(starCountQuery, userID).Scan(&count)
	return
}

const repositoryCountQuery = `
	SELECT count(id)
	FROM repositories
	WHERE user_id = $1
`

// RepositoryCount returns the amount of followers stored for a given userID
func (db *Userdatastore) RepositoryCount(userID int64) (int, error) {
	var count int
	err := db.QueryRow(repositoryCountQuery, userID).Scan(&count)
	return count, err
}

// UserExist returns true if an user is already registered in the db
func (db *Userdatastore) UserExist(userID int64) (bool, error) {
	var result bool
	err := db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM tokens
			WHERE id = $1
		)
	`, userID).Scan(&result)
	return result, err
}

func (db *Userdatastore) WithTx(fn func(tx *sqlx.Tx) error) error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = multierror.Append(err, rerr)
		}
		return err
	}
	return tx.Commit()
}
