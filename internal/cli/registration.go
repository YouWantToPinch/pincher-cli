package cli

func makeBaseCommandHandlers() []*cmdHandler {
	mdAct := middlewareValidateAction

	handlers := []*cmdHandler{
		{
			cmdElement: cmdElement{
				name:        "exit",
				description: "exit the program",
				priority:    0,
			},
			callback: nil,
		},
		{
			cmdElement: cmdElement{
				name:        "help",
				description: "See usage of another command.",
				parameters:  []string{"command"},
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
				parameters:  []string{"action"},
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
				name:        "ready",
				description: "Get server readiness",
				priority:    20,
			},
			callback: handlerReady,
		},
		{
			cmdElement: cmdElement{
				name:        "user",
				description: "Create a new user, or log in",
				parameters:  []string{"action"},
				priority:    50,
			},
			callback: mdAct(handlerUser),
			actions: []cmdElement{
				{
					name:        "add",
					description: "create a new user",
					parameters:  []string{"new_username", "new_password", "retype_password"},
				},
				{
					name:        "login",
					description: "log in as an existing user",
					parameters:  []string{"username", "password"},
					options: []cmdElement{
						{
							name:         "view-budget",
							description:  "specify a budget to view on successful login",
							parameters:   []string{"budget_name"},
							useShorthand: true,
						},
					},
				},
				{
					name:        "update",
					description: "update credentials of the logged-in user",
					parameters:  []string{"username", "password"},
					options: []cmdElement{
						{
							name:        "username",
							description: "set a new username for the user",
							parameters:  []string{"new_value"},
						},
						{
							name:        "password",
							description: "set a new password for the user",
							parameters:  []string{"new_value", "retyped_value"},
						},
					},
				},
				{
					name:        "logout",
					description: "log out existing user",
					parameters:  []string{},
				},
				{
					name:        "delete",
					description: "delete the logged-in user by first entering its credentials",
					parameters:  []string{"username", "password", "retype_password"},
				},
			},
		},
	}

	return handlers
}

func makeBudgetCommandHandler() *cmdHandler {
	handler := &cmdHandler{
		cmdElement: cmdElement{
			name:        "budget",
			description: "Manage budgets associated with logged-in user",
			parameters:  []string{"action"},
			priority:    100,
		},
		nonRegMsg: "login required",
		callback:  middlewareValidateAction(handlerBudget),
		actions: []cmdElement{
			{
				name:       "add",
				parameters: []string{"name"},
				options: []cmdElement{
					{
						name:        "notes",
						description: "Give your budget some notes",
						parameters:  []string{"notes_value"},
					},
				},
			},
			{
				name: "list",
				options: []cmdElement{
					{
						name:        "roles",
						description: "Filter results by user role. Can be ADMIN, MANAGER, CONTRIBUTOR, or VIEWER.",
						parameters:  []string{"role_title"},
					},
				},
			},
			{
				name:        "view",
				description: "specify a budget to interact with using other commands",
				parameters:  []string{"budget_name"},
			},
			{
				name:        "update",
				description: "update budget information, IE name, notes",
				parameters:  []string{"name"},
				options: []cmdElement{
					{
						name:        "name",
						description: "Update budget name",
						parameters:  []string{"name_value"},
					},
					{
						name:        "notes",
						description: "Update budget name",
						parameters:  []string{"notes_value"},
					},
				},
			},
			{
				name:        "delete",
				description: "delete an existing budget by name",
				parameters:  []string{"budget_name"},
			},
		},
	}

	return handler
}

