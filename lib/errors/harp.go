package errors

const (
	HarpInternalInitFail      uint64 = harp + internal + initFail
	HarpInternalOperationFail uint64 = harp + internal + operationFail

	HarpGrpcRequestError uint64 = harp + grpc + requestError

	HarpSQLOperationFail uint64 = harp + sql + operationFail
	HarpSQLNoResult      uint64 = harp + sql + noResult
	HarpSQLArgumentError uint64 = harp + sql + argumentError
)
