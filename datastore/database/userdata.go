package database

import (
	"fmt"

	"github.com/apex/log"
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
			if err := tx.QueryRow(`
				INSERT INTO users(login)
				VALUES($1)
				ON CONFLICT(login) DO
					UPDATE SET login = $1, updated_at = current_timestamp
				RETURNING id
			`, login).Scan(&followerID); err != nil {
				return fmt.Errorf("insert users failed: %w", err)
			}

			if _, err := tx.Exec(`
				INSERT INTO user_followers(user_id, follower_id)
				VALUES ($1, $2)
				ON CONFLICT DO NOTHING
			`, userID, followerID); err != nil {
				return fmt.Errorf("insert user_followers failed: %w", err)
			}
		}
		return nil
	})
}

// GetStars of a given userID
func (db *Userdatastore) GetStars(userID int64) ([]model.Star, error) {
	var repos []string
	if err := db.Select(&repos, `
		SELECT name
		FROM repositories
		WHERE user_id = $1
	`, userID); err != nil {
		return []model.Star{}, err
	}

	var stars []model.Star

	for _, repo := range repos {
		var stargazers []string
		if err := db.Select(&stargazers, `
			SELECT s.login
			FROM stargazers s
			JOIN repositories r ON
				r.user_id= $1 AND
				r.name = $2
				s.repository_id = r.id
		`, userID, repo); err != nil {
			return []model.Star{}, err
		}

		stars = append(stars, model.Star{
			RepoName:   repo,
			Stargazers: stargazers,
		})
	}

	return stars, nil
}

// SaveStars for a given userID
func (db *Userdatastore) SaveStars(userID int64, stars []model.Star) error {
	return db.WithTx(func(tx *sqlx.Tx) error {
		for _, star := range stars {
			var repoID int64
			// TODO: what if multiple users have the same repository? e.g. from orgs
			if err := tx.QueryRow(`
				INSERT INTO repositories(name, user_id)
				VALUES($1, $2)
				ON CONFLICT DO UPDATE
					SET name = $1,
						updated_at = current_timestamp
				RETURNING id
			`, star.RepoName, userID).Scan(&repoID); err != nil {
				return fmt.Errorf("insert repositories failed: %w", err)
			}

			for _, login := range star.Stargazers {
				var stargazerID int64
				if err := tx.QueryRow(`
					INSERT INTO users(login)
					VALUES($1)
					ON CONFLICT DO UPDATE
						SET login = $1,
							updated_at = current_timestamp
					RETURNING id
				`, login).Scan(&stargazerID); err != nil {
					return fmt.Errorf("insert users failed: %w", err)
				}

				if _, err := tx.Exec(`
					INSERT INTO starred_repositories(repository_id, stargazer_id)
					VALUES ($1, $2)
					ON CONFLICT DO NOTHING
				`, repoID, stargazerID); err != nil {
					return fmt.Errorf("insert starred_repositories failed: %w", err)
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
	log.Info("beginning tx")
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		log.Info("rolling back tx")
		if rerr := tx.Rollback(); rerr != nil {
			err = multierror.Append(err, rerr)
		}
		return err
	}

	log.Info("commiting tx")
	return tx.Commit()
}
