// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// ChatsColumns holds the columns for the "chats" table.
	ChatsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "client_id", Type: field.TypeUUID, Unique: true},
		{Name: "created_at", Type: field.TypeTime},
	}
	// ChatsTable holds the schema information for the "chats" table.
	ChatsTable = &schema.Table{
		Name:       "chats",
		Columns:    ChatsColumns,
		PrimaryKey: []*schema.Column{ChatsColumns[0]},
	}
	// FailedJobsColumns holds the columns for the "failed_jobs" table.
	FailedJobsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "name", Type: field.TypeString, Size: 2147483647},
		{Name: "payload", Type: field.TypeString, Size: 2147483647},
		{Name: "reason", Type: field.TypeString, Size: 2147483647},
		{Name: "created_at", Type: field.TypeTime},
	}
	// FailedJobsTable holds the schema information for the "failed_jobs" table.
	FailedJobsTable = &schema.Table{
		Name:       "failed_jobs",
		Columns:    FailedJobsColumns,
		PrimaryKey: []*schema.Column{FailedJobsColumns[0]},
	}
	// JobsColumns holds the columns for the "jobs" table.
	JobsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "name", Type: field.TypeString, Size: 2147483647},
		{Name: "payload", Type: field.TypeString, Size: 2147483647},
		{Name: "attempts", Type: field.TypeInt, Default: 0},
		{Name: "available_at", Type: field.TypeTime},
		{Name: "reserved_until", Type: field.TypeTime},
		{Name: "created_at", Type: field.TypeTime},
	}
	// JobsTable holds the schema information for the "jobs" table.
	JobsTable = &schema.Table{
		Name:       "jobs",
		Columns:    JobsColumns,
		PrimaryKey: []*schema.Column{JobsColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "job_available_at_reserved_until",
				Unique:  false,
				Columns: []*schema.Column{JobsColumns[4], JobsColumns[5]},
			},
		},
	}
	// MessagesColumns holds the columns for the "messages" table.
	MessagesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "author_id", Type: field.TypeUUID, Nullable: true},
		{Name: "initial_request_id", Type: field.TypeUUID},
		{Name: "is_visible_for_client", Type: field.TypeBool, Default: false},
		{Name: "is_visible_for_manager", Type: field.TypeBool, Default: false},
		{Name: "body", Type: field.TypeString, Size: 3000},
		{Name: "checked_at", Type: field.TypeTime, Nullable: true},
		{Name: "is_blocked", Type: field.TypeBool, Default: false},
		{Name: "is_service", Type: field.TypeBool, Default: false},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "chat_id", Type: field.TypeUUID},
		{Name: "problem_id", Type: field.TypeUUID},
	}
	// MessagesTable holds the schema information for the "messages" table.
	MessagesTable = &schema.Table{
		Name:       "messages",
		Columns:    MessagesColumns,
		PrimaryKey: []*schema.Column{MessagesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "messages_chats_messages",
				Columns:    []*schema.Column{MessagesColumns[10]},
				RefColumns: []*schema.Column{ChatsColumns[0]},
				OnDelete:   schema.NoAction,
			},
			{
				Symbol:     "messages_problems_messages",
				Columns:    []*schema.Column{MessagesColumns[11]},
				RefColumns: []*schema.Column{ProblemsColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "message_chat_id",
				Unique:  false,
				Columns: []*schema.Column{MessagesColumns[10]},
			},
			{
				Name:    "message_created_at_is_visible_for_client",
				Unique:  false,
				Columns: []*schema.Column{MessagesColumns[9], MessagesColumns[3]},
			},
			{
				Name:    "message_initial_request_id",
				Unique:  true,
				Columns: []*schema.Column{MessagesColumns[2]},
				Annotation: &entsql.IndexAnnotation{
					Where: "NOT is_service",
				},
			},
		},
	}
	// ProblemsColumns holds the columns for the "problems" table.
	ProblemsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "manager_id", Type: field.TypeUUID, Nullable: true},
		{Name: "resolved_at", Type: field.TypeTime, Nullable: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "chat_id", Type: field.TypeUUID},
	}
	// ProblemsTable holds the schema information for the "problems" table.
	ProblemsTable = &schema.Table{
		Name:       "problems",
		Columns:    ProblemsColumns,
		PrimaryKey: []*schema.Column{ProblemsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "problems_chats_problems",
				Columns:    []*schema.Column{ProblemsColumns[4]},
				RefColumns: []*schema.Column{ChatsColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "problem_id_manager_id",
				Unique:  true,
				Columns: []*schema.Column{ProblemsColumns[0], ProblemsColumns[1]},
				Annotation: &entsql.IndexAnnotation{
					Where: "resolved_at IS NULL",
				},
			},
			{
				Name:    "problem_chat_id",
				Unique:  false,
				Columns: []*schema.Column{ProblemsColumns[4]},
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		ChatsTable,
		FailedJobsTable,
		JobsTable,
		MessagesTable,
		ProblemsTable,
	}
)

func init() {
	MessagesTable.ForeignKeys[0].RefTable = ChatsTable
	MessagesTable.ForeignKeys[1].RefTable = ProblemsTable
	ProblemsTable.ForeignKeys[0].RefTable = ChatsTable
}
