#args = $(foreach a,$($(subst -,_,$1)_args),$(if $(value $a),$a="$($a)"))
#
#env_args = production
#rule2_args = version name
#
#TASKS = env rule2
#
#.PHONY: $(TASKS)
#$(TASKS):
#	@go run cli/cli.go $@ $(call args,$@)

test:
	go test ./test