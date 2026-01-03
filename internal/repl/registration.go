package repl

func batchRegisterCommands(r *commandRegistry, handlers []*cmdHandler, preregister bool) {
	for _, handler := range handlers {
		r.register(handler.name, handler, preregister)
	}
}

func registerBaseCommands(s *State, preregister bool) {
	mdAct := middlewareValidateAction

	handlers := []*cmdHandler{
		{
			cmdElement: cmdElement{
				name:        "exit",
				description: "exit the program",
				priority:    0,
			},
			callback: handlerExit,
		},
		{
			cmdElement: cmdElement{
				name:        "help",
				description: "See usage of another command.",
				arguments:   []string{"command"},
				priority:    1,
				options: []cmdElement{
					{
						name:         "verbose",
						description:  "show unregistered commands (those not available for use in the current CLI context)",
						useShorthand: true,
					},
				},
			},
			callback: handlerHelp,
		},
		{
			cmdElement: cmdElement{
				name:        "clear",
				description: "clear the terminal",
				priority:    2,
			},
			callback: handlerClear,
		},
		{
			cmdElement: cmdElement{
				name:        "config",
				description: "Add, Load, or Save a local user configuration for the Pincher-CLI",
				arguments:   []string{"action"},
				priority:    10,
			},
			actions: []cmdElement{
				{
					name:        "edit",
					description: "edit current user configuration",
				},
				{
					name:        "load",
					description: "load user configuration from the local machine",
				},
			},
			callback: mdAct(handlerConfig),
		},
		{
			cmdElement: cmdElement{
				name:        "log",
				description: "see Pincher-CLI logs",
				priority:    15,
			},
			callback: handlerLog,
		},
		{
			cmdElement: cmdElement{
				name:        "ready",
				description: "Get server readiness",
				priority:    25,
			},
			callback: handlerReady,
		},
		{
			cmdElement: cmdElement{
				name:        "user",
				description: "Create a new user, or log in",
				arguments:   []string{"action"},
				priority:    50,
			},
			callback: mdAct(handlerUser),
			actions: []cmdElement{
				{
					name:        "add",
					description: "create a new user",
					arguments:   []string{"new_username", "new_password", "retype_password"},
				},
				{
					name:        "login",
					description: "log in as an existing user",
					arguments:   []string{"username", "password"},
				},
				{
					name:        "update",
					description: "update credentials of the logged-in user",
					arguments:   []string{"username", "password"},
					options: []cmdElement{
						{
							name:        "username",
							description: "set a new username for the user",
							arguments:   []string{"new_value"},
						},
						{
							name:        "password",
							description: "set a new password for the user",
							arguments:   []string{"new_value", "retyped_value"},
						},
					},
				},
				{
					name:        "delete",
					description: "delete the logged-in user by first entering its credentials",
					arguments:   []string{"username", "password", "retype_password"},
				},
			},
		},
	}

	batchRegisterCommands(s.CommandRegistry, handlers, preregister)
}

func registerBudgetCommand(s *State, preregister bool) {
	handler := &cmdHandler{
		cmdElement: cmdElement{
			name:        "budget",
			description: "Manage budgets associated with user '" + s.Client.LoggedInUser.Username + "'",
			arguments:   []string{"action"},
			priority:    190,
		},
		nonRegMsg: "login required",
		callback:  middlewareValidateAction(handlerBudget),
		actions: []cmdElement{
			{
				name:      "add",
				arguments: []string{"name"},
				options: []cmdElement{
					{
						name:        "notes",
						description: "Give your budget some notes",
						arguments:   []string{"notes_value"},
					},
				},
			},
			{
				name: "list",
				options: []cmdElement{
					{
						name:        "role",
						description: "Filter results by user role. Can be ADMIN, MANAGER, CONTRIBUTOR, or VIEWER.",
						arguments:   []string{"role_title"},
					},
				},
			},
			{
				name:        "view",
				description: "specify a budget to interact with using other commands",
				arguments:   []string{"budget_name"},
			},
			{
				name:        "update",
				description: "update budget information, IE name, notes",
				arguments:   []string{"name"},
				options: []cmdElement{
					{
						name:        "name",
						description: "Update budget name",
						arguments:   []string{"name_value"},
					},
					{
						name:        "notes",
						description: "Update budget name",
						arguments:   []string{"notes_value"},
					},
				},
			},
			{
				name:        "delete",
				description: "delete an existing budget by name",
				arguments:   []string{"budget_name"},
			},
		},
	}

	s.CommandRegistry.register("budget", handler, preregister)
}

func registerResourceCommands(s *State, preregister bool) {
	mdAct := middlewareValidateAction

	handlers := []*cmdHandler{
		{
			cmdElement: cmdElement{
				name:        "account",
				arguments:   []string{"action"},
				description: "Manage accounts under budget '" + s.Client.ViewedBudget.Name + "'",
				priority:    180,
			},
			nonRegMsg: "first view a budget to see its accounts",
			callback:  mdAct(handlerAccount),
			actions: []cmdElement{
				{
					name:        "add",
					description: "Add a new account to budget",
					arguments:   []string{"name", "account_type"},
					options: []cmdElement{
						{
							name:        "notes",
							description: "give the new account some notes",
							arguments:   []string{"notes_value"},
						},
					},
				},
				{
					name:        "update",
					description: "update information on account by name",
					arguments:   []string{"name"},
					options: []cmdElement{
						{
							name:        "name",
							description: "rewrite account name",
							arguments:   []string{"new_name"},
						},
						{
							name:        "notes",
							description: "rewrite account notes",
							arguments:   []string{"new_notes"},
						},
						{
							name:        "type",
							description: "choose different account type",
							arguments:   []string{"new_type"},
						},
					},
				},
				{
					name:        "restore",
					description: "restore a soft-deleted account",
					arguments:   []string{"account_name"},
				},
				{
					name:        "list",
					description: "see a list of all accounts belonging to budget",
					options: []cmdElement{
						{
							name:         "include",
							description:  "Include accounts usually excluded with qualities like: 'deleted'",
							arguments:    []string{"quality"},
							useShorthand: true,
						},
					},
				},
				{
					name:        "delete",
					description: "Delete an account",
					arguments:   []string{"account_name"},
					options: []cmdElement{
						{
							name:         "hard",
							description:  "as opposed to a reversible soft deletion (default)",
							useShorthand: true,
						},
					},
				},
			},
		},
		{
			cmdElement: cmdElement{
				name:        "group",
				arguments:   []string{"action"},
				description: "Manage groups under budget '" + s.Client.ViewedBudget.Name + "'",
				priority:    180,
			},
			nonRegMsg: "first view a budget to see its groups",
			callback:  mdAct(handlerGroup),
			actions: []cmdElement{
				{
					name:        "add",
					description: "Add a new group to budget",
					arguments:   []string{"name"},
					options: []cmdElement{
						{
							name:        "notes",
							description: "give the new group some notes",
							arguments:   []string{"notes_value"},
						},
					},
				},
				{
					name:        "update",
					description: "update information on a group by name",
					arguments:   []string{"name"},
					options: []cmdElement{
						{
							name:        "name",
							description: "rewrite group name",
							arguments:   []string{"new_name"},
						},
						{
							name:        "notes",
							description: "rewrite group notes",
							arguments:   []string{"new_notes"},
						},
					},
				},
				{
					name:        "list",
					description: "see a list of all groups belonging to budget",
				},
				{
					name:        "delete",
					description: "Delete a group",
					arguments:   []string{"group_name"},
				},
			},
		},
	}

	batchRegisterCommands(s.CommandRegistry, handlers, preregister)
}
