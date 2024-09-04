// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdaoMapper

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"testing"
)

func Test_check(t *testing.T) {
	file, err := os.Open("mapper.xml")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	var root rootElement
	decoder := xml.NewDecoder(file)
	err = decoder.Decode(&root)
	if err != nil {
		fmt.Printf("Error decoding XML: %v\n", err)
		return
	}

	fmt.Printf("Root element: %s\n", root.XMLName.Local)

	if root.XMLName.Local == "mappers" {
		file.Seek(0, 0)
		var mappers mappers
		err = decoder.Decode(&mappers)
		if err != nil {
			fmt.Printf("Error decoding mappers: %v\n", err)
			return
		}
		fmt.Println("Mappers:")
		for _, mapper := range mappers.MapperList {
			fmt.Printf("  Resource: %s\n", mapper.Resource)
		}
	} else if root.XMLName.Local == "mapper" {
		file.Seek(0, 0)
		var mapper mapperCell
		err = decoder.Decode(&mapper)
		if err != nil {
			fmt.Printf("Error decoding mapper: %v\n", err)
			return
		}
		fmt.Printf("Mapper resource: %s\n", mapper.Resource)
	} else {
		fmt.Printf("Unknown root element: %s\n", root.XMLName.Local)
	}
}

func Test_mappers(t *testing.T) {
	file, err := os.Open("mappers.xml")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Parse the XML file
	var mappers mappers
	decoder := xml.NewDecoder(file)
	err = decoder.Decode(&mappers)
	if err != nil {
		fmt.Printf("Error decoding XML: %v\n", err)
		return
	}

	// Print the parsed data
	fmt.Println("Mapper resources:")
	for _, mapper := range mappers.MapperList {
		fmt.Printf("  Resource: %s\n", mapper.Resource)
	}
}

//func Test_mapper(t *testing.T) {
//	// Open the XML file
//	file, err := os.Open("mapper.xml")
//	if err != nil {
//		fmt.Printf("Error opening file: %v\n", err)
//		return
//	}
//	defer file.Close()
//
//	// Parse the XML file
//	var mapper mapper
//	decoder := xml.NewDecoder(file)
//	err = decoder.Decode(&mapper)
//	if err != nil {
//		fmt.Printf("Error decoding XML: %v\n", err)
//		return
//	}
//
//	// Print the parsed data
//	fmt.Printf("Namespace: %s\n", mapper.Namespace)
//	fmt.Println("Select statements:")
//	for _, s := range mapper.Selects {
//		fmt.Printf("  ID: %s, ResultType: %s, Query: %s\n", s.ID, s.ResultType, s.Query)
//	}
//
//	fmt.Println("Insert statements:")
//	for _, ins := range mapper.Inserts {
//		fmt.Printf("  ID: %s, ParameterType: %s, Query: %s\n", ins.ID, ins.ParameterType, ins.Query)
//	}
//
//	fmt.Println("Update statements:")
//	for _, upd := range mapper.Updates {
//		fmt.Printf("  ID: %s, ParameterType: %s, Query: %s\n", upd.ID, upd.ParameterType, upd.Query)
//	}
//
//	fmt.Println("Delete statements:")
//	for _, del := range mapper.Deletes {
//		fmt.Printf("  ID: %s, ParameterType: %s, Query: %s\n", del.ID, del.ParameterType, del.Query)
//	}
//}

func Test_parser(t *testing.T) {
	mp := newMapperParser()
	mp.parser("mappers.xml")
}

func Test_mapper2(t *testing.T) {
	exampleXML := `
	<mapper namespace="exampleMapper">
		<select id="findUser" resultType="User">
			SELECT * FROM users
		</select>
	</mapper>
	`
	var mapper mapper
	decoder := xml.NewDecoder(strings.NewReader(exampleXML))
	decoder.Decode(&mapper)
	for _, selectStmt := range mapper.CrudNodes {
		fmt.Println(selectStmt.XMLName.Local, ",", selectStmt.XMLName.Space)
		fmt.Println(strings.TrimSpace(selectStmt.Query))
		fmt.Println(len(selectStmt.Dynamics))
	}
}

func Test_dynamic(t *testing.T) {
	exampleXML := `
	<mapper namespace="exampleMapper">
		<select id="findUser" resultType="User">
			SELECT * FROM users
			<where>
				<if test="isActive">
					AND active = 1
					<if test="isVerified">
						AND verified = 1
					</if>
				</if>
			</where>
			<choose>
				<when test="isAdmin">
					AND role = 'admin'
				</when>
				<otherwise>
					AND role = 'user'
				</otherwise>
			</choose>
			<set>
				<if test="name != null">
					name = #{name},
				</if>
				<if test="age != null">
					age = #{age}
				</if>
			</set>
			<foreach collection="ids" item="id" open="AND id IN (" close=")" separator=",">
				#{id}
			</foreach>
		</select>
	</mapper>
	`
	var mapper mapper
	decoder := xml.NewDecoder(strings.NewReader(exampleXML))
	decoder.Decode(&mapper)

	// 输出解析结果
	fmt.Printf("Namespace: %s\n", mapper.Namespace)
	for _, selectStmt := range mapper.CrudNodes {
		fmt.Printf("Select ID: %s, ResultType: %s\n", selectStmt.ID, selectStmt.ResultType)
		for _, part := range selectStmt.Dynamics {
			fmt.Println("Query Part:", part)
		}
	}
}
