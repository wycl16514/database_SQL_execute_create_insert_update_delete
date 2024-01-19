module simple_db

replace file_manager => ./file_manager

replace log_manager => ./log_manager

replace buffer_manager => ./buffer_manager

replace tx => ./tx

replace record_manager => ./record_manager

replace metadata_management => ./metadata_manager

replace lexer => ./lexer

replace parser => ./parser

replace query => ./query

replace planner => ./planner

go 1.17

require query v0.0.0-00010101000000-000000000000

require (
	buffer_manager v0.0.0-00010101000000-000000000000
	file_manager v0.0.0-00010101000000-000000000000
	log_manager v0.0.0-00010101000000-000000000000
	parser v0.0.0-00010101000000-000000000000
	planner v0.0.0-00010101000000-000000000000
	record_manager v0.0.0-00010101000000-000000000000
	tx v0.0.0-00010101000000-000000000000
)

require (
	lexer v0.0.0-00010101000000-000000000000 // indirect
	metadata_management v0.0.0-00010101000000-000000000000
)
