run:
	@echo "用法: make run DIR=hello (DIR为code下的子目录)"
	@if [ -z "$(DIR)" ]; then \
		echo "请指定 DIR 变量，如 make run DIR=hello"; \
		exit 1; \
	 else \
		cd code/$(DIR) && go run main.go; \
	fi
