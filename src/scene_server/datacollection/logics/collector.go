/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package logics

import (
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
)

func (lgc *Logics) SearchCollector() ([]metadata.NetcollectorSearchResult, error) {
	result := []metadata.NetcollectorSearchResult{}

	mock := metadata.NetcollectorSearchResult{
		CloudID:   0,
		CloudName: "default area",
		InnerIP:   "192.168.1.1",
		Status: metadata.NetcollectorStatus{
			CollectorStatus: metadata.CollectorStatusNormal,
			ConfigStatus:    metadata.CollectorConfigStatusNormal,
			ReportStatus:    metadata.CollectorReportStatusNormal,
		},
		DeployTime: time.Now(),
		Version:    "",
		Config: metadata.NetcollectConfig{
			ScanRange: nil,
			Period:    "",
			Community: "",
		},
	}
	result = append(result, mock)
	return result, nil
}

func (lgc *Logics) UpdateCollector(config metadata.NetcollectorConfig) error {
	cond := condition.CreateCondition()
	cond.Field(common.BKCloudIDField).Eq(config.CloudID)
	cond.Field(common.BKHostInnerIPField).Eq(config.InnerIP)

	count, err := lgc.Instance.GetCntByCondition(common.BKTableNameNetcollectConfig, cond.ToMapStr())
	if err != nil {
		blog.Errorf("[UpdateCollector] count error: %v", err)
		return err
	}
	if count > 0 {
		err = lgc.Instance.UpdateByCondition(common.BKTableNameNetcollectConfig, config, cond)
		if err != nil {
			blog.Errorf("[UpdateCollector] UpdateByCondition error: %v", err)
			return err
		}
		return nil
	}

	_, err = lgc.Instance.Insert(common.BKTableNameNetcollectConfig, config)
	if err != nil {
		blog.Errorf("[UpdateCollector] UpdateByCondition error: %v", err)
		return err
	}

	return lgc.DiscoverNetDevice(config)
}

func (lgc *Logics) DiscoverNetDevice(config metadata.NetcollectorConfig) error {
	return nil
}