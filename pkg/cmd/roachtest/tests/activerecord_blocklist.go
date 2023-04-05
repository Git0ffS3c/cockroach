// Copyright 2020 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package tests

// These are lists of known activerecord test errors and failures.
// When the activerecord test suite is run, the results are compared to this list.
// Any passed test that is not on this list is reported as PASS - expected
// Any passed test that is on this list is reported as PASS - unexpected
// Any failed test that is on this list is reported as FAIL - expected
// Any failed test that is not on this list is reported as FAIL - unexpected
// Any test on this list that is not run is reported as FAIL - not run
//
// Please keep these lists alphabetized for easy diffing.
// After a failed run, an updated version of this blocklist should be available
// in the test log.
var activeRecordBlocklist = blocklist{
	"ActiveRecord::CockroachDB::Migration::PGChangeSchemaTest#test_change_string_to_date":                                               "49329",
	"ActiveRecord::CockroachDB::Migration::PGChangeSchemaTest#test_change_type_with_array":                                              "49329",
	"ActiveRecord::CockroachDB::Migration::PGChangeSchemaTest#test_change_type_with_symbol":                                             "49329",
	"ActiveRecord::CockroachDB::Migration::PGChangeSchemaTest#test_change_type_with_symbol_using_datetime":                              "49329",
	"ActiveRecord::CockroachDB::Migration::PGChangeSchemaTest#test_change_type_with_symbol_using_datetime_with_timestamptz_as_default":  "49329",
	"ActiveRecord::CockroachDB::Migration::PGChangeSchemaTest#test_change_type_with_symbol_using_timestamp_with_timestamptz_as_default": "49329",
	"ActiveRecord::CockroachDB::Migration::PGChangeSchemaTest#test_change_type_with_symbol_with_timestamptz_as_default":                 "49329",
	"ActiveRecord::CockroachDB::Migration::PGChangeSchemaTest#test_change_type_with_symbol_with_timestamptz":                            "49329",
}

