PROTOC_IMAGE := frontend-protoc
PROTO_DIR := ../backend/proto
OUT_DIR := src/lib/proto

.PHONY: dev
dev:
	npm run dev -- --open

.PHONY: generate/proto
generate/proto:
	@mkdir -p $(OUT_DIR)
	docker build -t $(PROTOC_IMAGE) -f tools/protoc/Dockerfile tools/protoc
	docker run --rm \
		-v $(PWD)/..:/workspace \
		-w /workspace/frontend \
		$(PROTOC_IMAGE) \
		protoc \
			-I $(PROTO_DIR) \
			-I /usr/include \
			--plugin=/usr/local/lib/node_modules/ts-proto/protoc-gen-ts_proto \
			--ts_proto_out=$(OUT_DIR) \
			--ts_proto_opt=outputServices=grpc-web,env=browser,esModuleInterop=true \
			--grpc-web_out=import_style=typescript,mode=grpcwebtext:$(OUT_DIR) \
			$(PROTO_DIR)/monitor/v1/monitor.proto

.PHONY: clean
clean:
	rm -rf $(OUT_DIR)
