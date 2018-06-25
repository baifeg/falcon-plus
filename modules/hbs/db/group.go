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
	gid, err := getGroupId(group)
	if err != nil || gid == 0 {
		return false, err
	}

	hid, err2 := getHostId
	if err2 != nil || hid == 0 {
		return false, err
	}

	url = "select count(*) as aCount from grp_host where grp_id=? and host_id=?"
	rows3, err4 := DB.Query(url, gid, hid)
	if err4 != nil {
		log.Println("ERROR:", err4)
		return false, err4
	}
	defer rows3.Close()
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

func getGroupId(group string) (int, error) {
	// 获取group的ID，不存在则创建
	url := "select id as gid from grp where grp_name=?"
	rows, err := DB.Query(url, group)
	if err != nil {
		log.Println("ERROR:", err)
		return 0, err
	}

	defer rows.Close()
	var gid int
	for rows.Next() {
		rows.Scan(&gid)
		break
	}
	if &gid == nil || gid <= 0 {
		url = "insert into grp(grp_name, create_user) value (?, ?)"
		result, err2 := DB.Exec(url, group, "hbs")
		if err2 != nil {
			log.Println("ERROR:", err2)
			return 0, err2
		}
		lastId, _ := result.LastInsertId()
		gid = int(lastId)
		log.Printf("Info: Add new hostgroup [%s], id [%d]\n", group, gid)
	}
	return gid, nil
}

func getHostId(host string) (int, error) {
	// 获取host的ID，必定存在
	url := "select id as hid from host where hostname=?"
	rows, err := DB.Query(url, host)
	if err != nil {
		log.Println("ERROR:", err)
		return 0, err
	}

	defer rows.Close()
	var hid int
	for rows.Next() {
		rows.Scan(&hid)
		break
	}

	if &hid == nil {
		return 0, nil
	}
	return hid, nil

}
