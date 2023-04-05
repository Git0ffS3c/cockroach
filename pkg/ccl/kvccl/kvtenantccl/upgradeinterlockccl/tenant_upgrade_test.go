// Copyright 2021 The Cockroach Authors.
//
// Licensed as a CockroachDB Enterprise file under the Cockroach Community
// License (the "License"); you may not use this file except in compliance with
// the License. You may obtain a copy of the License at
//
//     https://github.com/cockroachdb/cockroach/blob/master/licenses/CCL.txt

package upgradeinterlockccl_test

import (
	"context"
	gosql "database/sql"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/clusterversion"
	"github.com/cockroachdb/cockroach/pkg/jobs"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/security/username"
	"github.com/cockroachdb/cockroach/pkg/server"
	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/spanconfig"
	"github.com/cockroachdb/cockroach/pkg/sql"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/lease"
	"github.com/cockroachdb/cockroach/pkg/sql/sessiondatapb"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlinstance/instancestorage"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlliveness/slinstance"
	"github.com/cockroachdb/cockroach/pkg/sql/stats"
	"github.com/cockroachdb/cockroach/pkg/testutils/serverutils"
	"github.com/cockroachdb/cockroach/pkg/testutils/skip"
	"github.com/cockroachdb/cockroach/pkg/testutils/sqlutils"
	"github.com/cockroachdb/cockroach/pkg/testutils/testcluster"
	"github.com/cockroachdb/cockroach/pkg/upgrade/upgradebase"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/stop"
	"github.com/stretchr/testify/require"
)

