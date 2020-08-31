package errors

const (
	ClarinetInternalInitFail       uint64 = clarinet + internal + initFail
	ClarinetInternalConnectionFail uint64 = clarinet + internal + connectionFail
	ClarinetInternalParsingError   uint64 = clarinet + internal + parsingError

	ClarinetDriverRequestError       uint64 = clarinet + driver + requestError
	ClarinetDriverResponseError      uint64 = clarinet + driver + responseError
	ClarinetDriverReceiveError       uint64 = clarinet + driver + receiveError
	ClarinetDriverParsingError       uint64 = clarinet + driver + parsingError
	ClarinetDriverJsonUnmarshalError uint64 = clarinet + driver + jsonUnmarshalError

	ClarinetGraphQLArgumentError      uint64 = clarinet + graphql + argumentError
	ClarinetGraphQLRequestError       uint64 = clarinet + graphql + requestError
	ClarinetGraphQLSendError          uint64 = clarinet + graphql + sendError
	ClarinetGraphQLJsonUnmarshalError uint64 = clarinet + graphql + jsonUnmarshalError
	ClarinetGraphQLParsingError       uint64 = clarinet + graphql + parsingError
)
