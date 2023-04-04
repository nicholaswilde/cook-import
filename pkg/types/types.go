package types

type Config struct {
	DryRun						bool		`mapstructure:"dry-run"`
	Jsonify						bool		`mapstructure:"jsonify"`
	IgnoreFile				string		`mapstructure:"ignore-file"`
	RecipeSearchRoot	string		`mapstructure:"recipe-search-root"`
	LogLevel					string		`mapstructure:"log-level"`
	TemplateFiles			[]string	`mapstructure:"template-files"`
	WordWrap					int			`mapstructure:"word-wrap"`
}
