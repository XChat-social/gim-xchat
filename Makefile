PROTO_DIR := ./pkg/protocol/proto  # 定义 .proto 文件所在目录
OUT_DIR := ../                      # 定义输出目录

# 查找所有的 .proto 文件
PROTO_FILES := $(shell find $(PROTO_DIR) -type f -name "*.proto")

# 默认目标：生成所有 proto 文件的代码
all: generate-protos

# 生成目标
generate-protos: $(PROTO_FILES)
	@echo "Generating code for all .proto files..."
	@for proto_file in $(PROTO_FILES); do \
		echo "Processing $$proto_file..."; \
		protoc -I $(PROTO_DIR) \
		       --go_out=$(OUT_DIR) \
		       --go-grpc_out=$(OUT_DIR) \
		       $$proto_file; \
	done
	@echo "All .proto files have been processed."

# 清理生成的代码（可选）
clean:
	@echo "Cleaning generated files..."
	@find $(OUT_DIR) -type f -name "*.pb.go" -delete
	@echo "Cleanup completed."

# 帮助信息
help:
	@echo "Usage:"
	@echo "  make            - Generate Go and gRPC code for all .proto files"
	@echo "  make clean      - Remove all generated files"
	@echo "  make help       - Show this help message"