package graphql

import (
	"github.com/graphql-go/graphql"
	"hcc/harp/logger"
	"hcc/harp/mysql"
	"hcc/harp/types"
	"hcc/harp/uuidgen"
)

var mutationTypes = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		/* Create new volume */
		// http://localhost:8001/graphql?query=mutation+_{create_volume(size:1024000,type:"ext4",server_uuid:"[server_uuid]"){size,type,server_uuid}}
		"create_subnet": &graphql.Field{
			Type:        subnetType,
			Description: "Create new subnet",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"network_ip": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"netmask": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"os": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: create_subnet")

				uuid, err := uuidgen.Uuidgen()
				if err != nil {
					logger.Logger.Println("Failed to generate uuid!")
					return nil, nil
				}

				subnet := types.Subnet{
					UUID:      uuid,
					Name:      params.Args["name"].(string),
					NetworkIP: params.Args["network_ip"].(string),
					Netmask:   params.Args["netmask"].(string),
					Os:        params.Args["os"].(string),
				}

				//err = CheckServerUUID(subnet.ServerUUID)
				//if err != nil {
				//	logger.Logger.Println(err)
				//	return nil, nil
				//}

				sql := "insert into subnet(uuid, name, network_ip, netmask, os, created_at) values (?, ?, ?, ?, ?, now())"
				stmt, err := mysql.Db.Prepare(sql)
				if err != nil {
					logger.Logger.Println(err.Error())
					return nil, nil
				}
				defer stmt.Close()
				result, err2 := stmt.Exec(subnet.UUID, subnet.Name, subnet.NetworkIP, subnet.Netmask, subnet.Os)
				if err2 != nil {
					logger.Logger.Println(err2)
					return nil, nil
				}
				logger.Logger.Println(result.LastInsertId())

				return subnet, nil
			},
		},

		/* Update volume by uuid */
		// http://localhost:8001/graphql?query=mutation+_{update_subnet(uuid:"0ac56231-a0ee-4323-55ad-37c08c2d4a78",name:"aaaa",ip:"1234",netmask:"1234",os:"centos"){uuid,name,ip,netmask,os}}
		"update_subnet": &graphql.Field{
			Type:        subnetType,
			Description: "Update subnet by uuid",
			Args: graphql.FieldConfigArgument{
				"uuid": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"network_ip": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"netmask": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"os": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: update_subnet")

				requestedUUID, _ := params.Args["uuid"].(string)
				name := params.Args["name"].(string)
				ip, _ip := params.Args["network_ip"].(string)
				netmask, _netmask := params.Args["netmask"].(string)
				os, _os := params.Args["os"].(string)

				subnet := new(types.Subnet)

				if _ip && _netmask && _os {
					subnet.UUID = requestedUUID
					subnet.Name = name
					subnet.NetworkIP = ip
					subnet.Netmask = netmask
					subnet.Os = os

					sql := "update subnet set name = ?, network_ip = ?, netmask = ?, os = ? where uuid = ?"
					stmt, err := mysql.Db.Prepare(sql)
					if err != nil {
						logger.Logger.Println(err.Error())
						return nil, nil
					}
					defer stmt.Close()
					result, err2 := stmt.Exec(subnet.Name, subnet.NetworkIP, subnet.Netmask, subnet.Os, subnet.UUID)
					if err2 != nil {
						logger.Logger.Println(err2)
						return nil, nil
					}
					logger.Logger.Println(result.LastInsertId())

					return subnet, nil
				}
				return nil, nil
			},
		},

		/* Delete subnet by id */
		// http://localhost:8001/graphql?query=mutation+_{delete_subnet(uuid:"cccc"){uuid}}
		"delete_subnet": &graphql.Field{
			Type:        subnetType,
			Description: "Delete subnet by uuid",
			Args: graphql.FieldConfigArgument{
				"uuid": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: delete_subnet")

				requestedUUID, ok := params.Args["uuid"].(string)
				if ok {
					sql := "delete from subnet where uuid = ?"
					stmt, err := mysql.Db.Prepare(sql)
					if err != nil {
						logger.Logger.Println(err.Error())
						return nil, nil
					}
					defer stmt.Close()
					result, err2 := stmt.Exec(requestedUUID)
					if err2 != nil {
						logger.Logger.Println(err2)
						return nil, nil
					}
					logger.Logger.Println(result.RowsAffected())

					return requestedUUID, nil
				}
				return nil, nil
			},
		},
	},
})
