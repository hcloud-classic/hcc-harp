package graphql

import (
	"github.com/graphql-go/graphql"
	"hcc/harp/floatingip"
	"hcc/harp/logger"
	"hcc/harp/mysql"
	"hcc/harp/subnet"
	"hcc/harp/types"
	"time"
)

var queryTypes = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{

			////////////////////////////// Action ///////////////////////////////
			// http://localhost:8001/graphql?query={createSubnet(uuid:"6b18ae6c-d834-479b-62e0-80b04f5deed7"){uuid}}
			"updateSubnet": &graphql.Field{
				Type:        subnetType,
				Description: "Create subnet by uuid",
				Args: graphql.FieldConfigArgument{
					"uuid": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					subnet.UpdateSubnet()
					return nil, nil
				},
			},

			// http://localhost:8001/graphql?query={createFloatingip(uuid:"6b18ae6c-d834-479b-62e0-80b04f5deed7"){uuid}}
			"createFloatingip": &graphql.Field{
				Type:        subnetType,
				Description: "Create floating IP by uuid",
				Args: graphql.FieldConfigArgument{
					"uuid": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					floatingip.CreateFloatingIP()
					return nil, nil
				},
			},

			////////////////////////////// Subnet Read///////////////////////////////
			/* Get (read) single subnet by uuid
			   http://localhost:8001/graphql?query={subnet(uuid:"[volume_uuid]]"){uuid,size,type,server_uuid}}
			*/
			"subnet": &graphql.Field{
				Type:        subnetType,
				Description: "Get subnet by uuid",
				Args: graphql.FieldConfigArgument{
					"uuid": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: subnet")

					requestedUUID, ok := p.Args["uuid"].(string)
					if ok {
						subnet := new(types.Subnet)

						var uuid string
						var name string
						var networkIP string
						var netmask string
						var os string
						var createdAt time.Time

						sql := "select * from subnet where uuid = ?"
						err := mysql.Db.QueryRow(sql, requestedUUID).Scan(&uuid, &name, &networkIP, &netmask, &os, &createdAt)
						if err != nil {
							logger.Logger.Println(err)
							return nil, nil
						}

						subnet.UUID = uuid
						subnet.Name = name
						subnet.NetworkIP = networkIP
						subnet.Netmask = netmask
						subnet.Os = os
						subnet.CreatedAt = createdAt

						return subnet, nil
					}
					return nil, nil
				},
			},

			/* Get the number of subnet */
			// http://localhost:8001/graphql?query={createSubnet(uuid:"6b18ae6c-d834-479b-62e0-80b04f5deed7"){uuid}}
			//"num_subnet": &graphql.Field{
			//	Type:        subnetType,
			//	Description: "Create subnet by uuid",
			//	Args: graphql.FieldConfigArgument{
			//		"uuid": &graphql.ArgumentConfig{
			//			Type: graphql.String,
			//		},
			//	},
			//	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			//
			//		sql := "select count(*) from subnet"
			//		err := mysql.Db.QueryRow(sql).Scan()
			//		if err != nil {
			//			logger.Logger.Println(err)
			//			return nil, nil
			//		}
			//
			//		//subnet.CreateSubnet()
			//		//return nil, nil
			//		return num
			//	},
			//},

			/* Get (read) subnet list
			   http://localhost:8001/graphql?query={list_volume{uuid,size,type,server_uuid}}
			*/
			"list_subnet": &graphql.Field{
				Type:        graphql.NewList(subnetType),
				Description: "Get subnet list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: list_subnet")

					var subnets []types.Subnet
					var uuid string
					var name string
					var networkIP string
					var netmask string
					var os string
					var createdAt time.Time

					sql := "select * from subnet"
					stmt, err := mysql.Db.Query(sql)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}
					defer stmt.Close()

					for stmt.Next() {
						err := stmt.Scan(&uuid, &name, &networkIP, &netmask, &os, &createdAt)
						if err != nil {
							logger.Logger.Println(err)
						}

						subnet := types.Subnet{UUID: uuid, Name: name, NetworkIP: networkIP, Netmask: netmask, Os: os, CreatedAt: createdAt}

						logger.Logger.Println(subnet)
						subnets = append(subnets, subnet)
					}

					return subnets, nil
				},
			},
		},
	})