func makeResourceCommandHandlers() []*cmdHandler {
	mdAct := middlewareValidateAction

	handlers := []*cmdHandler{
		{
			cmdElement: cmdElement{
				name:        "account",
				parameters:  []string{"action"},
				description: "Manage accounts under budget in view",
				priority:    210,
			},
			nonRegMsg: "first view a budget to see its accounts",
			callback:  mdAct(handlerAccount),
			actions: []cmdElement{
				{
					name:        "add",
					description: "Add a new account to budget",
					parameters:  []string{"name"},
					options: []cmdElement{
						{
							name:        "notes",
							description: "give the new account some notes",
							parameters:  []string{"notes_value"},
						},
						{
							name:         "off-budget",
							description:  "create the account for tracking purposes only, seperating it from any categorization",
							useShorthand: true,
						},
					},
				},
				{
					name:        "update",
					description: "update information on account by name",
					parameters:  []string{"name"},
					options: []cmdElement{
						{
							name:        "name",
							description: "rewrite account name",
							parameters:  []string{"new_name"},
						},
						{
							name:        "notes",
							description: "rewrite account notes",
							parameters:  []string{"new_notes"},
						},
						{
							name:        "type",
							description: "choose different account type",
							parameters:  []string{"new_type"},
						},
					},
				},
				{
					name:        "restore",
					description: "restore a soft-deleted account",
					parameters:  []string{"account_name"},
				},
				{
					name:        "list",
					description: "see a list of all accounts belonging to budget",
					options: []cmdElement{
						{
							name:         "deleted",
							description:  "view only soft-deleted accounts",
							useShorthand: true,
						},
					},
				},
				{
					name:        "delete",
					description: "Delete an account",
					parameters:  []string{"account_name"},
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
				name:        "category",
				parameters:  []string{"action"},
				description: "Manage spending categories under budget in view",
				priority:    220,
			},
			nonRegMsg: "first view a budget to see its categories",
			callback:  mdAct(handlerCategory),
			actions: []cmdElement{
				{
					name:        "add",
					description: "Add a new category to budget",
					parameters:  []string{"name"},
					options: []cmdElement{
						{
							name:        "notes",
							description: "give the new category some notes",
							parameters:  []string{"notes_value"},
						},
						{
							name:         "group",
							description:  "assign the category to a group",
							parameters:   []string{"group_name"},
							useShorthand: true,
						},
					},
				},
				{
					name:        "update",
					description: "update information on a category by name",
					parameters:  []string{"name"},
					options: []cmdElement{
						{
							name:        "name",
							description: "rewrite category name",
							parameters:  []string{"new_name"},
						},
						{
							name:        "notes",
							description: "rewrite category notes",
							parameters:  []string{"new_notes"},
						},
						{
							name:         "group",
							description:  "assign the category to a group",
							parameters:   []string{"group_name"},
							useShorthand: true,
						},
					},
				},
				{
					name:        "list",
					description: "list all categories belonging to budget",
					options: []cmdElement{
						{
							name:         "group",
							description:  "list only categories grouped by given name",
							parameters:   []string{"group_name"},
							useShorthand: true,
						},
					},
				},
				{
					name:        "delete",
					description: "Delete a category",
					parameters:  []string{"category_name"},
				},
			},
		},
		{
			cmdElement: cmdElement{
				name:        "group",
				parameters:  []string{"action"},
				description: "Manage category groups under budget in view",
				priority:    230,
			},
			nonRegMsg: "first view a budget to see its groups",
			callback:  mdAct(handlerGroup),
			actions: []cmdElement{
				{
					name:        "add",
					description: "Add a new group to budget",
					parameters:  []string{"name"},
					options: []cmdElement{
						{
							name:        "notes",
							description: "give the new group some notes",
							parameters:  []string{"notes_value"},
						},
					},
				},
				{
					name:        "update",
					description: "update information on a group by name",
					parameters:  []string{"name"},
					options: []cmdElement{
						{
							name:        "name",
							description: "rewrite group name",
							parameters:  []string{"new_name"},
						},
						{
							name:        "notes",
							description: "rewrite group notes",
							parameters:  []string{"new_notes"},
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
					parameters:  []string{"group_name"},
				},
			},
		},
		{
			cmdElement: cmdElement{
				name:        "payee",
				parameters:  []string{"action"},
				description: "Manage payees under budget in view",
				priority:    240,
			},
			nonRegMsg: "first view a budget to see its payees",
			callback:  mdAct(handlerPayee),
			actions: []cmdElement{
				{
					name:        "add",
					description: "Add a new payee to budget",
					parameters:  []string{"name"},
					options: []cmdElement{
						{
							name:        "notes",
							description: "give the new payee some notes",
							parameters:  []string{"notes_value"},
						},
					},
				},
				{
					name:        "update",
					description: "update information on a payee by name",
					parameters:  []string{"name"},
					options: []cmdElement{
						{
							name:        "name",
							description: "rewrite payee name",
							parameters:  []string{"new_name"},
						},
						{
							name:        "notes",
							description: "rewrite payee notes",
							parameters:  []string{"new_notes"},
						},
					},
				},
				{
					name:        "list",
					description: "see a list of all payees belonging to budget",
				},
				{
					name:        "delete",
					description: "Delete a payee",
					parameters:  []string{"payee_name"},
					options: []cmdElement{
						{
							name:         "replacement",
							description:  "name of a payee to replace payee to delete for where it is still in use",
							parameters:   []string{"new_payee_name"},
							useShorthand: true,
						},
					},
				},
			},
		},
		{
			cmdElement: cmdElement{
				name:        "txn",
				parameters:  []string{"action"},
				description: "Manage category groups under budget in view",
				priority:    230,
			},
			nonRegMsg: "first view a budget to manage transactions within its accounts",
			callback:  mdAct(handlerTxn),
			actions: []cmdElement{
				{
					name:        "list",
					description: "see a list of transactions",
					options: []cmdElement{
						{
							name:         "account",
							description:  "filter by account",
							parameters:   []string{"account_name"},
							useShorthand: true,
						},
						{
							name:         "payee",
							description:  "filter by payee",
							parameters:   []string{"payee_name"},
							useShorthand: true,
						},
						{
							name:         "category",
							description:  "filter by category",
							parameters:   []string{"category=category_name"},
							useShorthand: true,
						},
						{
							name:         "dates",
							description:  "filter by time frame",
							parameters:   []string{"start_date", "end_date"},
							useShorthand: true,
						},
					},
				},
				{
					name:        "log",
					description: "log a deposit or withdrawal transaction to the budget in view.",
					parameters:  []string{"account", "payee", "amount", "category"},
					options: []cmdElement{
						{
							name:         "date",
							description:  "specify a date date for this transaction (defaults to present day)",
							parameters:   []string{"date"},
							useShorthand: true,
						},
						{
							name:         "notes",
							description:  "give the transaction some notes",
							parameters:   []string{"new_notes"},
							useShorthand: true,
						},
						{
							name:         "cleared",
							description:  "whether or not to represent this transaction as complete (false by default)",
							useShorthand: true,
						},
						{
							name:         "split",
							description:  "split this transaction into 2+ categories(writing 'split' for the main category argument). Let their amounts total passed to the amount argument.",
							parameters:   []string{"category=amount,..."},
							useShorthand: true,
						},
					},
				},
				{
					name:        "transfer",
					description: "log a transfer transaction between two accounts within the budget in view",
					parameters:  []string{"from_account", "to_account", "amount"},
					options: []cmdElement{
						{
							name:         "date",
							description:  "specify a date date for this transfer (defaults to present day)",
							parameters:   []string{"date"},
							useShorthand: true,
						},
						{
							name:         "notes",
							description:  "give the transfer some notes",
							parameters:   []string{"notes_value"},
							useShorthand: true,
						},
						{
							name:         "cleared",
							description:  "whether or not to represent this transfer as complete (false by default)",
							useShorthand: true,
						},
					},
				},
			},
		},
	}

	return handlers
}
