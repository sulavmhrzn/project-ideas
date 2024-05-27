package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/sulavmhrzn/projectideas/internal/validator"
)

type Idea struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UserId      int       `json:"-"`
	Tags        []Tag     `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
}

type Tag struct {
	Id    int    `json:"-"`
	Title string `json:"title"`
}

func ValidateIdea(v *validator.Validator, idea *Idea) {
	v.Check(idea.Title != "", "title", "must be provided")
	v.Check(len(idea.Title) < 100, "title", "must be smaller than 100 characters")
	v.Check(idea.Description != "", "description", "must be provided")
	v.Check(len(idea.Tags) != 0, "tags", "must be provided")
	for _, i := range idea.Tags {
		v.Check(i.Title != "", "title", "tags title must be provided")
	}
	var tagTitles []string
	for _, i := range idea.Tags {
		tagTitles = append(tagTitles, i.Title)
	}
	v.Check(validator.Unique(tagTitles...), "title", "tag title must be unique")
}

type IdeaModel struct {
	DB *sql.DB
}

func (m IdeaModel) Insert(input *Idea) (*Idea, error) {
	insertIdeaQuery := `
	INSERT INTO ideas 
	(title, description, user_id)
	VALUES 
	($1, $2, $3)
	RETURNING id, title, description, created_at`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var idea Idea
	err := m.DB.QueryRowContext(
		ctx,
		insertIdeaQuery,
		[]any{input.Title, input.Description, input.UserId}...,
	).Scan(&idea.Id, &idea.Title, &idea.Description, &idea.CreatedAt)

	if err != nil {
		return nil, err
	}

	for _, t := range input.Tags {
		selectTagQuery := `
		SELECT id, title FROM tags WHERE title = $1`
		var tag Tag
		err := m.DB.QueryRowContext(ctx, selectTagQuery, t.Title).Scan(&tag.Id, &tag.Title)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				err := m.DB.QueryRowContext(ctx, `INSERT INTO tags (title) VALUES ($1) RETURNING id, title`, t.Title).Scan(&tag.Id, &tag.Title)
				if err != nil {
					return nil, err
				}
			default:
				return nil, err
			}
		}
		idea.Tags = append(idea.Tags, tag)
		insertIdeasTagsQuery := `
		INSERT INTO ideas_tags (idea_id, tag_id)
		VALUES 
		($1, $2)`
		_, err = m.DB.ExecContext(ctx, insertIdeasTagsQuery, []any{idea.Id, tag.Id}...)
		if err != nil {
			return nil, err
		}
	}
	return &idea, nil
}

func (m IdeaModel) List() ([]Idea, error) {
	query := `
	SELECT ideas.id, ideas.title, ideas.description, ideas.created_at, tags.id, tags.title
	FROM ideas
	JOIN ideas_tags ON ideas_tags.idea_id = ideas.id
	JOIN tags ON ideas_tags.tag_id = tags.id
	ORDER BY created_at DESC`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ideasMap := make(map[int]*Idea)
	for rows.Next() {
		var ideaID int
		var ideaTitle, ideaDescription string
		var createdAt time.Time
		var tagID int
		var tagTitle string

		err := rows.Scan(&ideaID, &ideaTitle, &ideaDescription, &createdAt, &tagID, &tagTitle)
		if err != nil {
			return nil, err
		}
		if idea, exist := ideasMap[ideaID]; exist {
			idea.Tags = append(idea.Tags, Tag{Id: tagID, Title: tagTitle})
		} else {
			ideasMap[ideaID] = &Idea{
				Id:          ideaID,
				Title:       ideaTitle,
				Description: ideaDescription,
				CreatedAt:   createdAt,
				Tags:        []Tag{{Id: tagID, Title: tagTitle}},
			}
		}
	}

	var ideas []Idea
	for _, idea := range ideasMap {
		ideas = append(ideas, *idea)
	}
	return ideas, nil
}
