generate-swagger:
	oapi-codegen --package=rest --generate types openapi.yml > ./transportlayer/rest/models.gen.go
	oapi-codegen --package=rest --generate server openapi.yml > ./transportlayer/rest/server.gen.go
	oapi-codegen --package=rest --generate spec openapi.yml > ./transportlayer/rest/spec.gen.go