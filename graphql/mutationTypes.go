package graphql

import (
	"errors"
	"github.com/graphql-go/graphql"
	"hcc/harp/logger"
	"hcc/harp/mysql"
	"hcc/harp/rabbitmq"
	"hcc/harp/types"
)

var mutationTypes = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		/* Create new subnet */
		// http://192.168.110.240:7400/graphql?query=mutation+_{create_subnet(network_ip:"192.168.110.0",netmask:"255.255.255.0",gateway:"192.168.110.254",next_server: "192.168.110.240",name:"hcc",name_server:"8.8.8.8",domain_name:"google.com"){network_ip,netmask,gateway,next_server,name,name_server,domain_name}}
		"create_subnet": &graphql.Field{
			Type:        graphql.String,
			Description: "Create new subnet",
			Args: graphql.FieldConfigArgument{
				"network_ip": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"netmask": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"gateway": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"next_server": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"name_server": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"domain_name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: create_subnet")

				networkIP, networkIPOk := params.Args["network_ip"].(string)
				netmask, netmaskOk := params.Args["netmask"].(string)
				gateway, gatewayOk := params.Args["gateway"].(string)
				nextServer, nextServerOk := params.Args["next_server"].(string)
				name, nameOk := params.Args["name"].(string)

				if !networkIPOk {
					return nil, errors.New("need network_ip argument")
				}
				if !netmaskOk {
					return nil, errors.New("need netmask argument")
				}
				if !gatewayOk {
					return nil, errors.New("need gateway argument")
				}
				if !nextServerOk {
					return nil, errors.New("need next_server argument")
				}
				if !nameOk {
					return nil, errors.New("need name argument")
				}

				nameServer, nameServerOk := params.Args["name_server"].(string)
				if !nameServerOk {
					nameServer = ""
				}

				domainName, domainNameOk := params.Args["domain_name"].(string)
				if !domainNameOk {
					domainName = ""
				}

				subnet := types.Subnet{
					NetworkIP:  networkIP,
					Netmask:    netmask,
					Gateway:    gateway,
					NextServer: nextServer,
					Name:       name,
					NameServer: nameServer,
					DomainName: domainName,
				}

				err := rabbitmq.CreateSubnet(subnet)
				if err != nil {
					return nil, err
				}

				return "create_subnet queued successfully", nil
			},
		},
		////////////////////////////////////////////////////////////////////////////////
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
						logger.Logger.Println(err)
						return nil, err
					}
					defer func() {
						_ = stmt.Close()
					}()
					result, err2 := stmt.Exec(subnet.Name, subnet.NetworkIP, subnet.Netmask, subnet.Os, subnet.UUID)
					if err2 != nil {
						logger.Logger.Println(err2)
						return nil, err2
					}
					logger.Logger.Println(result.LastInsertId())

					return subnet, nil
				}
				return nil, errors.New("need ................... arguments")
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
						logger.Logger.Println(err)
						return nil, err
					}
					defer func() {
						_ = stmt.Close()
					}()
					result, err2 := stmt.Exec(requestedUUID)
					if err2 != nil {
						logger.Logger.Println(err2)
						return nil, err2
					}
					logger.Logger.Println(result.RowsAffected())

					return requestedUUID, nil
				}
				return nil, errors.New("need uuid argument")
			},
		},
	},
})