// TestTenantUpgradeInterlock validates the interlock between upgrading SQL
// servers and other SQL servers which may be running (and/or in the process
// of starting up). It runs with two SQL servers for the same tenant and starts
// one server which performs the upgrade, while the second server is starting
// up. It steps through all the phases of the interlock and validates that the
// system performs as expected for both starting SQL servers which are at the
// correct binary version, and those that are using a binary version that
// is too low.
func TestTenantUpgradeInterlock(t *testing.T) {
	defer leaktest.AfterTest(t)()
	// Times out under stress race.
	skip.UnderStressRace(t)
	// Test takes 100s+ to run.
	skip.UnderShort(t)

	defer log.Scope(t).Close(t)
	ctx := context.Background()

	const (
		currentBinaryVersion = iota
		laggingBinaryVersion
		numConfigs
	)

	type interlockTestVariant int

	var variants = map[interlockTestVariant]string{
		currentBinaryVersion: "current binary version",
		laggingBinaryVersion: "lagging binary version",
	}

	type interlockTestConfig struct {
		name          string
		expUpgradeErr [numConfigs][]string // empty if expecting a nil error
		expStartupErr [numConfigs]string   // empty if expecting a nil error
		pausePoint    upgradebase.PausePoint
	}

	var tests = []interlockTestConfig{
		{
			// In this case we won't see the new server in the first
			// transaction, and will instead see it when we try and commit the
			// first transaction.
			name:       "pause after first check for instances",
			pausePoint: upgradebase.AfterFirstCheckForInstances,
			expUpgradeErr: [numConfigs][]string{
				{""},
				{"pq: upgrade failed due to active SQL servers with incompatible binary version",
					fmt.Sprintf("sql server 2 is running a binary version %s which is less than the attempted upgrade version", clusterversion.TestingBinaryMinSupportedVersion.String())},
			},
			expStartupErr: [numConfigs]string{
				"",
				"",
			},
		},
		{
			name:       "pause after fence RPC",
			pausePoint: upgradebase.AfterFenceRPC,
			expUpgradeErr: [numConfigs][]string{
				{""},
				{"pq: upgrade failed due to active SQL servers with incompatible binary version"},
			},
			expStartupErr: [numConfigs]string{
				"",
				"",
			},
		},
		{
			name:       "pause after fence write to settings table",
			pausePoint: upgradebase.AfterFenceWriteToSettingsTable,
			expUpgradeErr: [numConfigs][]string{
				{""},
				{""},
			},
			expStartupErr: [numConfigs]string{
				"",
				"preventing SQL server from starting because its binary version is too low for the tenant active version",
			},
		},
		{
			name:       "pause after second check of instances",
			pausePoint: upgradebase.AfterSecondCheckForInstances,
			expUpgradeErr: [numConfigs][]string{
				{""},
				{""},
			},
			expStartupErr: [numConfigs]string{
				"",
				"preventing SQL server from starting because its binary version is too low for the tenant active version",
			},
		},
		{
			name:       "pause after migration",
			pausePoint: upgradebase.AfterMigration,
			expUpgradeErr: [numConfigs][]string{
				{""},
				{""},
			},
			expStartupErr: [numConfigs]string{
				"",
				"preventing SQL server from starting because its binary version is too low for the tenant active version",
			},
		},
		{
			name:       "pause after version bump RPC",
			pausePoint: upgradebase.AfterVersionBumpRPC,
			expUpgradeErr: [numConfigs][]string{
				{""},
				{""},
			},
			expStartupErr: [numConfigs]string{
				"",
				"preventing SQL server from starting because its binary version is too low for the tenant active version",
			},
		},
		{
			name:       "pause after write to settings table",
			pausePoint: upgradebase.AfterVersionWriteToSettingsTable,
			expUpgradeErr: [numConfigs][]string{
				{""},
				{""},
			},
			expStartupErr: [numConfigs]string{
				"",
				"preventing SQL server from starting because its binary version is too low for the tenant active version",
			},
		},
	}

	runTest := func(t *testing.T, variant interlockTestVariant, test interlockTestConfig) {
		logf := func(format string, args ...interface{}) {
			t.Helper()
			newArgs := make([]interface{}, len(args)+1)
			newArgs[0] = t.Name()
			copy(newArgs[1:], args)
			t.Logf("(%s) "+format, newArgs...)
		}
		logf(`upgrade interlock test: running variant "%s", configuration: "%s"`, variants[variant], test.name)

		reachedChannel := make(chan struct{})
		resumeChannel := make(chan struct{})
		completedChannel := make(chan struct{})
		defer close(reachedChannel)
		defer close(resumeChannel)
		defer close(completedChannel)

		bv := clusterversion.TestingBinaryVersion
		msv := clusterversion.TestingBinaryMinSupportedVersion

		// If there are any non-empty errors expected on the upgrade, then we're
		// expecting it to fail.
		expectingUpgradeToFail := test.expUpgradeErr[variant][0] != ""
		finalUpgradeVersion := clusterversion.TestingBinaryVersion
		if expectingUpgradeToFail {
			finalUpgradeVersion = msv
		}

		disableBackgroundTasks := func(s *cluster.Settings) {
			// This test makes use of inter-SQL server gRPCs, some of which are
			// going between SQL pods of incompatible version numbers. In cases
			// where a gRPC handshake detects an invalid versioned SQL server,
			// it will fail, and also trip the gRPC breaker. This will then
			// cause the test to fail with a "breaker open" error, instead of
			// the error we're expecting. As a result, we disable everything
			// that could run gRPCs in the background to ensure that the upgrade
			// gRPCs detect the first version mismatch error.
			stats.AutomaticStatisticsClusterMode.Override(ctx, &s.SV, false)
			stats.UseStatisticsOnSystemTables.Override(ctx, &s.SV, false)
			stats.AutomaticStatisticsOnSystemTables.Override(ctx, &s.SV, false)
			sql.DistSQLClusterExecMode.Override(ctx, &s.SV, int64(sessiondatapb.DistSQLOff))
		}

		reduceLeaseDurationAndReclaimLoopInterval := func(s *cluster.Settings) {
			// In cases where the other SQL server fails, it may go down holding
			// a descriptor lease. In that case, we want to decrease the lease
			// duration so that we can reclaim the lease faster than the default
			// (5 minutes). This allows the test to complete faster.
			lease.LeaseDuration.Override(ctx, &s.SV, 2000*time.Millisecond)
			lease.LeaseRenewalDuration.Override(ctx, &s.SV, 0)
			// Also in the failure case, we need the SQL instance (stored in
			// the system.sql_instances table) for the failed SQL server to
			// be removed as quickly as possible. This is because a failed
			// instance in that table can cause RPC attempts to that server
			// to fail. The changes below in the other server setup will
			// ensure that the failed SQL server expires quickly. Here we
			// need to also ensure that it will be reclaimed quickly.
			instancestorage.ReclaimLoopInterval.Override(ctx, &s.SV, 500*time.Millisecond)
		}

		reduceDefaultTTLAndHeartBeat := func(s *cluster.Settings) {
			// In the case where a SQL server is expected to fail, we need the
			// SQL instance (stored in the system.sql_instances table) for the
			// failed SQL server to be removed as quickly as possible. This is
			// because a failed instance in that table can cause RPC attempts to
			// that server to fail. To accomplish this, we decrease the default
			// TTL and heartbeat for the sessions so that they expire and are
			// cleaned up faster.
			slinstance.DefaultTTL.Override(ctx, &s.SV, 250*time.Millisecond)
			slinstance.DefaultHeartBeat.Override(ctx, &s.SV, 50*time.Millisecond)
		}

		// Initialize the version to the BinaryMinSupportedVersion so that
		// we can perform upgrades.
		settings := cluster.MakeTestingClusterSettingsWithVersions(bv, msv, false /* initializeVersion */)
		disableBackgroundTasks(settings)
		require.NoError(t, clusterversion.Initialize(ctx, msv, &settings.SV))

		tc := testcluster.StartTestCluster(t, 1, base.TestClusterArgs{
			ServerArgs: base.TestServerArgs{
				// Test validates tenant behavior. No need for the default test
				// tenant.
				DefaultTestTenant: base.TestTenantDisabled,
				Settings:          settings,
				Knobs: base.TestingKnobs{
					JobsTestingKnobs: jobs.NewTestingKnobsWithShortIntervals(),
					SpanConfig: &spanconfig.TestingKnobs{
						ManagerDisableJobCreation: true,
					},
					Server: &server.TestingKnobs{
						DisableAutomaticVersionUpgrade: make(chan struct{}),
						// Initialize to the minimum supported version
						// so that we can perform the upgrade below.
						BinaryVersionOverride: msv,
					},
				},
			},
		})
		defer tc.Stopper().Stop(ctx)

		connectToTenant := func(t *testing.T, addr string) (_ *gosql.DB, cleanup func()) {
			pgURL, cleanupPGUrl := sqlutils.PGUrl(t, addr, "Tenant", url.User(username.RootUser))
			tenantDB, err := gosql.Open("postgres", pgURL.String())
			require.NoError(t, err)
			return tenantDB, func() {
				tenantDB.Close()
				cleanupPGUrl()
			}
		}

		mkTenant := func(t *testing.T, id roachpb.TenantID, bv roachpb.Version, minBv roachpb.Version) (tenantDB *gosql.DB, cleanup func()) {
			settings := cluster.MakeTestingClusterSettingsWithVersions(bv, minBv, false /* initializeVersion */)
			disableBackgroundTasks(settings)
			if test.expStartupErr[variant] != "" {
				reduceLeaseDurationAndReclaimLoopInterval(settings)
			}
			require.NoError(t, clusterversion.Initialize(ctx, minBv, &settings.SV))
			tenantArgs := base.TestTenantArgs{
				TenantID: id,
				TestingKnobs: base.TestingKnobs{
					JobsTestingKnobs: jobs.NewTestingKnobsWithShortIntervals(),
					UpgradeManager: &upgradebase.TestingKnobs{
						InterlockPausePoint:               test.pausePoint,
						InterlockResumeChannel:            &resumeChannel,
						InterlockReachedPausePointChannel: &reachedChannel,
					},
				},
				Settings: settings,
			}
			tenant, err := tc.Server(0).StartTenant(ctx, tenantArgs)
			require.NoError(t, err)
			return connectToTenant(t, tenant.SQLAddr())
		}

		logf("creating an initial tenant server")
		// Create a tenant before upgrading anything, and verify its
		// version.
		tenantID := serverutils.TestTenantID()
		tenant, cleanup := mkTenant(t, tenantID, bv, msv)
		defer cleanup()
		initialTenantRunner := sqlutils.MakeSQLRunner(tenant)

		logf("verifying the tenant version")
		// Ensure that the tenant works.
		initialTenantRunner.CheckQueryResults(t, "SHOW CLUSTER SETTING version",
			[][]string{{msv.String()}})
		logf("verifying basic SQL functionality")
		initialTenantRunner.Exec(t, "CREATE TABLE t (i INT PRIMARY KEY)")
		initialTenantRunner.Exec(t, "INSERT INTO t VALUES (1), (2)")
		initialTenantRunner.CheckQueryResults(t, "SELECT * FROM t", [][]string{{"1"}, {"2"}})

		logf("verifying the version of the storage cluster")
		// Validate that the host cluster is at the expected version, and
		// upgrade it to the binary version.
		hostClusterRunner := sqlutils.MakeSQLRunner(tc.ServerConn(0))
		hostClusterRunner.CheckQueryResults(t, "SHOW CLUSTER SETTING version",
			[][]string{{msv.String()}})

		logf("upgrading the storage cluster")
		hostClusterRunner.Exec(t, "SET CLUSTER SETTING version = $1", bv.String())

		logf("checking the tenant after the storage cluster upgrade")
		// Ensure that the tenant still works.
		initialTenantRunner.CheckQueryResults(t, "SELECT * FROM t", [][]string{{"1"}, {"2"}})

		logf("start upgrading the tenant")
		// Start upgrading the tenant. This call will pause at the specified
		// pause point, and once resumed (below, outside this function),
		// will complete the upgrade.
		go func() {
			release := func() {
				log.Infof(ctx, "upgrade completed, resuming post-upgrade work")
				completedChannel <- struct{}{}
			}
			defer release()

			// Since we're in a separate go routine from the main test, we can't
			// use the SQLRunner methods for running queries, since they Fatal
			// under the covers. See https://go.dev/doc/go1.16#vet-testing-T for
			// more details.
			if expectingUpgradeToFail {
				getPossibleUpgradeErrorsString := func(variant interlockTestVariant) string {
					possibleErrorsString := ""
					for i := range test.expUpgradeErr[variant] {
						expErrString := test.expUpgradeErr[variant][i]
						if i > 0 {
							possibleErrorsString += " OR "
						}
						possibleErrorsString += `"` + expErrString + `"`
					}
					return possibleErrorsString
				}
				_, err := tenant.Exec("SET CLUSTER SETTING version = $1", bv.String())
				if err == nil {
					t.Errorf("expected %s, got: success", getPossibleUpgradeErrorsString(variant))
				} else {
					foundExpErr := false
					expErrString := ""
					for i := range test.expUpgradeErr[variant] {
						expErrString = test.expUpgradeErr[variant][i]
						if strings.Contains(err.Error(), expErrString) {
							foundExpErr = true
							break
						}
					}
					if !foundExpErr {
						t.Errorf("expected %s, got: %s", getPossibleUpgradeErrorsString(variant), err.Error())
					}
				}
			} else {
				_, err := tenant.Exec("SET CLUSTER SETTING version = $1", bv.String())
				if err != nil {
					t.Error("unexpected error: ", err.Error())
				}
			}
		}()

		// Wait until the upgrader reaches the pause point before
		// starting.
		<-reachedChannel
		logf("upgrader is ready")

		logf("starting another tenant server")
		// Now start a second SQL server for this tenant to see how the
		// two SQL servers interact.
		otherMsv := msv
		otherBv := bv
		if variant == laggingBinaryVersion {
			// If we're in "lagging binary" mode, we want the server to
			// startup with a binary that is too old for the upgrade to
			// succeed. To make this happen we set the binary version to
			// the tenant's minimum binary version, and then set the
			// server's minimum binary version to one major release
			// earlier. The last step isn't strictly required, but we
			// do it to prevent any code from tripping which validates
			// that the binary version cannot equal the minimum binary
			// version.
			otherBv = msv
			otherMsv.Major = msv.Major - 1
		}
		otherServerSettings := cluster.MakeTestingClusterSettingsWithVersions(otherBv, otherMsv, false /* initializeVersion */)
		disableBackgroundTasks(otherServerSettings)
		if test.expStartupErr[variant] != "" {
			reduceLeaseDurationAndReclaimLoopInterval(otherServerSettings)
			reduceDefaultTTLAndHeartBeat(otherServerSettings)
		}
		require.NoError(t, clusterversion.Initialize(ctx, otherMsv, &otherServerSettings.SV))
		otherServerStopper := stop.NewStopper()
		otherServer, otherServerStartError := tc.Server(0).StartTenant(ctx,
			base.TestTenantArgs{
				Stopper:  otherServerStopper,
				TenantID: tenantID,
				Settings: otherServerSettings,
			})

		var otherTenantRunner *sqlutils.SQLRunner
		expectingStartupToFail := test.expStartupErr[variant] != ""
		numTenantsStr := "1"
		if expectingStartupToFail {
			logf("shutting down the other tenant server")
			otherServerStopper.Stop(ctx)
		} else if otherServerStartError == nil {
			defer otherServer.Stopper().Stop(ctx)
			otherTenant, otherCleanup := connectToTenant(t, otherServer.SQLAddr())
			defer otherCleanup()
			otherTenantRunner = sqlutils.MakeSQLRunner(otherTenant)
			numTenantsStr = "2"
		}

		logf("waiting for the instance table to get in the right state")
		// Based on the success or failure of the "other" SQL server startup,
		// we'll either have 1 or 2 SQL servers running at this point. Confirm
		// that we're in the desired state by querying the sql_instances table
		// directly. We do this before we continue to ensure that the upgrade
		// doesn't encounter any stale SQL instances.
		initialTenantRunner.CheckQueryResultsRetry(t,
			"SELECT count(*) FROM system.sql_instances WHERE session_id IS NOT NULL", [][]string{{numTenantsStr}})

		// With tenant started, resume the upgrade and wait for it to
		// complete.
		logf("resuming upgrade")
		resumeChannel <- struct{}{}
		logf("waiting for upgrade to complete")
		<-completedChannel
		logf("upgrade completed")
		log.Infof(ctx, "continuing after upgrade completed")

		// Now that we've resumed the upgrade, process any errors. We must do
		// that after the upgrade completes because any unexpected errors above
		// will cause the test to FailNow, thus bringing down the cluster and
		// the tenant, which could cause the upgrade to hang if it's not resumed
		// first.
		if expectingStartupToFail {
			require.ErrorContains(t, otherServerStartError, test.expStartupErr[variant])
		} else {
			require.NoError(t, otherServerStartError)
		}

		// Handle errors and any subsequent processing once the upgrade
		// is completed.
		if test.expStartupErr[variant] == "" {
			// Validate that the version is as expected, and that the
			// second sql pod still works.
			logf("check the second server still works")
			otherTenantRunner.CheckQueryResults(t, "SELECT * FROM t", [][]string{{"1"}, {"2"}})
			logf("waiting for second server to reach target final version")
			otherTenantRunner.CheckQueryResultsRetry(t, "SHOW CLUSTER SETTING version",
				[][]string{{finalUpgradeVersion.String()}})
		}
	}

	for variant := range variants {
		variantName := variants[variant]
		t.Run(variantName, func(t *testing.T) {
			for i := range tests {
				test := tests[i]
				t.Run(test.name, func(t *testing.T) {
					runTest(t, variant, test)
				})
			}
		})
	}
}
