package migrate_db

import (
	"database/sql"
	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/forbole/juno/v3/logging"
	"github.com/forbole/juno/v3/types/config"
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
	"gitlab.com/rarimo/bdjuno/v3/database"
)

var migrations = &migrate.EmbedFileSystemMigrationSource{
	FileSystem: database.Migrations,
	Root:       "schema",
}

func NewMigrateDBCmd(cmdCfg *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "migrate-db",
		Short:   "Migrate the database schema",
		PreRunE: parsecmdtypes.ReadConfigPreRunE(cmdCfg),
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "up",
			Short: "migrate db up",
			RunE: func(cmd *cobra.Command, args []string) error {
				context, err := parsecmdtypes.GetParserContext(config.Cfg, cmdCfg)
				if err != nil {
					return err
				}
				db := database.Cast(context.Database)
				return migrateUp(db.Sql, context.Logger)
			},
		},
	)

	cmd.AddCommand(
		&cobra.Command{
			Use:   "down",
			Short: "migrate db down",
			RunE: func(cmd *cobra.Command, args []string) error {
				context, err := parsecmdtypes.GetParserContext(config.Cfg, cmdCfg)
				if err != nil {
					return err
				}
				db := database.Cast(context.Database)
				return migrateDown(db.Sql, context.Logger)
			},
		},
	)

	return cmd
}

func migrateUp(rawDB *sql.DB, log logging.Logger) error {
	applied, err := migrate.Exec(rawDB, "postgres", migrations, migrate.Up)
	if err != nil {
		return errors.Wrap(err, "failed to apply migrations")
	}
	log.Info("migrations applied", map[string]interface{}{
		"applied": applied,
	})
	return nil
}

func migrateDown(rawDB *sql.DB, log logging.Logger) error {
	applied, err := migrate.Exec(rawDB, "postgres", migrations, migrate.Down)
	if err != nil {
		return errors.Wrap(err, "failed to apply migrations")
	}
	log.Info("migrations applied", map[string]interface{}{
		"applied": applied,
	})
	return nil
}
