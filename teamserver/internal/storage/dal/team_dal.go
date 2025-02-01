package dal

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ksel172/Meduza/teamserver/models"
)

type ITeamDAL interface {
	CreateTeam(ctx context.Context, team *models.Team) error
	UpdateTeam(ctx context.Context, team *models.Team) error
	DeleteTeam(ctx context.Context, teamID string) error
	GetTeams(ctx context.Context) ([]models.Team, error)
	AddTeamMember(ctx context.Context, member *models.TeamMember) error
	RemoveTeamMember(ctx context.Context, memberID string) error
	GetTeamMembers(ctx context.Context, teamID string) ([]models.TeamMember, error)
}

type TeamDAL struct {
	db     *sql.DB
	schema string
}

func NewTeamDAL(db *sql.DB, schema string) *TeamDAL {
	return &TeamDAL{db: db, schema: schema}
}

func (dal *TeamDAL) CreateTeam(ctx context.Context, team *models.Team) error {
	// Check if a team with the same name already exists
	var existingTeamID string
	checkQuery := fmt.Sprintf(`SELECT id FROM %s.teams WHERE name=$1`, dal.schema)
	err := dal.db.QueryRowContext(ctx, checkQuery, team.Name).Scan(&existingTeamID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if existingTeamID != "" {
		return fmt.Errorf("a team with the name '%s' already exists", team.Name)
	}

	// Create the new team
	query := fmt.Sprintf(`INSERT INTO %s.teams(name) VALUES($1)`, dal.schema)
	_, err = dal.db.ExecContext(ctx, query, team.Name)
	return err
}

func (dal *TeamDAL) UpdateTeam(ctx context.Context, team *models.Team) error {
	query := fmt.Sprintf(`UPDATE %s.teams SET name=$1 WHERE id=$2`, dal.schema)
	_, err := dal.db.ExecContext(ctx, query, team.Name, team.ID)
	return err
}

func (dal *TeamDAL) DeleteTeam(ctx context.Context, teamID string) error {
	query := fmt.Sprintf(`DELETE FROM %s.teams WHERE id=$1`, dal.schema)
	_, err := dal.db.ExecContext(ctx, query, teamID)
	return err
}

func (dal *TeamDAL) GetTeams(ctx context.Context) ([]models.Team, error) {
	query := fmt.Sprintf(`SELECT id, name, created_at, updated_at FROM %s.teams`, dal.schema)
	rows, err := dal.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []models.Team
	for rows.Next() {
		var team models.Team
		if err := rows.Scan(&team.ID, &team.Name, &team.CreatedAt, &team.UpdatedAt); err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	return teams, nil
}

func (dal *TeamDAL) AddTeamMember(ctx context.Context, member *models.TeamMember) error {
	// Check if the user is already a member of the team
	var existingMemberID string
	checkQuery := fmt.Sprintf(`SELECT id FROM %s.team_members WHERE team_id=$1 AND user_id=$2`, dal.schema)
	err := dal.db.QueryRowContext(ctx, checkQuery, member.TeamID, member.UserID).Scan(&existingMemberID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if existingMemberID != "" {
		return fmt.Errorf("user '%s' is already a member of team '%s'", member.UserID, member.TeamID)
	}

	// Add the new team member
	query := fmt.Sprintf(`INSERT INTO %s.team_members(team_id, user_id, added_at) VALUES($1, $2, $3)`, dal.schema)
	_, err = dal.db.ExecContext(ctx, query, member.TeamID, member.UserID, member.AddedAt)
	return err
}

func (dal *TeamDAL) RemoveTeamMember(ctx context.Context, memberID string) error {
	query := fmt.Sprintf(`DELETE FROM %s.team_members WHERE id=$1`, dal.schema)
	_, err := dal.db.ExecContext(ctx, query, memberID)
	return err
}

func (dal *TeamDAL) GetTeamMembers(ctx context.Context, teamID string) ([]models.TeamMember, error) {
	query := fmt.Sprintf(`
        SELECT tm.id, tm.team_id, tm.user_id, u.username, u.role, tm.added_at
        FROM %s.team_members tm
        JOIN %s.users u ON tm.user_id = u.id
        WHERE tm.team_id = $1`, dal.schema, dal.schema)
	rows, err := dal.db.QueryContext(ctx, query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.TeamMember
	for rows.Next() {
		var member models.TeamMember
		if err := rows.Scan(&member.ID, &member.TeamID, &member.UserID, &member.Username, &member.Role, &member.AddedAt); err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, nil
}
