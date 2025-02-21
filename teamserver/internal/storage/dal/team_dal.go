package dal

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
	"github.com/ksel172/Meduza/teamserver/utils"
)

type ITeamDAL interface {
	CreateTeam(ctx context.Context, team *models.Team, creatorID string) error
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

func (dal *TeamDAL) CreateTeam(ctx context.Context, team *models.Team, creatorID string) error {
	return utils.WithTransactionTimeout(ctx, dal.db, 10, sql.TxOptions{}, func(context.Context, *sql.Tx) error {
		// Check if a team with the same name already exists
		var existingTeamID string
		checkQuery := fmt.Sprintf(`SELECT id FROM %s.teams WHERE name=$1`, dal.schema)
		queryStmt, err := dal.db.PrepareContext(ctx, checkQuery)
		if err != nil {
			logger.Error(logLevel, logDetailTeam, fmt.Sprintf("Failed to prepare read team query: %v", err))
			return fmt.Errorf("failed to prepare read team query: %w", err)
		}

		// Expected behavior is no rows returned
		if err = queryStmt.QueryRowContext(ctx, team.Name).Scan(&existingTeamID); err != sql.ErrNoRows {
			if err != nil {
				logger.Error(logLevel, logDetailTeam, fmt.Sprintf("Failed to query team: %v", err))
				return fmt.Errorf("failed to query team: %w", err)
			}
			return fmt.Errorf("team with the name '%s' already exists", team.Name)
		}

		// Create the new team
		createQuery := fmt.Sprintf(`INSERT INTO %s.teams(name) VALUES($1) RETURNING id`, dal.schema)
		createStmt, err := dal.db.PrepareContext(ctx, createQuery)
		if err != nil {
			logger.Error(logLevel, logDetailTeam, fmt.Sprintf("failed to prepare create team query: %v", err))
			return fmt.Errorf("failed to prepare create team query: %w", err)
		}
		if err = createStmt.QueryRowContext(ctx, team.Name).Scan(&team.ID); err != nil {
			logger.Error(logLevel, logDetailTeam, fmt.Sprintf("Failed to create team: %v", err))
			return fmt.Errorf("failed to create team: %w", err)
		}

		// Add the creator as a team member
		addMemberQuery := fmt.Sprintf(`INSERT INTO %s.team_members(team_id, user_id) VALUES($1, $2)`, dal.schema)
		addMemberStmt, err := dal.db.PrepareContext(ctx, addMemberQuery)
		if err != nil {
			logger.Error(logLevel, logDetailTeam, fmt.Sprintf("Failed to prepare create team query: %v", err))
			return fmt.Errorf("failed to prepare create team query: %w", err)
		}
		_, err = addMemberStmt.ExecContext(ctx, addMemberQuery, team.ID, creatorID)
		if err != nil {
			logger.Error(logLevel, logDetailTeam, fmt.Sprintf("Failed to add creator as team member: %v", err))
			return fmt.Errorf("failed to add creator as team member: %w", err)
		}

		return nil
	})
}

func (dal *TeamDAL) UpdateTeam(ctx context.Context, team *models.Team) error {
	query := fmt.Sprintf(`UPDATE %s.teams SET name=$1 WHERE id=$2`, dal.schema)
	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		_, err := stmt.ExecContext(ctx, team.Name, team.ID)
		if err != nil {
			logger.Error(logLevel, logDetailTeam, fmt.Sprintf("Failed to update team: %v", err))
			return fmt.Errorf("failed to update team: %w", err)
		}
		return nil
	})
}

func (dal *TeamDAL) DeleteTeam(ctx context.Context, teamID string) error {
	query := fmt.Sprintf(`DELETE FROM %s.teams WHERE id=$1`, dal.schema)
	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		_, err := stmt.ExecContext(ctx, teamID)
		if err != nil {
			logger.Error(logLevel, logDetailTeam, fmt.Sprintf("Failed to delete team: %v", err))
			return fmt.Errorf("failed to delete team: %w", err)
		}
		return nil
	})
}

func (dal *TeamDAL) GetTeams(ctx context.Context) ([]models.Team, error) {
	query := fmt.Sprintf(`SELECT id, name, created_at, updated_at FROM %s.teams`, dal.schema)

	return utils.WithResultTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) ([]models.Team, error) {
		rows, err := stmt.QueryContext(ctx)
		if err != nil {
			logger.Error(logLevel, logDetailTeam, fmt.Sprintf("Failed to get teams: %v", err))
			return nil, fmt.Errorf("failed to get teams: %w", err)
		}
		defer rows.Close()

		var teams []models.Team
		for rows.Next() {
			var team models.Team
			if err := rows.Scan(&team.ID, &team.Name, &team.CreatedAt, &team.UpdatedAt); err != nil {
				logger.Error(logLevel, logDetailTeam, fmt.Sprintf("Failed to scan team: %v", err))
				return nil, fmt.Errorf("failed to scan team: %w", err)
			}
			teams = append(teams, team)
		}
		return teams, nil
	})
}

func (dal *TeamDAL) AddTeamMember(ctx context.Context, member *models.TeamMember) error {
	// Add new team member, using upsert to check for conflicts
	query := fmt.Sprintf(`
		INSERT INTO %s.team_members(team_id, user_id)
		VALUES($1, $2)
		ON CONFLICT (team_id, user_id) DO NOTHING`, dal.schema)
	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		_, err := stmt.ExecContext(ctx, member.TeamID, member.UserID)
		if err != nil {
			logger.Error(logLevel, logDetailTeam, fmt.Sprintf("Failed to add team member: %v", err))
			return fmt.Errorf("failed to add team member: %w", err)
		}
		return nil
	})
}

func (dal *TeamDAL) RemoveTeamMember(ctx context.Context, memberID string) error {
	query := fmt.Sprintf(`DELETE FROM %s.team_members WHERE id=$1`, dal.schema)
	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		_, err := stmt.ExecContext(ctx, memberID)
		if err != nil {
			logger.Error(logLevel, logDetailTeam, fmt.Sprintf("Failed to remove team member: %v", err))
			return fmt.Errorf("failed to remove team member: %w", err)
		}
		return nil
	})
}

func (dal *TeamDAL) GetTeamMembers(ctx context.Context, teamID string) ([]models.TeamMember, error) {
	query := fmt.Sprintf(`
        SELECT tm.id, tm.team_id, tm.user_id, u.username, u.role, tm.added_at
        FROM %s.team_members tm
        JOIN %s.users u ON tm.user_id = u.id
        WHERE tm.team_id = $1`, dal.schema, dal.schema)
	return utils.WithResultTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) ([]models.TeamMember, error) {
		rows, err := stmt.QueryContext(ctx, teamID)
		if err != nil {
			logger.Error(logLevel, logDetailTeam, fmt.Sprintf("Failed to query team members: %v", err))
			return nil, fmt.Errorf("failed to query team members: %w", err)
		}
		defer rows.Close()

		var members []models.TeamMember
		for rows.Next() {
			var member models.TeamMember
			if err := rows.Scan(&member.ID, &member.TeamID, &member.UserID, &member.Username, &member.Role, &member.AddedAt); err != nil {
				logger.Error(logLevel, logDetailTeam, fmt.Sprintf("Failed to scan team member: %v", err))
				return nil, fmt.Errorf("failed to scan team member: %w", err)
			}
			members = append(members, member)
		}

		if err := rows.Err(); err != nil {
			logger.Error(logLevel, logDetailTeam, fmt.Sprintf("Row iteration error: %v", err))
			return nil, fmt.Errorf("row iteration error: %w", err)
		}

		return members, nil
	})
}
