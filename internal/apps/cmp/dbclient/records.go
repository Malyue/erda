// Copyright (c) 2021 Terminus, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dbclient

import (
	"github.com/jinzhu/gorm"

	"github.com/erda-project/erda/pkg/database/dbengine"
	"github.com/erda-project/erda/pkg/strutil"
)

type recordsReader struct {
	db     *dbengine.DBEngine
	limit  int
	offset int
}

type recordsWriter struct {
	db *dbengine.DBEngine
}

func (c *DBClient) RecordsReader() *recordsReader {
	return &recordsReader{db: c.DBEngine, limit: 0, offset: -1}
}

func (r *recordsReader) PageNum(n int) *recordsReader {
	r.offset = n
	return r
}

func (r *recordsReader) PageSize(n int) *recordsReader {
	r.limit = n
	return r
}

func (r *recordsReader) updateDBEngine(db *gorm.DB) *recordsReader {
	r.db = &dbengine.DBEngine{DB: db}
	return r
}

func (r *recordsReader) ByIDs(ids ...string) *recordsReader {
	return r.updateDBEngine(r.db.Where("id in (?)", ids))
}

func (r *recordsReader) ByPipelineIDs(ids ...string) *recordsReader {
	return r.updateDBEngine(r.db.Where("pipeline_id in (?)", ids))
}

func (r *recordsReader) ByRecordTypes(tps ...string) *recordsReader {
	tpsStr := strutil.Map(tps, func(s string) string { return "\"" + s + "\"" })
	return r.updateDBEngine(r.db.Where("record_type in (?)", tpsStr))
}

func (r *recordsReader) ByStatuses(statuses ...string) *recordsReader {
	statusesStr := strutil.Map(statuses, func(s string) string { return "\"" + s + "\"" })
	return r.updateDBEngine(r.db.Where("status in (?)", statusesStr))
}

func (r *recordsReader) ByOrgID(orgid string) *recordsReader {
	return r.updateDBEngine(r.db.Where("org_id = ?", orgid))
}

func (r *recordsReader) ByClusterNames(clusternames ...string) *recordsReader {
	clusternamesStr := strutil.Map(clusternames, func(s string) string { return "\"" + s + "\"" })
	return r.updateDBEngine(r.db.Where("cluster_name in (?)", clusternamesStr))
}

func (r *recordsReader) ByUserIDs(userids ...string) *recordsReader {
	useridsStr := strutil.Map(userids, func(s string) string { return "\"" + s + "\"" })
	return r.updateDBEngine(r.db.Where("user_id in (?)", useridsStr))
}

func (r *recordsReader) ByCreateTime(beforeNSecs int) *recordsReader {
	return r.updateDBEngine(r.db.Where("created_at > now() - interval %d second", beforeNSecs))
}

func (r *recordsReader) ByUpdateTime(beforeNSecs int) *recordsReader {
	return r.updateDBEngine(r.db.Where("updated_at > now() - interval %d second", beforeNSecs))
}

func (r *recordsReader) Limit(n int) *recordsReader {
	r.limit = n
	return r
}
func (r *recordsReader) Count() (int64, error) {
	var count int64
	err := r.db.Model(&Record{}).Count(&count).Error
	return count, err
}

func (r *recordsReader) Do() ([]Record, error) {
	records := []Record{}
	expr := r.db.Order("created_at desc")
	if r.limit != 0 {
		expr = expr.Limit(r.limit)
	}
	if r.offset != -1 {
		expr = expr.Offset(r.offset)
	}
	if err := expr.Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func (c *DBClient) RecordsWriter() *recordsWriter {
	return &recordsWriter{db: c.DBEngine}
}

func (w *recordsWriter) Create(s *Record) (uint64, error) {
	db := w.db.Save(s)
	return s.ID, db.Error
}
func (w *recordsWriter) Update(s Record) error {
	return w.db.Model(&s).Updates(s).Error
}
func (w *recordsWriter) Delete(ids ...uint64) error {
	return w.db.Delete(Record{}, "id in (?)", ids).Error
}
