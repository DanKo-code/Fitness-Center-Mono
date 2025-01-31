generate_sso:
	@protoc -I ./proto sso.proto --go_out=./gen --go-grpc_out=./gen --go_opt=Muser.proto=github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.user

generate_user:
	@protoc -I ./proto user.proto --go_out=./gen --go-grpc_out=./gen

generate_coach:
	@protoc -I ./proto coach.proto --go_out=./gen --go-grpc_out=./gen --go_opt=Mservice.proto=github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.service --go_opt=Mreview.proto=github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.review --go_opt=Muser.proto=github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.user

generate_service:
	@protoc -I ./proto service.proto --go_out=./gen --go-grpc_out=./gen

generate_abonement:
	@protoc -I ./proto abonement.proto --go_out=./gen --go-grpc_out=./gen --go_opt=Mservice.proto=github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.service

generate_review:
	@protoc -I ./proto review.proto --go_out=./gen --go-grpc_out=./gen

generate_order:
	@protoc -I ./proto order.proto --go_out=./gen --go-grpc_out=./gen --go_opt=Mabonement.proto=github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.abonement --go_opt=Mservice.proto=github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.service

generate_all: generate_sso generate_user generate_coach generate_service generate_abonement generate_review generate_order
	@echo "All proto file have been generated"
