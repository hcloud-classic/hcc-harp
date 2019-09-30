package graphql

import (
	"github.com/graphql-go/graphql"
	"hcc/harp/logger"
	"hcc/harp/mysql"
	"hcc/harp/types"
)

var queryTypes = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			////////////////////////////// volume ///////////////////////////////
			/* Get (read) single volume by uuid
			   http://localhost:8001/graphql?query={volume(uuid:"[volume_uuid]]"){uuid,size,type,server_uuid}}
			*/
			"volume": &graphql.Field{
				Type:        volumeType,
				Description: "Get volume by uuid",
				Args: graphql.FieldConfigArgument{
					"uuid": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: volume")

					requestedUUID, ok := p.Args["uuid"].(string)
					if ok {
						volume := new(types.Volume)

						var uuid string
						var size int
						var _type string
						var serverUUID string

						sql := "select * from volume where uuid = ?"
						err := mysql.Db.QueryRow(sql, requestedUUID).Scan(&uuid, &size, &_type, &serverUUID)
						if err != nil {
							logger.Logger.Println(err)
							return nil, nil
						}

						volume.UUID = uuid
						volume.Size = size
						volume.Type = _type
						volume.ServerUUID = serverUUID

						return volume, nil
					}
					return nil, nil
				},
			},

			/* Get (read) volume list
			   http://localhost:8001/graphql?query={list_volume{uuid,size,type,server_uuid}}
			*/
			"list_volume": &graphql.Field{
				Type:        graphql.NewList(volumeType),
				Description: "Get volume list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: list_volume")

					var volumes []types.Volume
					var uuid string
					var size int
					var _type string
					var serverUUID string

					sql := "select * from volume"
					stmt, err := mysql.Db.Query(sql)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}
					defer stmt.Close()

					for stmt.Next() {
						err := stmt.Scan(&uuid, &size, &_type, &serverUUID)
						if err != nil {
							logger.Logger.Println(err)
						}

						volume := types.Volume{UUID: uuid, Size: size, Type: _type, ServerUUID: serverUUID}

						logger.Logger.Println(volume)
						volumes = append(volumes, volume)
					}

					return volumes, nil
				},
			},
		},
	})