var activeRecordIgnoreList = blocklist{
	"ActiveRecord::Coders::YAMLColumnTestWithSafeLoad#test_returns_string_unless_starts_with_dash":                    "flaky - sometimes attempts to call a method on nil",
	"ActiveRecord::ConnectionAdapters::ConnectionSwappingNestedTest#test_prevent_writes_can_be_changed_granularly":    "flaky - sometimes attempts to call a method on nil",
	"ActiveRecord::InvertibleMigrationTest#test_down":                                                                 "flaky - sometimes complains that a relation already exists",
	"ActiveRecord::Migration::ForeignKeyChangeColumnTest#test_change_column_of_parent_table":                          "flaky - sometimes complains a file is not open for writing",
	"ActiveRecord::Migration::ForeignKeyChangeColumnTest#test_remove_foreign_key_by_column":                           "flaky - sometimes complains a file is not open for writing",
	"ActiveRecord::Migration::ForeignKeyChangeColumnTest#test_remove_foreign_key_by_column_in_change_table":           "flaky - sometimes complains a file is not open for writing",
	"ActiveRecord::Migration::ForeignKeyChangeColumnTest#test_remove_reference_column_of_child_table":                 "flaky - sometimes complains a file is not open for writing",
	"ActiveRecord::Migration::ForeignKeyChangeColumnTest#test_rename_column_of_child_table":                           "flaky - sometimes complains a file is not open for writing",
	"ActiveRecord::Migration::ForeignKeyChangeColumnTest#test_rename_reference_column_of_child_table":                 "flaky - sometimes complains a file is not open for writing",
	"ActiveRecord::Migration::ForeignKeyChangeColumnWithPrefixTest#test_change_column_of_parent_table":                "flaky - sometimes complains a file is not open for writing",
	"ActiveRecord::Migration::ForeignKeyChangeColumnWithPrefixTest#test_remove_foreign_key_by_column":                 "flaky - sometimes complains a file is not open for writing",
	"ActiveRecord::Migration::ForeignKeyChangeColumnWithPrefixTest#test_remove_foreign_key_by_column_in_change_table": "flaky - sometimes complains a file is not open for writing",
	"ActiveRecord::Migration::ForeignKeyChangeColumnWithPrefixTest#test_remove_reference_column_of_child_table":       "flaky - sometimes complains a file is not open for writing",
	"ActiveRecord::Migration::ForeignKeyChangeColumnWithPrefixTest#test_rename_column_of_child_table":                 "flaky - sometimes complains a file is not open for writing",
	"ActiveRecord::Migration::ForeignKeyChangeColumnWithPrefixTest#test_rename_reference_column_of_child_table":       "flaky - sometimes complains a file is not open for writing",
	"ActiveRecord::Migration::ForeignKeyChangeColumnWithSuffixTest#test_change_column_of_parent_table":                "flaky - sometimes complains a file is not open for writing",
	"ActiveRecord::Migration::ForeignKeyChangeColumnWithSuffixTest#test_remove_foreign_key_by_column":                 "flaky - sometimes complains a file is not open for writing",
	"ActiveRecord::Migration::ForeignKeyChangeColumnWithSuffixTest#test_remove_foreign_key_by_column_in_change_table": "flaky - sometimes complains a file is not open for writing",
	"ActiveRecord::Migration::ForeignKeyChangeColumnWithSuffixTest#test_remove_reference_column_of_child_table":       "flaky - sometimes complains a file is not open for writing",
	"ActiveRecord::Migration::ForeignKeyChangeColumnWithSuffixTest#test_rename_column_of_child_table":                 "flaky - sometimes complains a file is not open for writing",
	"ActiveRecord::Migration::ForeignKeyChangeColumnWithSuffixTest#test_rename_reference_column_of_child_table":       "flaky - sometimes complains a file is not open for writing",
	"ActiveRecord::PostgresqlConnectionTest#test_table_exists_logs_name":                                              "flaky - sometimes attempts to call a method on nil",
	"ActiveRecord::WhereChainTest#test_chaining_multiple":                                                             "flaky - sometimes complains that a relation does not exist",
	"BasicsTest#test_default_values_are_deeply_dupped":                                                                "flaky - sometimes attempts to call a method on nil",
	"BinaryTest#test_mixed_encoding":                                                                                                                           "flaky - sometimes complains that a relation does not exist",
	"CockroachDB::FixturesTest#test_bulk_insert":                                                                                                               "flaky",
	"CockroachDB::PostgresqlIntervalTest#test_column":                                                                                                          "flaky - sometimes complains that a relation does not exist",
	"CockroachDB::PostgresqlIntervalTest#test_simple_interval":                                                                                                 "flaky - sometimes time format is wrong",
	"CockroachDB::PostgresqlIntervalTest#test_interval_type":                                                                                                   "flaky - sometimes complains that a relation does not exist",
	"ConcurrentTransactionTest#test_transaction_isolation__read_committed":                                                                                     "flaky - https://github.com/cockroachdb/activerecord-cockroachdb-adapter/issues/237",
	"DefaultTextTest#test_default_texts_containing_single_quotes":                                                                                              "flaky - sometimes complains that a relation does not exist",
	"EagerAssociationTest#test_eager_association_loading_with_belongs_to_and_order_string_with_unquoted_table_name":                                            "flaky - sometimes complains that a relation does not exist",
	"EagerAssociationTest#test_eager_association_loading_with_habtm":                                                                                           "flaky - sometimes complains that a relation does not exist",
	"EagerAssociationTest#test_eager_load_has_many_with_string_keys":                                                                                           "flaky - sometimes complains that a relation does not exist",
	"EagerAssociationTest#test_eager_with_has_many_through_with_conditions_join_model_with_include":                                                            "flaky - sometimes complains that a relation does not exist",
	"EagerAssociationTest#test_eager_with_multiple_associations_with_same_table_has_many_and_habtm":                                                            "flaky - sometimes complains that a relation does not exist",
	"EagerAssociationTest#test_has_many_association_ignores_the_scoping":                                                                                       "flaky - sometimes complains that a relation does not exist",
	"EagerAssociationTest#test_nested_loading_does_not_raise_exception_when_association_does_not_exist":                                                        "flaky - sometimes complains that a relation does not exist",
	"EagerAssociationTest#test_nested_loading_through_has_one_association_with_order_on_association":                                                           "flaky - sometimes complains that a relation does not exist",
	"EagerAssociationTest#test_preconfigured_includes_with_habtm":                                                                                              "flaky - sometimes complains that a relation does not exist",
	"EagerAssociationTest#test_preloading_has_many_through_with_custom_scope":                                                                                  "flaky - sometimes complains that a relation does not exist",
	"EagerAssociationTest#test_three_level_nested_preloading_does_not_raise_exception_when_association_does_not_exist":                                         "flaky - sometimes complains that a relation does not exist",
	"FixturesTest#test_create_fixtures":                                                                                                                        "flaky - FK constraint violated sometimes when loading all fixture data",
	"FixturesTest#test_inserts":                                                                                                                                "flaky - sometimes complains that a relation does not exist",
	"FixturesTest#test_nonexistent_fixture_file":                                                                                                               "flaky - sometimes complains that a relation does not exist",
	"FixturesTest#test_subsubdir_file_with_arbitrary_name":                                                                                                     "flaky - sometimes complains that a relation does not exist",
	"FixturesTest#test_yaml_file_with_one_invalid_fixture":                                                                                                     "flaky - sometimes complains that a relation does not exist",
	"FoxyFixturesTest#test_identifies_strings":                                                                                                                 "flaky - sometimes complains that a relation does not exist",
	"FoxyFixturesTest#test_populates_timestamp_column":                                                                                                         "flaky - sometimes complains that a relation does not exist",
	"HasManyAssociationsTest#test_adding_array_and_collection":                                                                                                 "flaky - sometimes complains that a relation does not exist",
	"HasManyAssociationsTest#test_building_the_associated_object_with_an_invalid_type":                                                                         "flaky - sometimes complains that a relation does not exist",
	"HasManyAssociationsTest#test_calling_empty_on_an_association_that_has_not_been_loaded_performs_a_query":                                                   "flaky - sometimes complains that a relation does not exist",
	"HasManyAssociationsTest#test_calling_one_should_count_instead_of_loading_association":                                                                     "flaky - sometimes complains that a relation does not exist",
	"HasManyAssociationsTest#test_clear_collection_should_not_change_updated_at":                                                                               "flaky - sometimes complains that a relation does not exist",
	"HasManyAssociationsTest#test_clearing_an_exclusively_dependent_association_collection":                                                                    "flaky - sometimes complains that a relation does not exist",
	"HasManyAssociationsTest#test_counter_cache_updates_in_memory_after_update_with_inverse_of_disabled":                                                       "flaky - sometimes complains that a relation does not exist",
	"HasManyAssociationsTest#test_delete_on_association_clears_scope":                                                                                          "flaky - sometimes complains that a relation does not exist",
	"HasManyAssociationsTest#test_deleting_before_save":                                                                                                        "flaky - sometimes complains that a relation does not exist",
	"HasManyAssociationsTest#test_dependent_association_respects_optional_conditions_on_delete":                                                                "flaky - sometimes complains that a relation does not exist",
	"HasManyAssociationsTest#test_destroy_all_on_association_clears_scope":                                                                                     "flaky - sometimes complains that a relation does not exist",
	"HasManyAssociationsTest#test_destroying_by_string_id":                                                                                                     "flaky - sometimes complains that a relation does not exist",
	"HasManyAssociationsTest#test_get_ids_for_ordered_association":                                                                                             "flaky - sometimes complains that a relation does not exist",
	"HasManyAssociationsTest#test_ids_reader_cache_not_used_for_size_when_association_is_dirty":                                                                "flaky - sometimes complains that a relation does not exist",
	"HasManyAssociationsTest#test_in_memory_replacement_executes_no_queries":                                                                                   "flaky - sometimes complains that a relation does not exist",
	"HasManyAssociationsTest#test_included_in_collection":                                                                                                      "flaky - sometimes complains that a relation does not exist",
	"HasManyAssociationsTest#test_taking_with_a_number":                                                                                                        "flaky - sometimes complains that a relation does not exist",
	"HasOneThroughAssociationsTest#test_has_one_through_do_not_cache_association_reader_if_the_though_method_has_default_scopes":                               "flaky - sometimes complains that a relation does not exist",
	"HasOneThroughAssociationsTest#test_save_of_record_with_loaded_has_one_through":                                                                            "flaky - sometimes complains that a relation does not exist",
	"HasManyAssociationsTest#test_prevent_double_firing_the_before_save_callback_of_new_object_when_the_parent_association_saved_in_the_callback":              "flaky - sometimes complains that a relation does not exist",
	"HasManyAssociationsTest#test_transaction_when_deleting_persisted":                                                                                         "flaky - sometimes complains that a relation does not exist",
	"HasManyThroughAssociationsTest#test_add_two_instance_and_then_deleting":                                                                                   "flaky - sometimes complains that a relation does not exist",
	"HasManyThroughAssociationsTest#test_associate_existing_record_twice_should_add_to_target_twice":                                                           "flaky - sometimes complains that a relation does not exist",
	"HasManyThroughAssociationsTest#test_belongs_to_with_dependent_delete_all":                                                                                 "flaky - sometimes complains that a relation does not exist",
	"HasManyThroughAssociationsTest#test_build_a_model_from_hm_through_association_with_where_clause":                                                          "flaky - sometimes complains that a relation does not exist",
	"HasManyThroughAssociationsTest#test_delete_through_belongs_to_with_dependent_delete_all":                                                                  "flaky - sometimes complains that a relation does not exist",
	"HasManyThroughAssociationsTest#test_destroy_association":                                                                                                  "flaky - sometimes complains that a relation does not exist",
	"HasManyThroughAssociationsTest#test_has_many_through_polymorphic_with_primary_key_option":                                                                 "flaky - sometimes complains that a relation does not exist",
	"HasManyThroughAssociationsTest#test_has_many_through_with_source_scope":                                                                                   "flaky - sometimes complains that a relation does not exist",
	"HasManyThroughAssociationsTest#test_replace_by_id_order_is_preserved":                                                                                     "flaky - sometimes complains that a relation does not exist",
	"HasManyThroughAssociationsTest#test_update_counter_caches_on_delete_with_dependent_destroy":                                                               "flaky - sometimes complains that a relation does not exist",
	"HasManyThroughAssociationsTest#test_update_counter_caches_on_delete_with_dependent_nullify":                                                               "flaky - sometimes complains that a relation does not exist",
	"HotCompatibilityTest#test_update_after_remove_column":                                                                                                     "flaky - sometimes attempts to call a method on nil",
	"InverseBelongsToTests#test_child_instance_should_be_shared_with_replaced_via_accessor_parent":                                                             "flaky - sometimes complains that a relation does not exist",
	"InverseHasOneTests#test_parent_instance_should_be_shared_with_newly_built_child":                                                                          "flaky - sometimes complains that a relation does not exist",
	"IgnoreFixturesTest#test_ignores_books_fixtures":                                                                                                           "flaky - FK constraint violated sometimes when loading all fixture data",
	"IgnoreFixturesTest#test_ignores_parrots_fixtures":                                                                                                         "flaky - FK constraint violated sometimes when loading all fixture data",
	"JsonAttributeTest#test_select_array_json_value":                                                                                                           "flaky - sometimes attempts to call a method on nil",
	"LengthValidationTest#test_validates_size_of_association_using_within":                                                                                     "flaky - sometimes complains that a relation does not exist",
	"NestedThroughAssociationsTest#test_distinct_has_many_through_a_has_many_through_association_on_source_reflection":                                         "flaky - sometimes complains that a relation does not exist",
	"NestedThroughAssociationsTest#test_has_many_through_has_many_through_with_belongs_to_source_reflection_preload_via_joins":                                 "flaky - sometimes complains that a relation does not exist",
	"NestedThroughAssociationsTest#test_nested_has_many_through_with_conditions_on_through_associations_preload_via_joins":                                     "flaky - sometimes complains that a relation does not exist",
	"NestedThroughAssociationsTest#test_polymorphic_has_many_through_when_through_association_has_already_loaded":                                              "flaky - sometimes complains that a relation does not exist",
	"PessimisticLockingTest#test_no_locks_no_wait":                                                                                                             "flaky - sometimes complains that a relation does not exist",
	"PessimisticLockingTest#test_with_lock_configures_transaction":                                                                                             "flaky - sometimes complains that a relation does not exist",
	"PostgresqlInfinityTest#test_where_clause_with_infinite_range_on_a_datetime_column":                                                                        "flaky - sometimes complains that a relation does not exist",
	"PostgresqlIntervalTest#test_interval_type":                                                                                                                "flaky",
	"PostgresqlTimestampFixtureTest#test_bc_timestamp":                                                                                                         "flaky - sometimes datetime format does not match",
	"PostgresqlTimestampWithAwareTypesTest#test_timestamp_with_zone_values_with_rails_time_zone_support_and_time_zone_set":                                     "flaky - sometimes complains given Time instead of ActiveSupport::TimeWithZone",
	"PostgresqlTimestampWithTimeZoneTest#test_timestamp_with_zone_values_with_rails_time_zone_support_and_timestamptz_and_time_zone_set":                       "flaky - sometimes complains given Time instead of ActiveSupport::TimeWithZone",
	"QueryCacheTest#test_select_all_with_cache":                                                                                                                "flaky - sometimes complains that a relation does not exist",
	"RelationTest#test_finding_last_with_arel_order":                                                                                                           "flaky - sometimes complains that a relation does not exist",
	"RelationTest#test_joins_with_select_custom_attribute":                                                                                                     "flaky - sometimes complains that a relation does not exist",
	"RelationTest#test_order_triggers_eager_loading_when_ordering_using_hash_syntax":                                                                           "flaky - sometimes complains that a relation does not exist",
	"RelationTest#test_select_takes_a_variable_list_of_args":                                                                                                   "flaky - sometimes complains that a relation does not exist",
	"RelationTest#test_size_with_limit":                                                                                                                        "flaky - sometimes complains that a relation does not exist",
	"RelationTest#test_to_yaml":                                                                                                                                "flaky - sometimes complains that a relation does not exist",
	"ReservedWordTest#test_calculations_work_with_reserved_words":                                                                                              "flaky - sometimes complains that a relation does not exist",
	"StrictLoadingTest#test_raises_on_lazy_loading_a_strict_loading_belongs_to_relation":                                                                       "flaky - sometimes complains that a relation does not exist",
	"StrictLoadingTest#test_raises_on_unloaded_relation_methods_if_strict_loading":                                                                             "flaky - sometimes complains that a relation does not exist",
	"StrictLoadingTest#test_strict_loading_can_be_turned_off_on_an_association_in_a_model_with_strict_loading_on":                                              "flaky - sometimes complains that a relation does not exist",
	"StrictLoadingTest#test_strict_loading_with_reflection_is_ignored_in_validation_context":                                                                   "flaky - sometimes complains that a relation does not exist",
	"TestHasOneAutosaveAssociationWhichItselfHasAutosaveAssociations#test_when_great-grandchild_changed_in_memory,_saving_parent_should_save_great-grandchild": "flaky - sometimes attempts to call a method on nil",
	"TimestampTest#test_changing_parent_of_a_record_touches_both_new_and_old_polymorphic_parent_record_changes_with_other_class":                               "flaky - sometimes complains that a relation does not exist",
	"TimestampTest#test_destroying_a_record_with_a_belongs_to_that_specifies_touching_the_parent_should_update_the_parent_updated_at":                          "flaky - sometimes complains that a relation does not exist",
	"TimestampTest#test_timestamp_column_values_are_present_in_the_callbacks":                                                                                  "flaky - sometimes complains that a relation does not exist",
	"TimestampTest#test_touching_a_record_touches_parent_record_and_grandparent_record":                                                                        "flaky - sometimes complains that a relation does not exist",
	"TimestampTest#test_touching_a_record_with_a_belongs_to_that_uses_a_counter_cache_should_update_the_parent":                                                "flaky - sometimes complains that a relation does not exist",
	"TimestampTest#test_touching_update_at_attribute_as_symbol_updates_timestamp":                                                                              "flaky - sometimes complains that a relation does not exist",
	"TouchTest#test_create_turned_off":                                                                                                                         "flaky - sometimes complains that a relation does not exist",
	"TouchTest#test_update":                                                                                                                                    "flaky - sometimes complains that a relation does not exist",
	"TransactionsWithTransactionalFixturesTest#test_no_automatic_savepoint_for_inner_transaction":                                                              "flaky - sometimes attempts to call a method on nil",
	"TypeTest#test_registering_a_new_type":                                                                                                                     "flaky - sometimes attempts to call a method on nil",
	"UpdateAllTest#test_update_all_cares_about_optimistic_locking":                                                                                             "flaky - sometimes complains that a relation does not exist",
	"UniquenessValidationTest#test_validate_case_insensitive_uniqueness_with_special_sql_like_chars":                                                           "flaky - sometimes complains that a relation does not exist",
}
