SHELL := /bin/bash -euo pipefail

config_file_name = Kube-manifest.yaml
example_config = docs/$(config_file_name)

.PHONY: all
all: docs

# In a branch, run 'make docs' to update docs with
# generated code, then merge it to master.
docs: $(example_config)

# Use kinflate to create the standard kinflate configuration
# file that appears in the website's documentation.
$(example_config): /tmp/bin/kinflate
	rm -f TMP
	echo "# This is a generated example; do not edit.  Rebuild with 'make docs'." >> TMP
	echo " " >> TMP
	/tmp/bin/kinflate init
	cat $(config_file_name) >> TMP
	mv TMP $(example_config)
	rm $(config_file_name)

/tmp/bin/kinflate:
	go build -o /tmp/bin/kinflate kinflate.go
