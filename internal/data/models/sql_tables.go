package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Movie struct {
	bun.BaseModel `bun:"table:movies"`
	ID            int64     `bun:"id,pk,type:bigserial,notnull"`
	Title         string    `bun:"title,type:text,notnull"`
	Year          int16     `bun:"year,type:smallint,notnull"`
	Synopsis      string    `bun:"synopsis,type:longtext"`
	ImdbRating    float32   `bun:"imdb_rating,type:decimal(3,1)"`
	UserRating    float32   `bun:"user_rating,type:decimal(3,1)"`
	CreatedAt     time.Time `bun:"createdAt,type:time,notnull,default:current_timestamp"`

	// Cast is composed with people for movie characters
	Cast []*People `bun:"m2m:movie_cast,join:Movie=Person"`
	// Crew is composed with people for movie jobs, like director
	Crew []*People `bun:"m2m:movie_crew,join:Movie=Person"`
	// Genres
	Genres []*Genre `bun:"m2m:movie_genres,join:Movie=Genre"`
	// Themes
	Themes []*Theme `bun:"m2m:movie_themes,join:Movie=Theme"`
}

// UserRating belongs to Movie
type UserRating struct {
	bun.BaseModel `bun:"table:user_ratings"`
	ID            int64 `bun:"id,pk,type:bigserial,notnull"`
	// Directing is the director’s responsibility to unify the artistic, emotional, and technical
	// dimensions of a film, guiding actors, crew, and the camera itself toward a single coherent
	// vision.
	// Good directing clarifies intention, shapes performance, sculpts pacing, and ensures that
	// every creative choice serves the story’s deepest truth.
	Directing int8 `bun:"directing,tinyint,notnull"`
	// Photography is the art of shaping time and emotion by transforming space into story:
	// It is the deliberate orchestration of illumination, shadow, color, texture, and perspective
	// to guide the audience’s gaze and their feeling.
	// A good photography shot does not exist to be admired on its own, but to serve the rhythm,
	// meaning, and soul of the film.
	// Photography is not always beautiful but always meaningful.
	Photography int8 `bun:"photography,tinyint,notnull"`
	// Acting is the art of embodying human truth within a character's fiction. It is the craft of
	// lending a character a breath, instincts, emotions, so that an invented life feels lived.
	// A good acting shapes behavior so that every gesture, silence, and shift of the eyes reflects
	// an inner world.
	Acting int8 `bun:"acting_dubbing,tinyint,notnull"`
	// Scenario is the narrative architecture of the film: the structured articulation of its story,
	// characters, conflicts, and dramatic progression.
	// A good scenario provides, sufficiently and no more no less, clarity of intention while
	// leaving space for creative interpretation. It defines how each moment connects to others
	// without making everything explicit.
	Scenario int8 `bun:"scenario,tinyint,notnull"`
	// Music and sound are the 4th dimension of a movie, it is the invisible actor that helps to
	// shape emotions. It is composed of musical score, sound effects, ambient noise and silence
	// itself.
	// Good film's music and sound is like a heartbeat. It reveals subtext, punctuates dramatic
	// beats, and enhances the storytelling in ways the eye alone cannot.
	MusicSound int8 `bun:"music_sound,tinyint,notnull"`
	// Bonus rating is arbitrary, and it is up to the user to provide one or none.
	// Comment section can be used to explain why
	Bonus int8 `bun:"bonus,tinyint,notnull"`
	// Loved rating reveals that the user fell for the movie, even if the rating is bad. We all
	// have this bad movie that we love for any reason that only our heart knows ! :)
	Loved bool `bun:"loved,boolean,notnull"`
	// Comment let the user provide any comment they want about the movie
	Comment   string    `bun:"comment,text"`
	CreatedAt time.Time `bun:"createdAt,type:time,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updatedAt,type:time,notnull,default:current_timestamp"`

	// movie relationship
	MovieID *int64 `bun:"movie_id"`
	Movie   *Movie `bun:"rel:belongs-to,join:movie_id=id"`
}

