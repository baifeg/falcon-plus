// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package db

import (
	"log"
)

func QueryHostGroups() (map[int][]int, error) {
	m := make(map[int][]int)

	sql := "select grp_id, host_id from grp_host"
	rows, err := DB.Query(sql)
	if err != nil {
		log.Println("ERROR:", err)
		return m, err
	}

	defer rows.Close()
	for rows.Next() {
		var gid, hid int
		err = rows.Scan(&gid, &hid)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		if _, exists := m[hid]; exists {
			m[hid] = append(m[hid], gid)
		} else {
			m[hid] = []int{gid}
		}
	}

	return m, nil
}

// 添加成功返回true, nil
// 已存在返回false, nil
// 异常返回false, err
func PutHostIntoGroupIfNecessary(host, group string) (bool, error) {
	// 获取group的ID，不存在则创建
	url := "select id as gid from grp where grp_name=?"
	rows, err := DB.Query(url, group)
	if err != nil {
		log.Println("ERROR:", err)
		return false, err
	}

	defer rows.Close()
	var gid int
	for rows.Next() {
		rows.Scan(&gid)
		break
	}
	if &gid == nil {
		url = "insert into grp(grp_name, create_user) value (?, ?)"
		result, err2 := DB.Exec(url, group, "hbs")
		if err2 != nil {
			log.Println("ERROR:", err2)
			return false, err2
		}
		if result.RowsAffected() == 0 {
			log.Printf("ERROR: Add hostgroup[%s] failed\n", group)
			return false, nil
		}
		lastId, _ := result.LastInsertId()
		gid = int(lastId)
		log.Printf("Info: Add new hostgroup [%s], id [%d]\n", group, gid)
	}

	// 获取host的ID，必定存在
	url = "select id as hid from host where hostname=?"
	rows2, err3 := DB.Query(url, host)
	if err3 != nil {
		log.Println("ERROR:", err3)
		return false, err3
	}

	defer rows2.Close()
	var hid int
	for rows2.Next() {
		rows2.Scan(&hid)
		break
	}

	if &hid == nil {
		return false, nil
	}

	url = "select count(*) as aCount from grp_host where grp_id=? and host_id=?"
	rows3, err4 := DB.Query(url, gid, hid)
	if err4 != nil {
		log.Println("ERROR:", err4)
		return false, err4
	}
	var aCount int
	for rows3.Next() {
		rows3.Scan(&aCount)
		break
	}

	// 已存在
	if aCount > 0 {
		return false, nil
	}

	url = "insert into grp_host(grp_id, host_id) value (?, ?)"
	_, err5 := DB.Exec(url, gid, hid)
	if err5 != nil {
		log.Println("ERROR:", err5)
		return false, err5
	}
	log.Printf("Info: Add host [%s] to group [%s]\n", host, group)

	return true, nil
}
