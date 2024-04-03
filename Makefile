
# Microservices list
SERVICES := censors comments news apigate

.PHONY: all $(SERVICES)

all: $(SERVICES)

$(SERVICES):
	@echo "Launching the service $@"
	@(start cmd /c "go run ./cmd/$@/main.go")