// People play as characters or work for a job in movies
type People struct {
	bun.BaseModel `bun:"table:user_ratings"`
	ID            int64  `bun:"id,pk,type:bigserial,notnull"`
	FirstName     string `bun:"first_name,type:text,notnull"`
	LastName      string `bun:"last_name,type:text,notnull"`
	BirthDate     string `bun:"birth_date,type:text,notnull"`

	// movies characters in castings
	Casts []*Movie `bun:"m2m:movie_cast,join:Person=Movie"`
	// movies roles in castings (director, sound, ...)
	Crews []*Movie `bun:"m2m:movie_crew,join:Person=Movie"`
}

// Cast :
//   - belongs-to one movie
//   - has many people
//
// It is also in essence the relation table between people and a movie for characters
// It doesn't have an ID because it is purely relational between people and movies
type Cast struct {
	bun.BaseModel `bun:"table:casts"`
	CharacterName string `bun:"character_name,type:text,notnull"`
	// movie relationship
	MovieID int64  `bun:"movie_id,pk"`
	Movie   *Movie `bun:"rel:belongs-to,join:movie_id=id"`
	// People relationship
	PeopleID int64   `bun:"people_id,pk"`
	People   *People `bun:"rel:belongs-to,join:people_id=id"`
}

// Crew :
//   - belongs-to one movie
//   - has many people
//
// It is in essence the relation table between people and a movie for jobs
// It doesn't have an ID because it is purely relational between people and movies
type Crew struct {
	bun.BaseModel `bun:"table:casts"`
	RoleName      string `bun:"role_name,type:text,notnull"`
	// movie relationship
	MovieID int64  `bun:"movie_id,pk"`
	Movie   *Movie `bun:"rel:belongs-to,join:movie_id=id"`
	// People relationship
	PeopleID int64   `bun:"people_id,pk"`
	People   *People `bun:"rel:belongs-to,join:people_id=id"`
}

// Genre is what comes first in your mind when you categorize a movie.
// Like 'Horror', 'Action' or 'Drama'.
// It is quite static and the users won't be able to add more.
type Genre struct {
	bun.BaseModel `bun:"table:genres"`
	ID            int64  `bun:"id,pk,type:bigserial,notnull"`
	Name          string `bun:"name,type:text,notnull"`

	Movies []*Movie `bun:"m2m:movie_genres,join:Genre=Movie"`
}

// MovieGenre is obviously the relation between a movie and a genre
type MovieGenre struct {
	bun.BaseModel `bun:"table:movie_genres"`

	MovieID int64  `bun:"movie_id,pk"`
	Movie   *Movie `bun:"rel:belongs-to,join:movie_id=id"`

	GenreID int64  `bun:"genre_id,pk"`
	Genre   *Genre `bun:"rel:belongs-to,join:genre_id=id"`
}

// Theme is like a custom genre that the user can set to add more vision or information
// they value about a movie.
// Like "vengeance", "cathartic" or "chess".
// In other words, a movie genre and a theme can compose 'Action' and 'Vengeance, Chess'.
type Theme struct {
	bun.BaseModel `bun:"table:themes"`
	ID            int64  `bun:"id,pk,type:bigserial,notnull"`
	Name          string `bun:"name,type:text,notnull"`

	Movies []*Movie `bun:"m2m:movie_themes,join:Theme=Movie"`
}

// MovieTheme is obviously the relation between a movie and a theme
type MovieTheme struct {
	bun.BaseModel `bun:"table:movie_themes"`

	MovieID int64  `bun:"movie_id,pk"`
	Movie   *Movie `bun:"rel:belongs-to,join:movie_id=id"`

	ThemeID int64  `bun:"theme_id,pk"`
	Theme   *Genre `bun:"rel:belongs-to,join:theme_id=id"`
}
