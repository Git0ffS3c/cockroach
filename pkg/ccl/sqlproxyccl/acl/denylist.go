// Copyright 2021 The Cockroach Authors.
//
// Licensed as a CockroachDB Enterprise file under the Cockroach Community
// License (the "License"); you may not use this file except in compliance with
// the License. You may obtain a copy of the License at
//
//     https://github.com/cockroachdb/cockroach/blob/master/licenses/CCL.txt

package acl

import (
	"strings"
	"time"

	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
	"github.com/cockroachdb/errors"
)

// File represents a on-disk version of the denylist config.
// This also serves as a spec of expected yaml file format.
type DenylistFile struct {
	Seq      int64        `yaml:"SequenceNumber"`
	Denylist []*DenyEntry `yaml:"denylist"`
}

// Denylist represents an in-memory cache for the current denylist.
// It also handles the logic of deciding what to be denied.
type Denylist struct {
	entries map[DenyEntity]*DenyEntry
}

func (dl *Denylist) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var f DenylistFile
	if err := unmarshal(&f); err != nil {
		return err
	}
	dl.entries = make(map[DenyEntity]*DenyEntry)
	for _, entry := range f.Denylist {
		dl.entries[entry.Entity] = entry
	}

	return nil
}

func (dl *Denylist) CheckConnection(
	connection ConnectionTags, timeSource timeutil.TimeSource,
) error {
	ip := DenyEntity{Item: connection.IP, Type: IPAddrType}
	if err := dl.denied(ip, timeSource); err != nil {
		return errors.Wrapf(err, "connection ip '%v' denied", connection.IP)
	}
	cluster := DenyEntity{Item: connection.Cluster, Type: ClusterType}
	if err := dl.denied(cluster, timeSource); err != nil {
		return errors.Wrapf(err, "connection cluster '%v' denied", connection.Cluster)
	}
	return nil
}

// denied returns an error if the entity is denied access. The error message
// describes the reason for the denial.
func (dl *Denylist) denied(entity DenyEntity, timeSource timeutil.TimeSource) error {
	if ent, ok := dl.entries[entity]; ok &&
		(ent.Expiration.IsZero() || !ent.Expiration.Before(timeSource.Now())) {
		return errors.Newf("%s", ent.Reason)
	}
	return nil
}

// DenyEntry records info about one denied entity,
// the reason and the expiration time.
// This also serves as spec for the yaml config format.
type DenyEntry struct {
	Entity     DenyEntity `yaml:"entity"`
	Expiration time.Time  `yaml:"expiration"`
	Reason     string     `yaml:"reason"`
}

// DenyEntity represent one denied entity.
// This also serves as the spec for the config format.
type DenyEntity struct {
	Item string   `yaml:"item"`
	Type DenyType `yaml:"type"`
}

// DenyType is the type of the denied entity.
type DenyType int

// Enum values for DenyType.
const (
	IPAddrType DenyType = iota + 1
	ClusterType
	UnknownType
)

var strToTypeMap = map[string]DenyType{
	"ip":      IPAddrType,
	"cluster": ClusterType,
}

var typeToStrMap = map[DenyType]string{
	IPAddrType:  "ip",
	ClusterType: "cluster",
}

// UnmarshalYAML implements yaml.Unmarshaler interface for type.
func (typ *DenyType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw string
	err := unmarshal(&raw)
	if err != nil {
		return err
	}

	normalized := strings.ToLower(raw)
	t, ok := strToTypeMap[normalized]
	if !ok {
		*typ = UnknownType
	} else {
		*typ = t
	}

	return nil
}

// MarshalYAML implements yaml.Marshaler interface for type.
func (typ DenyType) MarshalYAML() (interface{}, error) {
	return typ.String(), nil
}

// String implements Stringer interface for type.
func (typ DenyType) String() string {
	s, ok := typeToStrMap[typ]
	if !ok {
		return "UNKNOWN"
	}
	return s
}
